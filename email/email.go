package email

import (
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
