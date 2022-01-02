package common

import (
	"fmt"
	"github.com/casbin/casbin/v2"

	sqladapter "github.com/Blank-Xu/sql-adapter"
)

var (
	A *sqladapter.Adapter
	E *casbin.Enforcer
)

//InitCasbin
//@title		InitCasbin()
//@description	初始化RBAC
//@author		zy
//@param
//@return		*casbin.Enforcer, error
func InitCasbin() (*casbin.Enforcer, error){
	var err error
	A, err = sqladapter.NewAdapter(DB, "mysql", "casbin_rule")
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return nil, err
	}
	E, err = casbin.NewEnforcer("./model.conf", A)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return nil, err
	}
	if err = E.LoadPolicy(); err != nil {
		fmt.Printf("LoadPolicy failed, err:%v\n", err)
	}
	added, _ := E.AddPolicy("teacher", "/teacher", "POST")
	if added {
		fmt.Println("addPolicy success!")
	} else {
		fmt.Println("Policy already existed!")
	}
	return E, err
}

//AddTeacherRole
//@title		AddTeacherRole()
//@description	添加角色
//@author		zy
//@param		username string
//@return
func AddTeacherRole(username string) {
	E.AddGroupingPolicy(username, "teacher")
}


