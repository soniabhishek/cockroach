package plog

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"gitlab.com/playment-main/angel/app/config"
	"runtime/debug"
)

func ErrorMail(tag string, err error, args ...interface{}) {

	// gets the stack trace of current go routine
	stackTrace := string(debug.Stack())

	errString := fmt.Sprintf("%#v", err)
	argsString := ""

	if len(args) > 0 {
		argsString = fmt.Sprintf("%#v", args)
	}

	if config.IsDevelopment() {
		fmt.Println(tag)
		fmt.Println(err)
		fmt.Println(errString)
		fmt.Println(argsString)
		fmt.Println(stackTrace)
		return
	}

	go func() {

		cl := sendgrid.NewSendGridClientWithApiKey(config.Get(config.SENDGRID_API_KEY))

		mail := sendgrid.NewMail()

		mail.To = []string{"dev@playment.in"}

		mail.From = "no-reply@playment.in"

		mail.Subject = "Error | Angel | " + tag + " | " + err.Error() + " | " + config.GetEnv()

		mail.Text =
			`
			Error occured : ` + err.Error() + `
			---

			` + errString + `

			` + argsString + `

			` + stackTrace

		sendErr := cl.Send(mail)
		if sendErr != nil {
			fmt.Println(err, errString, tag, err, args, stackTrace, sendErr)
		}
	}()
}

func Error(tag string, err error, args ...interface{}) {
	ErrorMail(tag, err, args)
}
