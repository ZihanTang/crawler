package main

import (
	"cc/crawler/db"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

var mysqlHost string
var mysqlPort int
var mysqlUsername string
var mysqlPassword string
var mysqlDatabase string

var d db.Database
var hh db.HouseHandler
var rxOk = regexp.MustCompile(`http://sh\.lianjia\.com/ershoufang/(.*)?$`)

var rxUsedHouse = regexp.MustCompile(`http://sh\.lianjia\.com/ershoufang/[0-9]*\.html`)
var rxPage = regexp.MustCompile(`http://sh\.lianjia\.com/ershoufang/[a-z]*(/pg[0-9]*)?$`)
var rxPrice = regexp.MustCompile("[0-9]+")
var rxFloat = regexp.MustCompile(`[0-9]+\.[0-9]+`)
var pageMap = make(map[string]bool)

type Ext struct {
	*gocrawl.DefaultExtender
}

// Visit actions on each visit of url
func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	u := ctx.NormalizedURL().String()
	if rxPage.MatchString(u) {
		fmt.Printf("Visit page link: %s, add to used_houses queue.\n", u)
		return nil, true
	}
	fmt.Printf("Visit used house link: %s, craw information\n", u)

	h := db.UsedHouse{}
	h.Link = u

	content := doc.Find(".overview .content")

	// total price
	tp := content.Find(".price .total").Text()
	tpi, err := strconv.Atoi(tp)
	if err != nil {
		h.TotalPrice = -1
	} else {
		h.TotalPrice = tpi
	}
	// unit price
	up := content.Find(".unitPrice .unitPriceValue").Text()
	up = rxPrice.FindString(up)
	upi, err := strconv.Atoi(up)
	if err != nil {
		h.UnitPrice = -1
	} else {
		h.UnitPrice = upi
	}
	// house info
	h.Layout = content.Find(".room .mainInfo").Text()
	h.Floor = content.Find(".room .subInfo").Text()
	h.Direction = content.Find(".type .mainInfo").Text()
	h.DecorationStatus = content.Find(".type .subInfo").Text()
	area := content.Find(".area .mainInfo").Text()
	h.AreaString = area
	area = rxFloat.FindString(area)
	af, err := strconv.ParseFloat(area, 64)
	if err != nil {
		h.Area = -1.0
	} else {
		h.Area = af
	}
	a := content.Find(".area .subInfo").Text()
	h.AgeString = a
	a = rxPrice.FindString(a)
	ai, err := strconv.Atoi(a)
	if err != nil {
		h.Age = -1
	} else {
		h.Age = ai
	}
	// community info
	h.HousingEstate = content.Find(".communityName a.info").Text()
	info := content.Find(".areaName .info").Text()
	infoFields := strings.Fields(info)
	h.District = infoFields[0]
	h.Region = infoFields[1]
	if len(infoFields) > 2 {
		h.Location = infoFields[2]
	}
	h.Subway = content.Find(".areaName .supplement").Text()
	// ID
	id := content.Find(".houseRecord span.info").Text()
	h.UUID = rxPrice.FindString(id)
	hh.House = &h
	if err := hh.Save(); err != nil {
		fmt.Println(err)
	}
	return nil, true
}

// Filter we deal with links that is either pageable list or used house detail page.
func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	u := ctx.NormalizedURL().String()

	if isVisited {
		return false
	}
	if !rxOk.MatchString(u) {
		return false
	}
	if rxPage.MatchString(u) {
		_, ok := pageMap[u]
		if ok {
			return false
		}
		pageMap[u] = true
		return true
	}
	return rxUsedHouse.MatchString(u)
}

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 500 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 100000
	flag.StringVar(&mysqlHost, "mysql-host", "127.0.0.1", "hostname of mysql server")
	flag.IntVar(&mysqlPort, "mysql-port", 3306, "port of mysql server")
	flag.StringVar(&mysqlUsername, "mysql-username", "root", "user of mysql server")
	flag.StringVar(&mysqlPassword, "mysql-password", "root", "password of mysql server")
	flag.StringVar(&mysqlDatabase, "mysql-database", "lianjia", "database name of mysql server")
	flag.Parse()
	d = db.Database{
		Host:           mysqlHost,
		Port:           mysqlPort,
		DatabaseName:   mysqlDatabase,
		DatabaseType:   "mysql",
		ConnectionType: "tcp",
		Username:       mysqlUsername,
		Password:       mysqlPassword,
	}
	hh = db.HouseHandler{
		DatabaseConfig: d,
		House:          nil,
	}
	if err := hh.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://sh.lianjia.com/ershoufang/pudong/")
}
