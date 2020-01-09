package services

import (
	"database/sql"
)

type BaseService struct {
}



func (b *BaseService) GetSqlConn() *sql.DB {
	mySqlStr :="root:123456@(127.0.0.1:3306)/mysqlname?charset=utf8"
	var err error
	engine, err := sql.Open("mysql", mySqlStr)
	if err != nil {
		panic(err)
	}
	return engine
}





