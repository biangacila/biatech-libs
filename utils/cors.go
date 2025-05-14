package utils

import (
	"os"
	"strings"
)

func GetAllowedMethods() []string {
	if os.Getenv("ALLOWED_METHODS") != "" {
		return strings.Split(os.Getenv("ALLOWED_METHODS"), ",")
	}
	return []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}
}

func GetAllowedOrigins() []string {
	if os.Getenv("ORIGINS") != "" {
		return strings.Split(os.Getenv("ORIGINS"), ",")
	}
	return []string{
		"http://localhost",
		"https://cloudcalls.easipath.com",
		"https://scrutinize.biacibenga.com",
	}
}
func GetAllowedHeaders() []string {
	if os.Getenv("ALLOWED_HEADERS") != "" {
		return strings.Split(os.Getenv("ALLOWED_HEADERS"), ",")
	}
	return []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Origin",
		"Authorization", "Origin",
		"X-Requested-With",
		"Accept",
		"Content-Type",
		"user-code",
		"org-code",
	}
}
func GetElasticSearchServerUrl(fallbackUrl string) string {
	if os.Getenv("ELASTICSEARCH_SERVER_URL") != "" {
		return os.Getenv("ELASTICSEARCH_SERVER_URL")
	}
	return fallbackUrl
}
func GetKafkaServerBroker(fallbackUrl string) []string {
	if os.Getenv("KAFKA_SERVER_BROKER") != "" {
		return strings.Split(os.Getenv("KAFKA_SERVER_BROKER"), ",")
	}
	return []string{fallbackUrl}
}
func GetKafkaServerBrokerOne(fallbackUrl string) string {
	if os.Getenv("KAFKA_SERVER_BROKER") != "" {
		return os.Getenv("KAFKA_SERVER_BROKER")
	}
	return fallbackUrl
}

func GetAuth1ServiceUrl(fallbackUrl string) string {
	server := fallbackUrl
	if os.Getenv("AUTH1_SERVICE") != "" {
		server = os.Getenv("AUTH1_SERVICE")
	}
	return server
}
func GetWebsocketServiceUrl() string {
	server := "https://cloudcalls.easipath.com/backend-telcowebsocket/api"
	if os.Getenv("WEBSOCKET_SERVICE") != "" {
		server = os.Getenv("WEBSOCKET_SERVICE")
	}
	return server
}
func GetWebsocketServiceGRPC() string {
	server := "safer.easipath.com:50051"
	if os.Getenv("WEBSOCKET_GRPC_SERVICE") != "" {
		server = os.Getenv("WEBSOCKET_GRPC_SERVICE")
	}
	return server
}
