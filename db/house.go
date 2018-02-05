package db

import "fmt"

// HouseHandler handle db related issues
type HouseHandler struct {
	Houses         []UsedHouse
	DatabaseConfig Database
}

//Save save to mysql
func (hh *HouseHandler) Save() error {
	e, err := hh.DatabaseConfig.Client()
	if err != nil {
		return err
	}

	t, _ := e.DBMetas()
	fmt.Println(t[0])
	// fmt.Println(e)
	// e.Sync()
	// fmt.Println(hh.Houses)
	// fmt.Println(&hh.Houses)
	affected, err := e.Insert(&hh.Houses)
	fmt.Println(affected)
	if err != nil {
		return err
	}
	return nil
}
