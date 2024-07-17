package databaseshortener

import (
	"database/sql"
	"fmt"
	apcfg "go-prof-sprint-1/internal/app_config"
	"log"

	"github.com/caarlos0/env"
)

func ExampetDataBaseCreateShortURLPageCfg() {
	var cfg apcfg.Config
	err := env.Parse(&cfg)
	flagRunAddr, apiRunAddr, fileName, databaseDSN, enableHTTPS, config := apcfg.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress != "" {
		flagRunAddr = "8080"
	}
	if cfg.BaseURL != "" {
		apiRunAddr = cfg.BaseURL
	}
	if cfg.FileStoragePath != "" {
		fileName = cfg.FileStoragePath
	}
	if cfg.DatabaseDSN != "" {
		databaseDSN = cfg.DatabaseDSN
	}
	log.Println(enableHTTPS)
	log.Println(config)
	err = DataBaseStartConfig(databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if databaseDSN != "localhost" {
		err = DataBasePingHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	_ = DataBaseCfg(flagRunAddr, apiRunAddr, fileName)
	out1, _ := DataBaseCreateShortURLPageCfg()
	fmt.Println(out1)
	var db *sql.DB
	db, _ = sql.Open("sqlite3", "cfg.db")
	defer db.Close()
	quer := "UPDATE cfg SET apiRunAddr='9090' WHERE id=1;"
	_, _ = db.Exec(quer)
	out2, _ := DataBaseCreateShortURLPageCfg()
	fmt.Println(out2)
	quer1 := "UPDATE cfg SET apiRunAddr=':8080' WHERE id=1;"
	_, _ = db.Exec(quer1)
	out3, _ := DataBaseCreateShortURLPageCfg()
	fmt.Println(out3)
}
