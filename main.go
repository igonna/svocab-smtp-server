package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Struct to capture incoming JSON request body
type EmailRequest struct {
	SMTPServerAddr string `json:"smtp_server_addr"` // Full address including port (e.g., smtp.example.com:587)
	Username       string `json:"username"`
	Password       string `json:"password"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	To             string `json:"to"`
}

// Function to send the email
func sendEmail(emailRequest EmailRequest) error {
	// Split the full server address into server and port
	serverAndPort := strings.Split(emailRequest.SMTPServerAddr, ":")
	if len(serverAndPort) != 2 {
		return fmt.Errorf("invalid SMTP server address format, expected 'server:port'")
	}
	smtpServer, smtpPort := serverAndPort[0], serverAndPort[1]

	// Set up the authentication information.
	auth := smtp.PlainAuth(
		"",
		emailRequest.Username,
		emailRequest.Password,
		smtpServer,
	)

	// The message to be sent
	message := []byte("To: " + emailRequest.To + "\r\n" +
		"Subject: " + emailRequest.Subject + "\r\n" +
		"\r\n" +
		emailRequest.Message + "\r\n")

	// Send the email
	err := smtp.SendMail(
		smtpServer+":"+smtpPort,
		auth,
		emailRequest.Username,
		[]string{emailRequest.To},
		message,
	)

	return err
}

// POST handler to send an email
func sendEmailHandler(c *gin.Context) {
	var emailRequest EmailRequest

	// Bind incoming JSON request to struct
	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Send the email
	err := sendEmail(emailRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	// Respond with success
	c.JSON(200, gin.H{"status": "Email sent successfully"})
}

func main() {
	// Create a new Gin router
	r := gin.Default()

	// POST route to send email
	r.POST("/send-email", sendEmailHandler)

	// Start the server
	err := r.Run(":5556")
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
