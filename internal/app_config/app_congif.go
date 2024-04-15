package appconfig

import (
	"flag"
)

type Config struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Home            string `env:"HOME"`
	ServerAddress   string `env:"serverAddress"`
	BaseURL         string `env:"baseURL"`
}

func ParseFlags() (a string, b string, f string) {
	var flagRunAddr string
	var apiRunAddr string
	var fileName string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&apiRunAddr, "b", "http://localhost:8080", "api page existance url adress")
	flag.StringVar(&fileName, "f", "text.txt", "txt file with short and long urls")
	flag.Parse()
	if flagRunAddr != "localhost:8080" && apiRunAddr == "http://localhost:8080" {
		apiRunAddr = "http://" + flagRunAddr
	}
	if flagRunAddr == "localhost:8080" && apiRunAddr != "http://localhost:8080" {
		flagRunAddr = apiRunAddr[7:]
	}
	return flagRunAddr, apiRunAddr, fileName
}
