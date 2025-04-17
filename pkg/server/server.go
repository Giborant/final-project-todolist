package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort = ":8080"
)

func getPort() string {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("некорректный порт '%s', используется %s", portStr, defaultPort)
		return defaultPort
	}

	if port < 1 || port > 65535 {
		log.Printf("порт %d вне допустимого диапазона, используется %s", port, defaultPort)
		return defaultPort
	}

	return fmt.Sprintf(":%d", port)
}

func Run() {
	//api.Init()
	port := getPort()
	// Запуск приколов
	log.Printf("Сервер запущен на http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}

}
