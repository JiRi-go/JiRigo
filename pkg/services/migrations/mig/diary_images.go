package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddDiaryImagesMigrations 그림일기 관련 마이그레이션 추가
func AddDiaryImagesMigrations(mg *migrator.Migrator) {
	diaryImagesTable := schema.Table{  
		Name: "diary_images",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "diary_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "image_url", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "image_type", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "original_filename", Type: schema.DB_Text, Nullable: false},
			{Name: "file_size", Type: schema.DB_Bool, Nullable: false}, // boolean???
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"diary_id"}, Type: schema.IndexNormal},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(diaryImagesTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(diaryImagesTable, *diaryImagesTable.Indices[0]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_diary_images_diary_id",
		"ALTER TABLE diary_images ADD CONSTRAINT fk_diary_images_diary_id FOREIGN KEY (diary_id) REFERENCES diaries(uid) ON DELETE CASCADE;",
	))
}
