package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddAnonymousPostImagesMigrations 익명 게시판 이미지 관련 마이그레이션 추가
func AddAnonymousPostImagesMigrations(mg *migrator.Migrator) {
	anonymousPostImagesTable := schema.Table{
		Name: "anonymous_post_images",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "post_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "image_url", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "image_type", Type: schema.DB_BigInt, Nullable: false},
			{Name: "original_filename", Type: schema.DB_Text, Nullable: false},
			{Name: "file_size", Type: schema.DB_BigInt, Nullable: true},
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
			{Name: "updated_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"post_id"}, Type: schema.IndexNormal},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(anonymousPostImagesTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(anonymousPostImagesTable, *anonymousPostImagesTable.Indices[0]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_anonymous_post_images_post_id",
		"ALTER TABLE anonymous_post_images ADD CONSTRAINT fk_anonymous_post_images_post_id FOREIGN KEY (post_id) REFERENCES anonymous_posts(uid) ON DELETE CASCADE;",
	))
}
