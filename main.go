package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"gopkg.in/gomail.v2"
)

const (
	baseURL   = "https://booking.kai.id"
	endpoint  = "/search?origination=n4aQD3sKMNnw1kDOEJg3zQ%3D%3D&destination=ykzJ6TGEvaQymY%2BY2hXnmA%3D%3D&tanggal=pduPkj1YvIsCQzvmTZbKcULjz2L8MsfAGymVI9PsRno%3D&adult=8Ho6OKf8tFhcLcNlxENpVg%3D%3D&infant=ThM%2B2IHux7cEkbUR6e9yBg%3D%3D&book_type="
	targetURL = baseURL + endpoint
)

var (
	fromEmail = "gakpernahlupa@gmail.com"
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

		time.Sleep(2 * time.Second)
		i -= 1
	}

}

func scrapeStart(url string) {
	// var Green = "\033[32m"
	// var Reset = "\033[0m"

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

		fmt.Print("isi text = ", e.Text)

		currentTime := time.Now().Format("15:04:05")

		if isSold {
			fmt.Printf("\033[1mTicket KAI \033[33;1m[%s-%s %s] \033[31mSOLD \033[0m<<<<<<<<<<< %s\n", stasiunAsal, stasiunTujuan, tanggalPesanan, currentTime)
			return
		}

		if isAvail || isLimit {
			fmt.Printf("\033[1mTicket KAI \033[33;1m[%s-%s %s] \033[32;1mAVAILABLE \033[0m<<<<< %s\n", stasiunAsal, stasiunTujuan, tanggalPesanan, currentTime)
			sendEmail(subject, textBody)
			return
		}

		fmt.Println("System Error <<<<<<<<<<", currentTime)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

}

func sendEmail(subject, textBody string) {
	htmlBody := fmt.Sprintf("<b>Link to buy:</b> <a href=\"%s\">HERE</a>", targetURL)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fromEmail)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", textBody)
	mailer.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, "aztwwnpupmrhfhtk")

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Println("Error sending email:", err)
	}
}
