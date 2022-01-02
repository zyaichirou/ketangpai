package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"ktp/common"
	"ktp/routers"
)

var store = cookie.NewStore([]byte("zyaichirou"))

func main() {
	DB, _ := common.InitDB("user:123456@tcp(127.0.0.1:3306)/ktp?charset=utf8mb4&parseTime=True&loc=Local")
	defer DB.Close()
	E, err := common.InitCasbin()
	if err = E.LoadPolicy(); err != nil {
		fmt.Printf("LoadPolicy failed, err:%v\n", err)
	}
	if err != nil {
		fmt.Printf("casbin failed, err:%v\n", err)
	}


	r := gin.Default()
	r.Use(sessions.Sessions("sessionid", store))


	r = routers.CollectRouter(r)

	r.Run(":9090")
	//panic(r.Run(":9090"))
}
