package routers

import (
	"github.com/gin-gonic/gin"
	"ktp/controller"
	middleware "ktp/middlerware"
)

//CollectRouter
//@title		CollectRoute()
//@description	各种路由集合
//@author		zy
//@param		r *gin.Engine
//@return		*gin.Engine
func CollectRouter(r *gin.Engine) *gin.Engine {

	// 内存限制
	r.MaxMultipartMemory = 100 << 20

	// HTML模板渲染
	r.LoadHTMLFiles("./template/admin/Index.html", "./template/admin/Update.html", "./template/admin/Logoff1.html", "./template/admin/Logoff2.html")

	//注册
	r.POST("/register", controller.Register)

	//登录
	r.POST("/login", controller.Login)

	// 修改/找回密码
	r.GET("/updatePassword1", func(c *gin.Context) {
		c.HTML(200,"Index.html", nil)
	})
	r.POST("/updatePassword1", controller.UpdatePassword1)
	r.POST("/updatePassword2", controller.UpdatePassword2)

	// 注销账号
	r.GET("/logOff", func(c *gin.Context) {
		c.HTML(200, "Logoff1.html", nil)
	})
	r.POST("/logOff1", controller.LogOffMyAccount1)
	r.POST("/logOff2", controller.LogOffMyAccount2)


	userGroup := r.Group("/user", middleware.JWTAuthMiddleware())
	{

		// 创建课堂
		userGroup.POST("/createClass", controller.CreateClass)

		teacher := userGroup.Group("/teacher", middleware.CheckTeacherRole())
		{

			// 布置作业
			teacher.POST("/releaseAssignment", controller.ReleaseAssignment)
			// 老师查阅学生作业
			teacher.POST("/checkStudentsAssignment", controller.CheckStudentsAssignment)


			// 上传-下载课件资料
			teacher.POST("/uploadClassPPT", controller.UploadPPT)
			teacher.POST("/checkPPT", controller.CheckPPT)
			teacher.POST("/downloadPPT", controller.DownloadPPT)


			// 发布话题  话题讨论
			teacher.POST("/topic", controller.PostTopic)
			teacher.POST("/comment", controller.CommentTopic)


			// 上课签到
			teacher.POST("/checkStudents", controller.CheckAllStudents)
			teacher.POST("/getCheckResult", controller.CheckOver)

			// 成绩管理
			// 打分
			teacher.POST("/uploadStudentScore", controller.UploadStudentScore)
			// 成绩分布 90分以上 不及格  平均分
			teacher.POST("/GetStudentsScore", controller.GetStudentsScore)

		}

		student := userGroup.Group("/student")
		{
			// 学生加入课堂
			student.POST("/joinClass", controller.JoinClass)
			// 学生查询作业
			student.POST("/checkAssignment", controller.CheckTheAssignment)
			// 学生提交作业
			student.POST("/uploadAssignment", controller.UploadAssignment)

			// 上传-下载课件资料
			student.POST("/checkPPT", controller.CheckPPT)
			student.POST("/downloadPPT", controller.DownloadPPT)

			// 发布话题  话题讨论
			student.POST("/topic", controller.PostTopic)
			student.POST("/comment", controller.CommentTopic)

			// 学生签到
			student.POST("/studentCheck", controller.StudentCheck)

		}
	}






	// 3.
	//r.POST("/UploadClassPPTs", middleware.JWTAuthMiddleware(), controller.UploadSinglePPTs)


	// 4.


	// 5.


	// 6.
	// 课中提问
	// 发布抢答
	//r.POST("/quickAnswer", middleware.JWTAuthMiddleware(), controller.QuickAnswer)
	// 学生抢答
	//r.POST("/firstAnswer", middleware.JWTAuthMiddleware(), controller.FirstAnswer)


	// 7.


	return r
}