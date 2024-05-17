package internal

type (
	Message struct {
		Event         string      `json:"event_name"`
		ReplyEvent    string      `json:"reply_event"`
		Body          interface{} `json:"body"`
		Meta          any         `json:"meta"`
		CorrelationID string      `json:"correlation_id"`
		AppID         *string     `json:"app_id"`
		ApiToken      *string     `json:"api_token"`
	}

	MessageReply struct {
		Event         string      `json:"event_name"`
		Error         bool        `json:"error"`
		StatusCode    int         `json:"status_code"`
		Body          interface{} `json:"body"`
		Meta          any         `json:"meta"`
		CorrelationID string      `json:"correlation_id"`
	}

	SmsBody struct {
		App     string `json:"app"`
		Sender  string `json:"sender"`
		Message string `json:"message"`
		TO      string `json:"to"`
	}

	EmailBody struct {
		App     string `json:"app"`
		Sender  string `json:"sender"`
		TO      string `json:"to"`
		Message string `json:"message"`
		Subject string `json:"subject"`
		CC      string `json:"cc"`
		BCC     string `json:"bcc"`
	}
)
