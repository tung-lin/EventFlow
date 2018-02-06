package email

import (
	"EventFlow/common/tool/parametertool"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
)

type EmailPlugin struct {
	Setting SettingConfig
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (trigger *EmailPlugin) FireAction(messageFromTrigger *string, parameters *map[string]interface{}) {

	dynamicto := parametertool.ReplaceWithParameter(&trigger.Setting.DynamicTo, parameters)

	to := append(trigger.Setting.StaticTo, strings.Split(dynamicto, ",")...)
	to = checkEmailAddress(&to)

	subject := parametertool.ReplaceWithParameter(&trigger.Setting.Subject, parameters)
	body := parametertool.ReplaceWithParameter(&trigger.Setting.Body, parameters)
	smtpServer := trigger.Setting.Address + ":" + strconv.Itoa(trigger.Setting.Port)
	auth := smtp.PlainAuth("", trigger.Setting.Username, trigger.Setting.Password, trigger.Setting.Address)

	header := ""
	header += fmt.Sprintf("From: %s\r\n", trigger.Setting.From)
	header += fmt.Sprintf("To: %s\r\n", strings.Join(to, ";"))
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "\r\n" + body

	err := smtp.SendMail(smtpServer, auth, trigger.Setting.From, to, []byte(header))

	if err != nil {
		log.Printf("[action][email] send email failed: %v", err)
	}
}

func checkEmailAddress(emailAddresses *[]string) (results []string) {

	for _, email := range *emailAddresses {
		if email != "" && emailRegexp.MatchString(email) {
			results = append(results, email)
		}
	}

	return
}
