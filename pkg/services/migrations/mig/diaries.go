package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddDiariesMigrations 그림일기 관련 마이그레이션 추가
func AddDiariesMigrations(mg *migrator.Migrator) {
	diariesTable := schema.Table{  
		Name: "diaries",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "user_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "title", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "content", Type: schema.DB_Text, Nullable: false},
			{Name: "is_private", Type: schema.DB_Bool, Nullable: false},
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
			{Name: "updated_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"user_id"}, Type: schema.IndexNormal},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(diariesTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(diariesTable, *diariesTable.Indices[0]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diaries_user_id",
		"ALTER TABLE diaries ADD CONSTRAINT fk_diaries_user_id FOREIGN KEY (user_id) REFERENCES users(uid) ON DELETE CASCADE;",
	))
}
