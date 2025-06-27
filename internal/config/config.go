package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	SMTPHost         string
	SMTPPort         int
	SMTPUsername     string
	SMTPPassword     string
	SenderEmail      string
	OTPExpiryMinutes int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	otpExpiryMinutes, _ := strconv.Atoi(os.Getenv("OTP_EXPIRY_MINUTES"))

	return &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         smtpPort,
		SMTPUsername:     os.Getenv("SMTP_USERNAME"),
		SMTPPassword:     os.Getenv("SMTP_PASSWORD"),
		SenderEmail:      os.Getenv("SENDER_EMAIL"),
		OTPExpiryMinutes: otpExpiryMinutes,
	}
}
