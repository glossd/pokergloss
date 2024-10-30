package mail

import (
	"context"
	"github.com/gin-gonic/gin"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqmail"
	"github.com/glossd/pokergloss/mail/idp"
	"github.com/glossd/pokergloss/mail/sender"
	log "github.com/sirupsen/logrus"
)

func Run(c *gin.Engine) func(context.Context) {
	if conf.IsProd() {
		idp.Init()
	}
	go func() {
		err := mqmail.Subscribe("mail-service", func(ctx context.Context, email *mqmail.Email) error {
			uemail, err := idp.GetUserEmail(ctx, email.ToUserId)
			if err != nil {
				return err
			}
			err = sender.SendEmail(uemail, email.Subject, email.Body)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Fatalf("mqmail.Subscribe failed: %s", err)
		}
	}()
	return func(ctx context.Context) {}
}
