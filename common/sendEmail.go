package common

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

//SendTo
//@title		SendTo()
//@description	给目标邮箱发邮件
//@author		zy
//@param		receiver string, identifyNums string
//@return
func SendTo(receiver string, identifyNums string) {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	em := email.NewEmail()

	// sender
	em.From = "你的邮箱号"

	// receiver
	em.To = []string{receiver}

	// subject
	em.Subject = "您的验证码"

	// Text
	em.Text = []byte(identifyNums)

	// config of server
	err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", "你的邮箱号", "你的邮箱授权码", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("send successfully ...")
}

