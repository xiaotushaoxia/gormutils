# 一些gorm的小工具

## WhereStatementOfDB
db.Where设置后，从db里拿出where语句
```go
package gormutils

import (
	"fmt"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestWhereStatementOfDB(t *testing.T) {
	var db *gorm.DB

	fmt.Println(WhereStatementOfDB(db)) // "", [], <nil>

	db2 := db.Where("id>?", 10).Where("age<20").Where("score>60").Where("score<?", 100)
	fmt.Println(WhereStatementOfDB(db2)) // WHERE id>? AND age<20 AND score>60 AND score<?, [10 100], <nil>
	fmt.Println(WhereStatementOfDB(db))  // "", [], <nil>

	db3 := db.Where("id>?", 1000)
	fmt.Println(WhereStatementOfDB(db3)) // WHERE id>?,[1000],<nil>

	db4 := db.Table("xsdsd")
	fmt.Println(WhereStatementOfDB(db4)) //"", [], <nil>
}

func TestDemo(t *testing.T) {
	// 为什么要写这么个函数？ 用于避免db.Raw的时候需要手动拼接sql
	//type Student struct {
	//	ID    int
	//	Score int
	//	Age   int
	//  Sex   int
	//}
	// 如：我有一个student表如上，需要查询每个分数的学生的数量
	// 可按条件查，比如年纪大于10的男学生每个分数的学生的数量。
	// 因为用了Raw所以db.Where没用生效，然后我就**要手动拼sql**了（这里还有好一些的办法吗，感觉我这个做法也憨憨的）
	var db *gorm.DB

	// 查询参数
	age := ""
	sex := ""
	var result = make([]map[string]any, 0) //查询结果

	where := ""
	var conds []string
	if age != "" {
		conds = append(conds, "agv>"+age)
	}
	if sex != "" {
		conds = append(conds, "sex="+sex)
	}
	if len(conds) > 0 {
		where = "where" + strings.Join(conds, " and ")
	}
	sqlStr := "select `score`,count(*) from student %s group by `score`;"
	err := db.Raw(fmt.Sprintf(sqlStr, where)).Scan(&result).Error
	if err != nil {
		panic(err)
	}

	// 有这个方法可能会好一点
	dd := db
	if age != "" {
		dd = dd.Where("age>?", age) // 还可以用这个?,不用手动格式化字符串了
	}
	if sex != "" {
		dd = dd.Where("sex=?", sex)
	}

	where, args, err := WhereStatementOfDB(dd)
	if err != nil {
		panic(err)
	}
	err = db.Raw(fmt.Sprintf(sqlStr, where), args...).Scan(&result).Error
	if err != nil {
		panic(err)
	}
}
```


## MultiTx
多个数据库的事务操作，让代码好看一点点
```go
package gormutils

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewMultiTx(t *testing.T) {
	var db1, db2 *gorm.DB
	// 当在多个数据库上的操作要保持事务时，我之前写的代码比较丑，所以写了这个MultiTx，让程序稍微好看一点
	// 旧
	err := db1.Transaction(func(tx1 *gorm.DB) error {
		return db2.Transaction(func(tx2 *gorm.DB) error {
			er := tx1.Updates(nil).Error
			if er != nil {
				return er
			}
			er = tx2.Updates(nil).Error
			if er != nil {
				return er
			}
			return nil
		})
	})
	if err != nil {
		panic(err)
	}

	// MultiTx
	mtx := NewMultiTx(db1, db2)
	err = mtx.TransactionN(func(txs ...*gorm.DB) error {
		tx1, tx2 := txs[0], txs[1]
		er := tx1.Updates(nil).Error
		if er != nil {
			return er
		}
		er = tx2.Updates(nil).Error
		if er != nil {
			return er
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
```
