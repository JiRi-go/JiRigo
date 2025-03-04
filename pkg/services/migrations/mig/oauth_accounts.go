package mig

import (
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
	"github.com/jslee/JiRigo/pkg/services/migrations/schema"
)

// AddOauthAccountsMigrations 인증 관련 마이그레이션 추가
func AddOauthAccountsMigrations(mg *migrator.Migrator) {
	oauthAccountsTable := schema.Table{
		Name: "oauth_accounts",
		Columns: []*schema.Column{
			{Name: "uid", Type: schema.DB_NVarchar, Length: 50, Nullable: false, IsPrimaryKey: true},
			{Name: "user_id", Type: schema.DB_NVarchar, Length: 50, Nullable: false},
			{Name: "provider", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "provider_id", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "provider_email", Type: schema.DB_NVarchar, Length: 255, Nullable: false},
			{Name: "access_token", Type: schema.DB_Text, Nullable: false},
			{Name: "refresh_token", Type: schema.DB_Text, Nullable: false},
			{Name: "created_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
			{Name: "updated_at", Type: schema.DB_DateTime, Nullable: false, Default: "NOW()"},
		},
		Indices: []*schema.Index{
			{Cols: []string{"provider_id"}, Type: schema.UniqueIndex},
		},
	}

	// 테이블 생성 마이그레이션 추가
	mg.AddMigration(migrator.NewAddTableMigration(oauthAccountsTable))

	// 인덱스 추가 마이그레이션 추가
	mg.AddMigration(migrator.NewAddIndexMigration(oauthAccountsTable, *oauthAccountsTable.Indices[0]))

	// 외래 키 추가 마이그레이션
	mg.AddMigration(migrator.NewRawSQLMigration(
		"add_foreign_key_oauth_accounts_user_id",
		"ALTER TABLE oauth_accounts ADD CONSTRAINT fk_oauth_accounts_user_id FOREIGN KEY (user_id) REFERENCES users(uid) ON DELETE CASCADE;",
	))
}
