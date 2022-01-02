package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ktp/common"
	"ktp/model"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Register
//@title		Register()
//@description	注册请求
//@author		zy
//@param		c *gin.Context
//@return
func Register(c *gin.Context) {
	var u model.UserInfo
	err := c.ShouldBind(&u)				//参数绑定
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"message": "无效的参数",
		})
		return
	}

	if u.Username == "" || u.Password == "" || u.Email == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "请将信息填充完整",
		})
		return
	}
	//判断用户名是否已经存在
	if common.QueryUserInfo(u) {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "username或Email已被注册",
		})
		return
	} else {
		if common.InsertUserInfo(u){
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"message": "注册成功",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"message": "some errors in sql",
			})
		}
	}
}


//Login
//@title		Login()
//@description	登录请求
//@author		zy
//@param		c *gin.Context
//@return
func Login(c *gin.Context) {
	var u model.UserInfo
	err := c.ShouldBind(&u)					//参数绑定
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"message": "无效的参数",
		})
		return
	}
	// 判断用户名密码是否正确
	if !common.QueryUserInfoExist(u) {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "用户名或密码有误!",
		})
		return
	}

	// 生成username对应的tokenString
	tokenString, err := common.GenToken(u.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message": "系统异常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"message": "login success",
		"data": gin.H{"token": tokenString},
	})
}

//LogOffMyAccount1
//@title		LogOffMyAccount1()
//@description	注销账号
//@author		zy
//@param		c *gin.Context
//@return
func LogOffMyAccount1(c *gin.Context) {
	var u model.UserInfo
	u.Username = c.PostForm("username")
	u.Password = c.PostForm("password")
	if u.Username == "" || u.Password == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "请将信息填充完整",
		})
		return
	}
	if common.QueryUserInfoExist(u) {
		common.GetEmail(&u)
		fmt.Println(u.Email)
		randNums := common.RandNumbers()
		common.SendTo(u.Email, randNums)
		session := sessions.Default(c)
		session.Set("username", u.Username)
		session.Set("identifiedNums", randNums)
		session.Save()
		c.HTML(http.StatusOK,"Logoff2.html",gin.H{})

	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "用户名或密码有误!",
		})
		return
	}
}


//LogOffMyAccount2
//@title		LogOffMyAccount2()
//@description	注销账号
//@author		zy
//@param		c *gin.Context
//@return
func LogOffMyAccount2(c *gin.Context) {
	var u model.UserInfo

	var err error
	var randNums string
	session := sessions.Default(c)
	u.Username = session.Get("username").(string)
	randNums = session.Get("identifiedNums").(string)
	if err != nil {
		fmt.Println("读取sessionid失败")
		return
	}
	identifiedNums := c.PostForm("identifiedNums")
	if identifiedNums == randNums {
		common.DeleteUser(u)
		c.JSON(http.StatusOK, gin.H{
			"message": "您已注销账号 期待与您下次相遇",
			"status":200,
		})
		return

	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "验证码有误!",
			"status": http.StatusForbidden,
		})
		return
	}
}

//UpdatePassword1
//@title		UpdatePassword1()
//@description	找回/更改密码请求
//@author		zy
//@param		c *gin.Context
//@return
func UpdatePassword1(c	 *gin.Context) {
	var u model.UserInfo

	u.Username = c.PostForm("username")
	u.Email = c.PostForm("email")

	if u.Username == "" || u.Email == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "请将信息填充完整",
		})
		return
	}

	if !common.QueryUserInfo(u) {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "username和Email不匹配",
		})
		return
	} else {
		fmt.Println(u.Username)
		fmt.Println(u.Email)
		randNums := common.RandNumbers()
		common.SendTo(u.Email, randNums)
		session := sessions.Default(c)
		session.Set("username", u.Username)
		session.Set("email", u.Email)
		session.Set("identifiedNums", randNums)
		session.Save()

		c.HTML(http.StatusOK,"Update.html",gin.H{})

	}
}

//UpdatePassword2
//@title		UpdatePassword2()
//@description	找回/更改密码请求
//@author		zy
//@param		c *gin.Context
//@return
func UpdatePassword2(c *gin.Context) {
	var u model.UserInfo

	var err error
	var randNums string
	session := sessions.Default(c)
	u.Username = session.Get("username").(string)
	randNums = session.Get("identifiedNums").(string)
	if err != nil {
		fmt.Println("读取sessionid失败")
		return
	}

	u.Password = c.PostForm("password")
	identifiedNums := c.PostForm("identifiedNums")
	if identifiedNums == randNums {
		if info := common.AlterUserInfo(u); info {
			c.JSON(http.StatusOK, gin.H{
				"message": "修改成功!",
				"status": 200,
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "修改失败!",
				"status": http.StatusForbidden,
			})
			return
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "验证码有误!",
			"status": http.StatusForbidden,
		})
		return
	}
}


//CreateClass
//@title		CreateClass()
//@description	创建课堂
//@author		zy
//@param		c *gin.Context
//@return
func CreateClass(c *gin.Context) {
	classId := common.RandClassId()
	for {
		if common.QueryClass(classId) {
			classId = common.RandClassId()
		} else {
			break
		}
	}
	username, _ := c.Get("username")
	className := c.PostForm("className")
	if className == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if common.FoundClass(classId, className, username.(string)) {
		common.AddTeacherRole(username.(string))
		err := os.Mkdir("D:/csa/classId/"+classId, os.ModePerm)
		if err != nil {
			fmt.Printf("dir already existed, err:%v\n", err)
		}
		err = os.Mkdir("D:/csa/classId/"+classId+"/PPT", os.ModePerm)
		if err != nil {
			fmt.Printf("dir already existed, err:%v\n", err)
		}
		err = os.Mkdir("D:/csa/classId/"+classId+"/assignment", os.ModePerm)
		if err != nil {
			fmt.Printf("dir already existed, err:%v\n", err)
		}


		c.JSON(http.StatusOK, gin.H{
			"message": "创建课堂成功！",
			"status": 200,
			"data": "课堂号为:" + classId,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":"服务器内部出现问题",
			"status": 500,
		})
		return
	}
}


//JoinClass
//@title		JoinClass()
//@description	加入课堂
//@author		zy
//@param		c *gin.Context
//@return
func JoinClass (c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	if classId == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if common.QueryClass(classId) {
		if common.JoinTheClass(classId, username.(string)) {
			c.JSON(http.StatusOK, gin.H{
				"message": "success to join class--" + classId,
				"status": 200,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "服务器内部出现问题",
				"status": 500,
			})
		}
		return
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "can't find the class",
			"status": 403,
		})
	}
}

//ReleaseAssignment
//@title		ReleaseAssignment()
//@description	发布作业
//@author		zy
//@param		c *gin.Context
//@return
func ReleaseAssignment(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	title := c.PostForm("title")
	content := c.PostForm("content")
	lastTime := c.PostForm("lastTime")   // 多少天
	if classId == "" || content == "" || lastTime == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	days, _ :=strconv.Atoi(lastTime)
	if common.ReleaseTheAssignment(classId, title, content, username.(string), days) {
		err := os.Mkdir("D:/csa/classId/"+classId+"/assignment/"+title, os.ModePerm)
		if err != nil {
			fmt.Printf("dir already existed, err:%v\n", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "作业发布成功!",
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
		return
	}
}

//CheckTheAssignment
//@title		CheckTheAssignment()
//@description	学生查询作业
//@author		zy
//@param		c *gin.Context
//@return
func CheckTheAssignment(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	if classId == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}

	ass1, ass2 := common.CheckAssignment(classId, username.(string))

	c.JSON(http.StatusOK, gin.H{
		"待完成": ass1,
		"已截止": ass2,
		"status": 200,
	})
}

//UploadAssignment
//@title		UploadAssignment()
//@description	学生提交作业
//@author		zy
//@param		c *gin.Context
//@return
func UploadAssignment(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	title := c.PostForm("title")
	if username == "" || classId == "" || title == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "上传文件有误!",
			"status": http.StatusBadRequest,
		})
		return
	}
	//询问是否有该title的作业
	if !common.QueryTitleExist(classId, title) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你在该班级下没有这个作业",
			"status": 403,
		})
		return
	}
	if common.QueryAssignmentLate(classId, title, username.(string)) {
		c.JSON(200, gin.H{
			"message": "已超过截止时间，禁止提交!",
			"status": 200,
		})
		return
	}

	i := 0
	for ; i < len(file.Filename); i++ {
		if file.Filename[i] == '.' {
			break
		}
	}
	var destPath string
	destPath = "D:/csa/classId/"+classId+"/assignment/"+title+"/"+ username.(string) + file.Filename[i:]
	err = c.SaveUploadedFile(file, destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
		return
	}

	if common.SubmitAssignment(classId, title, destPath, username.(string)) {
		c.JSON(200, gin.H{
			"message": "成功提交作业：" + title,
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
		return
	}
}

//CheckStudentsAssignment
//@title		CheckStudentsAssignment()
//@description	老师查阅作业
//@author		zy
//@param		c *gin.Context
//@return
func CheckStudentsAssignment(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	title := c.PostForm("title")
	studentName := c.PostForm("student")
	if classId == "" || title == "" || studentName == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}

	if !common.CheckTeacherAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的老师,没有权限",
			"status": 403,
		})
		return
	}
	destPath := common.GetAssignmentPath(classId, title, studentName)
	if destPath == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "该学生未提交作业",
			"status": 200,
		})
		return
	}
	fp, err := os.OpenFile(destPath, os.O_CREATE|os.O_APPEND, 6)
	defer fp.Close()
	if err != nil {
		fmt.Printf("open file failed, err:%v\n", err)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s",fp.Name()))
	c.File(destPath)

}




//UploadSinglePPT
//@title		UploadSinglePPT()
//@description	上传单个课件
//@author		zy
//@param		c *gin.Context
//@return
//func UploadSinglePPT(c *gin.Context) {
//	username, _ := c.Get("username")
//	classId := c.PostForm("classId")
//	file, err := c.FormFile("file")
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"message": "上传文件有误!",
//			"status": http.StatusBadRequest,
//		})
//		return
//	}
//	if classId == "" {
//		c.JSON(http.StatusForbidden, gin.H{
//			"message": "请将信息填充完整",
//			"status": 403,
//		})
//		return
//	}
//	if !common.CheckTeacherAndClass(classId, username.(string)) {
//		c.JSON(http.StatusForbidden, gin.H{
//			"message": "你没有权限上传课件",
//			"status": 403,
//		})
//		return
//	}
//	err = os.Mkdir("D:/csa/classId/PPT"+classId, os.ModePerm)
//	if err != nil {
//		fmt.Printf("dir already existed, err:%v\n", err)
//	}
//	err = c.SaveUploadedFile(file, "D:/csa/classId/PPT/"+classId+"/"+file.Filename)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": "服务器内部出现问题",
//			"status": 500,
//		})
//		return
//	}
//	c.JSON(200, gin.H{
//		"message": "成功上传课件:" + file.Filename,
//		"status": 200,
//	})
//}


//UploadPPT
//@title		UploadPPT()
//@description	上传课件
//@author		zy
//@param		c *gin.Context
//@return
func UploadPPT(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	form, err := c.MultipartForm()
	files := form.File["files"]
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "上传文件有误!",
			"status": http.StatusBadRequest,
		})
		return
	}
	if classId == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if !common.CheckTeacherAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你没有权限上传课件",
			"status": 403,
		})
		return
	}
	var name []string
	for _, file := range files {
		destPath := "D:/csa/classId/"+classId+"/PPT/"+file.Filename
		if common.UploadPPTs(classId, file.Filename, destPath) {
			fmt.Println("ppt ok")
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "服务器内部出现问题",
				"status": 500,
			})
			return
		}
		err = c.SaveUploadedFile(file, destPath)
		name = append(name, file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "unknown classId",
				"status": 500,
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"message": "成功上传课件!",
		"data": name,
		"status": 200,
	})
}

//CheckPPT
//@title		CheckPPT()
//@description	查看课件
//@author		zy
//@param		c *gin.Context
//@return
func CheckPPT(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	if classId == "" || username == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if common.CheckTeacherAndClass(classId, username.(string)) || common.CheckStudentAndClass(classId, username.(string)) {
		ppt := common.GetClassPPT(classId)
		c.JSON(http.StatusOK, gin.H{
			"message": "有以下课件",
			"ppt": ppt,
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的学生或老师",
			"status": 403,
		})
		return
	}
}


//DownloadPPT
//@title		DownloadPPT()
//@description	下载课件
//@author		zy
//@param		c *gin.Context
//@return
func DownloadPPT(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	pptName := c.PostForm("pptName")
	if classId == "" || username == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if common.CheckTeacherAndClass(classId, username.(string)) || common.CheckStudentAndClass(classId, username.(string)) {
		destPath := common.GetPPTPath(classId, pptName)
		fp, err := os.OpenFile(destPath, os.O_CREATE|os.O_APPEND, 6)
		defer fp.Close()
		if err != nil {
			fmt.Printf("open file failed, err:%v\n", err)
			return
		}
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s",fp.Name()))
		c.File(destPath)
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的学生或老师",
			"status": 403,
		})
		return
	}

}





//PostTopic
//@title		PostTopic()
//@description	发布话题
//@author		zy
//@param		c *gin.Context
//@return
func PostTopic(c *gin.Context) {
	username, _ := c.Get("username")
	topic := c.PostForm("topic")
	if topic == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if common.PostTheTopic(topic, username.(string)) {
		c.JSON(http.StatusOK, gin.H{
			"message": "发布成功!",
			"data": topic,
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
	}
}

//CommentTopic
//@title		CommentTopic()
//@description	讨论话题
//@author		zy
//@param		c *gin.Context
//@return
func CommentTopic(c *gin.Context) {
	username, _ := c.Get("username")
	topic := c.PostForm("topic")
	content := c.PostForm("content")
	if topic == "" || content == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}

	if !common.QueryTopicExist(topic) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "话题(" + topic + ")不存在",
			"status": 403,
		})
		return
	}

	if common.CommentTheTopic(topic, username.(string), content) {
		c.JSON(http.StatusOK, gin.H{
			"message": "评论成功!",
			"data": gin.H{
				"topic": topic,
				"content": content,
			},
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
	}
}

//CheckAllStudents
//@title		CheckAllStudents()
//@description	老师发布签到
//@author		zy
//@param		c *gin.Context
//@return
func CheckAllStudents(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	checkCode := common.RandClassId()
	for {
		if common.QueryClassCheckCode(checkCode) {
			classId = common.RandClassId()
		} else {
			break
		}
	}
	limitedSecond := c.PostForm("limitedSecond")
	limitedTime, err := strconv.Atoi(limitedSecond)
	if err != nil {
		fmt.Println(err)
	}
	if classId == "" || limitedSecond == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	fmt.Println("checkCode:", checkCode)
	if common.PostCheck(classId, checkCode, username.(string), limitedTime) {
		c.JSON(http.StatusOK, gin.H{
			"message": "发布签到成功!",
			"checkCode":  checkCode,
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题 or 你没有权限发布签到",
			"status": 500,
		})
	}
}

func CheckOver(c *gin.Context) {
	username, _ := c.Get("username")
	checkCode := c.PostForm("checkCode")
	classId := c.PostForm("classId")
	if checkCode == "" || classId == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}

	if !common.QueryClassIdAndCheckCode(classId, checkCode) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "课堂和签到码不匹配",
			"status": 403,
		})
		return
	}

	if !common.CheckTeacherAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "您不是该课堂的老师",
			"status": 403,
		})
		return
	}

	startTime, limitedTime := common.GetTheStartAndLimitedTime(classId, checkCode)
	for int(time.Now().Unix() - startTime.Unix()) < limitedTime {

	}
	students := common.AlreadyCheckInStu(classId, checkCode, startTime, limitedTime)  //找出所有时间相差在limitedTime以内的
	AllStudents := common.AllStudents(classId)
	m := make(map[string]int)
	for _, name := range AllStudents {
		m[name] = 0
	}
	for _, name := range students {
		m[name] = 1
	}
	var lateOrAbsentStu []string
	for name, value := range m {
		if value == 0 {
			lateOrAbsentStu = append(lateOrAbsentStu, name)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功签到学生名单",
		"已签到学生": students,
		"迟到或缺席学生": lateOrAbsentStu,
		"status": 200,
	})
	return
}

//StudentCheck
//@title		StudentCheck()
//@description	学生签到
//@author		zy
//@param		c *gin.Context
//@return
func StudentCheck(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	checkCode := c.PostForm("checkCode")
	if classId == "" || checkCode == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}
	if !common.QueryClassIdAndCheckCode(classId, checkCode) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "没有该签到",
			"status": 403,
		})
		return
	} else {
		if common.StudentCheckIn(classId, checkCode, username.(string)) {
			c.JSON(http.StatusOK, gin.H{
				"message": "签到!",
				"status": 200,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "服务器内部出现问题 or 你已签到",
				"status": 500,
			})
			return
		}
	}
}


//QuickAnswer
//@title		QuickAnswer()
//@description	抢答
//@author		zy
//@param		c *gin.Context
//@return
func QuickAnswer(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	question := c.PostForm("question")
	if classId == "" || question == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "请将信息填充完整",
		})
		return
	}
	// 如果username不是classId班级的老师
	if !common.CheckTeacherAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的老师,没有权限",
			"status": 403,
		})
		return
	}

	if common.StartQuickAnswer(classId, question) {
		c.JSON(http.StatusOK, gin.H{
			"message": "发布抢答成功!",
			"question": question,
			"status": 200,
		})
		//asd
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
		return
	}
}

//QuickAnswer
//@title		QuickAnswer()
//@description	抢答
//@author		zy
//@param		c *gin.Context
//@return
func FirstAnswer(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	question := c.PostForm("question")
	if classId == "" || question == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": http.StatusForbidden,
			"message": "请将信息填充完整",
		})
		return
	}
	if !common.CheckStudentAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的学生",
			"status": 403,
		})
		return
	}
	//asdas

}

//UploadStudentScore
//@title		UploadStudentScore()
//@description	老师上传学生成绩
//@author		zy
//@param		c *gin.Context
//@return
func UploadStudentScore(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	studentName := c.PostForm("studentName")
	score, _ := strconv.Atoi(c.PostForm("score"))

	// 如果username不是classId班级的老师
	if !common.CheckTeacherAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的老师,没有权限",
			"status": 403,
		})
		return
	}
	if score < 0 || score > 100 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "成绩有误!请重新上传",
			"status": 403,
		})
		return
	}
	if common.UploadStudentScore(classId, studentName, score) {
		c.JSON(http.StatusOK, gin.H{
			"message": "成功上传学生成绩!",
			"status": 200,
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部出现问题",
			"status": 500,
		})
		return
	}
}

//GetStudentsScore
//@title		GetStudentsScore()
//@description	老师获取学生成绩
//@author		zy
//@param		c *gin.Context
//@return
func GetStudentsScore(c *gin.Context) {
	username, _ := c.Get("username")
	classId := c.PostForm("classId")
	if classId == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "请将信息填充完整",
			"status": 403,
		})
		return
	}

	if !common.CheckTeacherAndClass(classId, username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "你不是该班级的老师,没有权限",
			"status": 403,
		})
		return
	}
	m := common.GetStudentsScore(classId)
	m1 := make(map[string]int)	//高分学生
	m2 := make(map[string]int)  //低分学生
	m3 := make(map[string]int)  //没有成绩的学生
	var totalScore int
	for student := range m {
		if m[student] == -1 {
			m3[student] = -1
		} else {
			if m[student] >= 90 {
				m1[student] = m[student]
			}
			if m[student] < 60 {
				m2[student] = m[student]
			}
			totalScore += m[student]
		}
	}
	var averageScore float64
	averageScore = float64(totalScore)/float64(len(m)-len(m3))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成绩成功",
		"90分以上": m1,
		"不及格": m2,
		"平均分": averageScore,
		"没有成绩的学生名单": m3,
	})
	return
}