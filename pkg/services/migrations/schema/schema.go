package schema

// 컬럼 타입 정의
type ColumnType string

const (
	DB_BigInt    ColumnType = "BIGINT"
	DB_Int       ColumnType = "INTEGER"
	DB_NVarchar  ColumnType = "VARCHAR"
	DB_Text      ColumnType = "TEXT"
	DB_Bool      ColumnType = "BOOLEAN"
	DB_DateTime  ColumnType = "TIMESTAMP WITH TIME ZONE"
	DB_TimeStamp ColumnType = "TIMESTAMP"
)

// 인덱스 타입 정의
type IndexType string

const (
	UniqueIndex IndexType = "UNIQUE"
	IndexNormal IndexType = "INDEX"
)

// 테이블 컬럼 정의
type Column struct {
	Name            string
	Type            ColumnType
	Length          int
	Nullable        bool
	IsPrimaryKey    bool
	IsAutoIncrement bool
	Default         interface{}
}

// 인덱스 정의
type Index struct {
	Cols []string
	Type IndexType
}

// 테이블 정의
type Table struct {
	Name    string
	Columns []*Column
	Indices []*Index
}
