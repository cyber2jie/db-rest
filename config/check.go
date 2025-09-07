package config

import (
	"db-rest/util"
	"errors"
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"path"
	"strings"
)

func IsSelectStatement(node ast.StmtNode) bool {
	//stmtLabel := ast.GetStmtLabel(node)
	//return stmtLabel == "Select"
	switch node.(type) {
	case *ast.SelectStmt:
		return true
	default:
		return false
	}
}
func parseSql(sql string) (*ast.StmtNode, error) {
	p := parser.New()
	stmts, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmts[0], nil
}
func CheckSql(sql string, useTidb bool) error {

	if useTidb {

		n, err := parseSql(sql)
		if err != nil {
			return errors.New(fmt.Sprintf("parse sql error,%v", err))
		}
		if !IsSelectStatement(*n) {
			return errors.New(fmt.Sprintf("sql is not select statement,it's a %s sql。", ast.GetStmtLabel(*n)))
		}
	} else {

		sql = strings.TrimSpace(sql)

		if !strings.HasPrefix(strings.ToLower(sql), "select ") {
			return errors.New("sql is not select statement")
		}

	}
	return nil
}

var checkFiles = []string{
	WORKSPACE_CONFIG,
	WORKSPACE_DB_NAME,
}

func CheckWorkSpaceExists(p string) bool {
	if !util.FileExists(p) {
		return false
	}

	for _, file := range checkFiles {
		if !util.FileExists(path.Join(p, file)) {
			return false
		}
	}
	return true
}
