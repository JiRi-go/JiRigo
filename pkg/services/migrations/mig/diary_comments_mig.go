package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddDiaryCommentsMigrations 그림일기 댓글 관련 마이그레이션 추가
func AddDiaryCommentsMigrations(mg *migrator.Migrator) {
	diaryCommentsTable := schema.Table{
		Name: "diary_comments",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "diary_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "user_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "parent_id", Type: schema.DB_NVarchar, Length: 50, Nullable: true},
			{Name: "content", Type: schema.DB_Text, Nullable: false},
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
			{Name: "updated_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"diary_id"}, Type: schema.IndexNormal},
			{Cols: []string{"user_id"}, Type: schema.IndexNormal},
			{Cols: []string{"parent_id"}, Type: schema.IndexNormal}, 
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(diaryCommentsTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(diaryCommentsTable, *diaryCommentsTable.Indices[0]))
	mg.AddMigration(migrator.NewAddIndexMigration(diaryCommentsTable, *diaryCommentsTable.Indices[1]))
	mg.AddMigration(migrator.NewAddIndexMigration(diaryCommentsTable, *diaryCommentsTable.Indices[2]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diary_comments_diary_id",
		"ALTER TABLE diary_comments ADD CONSTRAINT fk_diary_comments_diary_id FOREIGN KEY (diary_id) REFERENCES diaries(uid) ON DELETE CASCADE;",
	))

	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diary_comments_user_id",
		"ALTER TABLE diary_comments ADD CONSTRAINT fk_diary_comments_user_id FOREIGN KEY (user_id) REFERENCES users(uid) ON DELETE CASCADE;",
	))

	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diary_comments_parent_id",
		"ALTER TABLE diary_comments ADD CONSTRAINT fk_diary_comments_parent_id FOREIGN KEY (parent_id) REFERENCES diary_comments(uid) ON DELETE SET NULL;",
	))
}
