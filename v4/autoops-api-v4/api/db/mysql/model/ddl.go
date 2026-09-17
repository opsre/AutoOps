package model

// CreateDatabaseReq 创建数据库请求
type CreateDatabaseReq struct {
    DbName      string `json:"dbName" binding:"required"`
    Charset     string `json:"charset"`
    Collate     string `json:"collate"`
    IfNotExists bool   `json:"ifNotExists"`
}

// CreateTable 请求
type CreateTableColumn struct {
    Name          string  `json:"name" binding:"required"`
    Type          string  `json:"type" binding:"required"`
    Nullable      bool    `json:"nullable"`
    Default       *string `json:"default"`
    Comment       string  `json:"comment"`
    AutoIncrement bool    `json:"autoIncrement"`
    PrimaryKey    bool    `json:"pk"`
}

type CreateTableIndex struct {
    Name    string   `json:"name"`
    Columns []string `json:"columns" binding:"required"`
    Unique  bool     `json:"unique"`
}

type CreateTableReq struct {
    Db        string               `json:"db" binding:"required"`
    TableName string               `json:"tableName" binding:"required"`
    Engine    string               `json:"engine"`
    Charset   string               `json:"charset"`
    Columns   []CreateTableColumn  `json:"columns" binding:"required"`
    Indexes   []CreateTableIndex   `json:"indexes"`
}

