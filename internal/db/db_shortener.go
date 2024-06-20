package databaseshortener

import (
	"database/sql"
	"fmt"
	"sync"

	bn "go-prof-sprint-1/internal/bindata"

	flw "go-prof-sprint-1/internal/json_parser"

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

// NewDB создает новую sql сущность с конкретным необходимым нам драйвером
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

// RunMigrateScripts запускает скрипты миграции в зависимости от версии драйвера
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

// DataBaseCreateShortURLPageCfg возвращает адрес запуска api страницы для сокращения
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

// DatBaseDownloadFullURLPageGet по сокращенному url возращает его полную версию и flag, 1 - url успешно найден
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

// DataBaseDownloadFullURLPagePost добавляет в базу данных новую запись, содержащую короткий и длинный url с токеном пользователя, сделавшим запрос на сокращение
func DataBaseDownloadFullURLPagePost(id string, longURL string, token string) (err error) {
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
	quer := `INSERT INTO short_longURL(short_url, longURL, token) VALUES ('` + string(id) + `', '` + string(longURL) + `', '` + token + `')`

	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

// DataBaseCfg конфигурирует базу данных
func DataBaseCfg(flagRunAddr string, apiRunAddr string, fileName string) (err error) {
	db, err := NewDB()
	if err != nil {
		return err
	}

	defer db.Close()
	err = RunMigrateScripts(db)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO cfg(flagRunAddr, apiRunAddr, flnm) VALUES ('` + string(flagRunAddr) + `', '` + string(apiRunAddr) + `', '` + fileName + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

// DataBasePingHandler хендлер для проверки отклика базы данных
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

// DataBasePing проверяет, отвечает ли наша база данных
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

// DataBaseInsert вставляет в базу новую запись из формата json
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
	Consumer, err := flw.NewConsumer(id)
	if err != nil {
		return nil
	}
	for {
		readEvent, err_ := Consumer.ReadEvent()
		if err_ != nil {
			break
		}
		quer := `INSERT INTO short_longURL(short_url, longURL, token) VALUES ('` + string(readEvent.ShortURL) + `', '` + readEvent.LongURL + `', '` + readEvent.Token + `');`
		_, err = db.Exec(quer)
		if err != nil {
			return err
		}

	}
	return nil
}

// DataBaseFileNameSelect возвращает место хранения записей о сокращенныз url
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

// DataBaseJSONPage возвращает id записи с конкретными параметрами
func DataBaseJSONPage(shortURL string, longURL string, token string) (b int, err error) {
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

// DataBaseFilePost пишет запись о новом сокращенном url в файл, предназначенный для их хранения
func DataBaseFilePost(shortURL string, longURL string, token string) (err error) {
	fileName, err := DataBaseFileNameSelect()
	if err != nil {
		return err
	}
	Producer, err := flw.NewProducer(fileName)
	if err != nil {
		return err
	}
	defer Producer.Close()
	id, err := DataBaseJSONPage(shortURL, longURL, token)
	if err != nil {
		return err
	}
	var events = []*flw.Event{{ID: id, ShortURL: shortURL, LongURL: longURL, Token: token, DelFlag: 0}}
	err = Producer.WriteEvent(events[0])
	if err != nil {
		return err
	}
	return nil
}

// DataBaseCheckURLExistance проверяет наличие сокращенного url в базе данных
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

// DataBaseStartConfig конфигурирует изначаотный конфиг для записи
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

// DataBaseSelfConfigGet возвращает конфиг, на котором сейчас работает база данных
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

// DataBaseSelfConfigUpdate изменяет конфиг, записанный в базу данных
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

// DataBaseGetAllURLs возвращает все url, записанные пользователем с конкретным токеном
func DataBaseGetAllURLs(token string) (answb []AnswerBatch, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	apiRunAddr, err := DataBaseCreateShortURLPageCfg()
	if err != nil {
		return nil, err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return nil, err
	}
	quer := "SELECT short_url, longURL from short_longURL where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return nil, err
	}
	flag := 0
	for rows.Next() {
		answ := new(AnswerBatch)
		err = rows.Scan(&answ.ShortURL, &answ.OriginalURL)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		answ.ShortURL = apiRunAddr + "/" + answ.ShortURL
		answb = append(answb, *answ)
		flag = 1
	}
	if flag == 0 {
		return nil, nil
	}
	return
}

// DataBaseDeleteURL функция мягкого удаления записи о сокращенном url из базы данных
func DataBaseDeleteURL(longURL string, token string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	quer := "UPDATE short_longURL SET delFlag=1 WHERE short_url = '" + longURL + "' and token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil

}

// DataBaseDeleteURLs обертка для удаления url через горутины
func DataBaseDeleteURLs(ids flw.DeleteList, token string) (err error) {
	var wg sync.WaitGroup
	for _, produceItem := range ids {

		wg.Add(1)
		a := string(produceItem)
		go func(a string) {
			_ = DataBaseDeleteURL(a, token)
			wg.Done()
		}(a)

	}
	return nil
}

// DataBaseCheckURLDelition функция проверки удаления конкретного сокращенного url, 0 - url был удален
func DataBaseCheckURLDelition(shortURL string) (flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 1, err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 1, err
	}
	defer db.Close()
	quer := "SELECT delFlag FROM short_longURL where short_url='" + shortURL + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return 1, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return 1, rows.Err()
	}
	rows.Next()
	err = rows.Scan(&flag)
	if err != nil {
		return 1, err
	}
	return flag, nil
}

// AnswerBatch тип ответа от базы с сокращенным и полным url
type AnswerBatch struct {
	// ShortURL - сокращенный url
	ShortURL string `json:"short_url"`
	// OriginalURL - полный url
	OriginalURL string `json:"original_url"`
}
