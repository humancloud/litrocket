package apiv1

import (
	"fmt"
	"litrocket/utils/handlelog"
	"litrocket/utils/mail"
	"math/rand"
	"sync"
	"time"
)

var signupcode = make(map[string]string, 1024)
var forgotcode = make(map[string]string, 1024)
var codemutex sync.Mutex

func SignUpVerify(Mail string) {
	Verify(Mail, "欢迎注册Litrocket,", signupcode)
}

func ReloadPassVerify(Mail string) {
	Verify(Mail, "Litrocket 重置密码,", forgotcode)
}

func Verify(Mail, info string, mapp map[string]string) {
	code := GetVerifyCode()
	if err := mail.SendEmail(info+"您的验证码为: "+code, Mail); err != nil {
		handlelog.Handlelog("WARNING", "VerifyCode "+err.Error())
	}

	codemutex.Lock()
	mapp[Mail] = code
	codemutex.Unlock()
	go waitingcode(Mail, mapp) // 60s later , reload verify code
}

func GetVerifyCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func waitingcode(Mail string, mapp map[string]string) {
	time.Sleep(60 * time.Second)
	codemutex.Lock()
	mapp[Mail] = ""
	codemutex.Unlock()
}
