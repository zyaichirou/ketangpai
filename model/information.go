package model

//UserInfo 用户对象
type UserInfo struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}
