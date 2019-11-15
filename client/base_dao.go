package client

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/CalvinDjy/iteaGo/constant"
	"github.com/CalvinDjy/iteaGo/ilog"
	"reflect"
	"strings"
	"time"
)

type BaseDao struct {
	Ctx context.Context
	DbManager *DbManager `wired:"true"`
	connection *sql.DB
	debug bool
}

/**
 * 初始化Dao数据库连接
 */
func (bd *BaseDao) Init(dao interface{}) {
	v := reflect.ValueOf(dao)
	t := reflect.TypeOf(dao)
	if t == reflect.TypeOf(&BaseDao{}) {
		return
	}
	dbName := v.Elem().FieldByName("DbName")
	if !dbName.IsValid() || strings.EqualFold(dbName.String(), ""){
		ilog.Error("DbName of [", t.Elem().Name(), "] is not set or empty")
		return
	}
	bd.connection = bd.DbManager.GetDbConnection(dbName.String())
	bd.debug = bd.Ctx.Value(constant.DEBUG).(bool)
}

/**
 * 插入一条记录
 */
func (bd *BaseDao) Insert(table string, params map[string]interface{}) (int64, error) {
	var query, keys, values bytes.Buffer

	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql insert】耗时：", time.Since(start), "【", query.String(), "】")
		}()
	}

	var args []interface{}
	query.WriteString("insert into `")
	query.WriteString(table)
	query.WriteString("` (")
	for k, v := range params{
		keys.WriteString("`")
		keys.WriteString(k)
		keys.WriteString("`,")
		values.WriteString("?,")
		args = append(args, v)
	}

	query.WriteString(string(keys.Bytes()[0:keys.Len() - 1]))
	query.WriteString(") values (")
	query.WriteString(string(values.Bytes()[0:values.Len() - 1]))
	query.WriteString(")")

	r, err := bd.connection.Exec(query.String(), args...)
	if err != nil {
		ilog.Error("insert error : ", err)
		return -1, err
	}
	return r.LastInsertId()
}

/**
 * 插入多条记录
 */
func (bd *BaseDao) MultiInsert(table string, params []map[string]interface{}) (int, error) {
	if len(params) == 0 {
		return 0, nil
	}

	var query, keys, values bytes.Buffer
	var statement string
	var args []interface{}

	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql multi insert】耗时：", time.Since(start), "【", statement, "】")
			ilog.Info("params : ", args)
		}()
	}

	query.WriteString("insert into `")
	query.WriteString(table)
	query.WriteString("` (")
	for k := range params[0]{
		keys.WriteString("`")
		keys.WriteString(k)
		keys.WriteString("`,")

	}

	query.WriteString(string(keys.Bytes()[0:keys.Len() - 1]))
	query.WriteString(") values ")

	for _, item := range params {
		values.WriteString("(")
		for _, v := range item {
			values.WriteString("?,")
			args = append(args, v)
		}
		query.WriteString(string(values.Bytes()[0:values.Len() - 1]))
		query.WriteString("),")
		values.Reset()
	}

	statement = string(query.Bytes()[0:query.Len() - 1])
	_, err := bd.connection.Exec(statement, args...)
	if err != nil {
		ilog.Error("multi insert error : ", err)
		return -1, err
	}
	return len(params), nil
}

/**
 * 查询一条记录
 */
func (bd BaseDao) Find(sql string, p ...interface{}) *sql.Row{
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql find】耗时：", time.Since(start), "【", sql, "】")
		}()
	}
	return bd.connection.QueryRow(sql, p...)
}

/**
 * 查询记录
 */
func (bd BaseDao) Select(sql string, p ...interface{}) (*sql.Rows, error){
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql select】耗时：", time.Since(start), "【", sql, "】")
		}()
	}
	return bd.connection.Query(sql, p...)
}

/**
 * 更新记录
 */
func (bd BaseDao) Update(sql string, p ...interface{}) (int64, error) {
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql update】耗时：", time.Since(start), "【", sql, "】")
		}()
	}

	r, err := bd.connection.Exec(sql, p...)
	if err != nil {
		ilog.Error("update error : ", err)
		return -1, err
	}
	return r.RowsAffected()
}

/**
 * 根据主键id更新记录
 */
func (bd BaseDao) UpdateById(table string, id int32, params map[string]interface{}) (int64, error) {
	var query, set bytes.Buffer

	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql update】耗时：", time.Since(start), "【", query.String(), "】")
		}()
	}

	var args []interface{}
	query.WriteString("update `")
	query.WriteString(table)
	query.WriteString("` set ")
	for k, v := range params {
		set.WriteString(fmt.Sprintf("`%s`=?, ", k))
		args = append(args, v)
	}
	query.WriteString(string(set.Bytes()[0:set.Len() - 2]))
	query.WriteString(" where `id`=?")
	args = append(args, id)

	r, err := bd.connection.Exec(query.String(), args...)
	if err != nil {
		ilog.Error("update error : ", err)
		return -1, err
	}
	return r.RowsAffected()
}

/**
 * 删除记录
 */
func (bd BaseDao) Delete(sql string, p ...interface{}) (int64, error) {
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql delete】耗时：", time.Since(start), "【", sql, "】")
		}()
	}
	r, err := bd.connection.Exec(sql, p...)
	if err != nil {
		ilog.Error("delete error : ", err)
		return -1, err
	}
	return r.RowsAffected()
}

/********************************************************  事务方法  **************************************************/

/**
 * 事务
 */
func (bd BaseDao) Transaction(f func(*sql.Tx) (interface{}, error)) (interface{}, error){
	connection, err := bd.connection.Begin()
	if err != nil {
		ilog.Error("Transaction open error", err)
		return nil, err
	}
	result, err := f(connection)
	if err != nil {
		e := connection.Rollback()
		if e != nil {
			ilog.Error("Transaction rollback error", err)
			return nil, err
		}
		return nil, err
	} else {
		e := connection.Commit()
		if e != nil {
			ilog.Error("Transaction commit error", err)
			return nil, err
		}
		return result, err
	}
}

/**
 * 插入一条记录（事务内使用）
 */
func (bd *BaseDao) TInsert(conn *sql.Tx, table string, params map[string]interface{}) (int64, error) {
	var query, keys, values bytes.Buffer

	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql insert】耗时：", time.Since(start), "【", query.String(), "】")
		}()
	}

	var args []interface{}
	query.WriteString("insert into `")
	query.WriteString(table)
	query.WriteString("` (")
	for k, v := range params{
		keys.WriteString("`")
		keys.WriteString(k)
		keys.WriteString("`,")
		values.WriteString("?,")
		args = append(args, v)
	}

	query.WriteString(string(keys.Bytes()[0:keys.Len() - 1]))
	query.WriteString(") values (")
	query.WriteString(string(values.Bytes()[0:values.Len() - 1]))
	query.WriteString(")")

	r, err := conn.Exec(query.String(), args...)
	if err != nil {
		ilog.Error("insert error : ", err)
		return -1, err
	}
	return r.LastInsertId()
}

/**
 * 插入多条记录（事务内使用）
 */
func (bd *BaseDao) TMultiInsert(conn *sql.Tx, table string, params []map[string]interface{}) (int, error) {
	if len(params) == 0 {
		return 0, nil
	}

	var query, keys, values bytes.Buffer
	var statement string
	var args []interface{}
	var keySet []string

	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql multi insert】耗时：", time.Since(start), "【", statement, "】")
		}()
	}

	query.WriteString("insert into `")
	query.WriteString(table)
	query.WriteString("` (")
	for k := range params[0]{
		keys.WriteString("`")
		keys.WriteString(k)
		keys.WriteString("`,")
		keySet = append(keySet, k)
	}

	query.WriteString(string(keys.Bytes()[0:keys.Len() - 1]))
	query.WriteString(") values ")

	for _, item := range params {
		values.WriteString("(")
		for _, k := range keySet {
			values.WriteString("?,")
			args = append(args, item[k])
		}
		query.WriteString(string(values.Bytes()[0:values.Len() - 1]))
		query.WriteString("),")
		values.Reset()
	}

	statement = string(query.Bytes()[0:query.Len() - 1])
	_, err := conn.Exec(statement, args...)
	if err != nil {
		ilog.Error("multi insert error : ", err)
		return -1, err
	}
	return len(params), nil
}

/**
 * 查询一条记录（事务内使用）
 */
func (bd BaseDao) TFind(conn *sql.Tx, sql string, p ...interface{}) *sql.Row{
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql find】耗时：", time.Since(start), "【", sql, "】")
		}()
	}
	return conn.QueryRow(sql, p...)
}

/**
 * 查询记录
 */
func (bd BaseDao) TSelect(conn *sql.Tx, sql string, p ...interface{}) (*sql.Rows, error){
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql select】耗时：", time.Since(start), "【", sql, "】")
		}()
	}
	return conn.Query(sql, p...)
}

/**
 * 更新记录（事务内使用）
 */
func (bd BaseDao) TUpdate(conn *sql.Tx, sql string, p ...interface{}) (int64, error) {
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql update】耗时：", time.Since(start), "【", sql, "】")
		}()
	}

	r, err := conn.Exec(sql, p...)
	if err != nil {
		ilog.Error("update error : ", err)
		return -1, err
	}
	return r.RowsAffected()
}

/**
 * 根据主键id更新记录（事务内使用）
 */
func (bd BaseDao) TUpdateById(conn *sql.Tx, table string, id int32, params map[string]interface{}) (int64, error) {
	var query, set bytes.Buffer

	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql update】耗时：", time.Since(start), "【", query.String(), "】")
		}()
	}

	var args []interface{}
	query.WriteString("update `")
	query.WriteString(table)
	query.WriteString("` set ")
	for k, v := range params {
		set.WriteString(fmt.Sprintf("`%s`=?, ", k))
		args = append(args, v)
	}
	query.WriteString(string(set.Bytes()[0:set.Len() - 2]))
	query.WriteString(" where `id`=?")
	args = append(args, id)

	r, err := conn.Exec(query.String(), args...)
	if err != nil {
		ilog.Error("update error : ", err)
		return -1, err
	}
	return r.RowsAffected()
}

/**
 * 删除记录（事务内使用）
 */
func (bd BaseDao) TDelete(conn *sql.Tx, sql string, p ...interface{}) (int64, error) {
	if bd.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Mysql delete】耗时：", time.Since(start), "【", sql, "】")
		}()
	}
	r, err := conn.Exec(sql, p...)
	if err != nil {
		ilog.Error("delete error : ", err)
		return -1, err
	}
	return r.RowsAffected()
}
