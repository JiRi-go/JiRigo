package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddDiaryLikesMigrations 그림일기 관련 마이그레이션 추가
func AddDiaryLikesMigrations(mg *migrator.Migrator) {
	diaryLikesTable := schema.Table{  
		Name: "diary_likes",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "diary_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "user_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
			{Name: "updated_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"diary_id"}, Type: schema.IndexNormal}, 
			{Cols: []string{"user_id"}, Type: schema.IndexNormal},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(diaryLikesTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(diaryLikesTable, *diaryLikesTable.Indices[0]))
	mg.AddMigration(migrator.NewAddIndexMigration(diaryLikesTable, *diaryLikesTable.Indices[1]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diary_likes_diary_id",
		"ALTER TABLE diary_likes ADD CONSTRAINT fk_diary_likes_diary_id FOREIGN KEY (diary_id) REFERENCES diaries(uid) ON DELETE CASCADE;",
		))

	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diary_likes_user_id",
		"ALTER TABLE diary_likes ADD CONSTRAINT fk_diary_likes_user_id FOREIGN KEY (user_id) REFERENCES users(uid) ON DELETE CASCADE;",	
		))
}
	