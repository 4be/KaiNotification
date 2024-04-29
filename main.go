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
	}

}

func scrapeStart(url string) {
	c := colly.NewCollector()

	c.OnHTML("html", func(e *colly.HTMLElement) {
		isSold := strings.Contains(e.Text, "Habis")
		isAvail := strings.Contains(e.Text, "Tersedia")
		isLimit := strings.Contains(e.Text, "Tersisa")

		currentTime := time.Now().Format("15:04:05")

		if isSold {
			fmt.Println("Ticket 08/02/2024 Sold <<<<<<<<<<<", currentTime)
			return
		}

		if isAvail || isLimit {
			fmt.Printf("Tiket 08/02/2024 Available <<<<<<< %s\n", currentTime)
			sendEmail(subject, textBody)
			return
		}

		fmt.Println("System Error <<<<<<<<<<", currentTime)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

	// No need for ticker and select since we only run it once
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
