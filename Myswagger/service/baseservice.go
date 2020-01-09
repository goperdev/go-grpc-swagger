package services

import (
	"database/sql"
)

type BaseService struct {
}



func (b *BaseService) GetSqlConn() *sql.DB {
	mySqlStr :="root:Rp000000@(test.db.rpdns.com:3306)/educationcrm?charset=utf8"
	var err error
	engine, err := sql.Open("mysql", mySqlStr)
	if err != nil {
		panic(err)
	}
	return engine
}


//设置线上支付地址
func (b *BaseService) GetPayUrl(GroupId string) string{
	return "https://h5.petrvet.com/wx/index.html?OrderId=" + GroupId //同步回调地址
}



