package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/gocolly/colly"
	"gopkg.in/gomail.v2"
)

var (
	fromEmail = "Harumnyoo <ews@gmail.com>"
	coreEmail = "gakpernahlupa@gmail.com"
	toEmail   = "jaderizan@gmail.com"
	subject   = "Tiket Bandung - Jakarta Tersedia!"
	textBody  = "Beli gan buruan"
)

func main() {
	fmt.Print("Input URL KAI: ")

	var in1 string
	_, err := fmt.Scanln(&in1)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		scrapeStart(in1)

		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		i -= 1
	}

}

func scrapeStart(url string) {

	c := colly.NewCollector()

	var tanggalPesanan, stasiunAsal, stasiunTujuan string

	c.OnHTML("input[name='origination']#origination", func(e *colly.HTMLElement) {
		stasiunAsal = e.Attr("value")
	})

	c.OnHTML("input[name='destination']#destination", func(e *colly.HTMLElement) {
		stasiunTujuan = e.Attr("value")
	})

	c.OnHTML("input[name='tanggal']#departure_dateh", func(e *colly.HTMLElement) {
		tanggalPesanan = e.Attr("value")
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		isSold := strings.Contains(e.Text, "Habis")
		isAvail := strings.Contains(e.Text, "Tersedia")
		isLimit := strings.Contains(e.Text, "Tersisa")

		currentTime := time.Now().Format("15:04:05")

		if isSold {
			fmt.Printf("\033[1mTicket KAI \033[33;1m[%s-%s %s] \033[31mSOLD \033[0m<<<<<<<<<<< %s\n", stasiunAsal, stasiunTujuan, tanggalPesanan, currentTime)
			return
		}

		if isAvail || isLimit {
			fmt.Printf("\033[1mTicket KAI \033[33;1m[%s-%s %s] \033[32;1mAVAILABLE \033[0m<<<<< %s\n", stasiunAsal, stasiunTujuan, tanggalPesanan, currentTime)
			notifyBySound()
			sendEmail(subject, textBody, url)
			return
		}

		fmt.Println("System Error <<<<<<<<<<", currentTime)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

}

func notifyBySound() {
	err := beeep.Beep(1, beeep.DefaultDuration)
	if err != nil {
		fmt.Println("Error notifying:", err)
	}
}

func sendEmail(subject, textBody, url string) {
	htmlBody := fmt.Sprintf("<b>Link to buy:</b> <a href=\"%s\">HERE</a>", url)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fromEmail)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", textBody)
	mailer.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, coreEmail, "aztwwnpupmrhfhtk")

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Println("Error sending email:", err)
	}
}
