package gormutils

import (
	"fmt"

	"gorm.io/gorm"
)

type MultiTx struct {
	DBs []*gorm.DB
}

func NewMultiTx(dbs ...*gorm.DB) *MultiTx {
	return &MultiTx{DBs: dbs}
}

func (mtx *MultiTx) TransactionN(fc func(txs ...*gorm.DB) error) error {
	// 忽略这个奇怪的实现吧。。。
	if len(mtx.DBs) == 0 {
		return fmt.Errorf("no db input")
	}
	if len(mtx.DBs) == 1 {
		return mtx.DBs[0].Transaction(func(tx *gorm.DB) error {
			return fc(tx)
		})
	}
	if len(mtx.DBs) == 2 {
		db0, db1 := mtx.DBs[0], mtx.DBs[1]
		return db0.Transaction(func(tx0 *gorm.DB) error {
			return db1.Transaction(func(tx1 *gorm.DB) error {
				return fc(tx0, tx1)
			})
		})
	}
	if len(mtx.DBs) == 3 {
		db0, db1, db2 := mtx.DBs[0], mtx.DBs[1], mtx.DBs[2]
		return db0.Transaction(func(tx0 *gorm.DB) error {
			return db1.Transaction(func(tx1 *gorm.DB) error {
				return db2.Transaction(func(tx2 *gorm.DB) error {
					return fc(tx0, tx1, tx2)
				})
			})
		})
	}
	return fmt.Errorf("db number > 3 not support")
}
