// Модуль appconfig отвечает за изначальную конфигурацию приложения по сокращению url
package appconfig

import (
	"flag"
)

// Config франит в себе настройку проекта
type Config struct {
	// FileStoragePath содержит в себе путь хранения текстового файла с url
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// DatabaseDSN содержит в себе адрес psql сервера
	DatabaseDSN string `env:"DATABASE_DSN"`
	// Home содержит в себе адрес изначальной директории
	Home string `env:"HOME"`
	// ServerAddress содержит в себе адрес запуска нашего приложения
	ServerAddress string `env:"serverAddress"`
	// BaseURL содержит в себе начало пути адреса для сокращения url
	BaseURL string `env:"baseURL"`
	// При передаче флага EnableHTTPS сервер запуститься по методу http.ListenAndServeTLS
	EnableHTTPS bool `env:"ENABLE_HTTPS"`
	// Config содержит в себе имя файла в формате JSON с параметрами кофигурации приложения
	Config string `env:"CONFIG"`
	// TrustedSubnet содержит в себе адресмдоверенной подсети
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

// ParseFlags возвращает флаги, необходимын для работы приложения
func ParseFlags() (a string, b string, f string, v string, s bool, c string, t string) {
	var flagRunAddr string
	var apiRunAddr string
	var fileName string
	var dataBaseAddress string
	var enableHTTPS bool
	var config string
	var trustedSubnet string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&apiRunAddr, "b", "http://localhost:8080", "api page existance url adress")
	flag.StringVar(&fileName, "f", "text.txt", "txt file with short and long urls")
	flag.StringVar(&dataBaseAddress, "d", "localhost", "databaseport")
	flag.BoolVar(&enableHTTPS, "s", false, "if flag is on, using http.ListenAndServeTLS")
	flag.StringVar(&config, "c", "ns", "JSON file adress with exe config")
	flag.StringVar(&trustedSubnet, "t", "ns", "JSON file adress with exe config")
	flag.Parse()
	if flagRunAddr != "localhost:8080" && apiRunAddr == "http://localhost:8080" {
		apiRunAddr = "http://" + flagRunAddr
	}
	if flagRunAddr == "localhost:8080" && apiRunAddr != "http://localhost:8080" {
		flagRunAddr = apiRunAddr[7:]
	}
	return flagRunAddr, apiRunAddr, fileName, dataBaseAddress, enableHTTPS, config, trustedSubnet
}
