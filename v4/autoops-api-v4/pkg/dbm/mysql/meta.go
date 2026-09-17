package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"dodevops-api/pkg/dbm/dbi"
)

func init() {
	meta := new(Meta)
	dbi.Register(DbTypeMysql, meta)
	dbi.Register(DbTypeMariadb, meta)
}

const (
	DbTypeMysql   dbi.DbType = "mysql"
	DbTypeMariadb dbi.DbType = "mariadb"
)

type Meta struct {
}

func (mm *Meta) GetSqlDb(ctx context.Context, d *dbi.DbInfo) (*sql.DB, error) {
	// SSH tunnel functionality disabled
	if d.SshTunnelMachineId > 0 {
		return nil, fmt.Errorf("SSH tunnel not supported")
	}
	// 设置dataSourceName  -> 更多参数参考：https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true&timeout=8s", d.Username, d.Password, d.Network, d.Host, d.Port, d.Database)
	if d.Params != "" {
		dsn = fmt.Sprintf("%s&%s", dsn, d.Params)
	}
	const driverName = "mysql"
	return sql.Open(driverName, dsn)
}

func (mm *Meta) GetDialect(conn *dbi.DbConn) dbi.Dialect {
	return &MysqlDialect{dc: conn}
}

func (mm *Meta) GetMetadata(conn *dbi.DbConn) dbi.Metadata {
	return &MysqlMetadata{dc: conn}
}

func (mm *Meta) GetDbDataTypes() []*dbi.DbDataType {
	return []*dbi.DbDataType{
		UnsignedBigint, Bigint, Tinyint, Smallint, Int, Bit, Float, Double, Decimal,
		Varchar, Char, Text, Longtext, Mediumtext,
		Datetime, Date, Time, Timestamp,
		Enum, JSON, Set,
		Binary, Blob, Longblob, Mediumblob, Varbinary,
	}
}

func (mm *Meta) GetCommonTypeConverter() dbi.CommonTypeConverter {
	return &commonTypeConverter{}
}
