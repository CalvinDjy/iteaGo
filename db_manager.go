package itea

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"sync"
	"fmt"
	"context"
	"encoding/json"
)

const (
	MAX_OPEN_CONNS = 20
	MAX_IDLE_CONNS = 10
	CONN_MAX_LIFE_TIME = 14400 * time.Second
)

type DbManager struct {
	Ctx				context.Context
	databases 		map[string]DatabaseConf
	connections 	map[string]*sql.DB
	mutex 			*sync.Mutex
}

func (dm *DbManager) Construct() {
	var conf map[string]DatabaseConf
	if c := dm.Ctx.Value("connections").(*json.RawMessage); c != nil {
		err := json.Unmarshal(*c, &conf)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Can not find database config of connections!")
	}
	dm.databases = conf
	dm.connections = make(map[string]*sql.DB)
	dm.mutex = new(sync.Mutex)
}

func (dm *DbManager) GetDbConnection(name string) (db *sql.DB) {
	defer dm.mutex.Unlock()
	dm.mutex.Lock()
	if dm.connections[name] != nil {
		return dm.connections[name]
	}
	log.Println("DB connection not exist for [", name, "]")
	dm.connections[name] = dm.createConnection(name)
	log.Println("DB connection create success for [", name, "]")
	return dm.connections[name]
}

func (dm *DbManager) createConnection(name string) (db *sql.DB) {
	if dbconfig, ok := dm.databases[name]; ok {
		db, err := sql.Open(dbconfig.Driver, dm.dataSource(dbconfig))
		if err != nil {
			log.Println("databse [", name, "] open fail : ", err)
			return nil
		}

		if dbconfig.MaxConn != 0 {
			db.SetMaxOpenConns(dbconfig.MaxConn)
		} else {
			db.SetMaxOpenConns(MAX_OPEN_CONNS)
		}

		if dbconfig.MaxIdle != 0 {
			db.SetMaxIdleConns(dbconfig.MaxIdle)
		} else {
			db.SetMaxIdleConns(MAX_IDLE_CONNS)
		}

		if dbconfig.ConnMaxLift != 0 {
			db.SetConnMaxLifetime(time.Duration(dbconfig.ConnMaxLift) * time.Second)
		} else {
			db.SetConnMaxLifetime(CONN_MAX_LIFE_TIME)
		}

		return db
	} else {
		log.Println("can not find config of databse [", name, "]")
		return nil
	}
}

func (dm *DbManager) dataSource(dbconfig DatabaseConf) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbconfig.Username, dbconfig.Password, dbconfig.Ip, dbconfig.Port, dbconfig.Database, dbconfig.Charset)
}
