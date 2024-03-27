package dtbs

import (
	"database/sql"
	"fmt"
	"log"

	bn "go-prof-sprint-1/cmd/bindata"

	flw "go-prof-sprint-1/cmd/flwrk"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func NewDB(dbPath string) (*sql.DB, error) {
	sqliteDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sqlite DB")
	}

	return sqliteDB, nil
}

func RunMigrateScripts(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 db driver failed %s", err)
	}

	res := bindata.Resource(bn.AssetNames(),
		func(name string) ([]byte, error) {
			return bn.Asset(name)
		})

	d, _ := bindata.WithInstance(res)
	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}
	_ = m.Down()
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrating database failed %s", err)
	}

	return nil
}

func DBMainPageCfg() (a string, err error) {
	var db *sql.DB
	var vbn string
	db, err = sql.Open("sqlite3", "teeest.db")
	if err != nil {
		return "", err
	}
	defer db.Close()
	quer := "SELECT vbn FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", rows.Err()
	}
	rows.Next()
	err = rows.Scan(&vbn)
	if err != nil {
		return "", err
	}
	return vbn, nil
}

func DBAppgGt(id string) (a string, b int, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "teeest.db")
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

func DBAppgPst(id string, longURL string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "teeest.db")
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "INSERT INTO short_longURL(short_url, longURL) VALUES('" + string(id) + "', '" + string(longURL) + "');"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DBMnCf(flagRunAddr string, vbn string, fileName string) (err error) {
	//db, err := sql.Open("sqlite3", "shortlongurl.db")
	db, err := NewDB("teeest.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = RunMigrateScripts(db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	quer := `INSERT INTO cfg(flagRunAddr, vbn, flnm) VALUES ('` + string(flagRunAddr) + `', '` + string(vbn) + `', '` + fileName + `');`
	_, err = db.Exec(quer)

	if err != nil {
		return err
	}
	return nil
}

func DBIns(id string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "teeest.db")
	if err != nil {
		return err
	}
	defer db.Close()
	Consumer, err := flw.NewConsumer(id)
	if err != nil {
		return nil
	}
	var eventTable []flw.Event
	flag := 0
	for {
		readEvent, err_ := Consumer.ReadEvent()
		if err_ != nil {
			break
		}
		eventTable = append(eventTable, *readEvent)
		flag = 1

	}
	if flag == 1 {
		sqlStr := "INSERT INTO short_longURL(id, short_url, longURL) VALUES "
		vals := []interface{}{}

		for _, row := range eventTable {
			sqlStr += "(?, ?, ?),"
			vals = append(vals, row.ID, row.ShortURL, row.LongURL)
		}
		sqlStr = sqlStr[0 : len(sqlStr)-1]
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(vals...)
		if err != nil {
			return err
		}
	}
	return nil
}

func DBFileNameCfg() (a string, err error) {
	var db *sql.DB
	var vbn string
	db, err = sql.Open("sqlite3", "teeest.db")
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
	err = rows.Scan(&vbn)
	if err != nil {
		return "", err
	}
	return vbn, nil
}
func DBjsnpps(shortURL string, longURL string) (b int, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "teeest.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	quer := "SELECT id FROM short_longURL WHERE short_url = '" + string(shortURL) + "' and longURL ='" + longURL + "';"
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

func Flpst(shortURL string, longURL string) (err error) {
	fileName, err := DBFileNameCfg()
	if err != nil {
		log.Fatal(err)
	}
	Producer, err := flw.NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()
	id, err := DBjsnpps(shortURL, longURL)
	if err != nil {
		log.Fatal(err)
	}
	var events = []*flw.Event{{ID: id, ShortURL: shortURL, LongURL: longURL}}
	err = Producer.WriteEvent(events[0])
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
