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
var rxOk = regexp.MustCompile(`http://sh\.lianjia\.com/ershoufang/(.*)?$`)

var rxPrice = regexp.MustCompile("[0-9]+")

type Ext struct {
	*gocrawl.DefaultExtender
}

// Visit actions on each visit of url
func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())
	var hs []db.UsedHouse
	doc.Find(".sellListContent li").Each(func(i int, s *goquery.Selection) {
		h := db.UsedHouse{}
		// get outlinks
		title := s.Find(".title a")
		if l, ok := title.Attr("href"); ok {
			h.Link = l
		}
		if uuid, ok := title.Attr("data-housecode"); ok {
			h.UUID = uuid
		}
		// get basic info
		info := s.Find(".houseInfo")
		infoSlice := strings.Split(info.Text(), "|")
		h.HousingEstate = infoSlice[0]
		h.Layout = infoSlice[1]
		h.Direction = infoSlice[3]
		h.DecorationStatus = infoSlice[4]
		// get region info
		position := s.Find(".positionInfo")
		positionSlice := strings.Split(position.Text(), "-")
		h.Floor = strings.TrimSpace(positionSlice[0])
		h.Region = strings.TrimSpace(positionSlice[1])
		// get unit price
		up := s.Find(".priceInfo .unitPrice span").Text()
		h.UnitPriceString = up
		p := rxPrice.FindString(up)
		i, err := strconv.Atoi(p)
		if err != nil {
			h.UnitPrice = -1
		} else {
			h.UnitPrice = i
		}
		// get total price
		tp := s.Find(".priceInfo .totalPrice").Text()
		h.TotalPriceString = tp
		p = rxPrice.FindString(tp)
		i, err = strconv.Atoi(p)
		if err != nil {
			h.TotalPrice = -1
		} else {
			h.TotalPrice = i
		}
		// get tax free
		h.Subway = s.Find(".subway").Text()
		h.TaxFree = s.Find(".taxfree").Text()
		// followInfo
		h.Follow = s.Find(".followInfo").Text()
		hs = append(hs, h)
	})
	if len(hs) > 0 {
		hh := db.HouseHandler{
			Houses:         hs,
			DatabaseConfig: d,
		}
		if err := hh.Save(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return nil, true
}

// Filter only craw link that match rxOk pattern
func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	return rxOk.MatchString(ctx.NormalizedURL().String())
}

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 100
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

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://sh.lianjia.com/ershoufang/pudong/")
}
