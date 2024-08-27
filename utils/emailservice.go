package utils

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	gomail "gopkg.in/mail.v2"
)

func SendVerificationEmail(email, token string) error {
	// Set up the email
	m := gomail.NewMessage()
	m.SetHeader("From", "tarikuamisganaw@gmail.com") // Change to your sender email
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification")
	verificationLink := fmt.Sprintf("http://localhost:8080/users/verify-email?token=%s&email=%s", token, email)
	m.SetBody("text/plain", fmt.Sprintf("Click the link below to verify your email:\n\n%s", verificationLink))

	// Set up the SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-email-password-or-app-password")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
func ValidateEmailToken(tokenString, email string) bool {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})

	// Check if token is valid and email matches
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject == email
	}

	return false
}

// GenerateTokenHash generates a token hash to be used in password reset or email verification
func GenerateTokenHash(email string) string {
	hash := sha256.New()
	hash.Write([]byte(email + fmt.Sprintf("%d", time.Now().Unix())))
	return hex.EncodeToString(hash.Sum(nil))
}
