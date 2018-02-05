package db

import (
	"fmt"
	// for go mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

//UsedHouse basic struct holding hourse info
type UsedHouse struct {
	UUID             string `yaml:"uuid" xorm:"pk 'uuid'"`
	Region           string `yaml:"region" xorm:"'region'"`
	Layout           string `yaml:"layout" xorm:"'layout'"`
	Area             string `yaml:"area" xorm:"'area'"`
	Direction        string `yaml:"direction" xorm:"'direction'"`
	HousingEstate    string `yaml:"housing_estate" xorm:"'housing_estate'"`
	Floor            string `yaml:"floor" xorm:"'floor'"`
	DecorationStatus string `yaml:"decoration_status" xorm:"'decoration_status'"`
	TotalPrice       int    `yaml:"total_price" xorm:"'total_price'"`
	TotalPriceString string `yaml:"total_price_string" xorm:"'total_price_string'"`
	UnitPrice        int    `yaml:"unit_price" xorm:"'unit_price'"`
	UnitPriceString  string `yaml:"unit_price_string" xorm:"'unit_price_string'"`
	Link             string `yaml:"link" xorm:"'link'"`
	TaxFree          string `yaml:"tax_free" xorm:"'tax_free'"`
	Subway           string `yaml:"subway" xorm:"'subway'"`
	Follow           string `yaml:"follow" xorm:"'follow'"`
}

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
	fmt.Println(d.DSN())
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
