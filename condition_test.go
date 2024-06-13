package gormutils

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type areaQueryBody struct {
	nameIDQueryBody `cond:"embedded"`
	//Name  string `cond:"colum:name;opr:like;pattern:%%?%%"`
	//ID    []int  `cond:"colum:id;opr:in"`
	Type  string `cond:"colum:type;op:in;split:int;sep:,"`
	Start string `cond:"colum:created_at;opr:>"`
	Desc  string `cond:"colum:desc;opr:in"` // split默认值string,sep默认值,
}

type locationQueryBody struct {
	nameIDQueryBody `cond:"embedded"`
	//Name  string `cond:"colum:name;opr:like;pattern:%%?%%"`
	//ID    []int  `cond:"colum:id;opr:in"`
	Type  string `cond:"colum:type;opr:in;split:int;sep:,"`
	Start string `cond:"colum:created_at;opr:>"`
	Desc  string `cond:"colum:desc;opr:in"` // split默认值string,sep默认值,
}

type nameIDQueryBody struct {
	Name string `cond:"colum:name;opr:like;patten:%%?%%"`
	ID   []int  `cond:"colum:id;opr:in"`
}

func TestWithConditionUsage(t *testing.T) {
	withConditionUsage()
}

func withConditionUsage() {
	aa := areaQueryBody{
		nameIDQueryBody: nameIDQueryBody{
			Name: "区",
			ID:   []int{1, 2, 3},
		},
		Start: "2023-05-08 11:03:23",
		Type:  "1,3,4",
		Desc:  "缓存区1,缓存区2",
	}
	// SELECT * FROM `area`
	// WHERE `type` in (1,3,4) AND `created_at` > '2023-05-08 11:03:23' AND `desc` in ('缓存区1','缓存区2') AND `name` like '%%区%%' AND `id` in (1,2,3)

	db := getDB().Debug()

	db = db.Table("area")

	// 或者  db = db.Scopes(WithCondition(aa))，
	// Scopes是惰性的，到最后find的时候才会执行。因为我这里想打印一下where语句，所以用了db = WithCondition(aa)(db)
	db = WithCondition(aa)(db)
	fmt.Println(WhereStatementOfDB(db))
	//db = db.Scopes(WithCondition(aa))

	var result []map[string]any
	err := db.Find(&result).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("xxx")

	for _, m := range result {
		fmt.Println(m)
	}

	// 旧的方式可能需要 有人有更好的办法 谢谢分享
	//if aa.Name != "" {
	//	db = db.Where("name", aa.Name)
	//}
	//if len(aa.ID) > 0 {
	//	db = db.Where("id in ?", aa.ID)
	//}
	//if aa.Type != "" {
	//	db = db.Where("type in ?", trySplitInString(aa.Type, _queryConditionInt, ","))
	//}
	//...
}

func getDB() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/testquery?charset=utf8mb4&parseTime=True&loc=Local"
	vv, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return vv
}
