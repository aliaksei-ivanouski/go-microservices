package main

import (
	"log"
	"net/http"

	"github.com/tsawler/toolbox"
)

var tools toolbox.Tools

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	log.Println("msg", msg)

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	var payload toolbox.JSONResponse
	payload.Error = false
	payload.Message = "sent to " + requestPayload.To

	tools.WriteJSON(w, http.StatusAccepted, payload)
}
