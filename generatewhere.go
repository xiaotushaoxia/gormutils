package gormutils

import (
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

func WhereStatementOfDB(db *gorm.DB) (string, []any, error) {
	// 这里要copy一个db， 不然callbacks.BuildQuerySQL时候会添加Statement.Vars，然后最后执行的时候又会去去添加一遍
	tempDB, tempStmt := cloneDBAndStmtForWhere(db)
	defer back(tempDB, tempStmt)
	callbacks.BuildQuerySQL(tempDB)
	if tempDB.Error != nil {
		return "", nil, tempDB.Error
	}
	return tempDB.Statement.SQL.String(), tempDB.Statement.Vars, nil
}

func cloneDBAndStmtForWhere(db *gorm.DB) (*gorm.DB, *gorm.Statement) {
	// clone DB
	tempDB := tempDBPool.Get().(*gorm.DB)
	tempDB.Config = db.Config
	tempDB.Error = db.Error
	tempDB.Statement = nil

	// clone Statement
	tempStmt := tempStmtPool.Get().(*gorm.Statement)
	tempStmt.Clauses = map[string]clause.Clause{}

	if db.Statement != nil && db.Statement.Clauses != nil {
		for k, c := range db.Statement.Clauses {
			tempStmt.Clauses[k] = c
		}
	}
	tempStmt.Vars = nil
	tempStmt.SQL = strings.Builder{}
	tempStmt.BuildClauses = []string{"WHERE"}

	//这里循环引用了，有点怪。不过应该是没关系的，循环引用对现代gc来说可能已经不是问题了
	tempDB.Statement = tempStmt
	tempStmt.DB = tempDB // 某些操作会访问tempStmt.DB，所以要设置进来

	return tempDB, tempStmt
}

var tempStmtPool = sync.Pool{New: func() any {
	return &gorm.Statement{}
}}

var tempDBPool = sync.Pool{New: func() any {
	return &gorm.DB{}
}}

func back(db *gorm.DB, statement *gorm.Statement) {
	db.Statement = nil
	statement.DB = nil
	tempDBPool.Put(db)
	tempStmtPool.Put(statement)
}
