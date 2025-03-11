package migrator

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// 마이그레이션 인터페이스
type Migration interface {
	SQL(dialect Dialect) string
	ID() string
	Name() string
}

// 데이터베이스 방언 인터페이스
type Dialect interface {
	CreateTableSQL(table schema.Table) string
	AddIndexSQL(table string, index schema.Index) string
	Quote(name string) string
}

type PostgresDialect struct{}

// 테이블 생성 SQL 생성
func (d *PostgresDialect) CreateTableSQL(table schema.Table) string {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", d.Quote(table.Name))

	// 컬럼 정의
	for i, col := range table.Columns {
		sql += "  " + d.ColumnDefinition(col)

		if i < len(table.Columns)-1 {
			sql += ",\n"
		}
	}

	sql += "\n);"
	return sql
}

// 컬럼 정의 SQL 생성
func (d *PostgresDialect) ColumnDefinition(col *schema.Column) string {
	sql := d.Quote(col.Name) + " "

	// 타입 정의
	sql += string(col.Type)
	if col.Length > 0 && col.Type == schema.DB_NVarchar {
		sql += fmt.Sprintf("(%d)", col.Length)
	}

	// 제약 조건
	if col.IsPrimaryKey {
		sql += " PRIMARY KEY"
	}

	if col.IsAutoIncrement {
		sql += " GENERATED ALWAYS AS IDENTITY"
	}

	if !col.Nullable {
		sql += " NOT NULL"
	}

	if col.Default != nil {
		defaultValue := fmt.Sprintf("%v", col.Default) // interface{} → string 변환
		sql += " DEFAULT " + defaultValue
	}

	// if col.Default != "" {
	// 	sql += " DEFAULT " + col.Default
	// }

	return sql
}

// 인덱스 추가 SQL 생성
func (d *PostgresDialect) AddIndexSQL(table string, index schema.Index) string {
	cols := ""
	for i, col := range index.Cols {
		cols += d.Quote(col)
		if i < len(index.Cols)-1 {
			cols += ", "
		}
	}

	indexType := ""
	if index.Type == schema.UniqueIndex {
		indexType = "UNIQUE "
	}

	indexName := fmt.Sprintf("idx_%s_%s", table, index.Cols[0])
	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s);",
		indexType, d.Quote(indexName), d.Quote(table), cols)
}

// 식별자 인용 처리
// ex) where uid = "uid_xxx_xxxx" 의 "" 처리
func (d *PostgresDialect) Quote(name string) string {
	return "\"" + name + "\""
}

//============= Table =============

// 테이블 생성 마이그레이션
type AddTableMigration struct {
	Table schema.Table
	id    string
	name  string
}

// 새 테이블 생성 마이그레이션 생성
func NewAddTableMigration(table schema.Table) *AddTableMigration {
	return &AddTableMigration{
		Table: table,
		id:    fmt.Sprintf("%s_create_%s_table", time.Now().Format("20060102"), table.Name),
		name:  fmt.Sprintf("create %s table", table.Name),
	}
}

// SQL 문 생성
func (m *AddTableMigration) SQL(dialect Dialect) string {
	return dialect.CreateTableSQL(m.Table)
}

// 마이그레이션 ID 반환
func (m *AddTableMigration) ID() string {
	return m.id
}

// 마이그레이션 이름 반환
func (m *AddTableMigration) Name() string {
	return m.name
}

//============= Index =============

// 인덱스 추가 마이그레이션
type AddIndexMigration struct {
	Table schema.Table
	Index schema.Index
	id    string
	name  string
}

// 새 인덱스 추가 마이그레이션 생성
func NewAddIndexMigration(table schema.Table, index schema.Index) *AddIndexMigration {
	indexName := index.Cols[0]
	return &AddIndexMigration{
		Table: table,
		Index: index,
		id:    fmt.Sprintf("%s_add_index_%s_%s", time.Now().Format("20060102"), table.Name, indexName),
		name:  fmt.Sprintf("add index %s.%s", table.Name, indexName),
	}
}

// SQL 문 생성
func (m *AddIndexMigration) SQL(dialect Dialect) string {
	return dialect.AddIndexSQL(m.Table.Name, m.Index)
}

// 마이그레이션 ID 반환
func (m *AddIndexMigration) ID() string {
	return m.id
}

// 마이그레이션 이름 반환
func (m *AddIndexMigration) Name() string {
	return m.name
}

//============= Raw SQL(FK 사용) =============

// Raw SQL 실행 마이그레이션
type RawSQLMigration struct {
	id   string
	name string
	sql  string
}

// 새 SQL 실행 마이그레이션 생성
func NewRawSQLMigration(id, sql string) *RawSQLMigration {
	return &RawSQLMigration{
		id:   id,
		name: fmt.Sprintf("Execute raw SQL: %s", id),
		sql:  sql,
	}
}

// SQL 문 반환
func (m *RawSQLMigration) SQL(dialect Dialect) string {
	return m.sql
}

// 마이그레이션 ID 반환
func (m *RawSQLMigration) ID() string {
	return m.id
}

// 마이그레이션 이름 반환
func (m *RawSQLMigration) Name() string {
	return m.name
}

//============= Migrator =============

// 마이그레이션 매니저
type Migrator struct {
	db         *sql.DB
	dialect    Dialect
	migrations []Migration
}

// NewMigrator 새 마이그레이션 매니저 생성
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:         db,
		dialect:    &PostgresDialect{},
		migrations: []Migration{},
	}
}

// AddMigration 마이그레이션 추가
func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// RunMigrations 모든 마이그레이션 실행
func (m *Migrator) RunMigrations() error {
	// 마이그레이션 테이블 생성
	err := m.createMigrationsTable()
	if err != nil {
		return err
	}

	// 각 마이그레이션 실행
	for _, migration := range m.migrations {
		// 이미 적용된 마이그레이션인지 확인
		applied, err := m.isMigrationApplied(migration.ID())
		if err != nil {
			return err
		}

		if applied {
			log.Printf("마이그레이션 '%s'는 이미 적용되었습니다.", migration.Name())
			continue
		}

		log.Printf("마이그레이션 '%s' 실행 중...", migration.Name())

		// SQL 생성 및 실행
		sql := migration.SQL(m.dialect)
		_, err = m.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("마이그레이션 '%s' 실행 실패: %w", migration.Name(), err)
		}

		// 마이그레이션 기록
		err = m.recordMigration(migration.ID())
		if err != nil {
			return err
		}

		log.Printf("마이그레이션 '%s' 성공적으로 적용되었습니다.", migration.Name())
	}

	return nil
}

// createMigrationsTable 마이그레이션 이력 테이블 생성
func (m *Migrator) createMigrationsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        id VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    )`

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("마이그레이션 테이블 생성 실패: %w", err)
	}

	return nil
}

// isMigrationApplied 마이그레이션 적용 여부 확인
func (m *Migrator) isMigrationApplied(id string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE id = $1)"

	err := m.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("마이그레이션 확인 실패: %w", err)
	}

	return exists, nil
}

// recordMigration 마이그레이션 실행 기록
func (m *Migrator) recordMigration(id string) error {
	query := "INSERT INTO schema_migrations (id) VALUES ($1)"

	_, err := m.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("마이그레이션 기록 실패: %w", err)
	}

	return nil
}
