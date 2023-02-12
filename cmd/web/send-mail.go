package main

import (
	"fmt"
	"github.com/deenikarim/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//TODO: WHAT WE ARE GOING TO DO IN HERE

//ListenForMail listens all the times for incoming mail from the channel
func ListenForMail() {
	//todo: to create something that runs indefinitely and fires things off in the background
	//  to happen asynchronously.

	//how to information from a channel is just to declare a variable and make that equal to whatever we get
	// from the channel we listen to
	//m := <-app.MailChan

	//start a goroutine
	go func() {
		//todo: have a function that will execute in the background and will listen all the time for
		//  incoming data
		for {
			//we just listen to the channel indefinitely because it never goes out of the For loop
			msg := <-app.MailChan

			//let's send a email message when we get one by calling sendMsg() function and pass it
			//through channel
			sendMsg(msg)
		}
	}()
}

//sendMsg just send emails
func sendMsg(m models.MailData) {
	//todo: how to send email message by using the simple mail package we imported

	//TODO: STEP:1 is to tell the mail package we just imported what our server is
	//  [gives us the means to connect to a local mail server]
	server := mail.NewSMTPClient() //the server will have some information associated with it
	//where is our mail server
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false //do not want to keep connection to the mail server active all the time
	//so only make connection when i tell it to send mail
	server.ConnectTimeout = 10 * time.Second //can not connect within 10 seconds, then give up
	server.SendTimeout = 10 * time.Second
	//WHEN IN PRODUCTION MODE
	//server.Username = ""
	//server.Password = ""
	//server.Encryption =""

	//TODO: STEP-2: now we need an client
	//represents a SMTP Client for send email
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	//TODO: STEP-3:: need to construct our email message in a format that our client understands by "m variable"
	//this creates a new empty email message [This creates a struct and is a pointer to email from the package]
	// which means we can now setup some configurations
	email := mail.NewMSG()
	//set configuration
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	//todo add: if a template is specified, then use that but if not, follow normal way
	//if no template is specified
	if m.Template == "" {
		//set the body of the email message
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		//if a template is specified, then use that
		//load the template from disk
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-template/%s", m.Template)) //replace the placeholder with the template name
		if err != nil {
			app.ErrorLog.Println(err)
		}
		//need to convert the template that we read from disk because ReadFile() returns an array of bytes
		//so need to convert it into string
		mailTemplate := string(data) // now we have template in the memory

		//replace that body placeholder with the contents we specified in the variable called mailTemplate
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)

		//set the body of the email message
		email.SetBody(mail.TextHTML, msgToSend)
	}

	//TODO: STEP-4:: finally, at the point where we can send the email
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent successfully")
	}

	/********************************************************************************************************/
	//INFO SUMMARY:the client() has everything it needs, has server and client information, has the email body
	/********************************************************************************************************/
}
