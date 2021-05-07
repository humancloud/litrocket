package mail

import (
	"net/smtp"
	"strings"
)

const (
	MyEMail = "2144103614@qq.com"
	PWD     = "tubvfxzokhokcdad"
	HOST    = "smtp.qq.com:587" //例如还有 smtp.126.com:25
	Subject = "验证码"
)

//	Send Mail.
//	body 正文, destmail 接收端邮箱.
func SendEmail(body, destmail string) error {
	hp := strings.Split(HOST, ":")          // 多个HOST由:分割
	send_to := strings.Split(destmail, ";") // 目标邮箱可以有多个，以;分割
	content_type := "Content_Type: text/html; charset=UTF-8"

	auth := smtp.PlainAuth("", MyEMail, PWD, hp[0])

	// Message.
	msg := []byte("To:" + destmail + "\r\nFrom: Litrocket" + "<" +
		MyEMail + ">\r\nSubject: " + Subject + "\r\n" +
		content_type + "\r\n\r\n" + body)

	return smtp.SendMail(HOST, auth, MyEMail, send_to, msg)
}
