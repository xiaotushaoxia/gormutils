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
