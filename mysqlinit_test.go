package gormutils

import (
	"fmt"
	"testing"
)

func TestNewDB(t *testing.T) {
	a := Address{
		Address:  "127.0.0.1:3306",
		User:     "root",
		Password: "123456",
		DBName:   "testinit",
		//DialTimeout:     30,
		//ReadTimeout:     30,
		//WriteTimeout:    30,
		//MaxIdleConns:    10,
		//MaxOpenConns:    20,
		//MaxLifeTime:     300,
		//ConnMaxIdleTime: 100,
	}

	db, err := NewDB(&a)
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	err = db.Create(&User{Name: "tu", Age: 18}).Error
	if err != nil {
		panic(err)
	}

	var us []*User
	err = db.Find(&us).Error
	if err != nil {
		panic(err)
	}
	for _, u := range us {
		fmt.Println(u)
	}
}

type User struct {
	ID   int
	Name string
	Age  int
}
