package common

import (
	"ktp/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB


//InitDB
//@title		InitDB()
//@description	连接数据库
//@author		zy
//@param		dsn string
//@return		*sql.DB error
func InitDB(dsn string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("failed to open database, err:%v", err)
		return nil, err
	}
	err = DB.Ping()

	if err != nil {
		fmt.Printf("failed to connect database, err: %v", err)
		return nil, err
	}
	fmt.Println("success to open database")
	return DB, err
}

//QueryUserInfo
//@title		QueryUserInfo()
//@description	查询用户名u.Username和u.email是否已经存在
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func QueryUserInfo(u model.UserInfo) bool {
	sqlStr := "select username, email from user where username=? and email=?"
	var userTemp model.UserInfo
	err := DB.QueryRow(sqlStr, u.Username, u.Email).Scan(&userTemp.Username, &userTemp.Email)
	if err != nil {
		fmt.Printf("%s和%s还未被注册\n", u.Username, u.Email)
		return false
	}
	fmt.Printf("用户%s你好\n", userTemp.Username)
	return true
}


//InsertUserInfo
//@title		InsertUserInfo()
//@description	注册成功时插入到user表中
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func InsertUserInfo(u model.UserInfo) bool{
	sqlStr := "insert into user(username, password, email) values (?,?,?)"	//sql语句
	ret, err := DB.Exec(sqlStr, u.Username, u.Password, u.Email)				//插入操作
	if err != nil {
		fmt.Printf("insert failed, err:%v", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d.\n", id)
	return true
}

//QueryUsernameExist
//@title		QueryUsernameExist()
//@description	查询用户名和密码是否正确
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func QueryUsernameExist(u model.UserInfo) bool {
	sqlStr := "select username from user where username=?"	//sql语句
	var userTemp model.UserInfo

	err := DB.QueryRow(sqlStr, u.Username).Scan(&userTemp.Username) //调用QueryRow进行插入

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	fmt.Printf("用户:%s你好\n", userTemp.Username)
	return true
}

//QueryUserInfoExist
//@title		QueryUserInfoExist()
//@description	查询用户名和密码是否正确
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func QueryUserInfoExist(u model.UserInfo) bool {
	sqlStr := "select username, password from user where username=? and password=?"	//sql语句
	var userTemp model.UserInfo

	err := DB.QueryRow(sqlStr, u.Username, u.Password).Scan(&userTemp.Username, &userTemp.Password) //调用QueryRow进行插入

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	fmt.Printf("用户:%s你好\n", userTemp.Username)
	return true
}

//GetEmail
//@title		GetEmail()
//@description	得到用户邮箱
//@author		zy
//@param		u userinformation.UserInfo
//@return		bool
func GetEmail(u *model.UserInfo) bool {
	sqlStr := "select email from user where username=? and password=?"
	var userTemp model.UserInfo

	err := DB.QueryRow(sqlStr, u.Username, u.Password).Scan(&userTemp.Email) //调用QueryRow进行插入

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	u.Email = userTemp.Email
	fmt.Printf("用户:%s你好\n", userTemp.Username)
	return true
}

//AlterUserInfo
//@title		AlterUserInfo()
//@description	更改用户密码
//@author		zy
//@param		u model.UserInfo
//@return		bool
func AlterUserInfo(u model.UserInfo) bool{
	sqlStr := "update user set password=? where username=?"
	_, err := DB.Exec(sqlStr, u.Password, u.Username)
	if err != nil {
		fmt.Printf("failed to update, err:%v\n", err)
		return false
	}
	fmt.Printf("success to update password\n")
	return true
}

func DeleteUser(u model.UserInfo) {
	sqlStr1 := "delete from user where username=?"
	ret, err := DB.Exec(sqlStr1, u.Username)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)

	sqlStr1 = "delete from class where member=?"
	ret, err = DB.Exec(sqlStr1, u.Username)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err = ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)

	sqlStr1 = "delete from topics where username=?"
	ret, err = DB.Exec(sqlStr1, u.Username)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err = ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

//QueryClass
//@title		QueryClass()
//@description	classId的课堂是否已存在
//@author		zy
//@param		classId string, username string
//@return		bool
func QueryClass(classId string) bool {
	sqlStr := "select classId from class where classId=?"
	var tempId string

	err := DB.QueryRow(sqlStr, classId).Scan(&tempId)
	if err != nil {
		fmt.Printf("该classId--%s还未被使用:\n", classId)
		return false
	}
	fmt.Printf("classId:%s已被使用\n", tempId)
	return true
}


//FoundClass
//@title		FoundClass()
//@description	创建课堂
//@author		zy
//@param		classID string, username string
//@return		bool
func FoundClass(classId , className, username string) bool {
	sqlStr := "insert into class(classId, className, member, identity) values(?,?,?,?)"

	ret, err := DB.Exec(sqlStr, classId, className, username, "1")
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d\n", id)
	return true
}


//GetClassMembers
//@title		GetClassMembers()
//@description	获得班级的人数
//@author		zy
//@param		classId string
//@return		bool
func GetClassMembers(classId string) int {
	sqlStr := "select count(*) from class where classId=?"
	ret, err := DB.Exec(sqlStr, classId)
	if err != nil {
		fmt.Printf("Get classMembers failed, err:%v", err)
		return -1
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v", err)
		return -1
	}
	return int(n)
}

//CheckTeacherAndClass
//@title		CheckTeacherAndClass()
//@description	老师是否和课堂匹配
//@author		zy
//@param		classID string, username string
//@return		bool
func CheckTeacherAndClass(classId, username string) bool {
	sqlStr := "select classId, member from class where classId=? and member=? and identity!='0'"
	var tempClassId, tempUsername string
	err := DB.QueryRow(sqlStr, classId, username).Scan(&tempClassId, &tempUsername)
	fmt.Println("check Username:", tempUsername)
	if err != nil {
		fmt.Printf("checkTeacherAndClass failed, err:%v\n", err)
		return false
	}
	return true
}

//CheckStudentAndClass
//@title		CheckStudentAndClass()
//@description	学生是否和课堂匹配
//@author		zy
//@param		classID string, username string
//@return		bool
func CheckStudentAndClass(classId, username string) bool {
	sqlStr := "select classId, member from class where classId=? and member=?"
	var tempClassId, tempUsername string
	err := DB.QueryRow(sqlStr, classId, username).Scan(&tempClassId, &tempUsername)
	fmt.Println("check Username:", tempUsername)
	if err != nil {
		fmt.Printf("checkStudentAndClass failed, err:%v\n", err)
		return false
	}
	return true
}

//JoinTheClass
//@title		JoinTheClass()
//@description	加入课堂
//@author		zy
//@param		classID string, username string
//@return		bool
func JoinTheClass(classId, username string) bool {
	sqlStr1 := "select className from class where classId=?"
	var className string
	err := DB.QueryRow(sqlStr1, classId).Scan(&className)
	if err != nil {
		fmt.Printf("select className failed, err:%v\n", err)
		return false
	}
	sqlStr2 := "insert into class(classId, className, member, identity) values(?,?,?,?)"
	ret, err := DB.Exec(sqlStr2, classId, className, username, "0")
	if err != nil {
		fmt.Printf("student insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d.\n", id)
	return true
}

//ReleaseTheAssignment
//@title		ReleaseTheAssignment()
//@description	发布作业
//@author		zy
//@param		classID string, username string
//@return		bool
func ReleaseTheAssignment(classId, title, content, username string, lastTime int) bool {
	sqlStr1 := "insert into assignment(classId,title,content,username,beginTime,lastTime, isTeacher) values(?,?,?,?,?,?,?)"
	ret, err := DB.Exec(sqlStr1, classId, title, content, username, time.Now(), lastTime, 1)
	if err != nil {
		fmt.Printf("Release insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d\n", id)
	allStu := AllStudents(classId)
	sqlStr2 := "insert into assignment(classId,title,content,username,beginTime,lastTime) values(?,?,?,?,?,?)"
	for _, name := range allStu {
		_, err := DB.Exec(sqlStr2, classId, title, content, name, time.Now(), lastTime)
		if err != nil {
			fmt.Printf("Release insert failed, err:%v\n", err)
			return false
		}
	}
	return true
}


//CheckAssignment
//@title		CheckAssignment()
//@description	检查还有多少作业没写
//@author		zy
//@param		classId, username string
//@return		m1, m2 map[string]string
func CheckAssignment(classId, username string) (m1, m2 map[string]string){
	sqlStr := "select title, content, beginTime, lastTime from assignment where classId=? and username=? and assignmentPath is NULL"
	rows, err :=DB.Query(sqlStr, classId, username)
	if err != nil {
		fmt.Printf("CheckAssignment failed, err:%v\n", err)
		return nil, nil
	}
	defer rows.Close()
	m1 = make(map[string]string)
	m2 = make(map[string]string)
	for rows.Next()  {
		var tempBT time.Time
		var tempTitle, tempContent string
		var tempLastTime int64
		err := rows.Scan(&tempTitle, &tempContent, &tempBT, &tempLastTime)
		fmt.Printf("Got assignment:%s\n", tempTitle)
		if err != nil {
			fmt.Printf("AllStudents rows scan failed, err:%v\n", err)
			return nil, nil
		}
		if time.Now().Unix() - tempBT.Unix() < tempLastTime*3600*24 {
			m1[tempTitle] = tempContent
		} else {
			m2[tempTitle] = tempContent
		}
	}

	return m1, m2
}


//QueryTitleExist
//@title		QueryTitleExist()
//@description	检查是否有Title的作业存在
//@author		zy
//@param		classId, title string
//@return		bool
func QueryTitleExist(classId, title string) bool {
	sqlStr := "select distinct title from assignment where title=? and classId=?"
	var tempTitle string
	err := DB.QueryRow(sqlStr, title, classId).Scan(&tempTitle)
	if err != nil {
		fmt.Printf("QueryTitleExist failed, err:%v\n", err)
		return false
	}
	return true
}

//QueryAssignmentLate
//@title		QueryAssignmentLate()
//@description	查询作业为title的是否已过截止时间
//@author		zy
//@param		classId, title string
//@return		bool
func QueryAssignmentLate(classId, title, username string) bool {
	sqlStr := "select beginTime, lastTime from assignment where classId=? and title=? and username=?"
	var tempBT time.Time
	var lastTime int64
	err := DB.QueryRow(sqlStr, classId, title, username).Scan(&tempBT, &lastTime)
	if err != nil {
		fmt.Printf("CommitAssignment query failed, err:%v\n", err)
		return false
	}
	if time.Now().Unix() - tempBT.Unix() > lastTime*3600*24 {
		return true
	}
	return false
}

//SubmitAssignment
//@title		SubmitAssignment()
//@description	提交作业
//@author		zy
//@param		classId, title, destPath, username string
//@return		bool
func SubmitAssignment(classId, title, destPath, username string) bool {
	sqlStr := "update assignment set assignmentPath=? where classId=? and title=? and username=?"
	ret, err := DB.Exec(sqlStr, destPath, classId, title, username)
	if err != nil {
		fmt.Printf("commitAssignment failed, err:%v\n", err)
		return false
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return false
	}
	fmt.Printf("update success, affected rows:%d\n", n)
	return true
}


//GetAssignmentPath
//@title		GetAssignmentPath()
//@description	获取作业文件路径
//@author		zy
//@param		classId, title, destPath, username string
//@return		bool
func GetAssignmentPath(classId, title, studentName string) string {
	sqlStr := "select assignmentPath from assignment where classId=? and title=? and username=?"
	var temp string
	err := DB.QueryRow(sqlStr, classId, title, studentName).Scan(&temp)
	if err != nil {
		fmt.Printf("GetAssignmentPath failed, err:%v\n", err)
		return ""
	}
	return temp
}

//UploadPPTs
//@title		UploadPPTs()
//@description	上传ppt路径
//@author		zy
//@param		classId, title, destPath, username string
//@return		bool
func UploadPPTs(classId, Filename, destPath string) bool{
	sqlStr := "insert into classPPT(classId, pptName, pptPath) values(?,?,?)"
	ret, err := DB.Exec(sqlStr, classId, Filename, destPath)
	if err != nil {
		fmt.Printf("ppt insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d\n", id)
	return true
}

//GetClassPPT
//@title		GetClassPPT()
//@description	获取班级课件
//@author		zy
//@param		classId, title, destPath, username string
//@return		bool
func GetClassPPT(classId string) []string {
	sqlStr := "select pptName from classppt where classId=?"
	rows, err := DB.Query(sqlStr, classId)
	if err != nil {
		fmt.Printf("GetClassPPT failed, err:%v\n", err)
		return nil
	}
	defer rows.Close()

	var PPT []string

	for rows.Next()  {
		var tempPPT string
		err := rows.Scan(&tempPPT)
		fmt.Printf("Got PPTname:%s\n", tempPPT)
		if err != nil {
			fmt.Printf("GetClassPPT rows scan failed, err:%v\n", err)
			return PPT
		}
		PPT = append(PPT, tempPPT)
	}
	return PPT
}

//GetPPTPath
//@title		GetPPTPath()
//@description	获取课件文件路径
//@author		zy
//@param		classId, title, destPath, username string
//@return		bool
func GetPPTPath(classId, pptName string) string {
	sqlStr := "select pptPath from classppt where classId=? and pptName=?"
	var temp string
	err := DB.QueryRow(sqlStr, classId, pptName).Scan(&temp)
	if err != nil {
		fmt.Printf("GetAssignmentPath failed, err:%v\n", err)
		return ""
	}
	return temp
}

//PostTheTopic
//@title		PostTheTopic()
//@description	发布话题
//@author		zy
//@param		topic, username string
//@return		bool
func PostTheTopic(topic, username string) bool{
	sqlStr := "insert into topics(topic, username, content, postTime) values(?,?,?,?)"
	ret, err := DB.Exec(sqlStr, topic, username, "", time.Now())
	if err != nil {
		fmt.Printf("topic insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d\n", id)
	return true
}

//QueryTopicExist
//@title		QueryTopicExist()
//@description	话题是否存在
//@author		zy
//@param		topic, username string
//@return		bool
func QueryTopicExist(topic string) bool {
	sqlStr := "select topic from topics where topic=?"
	var tempTopic string

	err := DB.QueryRow(sqlStr, topic).Scan(&tempTopic)
	if err != nil {
		fmt.Printf("不存在话题(%s):\n", tempTopic)
		return false
	}
	fmt.Printf("话题%s存在\n", tempTopic)
	return true
}

//CommentTheTopic
//@title		CommentTheTopic()
//@description	话题讨论
//@author		zy
//@param		topic, username, content string
//@return		bool
func CommentTheTopic(topic, username, content string) bool{
	sqlStr := "insert into topics(topic, username, content, postTime) values(?,?,?,?)"
	ret, err := DB.Exec(sqlStr, topic, username, content, time.Now())
	if err != nil {
		fmt.Printf("content insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d\n", id)
	return true
}

//PostCheck
//@title		PostCheck()
//@description	发起签到
//@author		zy
//@param		classId, checkCode, username string, limitedSecond int
//@return		bool
func PostCheck(classId, checkCode, username string, limitedSecond int) bool{
	if !CheckTeacherAndClass(classId, username) {
		fmt.Printf("你没有权限!\n")
		return false
	}
	sqlStr := "insert into checkClass(classId, checkCode, beginTime, beginer, username, limitedSecond) values(?,?,?,?,?,?)"
	ret, err := DB.Exec(sqlStr, classId, checkCode, time.Now(), 1, username, limitedSecond)
	if err != nil {
		fmt.Printf("checkClass insert failed, err:%v\n", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("classCheck insert success, the id is %d\n", id)
	return true
}


//QueryClassCheckCode
//@title		QueryClassCheckCode()
//@description	签到码是否被使用
//@author		zy
//@param		checkCode string
//@return		bool
func QueryClassCheckCode(checkCode string) bool {
	sqlStr := "select checkCode from checkClass where checkCode=?"
	var tempCheckCode string
	err := DB.QueryRow(sqlStr, checkCode).Scan(&tempCheckCode)
	if err != nil {
		fmt.Printf("checkCode还未被使用, err:%v\n", err)
		return false
	}
	return true
}

//QueryClassIdAndCheckCode
//@title		QueryClassIdAndCheckCode()
//@description	签到码和课堂是否匹配
//@author		zy
//@param		classId, checkCode string
//@return		bool
func QueryClassIdAndCheckCode(classId, checkCode string) bool {
	sqlStr := "select classId, checkCode from checkClass where classId=? and checkCode=?"
	var tempClassId, tempCheckCode string

	err := DB.QueryRow(sqlStr, classId, checkCode).Scan(&tempClassId, &tempCheckCode)
	if err != nil {
		fmt.Printf("不存在该班级的签到:%v\n", err)
		return false
	}
	fmt.Println("存在签到")
	return true
}

//StudentCheckIn
//@title		StudentCheckIn()
//@description	学生签到
//@author		zy
//@param		classId, checkCode, username string
//@return		bool
func StudentCheckIn(classId, checkCode, username string) bool {
	sqlStr1 := "select classId from checkClass where classId=? and checkCode=? and username=?"
	var tempClassId string
	err := DB.QueryRow(sqlStr1, classId, checkCode, username).Scan(&tempClassId)
	if err == nil {
		fmt.Printf("studentCheck have checked, err:%v\n", err)
		return false
	}
	sqlStr := "insert into checkClass(classId, checkCode, username, clickTime) values (?,?,?,?)"
	ret, err := DB.Exec(sqlStr, classId, checkCode, username, time.Now())
	if err != nil {
		fmt.Printf("insert checkClass failed, err:%v", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert checkClass success, the id is %d.\n", id)
	return true
}

//GetTheStartAndLimitedTime
//@title		GetTheStartAndLimitedTime()
//@description	获得签到开始和限制时间
//@author		zy
//@param		classId, checkCode string
//@return		t time.Time, tl int
//有问题
func GetTheStartAndLimitedTime(classId, checkCode string) (t time.Time, tl int) {
	fmt.Printf("classId=%s  checkCode:%s\n", classId, checkCode)
	sqlStr := "select beginTime, limitedSecond from checkClass where classId=? and checkCode=? and limitedSecond != 0"
	var startTime time.Time
	var limitedTime int
	err := DB.QueryRow(sqlStr, classId, checkCode).Scan(&startTime, &limitedTime)
	if err != nil {
		fmt.Printf("time select failed, err:%v\n", err)
		return
	}
	return startTime, limitedTime
}

//AlreadyCheckInStu
//@title		AlreadyCheckInStu()
//@description	获得成功签到学生的名单
//@author		zy
//@param		classId, checkCode string, startTime time.Time, limitedTime int
//@return		[]string
func AlreadyCheckInStu(classId, checkCode string, startTime time.Time, limitedTime int) []string {
	sqlStr := "select username, clickTime from checkClass where classId=? and checkCode=? and beginer != 1"
	rows, err := DB.Query(sqlStr, classId, checkCode)
	if err != nil {
		fmt.Printf("CheckInStu failed, err:%v\n", err)
		return nil
	}
	defer rows.Close()

	var students []string

	for rows.Next()  {
		var tempUsername string
		var tempClickTime time.Time
		err := rows.Scan(&tempUsername, &tempClickTime)
		fmt.Printf("Got username:%s\n", tempUsername)
		if err != nil {
			fmt.Printf("CheckInStu rows scan failed, err:%v\n", err)
			return students
		}
		if int(tempClickTime.Unix() - startTime.Unix()) < limitedTime {
			fmt.Println("checked student:", tempUsername)
			students = append(students, tempUsername)
		}
	}
	fmt.Println(students)
	return students
}

//StartQuickAnswer
//@title		StartQuickAnswer()
//@description	开启抢答
//@author		zy
//@param		classId, question string
//@return		bool
func StartQuickAnswer(classId, question string) bool {
	sqlStr := "insert into quickAnswer(classId, question, beginTime) values (?,?,?)"	//sql语句
	ret, err := DB.Exec(sqlStr, classId, question, time.Now())				//插入操作
	if err != nil {
		fmt.Printf("insert failed, err:%v", err)
		return false
	}
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err: %v\n", err)
		return false
	}
	fmt.Printf("insert success, the id is %d.\n", id)
	return true
}

//AllStudents
//@title		AllStudents()
//@description	获取一个班学生的名单
//@author		zy
//@param		classId string
//@return		[]string
func AllStudents(classId string) []string {
	sqlStr := "select member from class where classId=? and identity=0"
	rows, err := DB.Query(sqlStr, classId)
	if err != nil {
		fmt.Printf("AllStudents failed, err:%v\n", err)
		return nil
	}
	defer rows.Close()

	var students []string

	for rows.Next()  {
		var tempUsername string
		err := rows.Scan(&tempUsername)
		fmt.Printf("Got username:%s\n", tempUsername)
		if err != nil {
			fmt.Printf("AllStudents rows scan failed, err:%v\n", err)
			return students
		}
		students = append(students, tempUsername)
	}
	return students
}

//UploadStudentScore
//@title		UploadStudentScore()
//@description	上传学生成绩
//@author		zy
//@param		classId, studentName string, score int
//@return		bool
func UploadStudentScore(classId, studentName string, score int) bool {
	sqlStr := "update class set score=? where classId=? and member=?"
	_, err := DB.Exec(sqlStr, score, classId, studentName)
	if err != nil {
		fmt.Printf("failed to uploadStudentScore, err:%v\n", err)
		return false
	}
	fmt.Printf("success to uploadStudentScore\n")
	return true

}

//GetStudentsScore
//@title		GetStudentsScore()
//@description	获取学生成绩
//@author		zy
//@param		classId, studentName string, score int
//@return		bool
func GetStudentsScore(classId string) (m map[string]int){
	sqlStr := "select member, score from class where classId=? and identity ='0'"
	rows, err := DB.Query(sqlStr, classId)
	if err != nil {
		fmt.Printf("GetStudentScore failed, err:%v\n", err)
		return nil
	}
	defer rows.Close()

	m = make(map[string]int)

	for rows.Next()  {
		var tempUsername string
		var tempScore int
		err := rows.Scan(&tempUsername, &tempScore)
		fmt.Printf("Got username:%s\n", tempUsername)
		if err != nil {
			fmt.Printf("GetStudentScore rows scan failed, err:%v\n", err)
			return m
		}
		m[tempUsername] = tempScore
	}
	fmt.Println(m)
	return m
}