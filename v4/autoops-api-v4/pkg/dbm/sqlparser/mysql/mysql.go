package mysql

import (
	"dodevops-api/pkg/dbm/sqlparser/base"
	mysqlparser "dodevops-api/pkg/dbm/sqlparser/mysql/antlr4"
	"dodevops-api/pkg/dbm/sqlparser/sqlstmt"
	"dodevops-api/pkg/logx"

	"github.com/antlr4-go/antlr/v4"
)

func GetMysqlParserTree(baseLine int, statement string) (antlr.ParseTree, *antlr.CommonTokenStream, error) {
	lexer := mysqlparser.NewMySqlLexer(antlr.NewInputStream(statement))
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := mysqlparser.NewMySqlParser(stream)

	lexerErrorListener := &base.ParseErrorListener{
		BaseLine: baseLine,
	}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrorListener)

	parserErrorListener := &base.ParseErrorListener{
		BaseLine: baseLine,
	}
	parser.RemoveErrorListeners()
	parser.AddErrorListener(parserErrorListener)
	parser.BuildParseTrees = true

	tree := parser.Root()

	if lexerErrorListener.Err != nil {
		return nil, nil, lexerErrorListener.Err
	}

	if parserErrorListener.Err != nil {
		return nil, nil, parserErrorListener.Err
	}

	return tree, stream, nil
}

type MysqlParser struct {
}

func (*MysqlParser) Parse(stmt string) (stmts []sqlstmt.Stmt, err error) {
	defer func() {
		if e := recover(); e != nil {
			logx.ErrorTrace("mysql sql parser err: ", e)
			err = e.(error)
		}
	}()
	tree, _, err := GetMysqlParserTree(1, stmt)
	if err != nil {
		return nil, err
	}

	return tree.Accept(new(MysqlVisitor)).([]sqlstmt.Stmt), nil
}
