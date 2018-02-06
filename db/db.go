package db

import (
	"fmt"
	// for go mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

// Database struct that hold db information
type Database struct {
	DatabaseName   string       `yaml:"database_name"`
	DatabaseType   string       `yaml:"database_type"`
	ConnectionType string       `yaml:"connection_type"`
	Host           string       `yaml:"host"`
	Port           int          `yaml:"port"`
	Username       string       `yaml:"username"`
	Password       string       `yaml:"password"`
	Engine         *xorm.Engine `yaml:"engine"`
}

//Client get database client
func (d *Database) Client() (*xorm.Engine, error) {
	if d.Engine == nil {
		e, err := xorm.NewEngine(d.DatabaseType, d.DSN())
		d.Engine = e
		return e, err
	}
	return d.Engine, nil
}

//DSN get data source name
func (d *Database) DSN() string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8", d.Username, d.Password, d.ConnectionType, d.Host, d.Port, d.DatabaseName)
}
