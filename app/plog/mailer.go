package plog

import (
	"gopkg.in/sendgrid/sendgrid-go.v2"
	"github.com/crowdflux/angel/app/config"
	"runtime/debug"
	"fmt"
	"github.com/Sirupsen/logrus"
)



func ErrorMail(tag string, err error, args ...interface{}) {

	if err == nil {
		return
	}
	// gets the stack trace of current go routine
	stackTrace := string(debug.Stack())

	errString := fmt.Sprintf("%+v", err)
	argsString := ""

	if len(args) > 0 {
		argsString = fmt.Sprintf("%+v", args)
	}

	if config.IsDevelopment() || config.IsStaging() {
		logr.WithFields(logrus.Fields{
			"err" : err,
			"err_string" : errString,
			"args_string": argsString,
			"stack_trace" : stackTrace,
		}).Debug("Mailer in Dev/Staging: "+ tag)
		fmt.Println(tag)
		fmt.Println(err)
		fmt.Println(errString)
		fmt.Println(argsString)
		fmt.Println(stackTrace)
		return
	}
	subject:= "Error | Angel | " + tag + " | " + err.Error() + " | " + config.GetEnv()

	text:=
		`
		Error occured : ` + err.Error() + `
			---

			` + errString + `

			` + argsString + `

			` + stackTrace

	sendErr:= sendMail(subject,text)
	if sendErr != nil {

		logr.WithFields(logrus.Fields{
			"error" : err,
			"error_string" : errString,
			"stack_trace" : stackTrace,
			"mail_error" : sendErr,
			"args" : fmt.Sprintf("%+v", args),
		}).Error("Mailer error : "+tag)
	}
}

func sendMail (subject string, text string) error {

	cl := sendgrid.NewSendGridClientWithApiKey(config.SENDGRID_API_KEY.Get())

	mail := sendgrid.NewMail()

	mail.To = []string{"dev@playment.in"}

	mail.From = "no-reply@playment.in"

	mail.Subject = subject

	mail.Text = text

	sendErr := cl.Send(mail)
	return sendErr

}

