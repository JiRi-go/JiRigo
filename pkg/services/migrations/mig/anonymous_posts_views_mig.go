package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddAnonymousPostViewsMigrations 익명 게시판 조회수 관련 마이그레이션 추가
func AddAnonymousPostViewsMigrations(mg *migrator.Migrator) {
	anonymousPostViewsTable := schema.Table{
		Name: "anonymous_post_views",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "post_id", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "user_id", Type: schema.DB_NVarchar, Length: 50, Nullable: true},
			{Name: "viewed_at", Type: schema.DB_DateTime, Nullable: false},
		},
		Indices: []*schema.Index{
			{Cols: []string{"post_id"}, Type: schema.IndexNormal},
			{Cols: []string{"user_id"}, Type: schema.IndexNormal},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(anonymousPostViewsTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(anonymousPostViewsTable, *anonymousPostViewsTable.Indices[0]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_anonymous_post_views_post_id",
		"ALTER TABLE anonymous_post_views ADD CONSTRAINT fk_anonymous_post_views_post_id FOREIGN KEY (post_id) REFERENCES anonymous_posts(uid) ON DELETE CASCADE;",
	))

	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_anonymous_post_views_user_id",
		"ALTER TABLE anonymous_post_views ADD CONSTRAINT fk_anonymous_post_views_user_id FOREIGN KEY (user_id) REFERENCES users(uid) ON DELETE SET NULL;",
	))
}
