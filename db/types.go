package db

//HouseInfo basic struct holding hourse info
type HouseInfo struct {
	UUID      string `yaml:"uuid"`
	Community string `yaml:"community"`
	Structure string `yaml:"structure"`
	Area      string `yaml:"area"`
	Direction string `yaml:"direction"`
	Status    string `yaml:"status"`
	UnitPrice string `yaml:"unit_price"`
	Link      string `yaml:"link"`
}
