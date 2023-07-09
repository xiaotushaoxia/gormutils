package gormutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/xiaotushaoxia/kvtag"
	"gorm.io/gorm"
)

const (
	queryConditionTagName = "cond"
	queryConditionTagSep  = ";"

	//"colum:user;opr:in;split:sting;sep:,"`
	queryConditionColum   = "colum"
	queryConditionOpr     = "opr"
	queryConditionSplit   = "split"
	queryConditionSep     = "sep"
	queryConditionPattern = "pattern"

	// pattern: 会把结构体的值替换pattern字符串中的?
	// pattern:?%
	// 如 a := Query{Name: "xx"}  ->  where name like 'xx%'

	// split 支持的字符串
	_queryConditionString = "string"
	_queryConditionInt    = "int"
)

//      LIKE Operator	                      Description
// WHERE CustomerName LIKE 'a%'	      Finds any values that start with "a"
// WHERE CustomerName LIKE '%a'	      Finds any values that end with "a"
// WHERE CustomerName LIKE '%or%'	  Finds any values that have "or" in any position
// WHERE CustomerName LIKE '_r%'	  Finds any values that have "r" in the second position
// WHERE CustomerName LIKE 'a_%'	  Finds any values that start with "a" and are at least 2 characters in length
// WHERE CustomerName LIKE 'a__%'	  Finds any values that start with "a" and are at least 3 characters in length
// WHERE ContactName LIKE 'a%o'	      Finds any values that start with "a" and ends with "o"   不支持！！

func WithCondition(cond any) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tagInfo := kvtag.ParserTag(cond, queryConditionTagName, queryConditionTagSep)
		v := reflect.Indirect(reflect.ValueOf(cond))
		// 遍历结构体
		for i := 0; i < v.NumField(); i++ {
			// 忽略掉没有tag的字段和零值字段
			ft, hasFieldTag := tagInfo.GetFieldTagByIndex(i)
			if !hasFieldTag {
				continue
			}
			fv := v.Field(i)
			if fv.IsZero() {
				continue
			}

			arg := fv.Interface()
			st := ft.TagSetting
			if isInOpr(st[queryConditionOpr]) {
				if argS, isStr := arg.(string); isStr {
					arg = trySplitInString(argS, st[queryConditionSplit], st[queryConditionSep])
				}
				if k := reflect.ValueOf(arg).Kind(); !(k == reflect.Slice || k == reflect.Array) {
					continue
				}
			} else if isLikeOpr(st[queryConditionOpr]) {
				pattern_ := st[queryConditionPattern]
				if pattern_ == "" {
					continue
				}
				arg = strings.ReplaceAll(pattern_, "?", fmt.Sprintf("%v", arg))
			}
			db = db.Where(fmt.Sprintf("`%s` %s ?", st[queryConditionColum], st[queryConditionOpr]), arg)
		}
		return db
	}

}

// trySplitInString 尝试把 1,2,3 分割成 [1,2,3], 看情况返回[]string或者[]int
func trySplitInString(v string, to string, sep string) any {
	if to == "" {
		to = _queryConditionString
	}
	if sep == "" {
		sep = ","
	}

	ss := strings.Split(v, sep)
	if to == _queryConditionString {
		return ss
	}
	if to == _queryConditionInt {
		var s []int
		for _, i2 := range ss {
			i, err := strconv.Atoi(i2)
			if err != nil {
				continue
			}
			s = append(s, i)
		}
		if len(s) == 0 {
			return nil
		}
		return s
	}
	return nil
}

func isLikeOpr(s string) bool {
	return strings.ToUpper(s) == "LIKE"
}

func isInOpr(s string) bool {
	upper := strings.ToUpper(s)
	return upper == "IN" || upper == "NOT IN"
}
