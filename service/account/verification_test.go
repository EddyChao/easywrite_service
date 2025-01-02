package account

import (
	"easywrite-service/tools"
	"fmt"
	"testing"
)

func TestCode(t *testing.T) {
	err := tools.SendCodeToEmail("weilizan71@gmail.com", "【xx】698927（设置本机号码验证码）。工作人员不会向您索要，请勿向任何人泄露，以免造成账户或资金损失。")
	if err != nil {
		fmt.Println(err)
	}
}

//找回密码
