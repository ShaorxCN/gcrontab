package email

import (
	"fmt"
	"gcrontab/entity/task"

	"github.com/jordan-wright/email"
)

func sendTest() error {
	em := email.NewEmail()
	em.Subject = "crontabAlert email test"
	em.To = []string{testAddress}
	em.Subject = "init send"
	content := "init test"
	em.Text = []byte(content)
	em.From = user
	return em.SendWithTLS(addr, auth, tlsConfig)
}

// SendCrontabAlert 发送定时任务警告邮件
func SendCrontabAlert(code int, body string, t *task.Task, timeStamp int64, emails []string) error {

	em := email.NewEmail()
	em.Subject = fmt.Sprintf("Task[%s] exec failed alert", t.Name)

	//TODO: 分批发 如果人太多的话
	em.To = emails
	content := fmt.Sprintf(crontabAlertHTML, t.Name, code, body, fmt.Sprintf(viewURL, t.ID.GetIDValue(), timeStamp))
	em.From = user
	em.HTML = []byte(content)
	return em.SendWithTLS(addr, auth, tlsConfig)
}
