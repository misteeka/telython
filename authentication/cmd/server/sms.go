package main

// test file

import (
	"github.com/vonage/vonage-go-sdk"
)

var smsClient *vonage.SMSClient

var API_KEY string
var API_SECRET string

func initSMSClient() {
	auth := vonage.CreateAuthFromKeySecret(API_KEY, API_SECRET)
	smsClient = vonage.NewSMSClient(auth)
}

func sendSms(from string, to string, text string) bool {
	_, response, err := smsClient.Send(from, to, text, vonage.SMSOpts{})
	if err != nil {
		return false
	}
	if response.Messages[0].Status != "0" {
		return false
	}
	return true
}
