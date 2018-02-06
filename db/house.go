package db

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// HouseHandler handle db related issues
type HouseHandler struct {
	House          *UsedHouse
	DatabaseConfig Database
}

// Init init house handler
func (hh *HouseHandler) Init() error {
	e, err := hh.DatabaseConfig.Client()
	if err != nil {
		return err
	}
	e.Sync(new(UsedHouse))
	return nil
}

//Save save to mysql
func (hh *HouseHandler) Save() error {
	e, err := hh.DatabaseConfig.Client()
	if err != nil {
		return err
	}
	_, err = e.Insert(hh.House)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//UsedHouse basic struct holding hourse info
type UsedHouse struct {
	UUID string `yaml:"uuid" xorm:"pk 'uuid'"`
	// basic info
	Layout           string  `yaml:"layout" xorm:"'layout'"`
	Area             float64 `yaml:"area" xorm:"'area'"`
	AreaString       string  `yaml:"areaString" xorm:"'area_string'"`
	Age              int     `yaml:"age" xorm:"'age'"`
	Floor            string  `yaml:"floor" xorm:"'floor'"`
	DecorationStatus string  `yaml:"decoration_status" xorm:"'decoration_status'"`
	TotalPrice       int     `yaml:"total_price" xorm:"'total_price'"`
	UnitPrice        int     `yaml:"unit_price" xorm:"'unit_price'"`
	AgeString        string  `yaml:"age_string" xorm:"'age_string'"`
	Direction        string  `yaml:"direction" xorm:"'direction'"`
	// community info
	Location      string `yaml:"location" xorm:"location"`
	District      string `yaml:"district" xorm:"district"`
	Region        string `yaml:"region" xorm:"'region'"`
	Subway        string `yaml:"subway" xorm:"'subway'"`
	HousingEstate string `yaml:"housing_estate" xorm:"'housing_estate'"`
	// others
	Link string `yaml:"link" xorm:"'link'"`
}

// Digest get digest
func (uh *UsedHouse) Digest() string {
	s := fmt.Sprintf("%v", uh)
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
