package databaseshortener

import (
	apcfg "go-prof-sprint-1/internal/app_config"
	"log"
	"testing"

	"github.com/caarlos0/env"
)

func TestNewDB(t *testing.T) {
	_, err := NewDB()
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestRunMigrateScripts(t *testing.T) {
	err := RunMigrateScripts(nil)
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseCreateShortURLPageCfg(t *testing.T) {
	_, err := DataBaseCreateShortURLPageCfg()
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
func TestDatBaseDownloadFullURLPageGet(t *testing.T) {
	_, _, err := DatBaseDownloadFullURLPageGet("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
func TestDataBaseDownloadFullURLPagePost(t *testing.T) {
	err := DataBaseDownloadFullURLPagePost("", "", "")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseCfg(t *testing.T) {
	err := DataBaseCfg("", "", "")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBasePingHandler(t *testing.T) {
	err := DataBasePingHandler()
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBasePing(t *testing.T) {
	err := DataBasePing("", "")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseInsert(t *testing.T) {
	err := DataBaseInsert("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseFileNameSelect(t *testing.T) {
	_, err := DataBaseFileNameSelect()
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseJSONPage(t *testing.T) {
	_, err := DataBaseJSONPage("", "", "")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
func TestDataBaseFilePost(t *testing.T) {
	err := DataBaseFilePost("", "", "")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseCheckURLExistance(t *testing.T) {
	_, _, err := DataBaseCheckURLExistance("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseStartConfig(t *testing.T) {
	err := DataBaseStartConfig("")
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseSelfConfigGet(t *testing.T) {
	_, _, err := DataBaseSelfConfigGet()
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseSelfConfigUpdate(t *testing.T) {
	_, _, err := DataBaseSelfConfigGet()
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseGetAllURLs(t *testing.T) {
	_, err := DataBaseGetAllURLs("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
func TestDataBaseDeleteURL(t *testing.T) {
	err := DataBaseDeleteURL("", "")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestDataBaseDeleteURLs(t *testing.T) {
	err := DataBaseDeleteURLs(nil, "")
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
func TestDataBaseCheckURLDelition(t *testing.T) {
	_, err := DataBaseCheckURLDelition("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestNewDBGood(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did  panic")
		}
	}()
	_ = DataBaseStartConfig(":8080")
	_, err := NewDB()
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestRunMigrateScriptsBad(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = DataBaseStartConfig(":8080")
	err := RunMigrateScripts(nil)
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestNewDBGood1(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did  panic")
		}
	}()
	var cfg apcfg.Config
	err := env.Parse(&cfg)
	flagRunAddr, apiRunAddr, fileName, databaseDSN := apcfg.ParseFlags()
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
	log.Println(cfg)
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
	err = DataBaseCfg(flagRunAddr, apiRunAddr, fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = DataBaseInsert(fileName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DataBaseCreateShortURLPageCfg()
	//_, err := NewDB()
	if err != nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
