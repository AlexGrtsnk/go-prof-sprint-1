package databaseshortener

import (
	"database/sql"
	"fmt"
	"log"

	bn "go-prof-sprint-1/internal/bindata"

	Flw "go-prof-sprint-1/internal/json_parser"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const drriver = "sqlite3"
const dbbName = "shortenerdbs.db"

func NewDB() (*sql.DB, error) {
	dbname, driverTemp, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	sqliteDB, err := sql.Open(driverTemp, dbname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sqlite DB")
	}

	return sqliteDB, nil
}

func RunMigrateScripts(db *sql.DB) error {
	var driver database.Driver
	var err error
	dbNameTemp, _, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	if dbNameTemp == dbbName {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	} else {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	}
	if err != nil {
		return fmt.Errorf("creating db driver failed %s", err)
	}

	res := bindata.Resource(bn.AssetNames(),
		func(name string) ([]byte, error) {
			return bn.Asset(name)
		})

	d, _ := bindata.WithInstance(res)
	m, err := migrate.NewWithInstance("go-bindata", d, dbNameTemp, driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}
	if dbNameTemp == dbbName {
		_ = m.Steps(-1)
		err = m.Steps(1)
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating database failed %s", err)
		}
	} else {
		_ = m.Down()
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating database failed %s", err)
		}
	}
	return nil
}

func DataBaseCreateShortURLPageCfg() (apiRunAddr_ string, err error) {
	var db *sql.DB
	var apiRunAddr string
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return "", err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return "", err
	}
	defer db.Close()
	quer := "SELECT apiRunAddr FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", rows.Err()
	}
	rows.Next()
	err = rows.Scan(&apiRunAddr)
	if err != nil {
		return "", err
	}
	return apiRunAddr, nil
}

func DatBaseDownloadFullURLPageGet(id string) (longURL_ string, flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return "", 0, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return "", 0, err
	}
	defer db.Close()
	quer := "SELECT longURL FROM short_longURL WHERE short_url = '" + string(id) + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", 0, rows.Err()
	}
	rows.Next()
	var longURL string
	err = rows.Scan(&longURL)
	if err != nil {
		return "", 0, err
	}
	return longURL, 1, nil
}

func DataBaseDownloadFullURLPagePost(id string, longURL string, tknm string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO short_longURL(short_url, longURL, tknm) VALUES ('` + string(id) + `', '` + string(longURL) + `', '` + tknm + `')`
	//quer := `INSERT INTO cfg(flagRunAddr, apiRunAddr, flnm) VALUES ('` + string(id) + `', '` + string(longURL) + `', '` + tknm[:10] + `')`
	//quer := `INSERT INTO cfg(flagRunAddr, apiRunAddr, flnm) VALUES ('` + string(id) + `', '` + string(longURL) + `', '` + tknm[:10] + `')`

	//quer := "INSERT INTO short_longURL(short_url, longURL, tknm) VALUES ('qwe', 'qwe1', 'qwe2')"
	//quer := "INSERT INTO short_longURL VALUES (1, 'qwe12', 'qwe', 'qwe1')"
	//quer := `INSERT INTO short_longURL(short_url, longURL, token) VALUES ('` + string(id) + `', '` + string(longURL) + `', '` + string(id) + `')`
	//quer := "INSERT INTO short_longURL(short_url, longURL) VALUES('" + string(id) + "', '" + string("dasda") + "');"
	fmt.Println(quer)
	_, err = db.Exec(quer)
	if err != nil {
		fmt.Println("adsadasdasdawdawdawe")
		return err
	}
	return nil
}

func DataBaseCfg(flagRunAddr string, apiRunAddr string, fileName string) (err error) {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	err = RunMigrateScripts(db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	quer := `INSERT INTO cfg(flagRunAddr, apiRunAddr, flnm) VALUES ('` + string(flagRunAddr) + `', '` + string(apiRunAddr) + `', '` + fileName + `')`
	fmt.Println(quer)
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}
func DataBasePingHandler() (err error) {
	_, driverTemp, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	dbName := fmt.Sprintf("host=%s port=%s  user=%s password=%s dbname=%s sslmode=disable",
		`postgres`, `5432`, `postgres`, `postgres`, `praktikum`)
	err = DataBasePing(dbName, driverTemp)
	if err != nil {
		err = DataBaseSelfConfigUpdate(dbbName, drriver)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func DataBasePing(dbbname string, driver string) (err error) {
	var db *sql.DB
	var res string
	db, err = sql.Open(driver, dbbname)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "SELECT 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}
	rows.Next()
	err = rows.Scan(&res)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseInsert(id string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	Consumer, err := Flw.NewConsumer(id)
	if err != nil {
		return nil
	}
	for {
		readEvent, err_ := Consumer.ReadEvent()
		if err_ != nil {
			break
		}
		quer := `INSERT INTO short_longURL(short_url, longURL, tknm) VALUES ('` + string(readEvent.ShortURL) + `', '` + readEvent.LongURL + `', '` + readEvent.Tknm + `');`
		_, err = db.Exec(quer)
		if err != nil {
			return err
		}

	}
	return nil
}

func DataBaseFileNameSelect() (flnm string, err error) {
	var db *sql.DB
	var apiRunAddr string
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return "", err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return "", err
	}
	defer db.Close()
	quer := "SELECT flnm FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", rows.Err()
	}
	rows.Next()
	err = rows.Scan(&apiRunAddr)
	if err != nil {
		return "", err
	}
	return apiRunAddr, nil
}
func DataBaseJSONPage(shortURL string, longURL string, tknm string) (b int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	//quer := "SELECT id FROM short_longURL WHERE short_url = '" + string(shortURL) + "' and longURL ='" + longURL + "' and tknm = '" + tknm + "';"
	quer := "SELECT id FROM short_longURL WHERE short_url = '" + string(shortURL) + "' and longURL ='" + longURL + "';"
	fmt.Println(quer)
	rows, err := db.Query(quer)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return 0, rows.Err()
	}
	rows.Next()
	var id int
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DataBaseFilePost(shortURL string, longURL string, tknm string) (err error) {
	fileName, err := DataBaseFileNameSelect()
	if err != nil {
		log.Fatal(err)
	}
	Producer, err := Flw.NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()
	id, err := DataBaseJSONPage(shortURL, longURL, tknm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("12344, ", tknm)
	var events = []*Flw.Event{{ID: id, ShortURL: shortURL, LongURL: longURL, Tknm: tknm, DelFlag: 0}}
	err = Producer.WriteEvent(events[0])
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func DataBaseCheckURLExistance(longURL string) (shortURL string, flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return "", 0, err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return "", 0, err
	}
	defer db.Close()
	var shoortURL string
	if err := db.QueryRow("SELECT short_url FROM short_longURL WHERE longURL = '" + string(longURL) + "';").Scan(&shoortURL); err != nil {
		if err == sql.ErrNoRows {
			return "", 0, nil
		}
		return "", 0, err
	}
	return shoortURL, 1, nil
}

func DataBaseStartConfig(dbName string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return err
	}
	defer db.Close()
	var driver string
	if dbName != "localhost" {
		driver = "pgx"
	} else {
		driver = "sqlite3"
		dbName = dbbName
	}
	sts1 := `
	DROP TABLE IF EXISTS cfg;
	CREATE TABLE cfg (id INTEGER PRIMARY KEY, dbbname TEXT, driver TEXT);
	INSERT INTO cfg(dbbname, driver) VALUES ('` + string(dbName) + `', '` + string(driver) + `');`
	_, err = db.Exec(sts1)

	if err != nil {
		return err
	}
	return nil
}

func DataBaseSelfConfigGet() (dbbname string, driver string, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return "", "", err
	}
	quer := "SELECT dbbname FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", "", rows.Err()
	}
	rows.Next()
	var dbNameTemp string
	err = rows.Scan(&dbNameTemp)
	if err != nil {
		return "", "", err
	}
	quer = "SELECT driver FROM cfg WHERE id = 1;"
	rows, err = db.Query(quer)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", "", rows.Err()
	}
	rows.Next()
	var driverTemp string
	err = rows.Scan(&driverTemp)
	if err != nil {
		return "", "", err
	}
	return dbNameTemp, driverTemp, nil
}
func DataBaseSelfConfigUpdate(dbbname string, driver string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "UPDATE cfg SET dbbname='" + dbbname + "', '" + "driver='" + driver + "' WHERE id=1;"
	_, err = db.Exec(quer)

	if err != nil {
		return err
	}
	return nil
}
