package email

import (
	"EventFlow/common/tool/arraytool"
	"EventFlow/common/tool/parametertool"
	"fmt"
	"strings"
)

type EmailPlugin struct {
	Setting SettingConfig
}

func (trigger *EmailPlugin) FireAction(messageFromTrigger *string, parameters *map[string]interface{}) {

	dynamicto := parametertool.ReplaceWithParameter(&trigger.Setting.DynamicTo, parameters)
	to := append(trigger.Setting.StaticTo, strings.Split(dynamicto, ",")...)
	to = arraytool.RemoveItem(to, "")

	subject := parametertool.ReplaceWithParameter(&trigger.Setting.Subject, parameters)
	body := parametertool.ReplaceWithParameter(&trigger.Setting.Body, parameters)
	//smtpServer := trigger.Setting.Address + ":" + strconv.Itoa(trigger.Setting.Port)
	//auth := smtp.PlainAuth("", trigger.Setting.Username, trigger.Setting.Password, trigger.Setting.Address)

	header := ""
	header += fmt.Sprintf("From: %s\r\n", trigger.Setting.From)
	header += fmt.Sprintf("To: %s\r\n", strings.Join(to, ";"))
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "\r\n" + body

	//err := smtp.SendMail(smtpServer, auth, trigger.Setting.From, to, []byte(header))

	//if err != nil {
	//	log.Printf("[action][email] send email failed: %v", err)
	//}
}
