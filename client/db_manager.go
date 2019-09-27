package client

import (
	"database/sql"
	"fmt"
	"github.com/CalvinDjy/iteaGo/constant"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/system"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

const (
	CONNECTION_KEY = "connections"
	MAX_OPEN_CONNS = 20
	MAX_IDLE_CONNS = 10
	CONN_MAX_LIFE_TIME = 14400 * time.Second
)

type DatabaseConf struct {
	Driver 			string
	Ip 				string
	Port 			string
	Database 		string
	Username 		string
	Password 		string
	Charset 		string
	MaxConn 		int
	MaxIdle 		int
	ConnMaxLift 	int
}

type DbManager struct {
	databases 		map[string]*DatabaseConf
	connections 	map[string]*sql.DB
	mutex 			*sync.Mutex
}

func (dm *DbManager) Construct() {
	if c := system.Conf.GetStructMap(fmt.Sprintf("%s.%s", constant.DATABASE_KEY, CONNECTION_KEY), DatabaseConf{}); c != nil {
		dbconf := make(map[string]*DatabaseConf)
		for db, dc := range c {
			dbconf[db] = dc.(*DatabaseConf)
		}
		dm.databases = dbconf
	} else {
		panic("Can not find database config of connections!")
	}
	dm.connections = make(map[string]*sql.DB)
	dm.mutex = new(sync.Mutex)
}

func (dm *DbManager) GetDbConnection(name string) (db *sql.DB) {
	defer dm.mutex.Unlock()
	dm.mutex.Lock()
	if dm.connections[name] != nil {
		return dm.connections[name]
	}
	dm.connections[name] = dm.createConnection(name)
	ilog.Info("DB connection create success for [", name, "]")
	return dm.connections[name]
}

func (dm *DbManager) createConnection(name string) (db *sql.DB) {
	if dbconfig, ok := dm.databases[name]; ok {
		db, err := sql.Open(dbconfig.Driver, dm.dataSource(dbconfig))
		if err != nil {
			ilog.Error("databse [", name, "] open fail : ", err)
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
		ilog.Error("can not find config of databse [", name, "]")
		return nil
	}
}

func (dm *DbManager) dataSource(dbconfig *DatabaseConf) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbconfig.Username, dbconfig.Password, dbconfig.Ip, dbconfig.Port, dbconfig.Database, dbconfig.Charset)
}
