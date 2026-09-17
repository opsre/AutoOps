package model

// TableInfo 表信息
type TableInfo struct {
    TableName    string `json:"tableName"`
    TableComment string `json:"tableComment"`
    TableType    string `json:"tableType"`
    CreateTime   string `json:"createTime"`
    UpdateTime   string `json:"updateTime"`
    TableRows    int64  `json:"tableRows"`
    DataLength   int64  `json:"dataLength"`
    IndexLength  int64  `json:"indexLength"`
    Engine       string `json:"engine"`
    Collation    string `json:"collation"`
    TotalLength  int64  `json:"totalLength"`   // DataLength + IndexLength
    Size         string `json:"size"`          // 可读化大小，如 "1.5 MB"
}

// ColumnInfo 列信息
type ColumnInfo struct {
	ColumnName    string `json:"columnName"`
	ColumnType    string `json:"columnType"`
	ColumnComment string `json:"columnComment"`
	DataType      string `json:"dataType"`
	IsNullable    string `json:"isNullable"`
	ColumnKey     string `json:"columnKey"`
	Extra         string `json:"extra"`
	DefaultValue  string `json:"defaultValue"`
	CharMaxLength int64  `json:"charMaxLength"`
	NumPrecision  int64  `json:"numPrecision"`
	NumScale      int64  `json:"numScale"`
}

// IndexInfo 索引信息
type IndexInfo struct {
	IndexName    string `json:"indexName"`
	ColumnName   string `json:"columnName"`
	IndexType    string `json:"indexType"`
	SeqInIndex   int    `json:"seqInIndex"`
	NonUnique    int    `json:"nonUnique"`
	IndexComment string `json:"indexComment"`
}

// DbServer 数据库服务器信息
type DbServer struct {
	Version string `json:"version"`
	Extra   string `json:"extra"`
}

// CopyTableDto 复制表DTO
type CopyTableDto struct {
	DbId      uint   `json:"dbId" validate:"required"`
	Db        string `json:"db" validate:"required"`
	TableName string `json:"tableName" validate:"required"`
	NewTable  string `json:"newTable" validate:"required"`
	CopyData  bool   `json:"copyData"`
}

// TableQuery 表查询条件
type TableQuery struct {
	DbId       uint   `json:"dbId"`
	Db         string `json:"db"`
	TableName  string `json:"tableName"`
	TableType  string `json:"tableType"`
}

// HintTable 提示表(用于SQL编辑器自动完成)
type HintTable struct {
	TableName string       `json:"tableName"`
	Columns   []HintColumn `json:"columns"`
}

// HintColumn 提示列
type HintColumn struct {
    ColumnName string `json:"columnName"`
    ColumnType string `json:"columnType"`
    ColumnComment string `json:"columnComment"`
}

// TableStatus 表状态信息(用于快速统计)
type TableStatus struct {
    TableName   string `json:"tableName"`
    Rows        int64  `json:"rows"`
    DataLength  int64  `json:"dataLength"`
    IndexLength int64  `json:"indexLength"`
    Comment     string `json:"comment"`
    CreateTime  string `json:"createTime"`
}
