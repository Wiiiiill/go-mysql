package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Wiiiiill/go-mysql/helper"
)

var (
	_configFile = "config.json"
	_Config     *DBConfig
	_once       sync.Once
	Repositorys []*Repository
	RMap        map[string]*Repository
)

type DBConfig struct {
	Database []DatabaseCluster
}

func Close() {
	for i := 0; i < len(Repositorys); i++ {
		Repositorys[i].Close()
	}
}

// Init config
func Init(configFileName string) {
	_once.Do(func() {
		c := &DBConfig{}
		if len(configFileName) > 0 {
			_configFile = configFileName
		}
		err := helper.ReadJSON(c, _configFile)
		if err != nil {
			panic(err)
		}
		if len(c.Database) == 0 {
			panic(fmt.Errorf("no database config"))
		}
		_Config = c
		Repositorys = make([]*Repository, len(c.Database))
		RMap = make(map[string]*Repository)
		for _, v := range c.Database {
			db := createRepository(&v)
			db.OpenReadDatabases()
			db.OpenWriteDatabase()
			Repositorys = append(Repositorys, db)
			if RMap[v.Name] != nil {
				panic(fmt.Errorf("duplicate db name"))
			}
			RMap[v.Name] = db
		}
	})
}

func DB(key string) *Repository {
	obj, ok := RMap[key]
	if ok {
		return obj
	}
	panic(fmt.Errorf("no such db"))
}

// WebConfig get WebConfig
func GetConfig() *DBConfig {
	return _Config
}

// createRepository return contract.Repository
func createRepository(dc *DatabaseCluster) *Repository {

	cat := &Repository{
		databaseCluster: dc,
	}

	return cat
}

// Repository struct
type Repository struct {
	databaseCluster *DatabaseCluster
	writeDatabase   *sql.DB
	readDatabases   []*sql.DB
}

type DatabaseCluster struct {
	Name      string
	Driver    string
	Database  string
	Username  string
	Password  string
	Charset   string
	Collation string
	Write     *DatabaseHostConfig
	Read      *[]DatabaseHostConfig
}

type DatabaseHostConfig struct {
	Host string
	Port int
}

// Close databases
func (r *Repository) Close() error {

	if r.writeDatabase != nil {
		r.writeDatabase.Close()
	}

	for _, rd := range r.readDatabases {
		if rd != nil {
			rd.Close()
		}
	}

	return nil
}

// OpenWriteDatabase of config.database.write
func (r *Repository) OpenWriteDatabase() error {
	db, err := sql.Open(r.databaseCluster.Driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		r.databaseCluster.Username,
		r.databaseCluster.Password,
		r.databaseCluster.Write.Host,
		r.databaseCluster.Write.Port,
		r.databaseCluster.Database,
		r.databaseCluster.Charset))

	if err != nil {
		return err
	}

	r.writeDatabase = db

	return nil
}

// OpenReadDatabases of config.database.read
func (r *Repository) OpenReadDatabases() error {

	var readDatabases []*sql.DB

	for _, rd := range *r.databaseCluster.Read {
		db, err := sql.Open(r.databaseCluster.Driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			r.databaseCluster.Username,
			r.databaseCluster.Password,
			rd.Host,
			rd.Port,
			r.databaseCluster.Database,
			r.databaseCluster.Charset))

		if err != nil {
			return err
		}

		readDatabases = append(readDatabases, db)
	}

	r.readDatabases = readDatabases

	return nil
}

// SelectDB select DB for Query
func (r *Repository) SelectDB() *sql.DB {
	return r.readDatabases[helper.RandMax(len(r.readDatabases))]
}

// Query ...
func (r *Repository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return r.SelectDB().Query(query, args...)
}

// QueryRow ...
func (r *Repository) QueryRow(query string, args ...interface{}) *sql.Row {
	return r.SelectDB().QueryRow(query, args...)
}

// Prepare ...
func (r *Repository) Prepare(query string) (*sql.Stmt, error) {
	return r.writeDatabase.Prepare(query)
}

// TxPrepare ...
func (r *Repository) TxPrepare(tx *sql.Tx, query string) (*sql.Stmt, error) {
	return tx.Prepare(query)
}

// StmtExec ...
func (r *Repository) StmtExec(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	return stmt.Exec(args...)
}

// Exec ...
func (r *Repository) Exec(query string, args ...interface{}) (sql.Result, error) {
	return r.writeDatabase.Exec(query, args...)
}

// Begin ...
func (r *Repository) Begin() (*sql.Tx, error) {
	return r.writeDatabase.Begin()
}

// TxExec ...
func (r *Repository) TxExec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	return tx.Exec(query, args...)
}

// Rollback ...
func (r *Repository) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// Commit ...
func (r *Repository) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// Now ...
func (r *Repository) Now() *time.Time {
	time := time.Now()
	return &time
}

// Max get max key value
func (r *Repository) Max(tableName string, key string, AppNum, AppID int) (uint64, error) {

	sqlx := "SELECT MAX(`" + key + "`) FROM `" + tableName + "` WHERE `" + key + "` % ? = ? "

	row := r.QueryRow(sqlx, AppNum, AppID)

	var val *uint64

	err := row.Scan(&val)

	if err != nil {
		return 0, err
	}

	if val == nil {
		return 0, nil
	}

	return *val, nil
}
