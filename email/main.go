package main

import (
	"log"

	"gopkg.in/gomail.v2"
)

func main() {
	d := gomail.NewDialer("smtp.gmail.com", 465, "ajudanpribadi50@gmail.com", "hrrmjatfwuuemrlz")
	m := gomail.NewMessage()
	m.SetHeader("From", "ajudanpribadi50@gmail.com")
	m.SetHeader("To", "aryadevaraj19@gmail.com")
	m.SetHeader("Subject", "Hey Kudo!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	log.Println("sending...")
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	log.Println("success")

	// Send emails using d.
}
