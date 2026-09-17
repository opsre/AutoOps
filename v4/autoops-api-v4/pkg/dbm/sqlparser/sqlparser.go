package sqlparser

import (
	"dodevops-api/pkg/dbm/sqlparser/sqlstmt"
)

type DbDialect string

type SqlParser interface {

	// sql解析
	// @param stmt sql语句
	Parse(stmt string) ([]sqlstmt.Stmt, error)
}

// SQLSplit sql切割 - disabled
// func SQLSplit(r io.Reader, delimiter rune, callback utils.StmtCallback) error {
// 	return utils.SplitStmts(r, delimiter, callback)
// }
