package utils

import (
	"github.com/joho/godotenv"
	"os"
)

type SmtpInfo struct {
	Host     string
	Port     string
	Username string
	Password string
}

func GetSmtpInfo() (info SmtpInfo, err error) {
	if os.Getenv("SMTP_SERVER") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return info, err
		}
	}

	info.Host = os.Getenv("SMTP_SERVER")
	info.Port = os.Getenv("SMTP_PORT")
	info.Username = os.Getenv("SMTP_USERNAME")
	info.Password = os.Getenv("SMTP_PASSWORD")

	return info, err
}
func GoogleAuthCallbackUri() string {
	if os.Getenv("GOOGLE_CALLBACK_URL") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return ""
		}
	}
	return os.Getenv("GOOGLE_CALLBACK_URL")
}
