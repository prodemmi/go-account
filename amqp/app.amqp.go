package amqp

import (
	"account/http/services"
	"account/internal"
	"fmt"
	"net/http"
)

type (
	AppAMQP struct {
		AMQP        internal.IAMQPService
		AuthService services.IAuthService
	}

	IAppAMQP interface {
		exists()
		Setup()
	}
)

func (a AppAMQP) exists() {
	messages, err := a.AMQP.Subscribe("CHECK_AUTH")
	if err != nil {
		fmt.Printf("cannot subscribe to CHECK_AUTH: %s", err.Error())
		return
	}

	for d := range messages {
		msg, _ := internal.StringToMessage(string(d.Body))
		valid, err := a.AuthService.Check(*msg.AppID, *msg.ApiToken)

		if err != nil {
			_, err := a.AMQP.ReplyTo(d.AppId, internal.MessageReply{
				Event:         msg.ReplyEvent,
				Body:          err.Error(),
				Error:         true,
				StatusCode:    http.StatusInternalServerError,
				CorrelationID: msg.CorrelationID,
			})
			if err != nil {
				return
			}
		} else {
			_, err := a.AMQP.ReplyTo(d.AppId, internal.MessageReply{
				Event:         msg.ReplyEvent,
				Body:          fmt.Sprint(valid),
				Error:         false,
				StatusCode:    http.StatusOK,
				CorrelationID: msg.CorrelationID,
			})
			if err != nil {
				return
			}
		}
	}
}

func (a AppAMQP) Setup() {
	go a.exists()
}

func NewAppAMQP(amqp internal.IAMQPService, authService services.IAuthService) IAppAMQP {
	return AppAMQP{AMQP: amqp, AuthService: authService}
}
