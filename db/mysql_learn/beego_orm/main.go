package main

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 连接数据库
func init() {
	err := orm.RegisterDataBase("default", "mysql",
		"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4")
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	orm.RegisterModel(new(PetParentUser))

}

// PetParentUser 用户信息orm对象
type PetParentUser struct {
	// 一些字段
}

// TableName 数据库表名
func (info *PetParentUser) TableName() string {
	return "czt_user_info"
}
func main() {

	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.QueryTable("czt_user_info").Filter("user_name", "雷磊").Values(&maps)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", maps)

}
