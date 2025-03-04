package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddUserMigrations 사용자 관련 마이그레이션 추가
func AddUserMigrations(mg *migrator.Migrator) {
	usersTable := schema.Table{
		Name: "users",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "name", Type: schema.DB_NVarchar, Length: 122, Nullable: false},
			{Name: "nick_name", Type: schema.DB_NVarchar, Length: 122, Nullable: false},
			{Name: "email", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "password", Type: schema.DB_Text, Nullable: false},
			{Name: "role", Type: schema.DB_NVarchar, Length: 5, Nullable: false},
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
			{Name: "updated_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"email"}, Type: schema.UniqueIndex},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(usersTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(usersTable, *usersTable.Indices[0]))
}
