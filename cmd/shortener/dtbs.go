package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func dbMnp() (a string, err error) {
	var db *sql.DB
	var vbn string
	db, err = sql.Open("sqlite3", "conf.db")
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

func dbAppgGt(id string) (a string, b int, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "shortlongurl.db")
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

func dbAppgPst(id string, longURL string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "shortlongurl.db")
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

func dbMnCf(flagRunAddr string, vbn string, fileName string) (err error) {
	db, err := sql.Open("sqlite3", "shortlongurl.db")
	if err != nil {
		return err
	}

	defer db.Close()

	sts := `
DROP TABLE IF EXISTS short_longURL;
CREATE TABLE short_longURL(id INTEGER PRIMARY KEY, short_url TEXT, longURL TEXT);`
	_, err = db.Exec(sts)

	if err != nil {
		return err
	}

	db1, err := sql.Open("sqlite3", "conf.db")
	if err != nil {
		return err
	}

	defer db1.Close()
	sts1 := `
DROP TABLE IF EXISTS cfg;
CREATE TABLE cfg (id INTEGER PRIMARY KEY, flagRunAddr TEXT, vbn TEXT, flnm TEXT);
INSERT INTO cfg(flagRunAddr, vbn, flnm) VALUES ('` + string(flagRunAddr) + `', '` + string(vbn) + `', '` + fileName + `');`
	_, err = db1.Exec(sts1)

	if err != nil {
		return err
	}
	return nil
}

func dbins(id string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "shortlongurl.db")
	if err != nil {
		return err
	}
	defer db.Close()
	Consumer, err := NewConsumer(id)
	if err != nil {
		return nil
	}
	for {
		readEvent, err_ := Consumer.ReadEvent()
		if err_ != nil {
			break
		}
		quer := `INSERT INTO short_longURL(id, short_url, longURL) VALUES (` + fmt.Sprint(readEvent.ID) + `, '` + string(readEvent.Short_URL) + `', '` + readEvent.Long_URL + `');`
		_, err = db.Exec(quer)
		if err != nil {
			return err
		}

	}
	return nil
}

func dbfln() (a string, err error) {
	var db *sql.DB
	var vbn string
	db, err = sql.Open("sqlite3", "conf.db")
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
func dbjsnpps(shortURL string, longURL string) (b int, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "shortlongurl.db")
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
