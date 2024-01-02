package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Username       string
	Password       string
	Hostname       string
	DatabaseName   string
	Port           int
	DataSourceName string
	Connection     *sql.DB
}

func (d *Database) Connect() {
	var err error
	log.Debug("Connecting to database, start")
	defer log.Debug("Connected to database, end")

	d.Connection, err = sql.Open("mysql", d.DataSourceName)
	if err != nil {
		log.Errorf("%q", err)
	}
}

func (d *Database) Close() {
	log.Debug("Closing database connection, start")
	defer log.Debug("Closed database connection, start")

	d.Connection.Close()
}

func (d *Database) SetDataSourceName() {
	d.DataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Hostname, d.Port, d.DatabaseName)
}

func (d *Database) SetUsername(username string) {
	d.Username = username
	d.SetDataSourceName()
}

func (d *Database) SetPassword(password string) {
	d.Password = password
	d.SetDataSourceName()
}

func (d *Database) SetHostname(hostname string) {
	d.Hostname = hostname
	d.SetDataSourceName()
}

func (d *Database) SetPort(port int) {
	d.Port = port
	d.SetDataSourceName()
}

func (d *Database) GetConnection() *sql.DB {
	return d.Connection
}

func (d *Database) DisplaySettings() string {
	retv := ""
	retv += fmt.Sprintf("Username: %s\n", d.Username)
	retv += fmt.Sprintf("Password: %s\n", d.Password)
	retv += fmt.Sprintf("Hostname: %s\n", d.Hostname)
	retv += fmt.Sprintf("Port: %d\n", d.Port)
	retv += fmt.Sprintf("DataSourceName: %s\n", d.DataSourceName)
	return retv
}

func NewDatabase(username, password, hostname, database string, port int) *Database {
	retv := &Database{}

	retv.Username = username
	retv.Password = password
	retv.Hostname = hostname
	retv.DatabaseName = database
	retv.Port = port

	retv.SetDataSourceName()

	retv.Connect()

	return retv
}
