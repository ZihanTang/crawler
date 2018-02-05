package main

import (
	"cc/crawler/db"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())
	var hs []*db.HouseInfo
	doc.Find(".sellListContent li").Each(func(i int, s *goquery.Selection) {
		h := db.HouseInfo{}
		a := s.Find(".title a")
		if l, ok := a.Attr("href"); ok {
			h.Link = l
		}
		if uuid, ok := a.Attr("data-housecode"); ok {
			h.UUID = uuid
		}
		info := s.Find(".houseInfo")
		infoSlice := strings.Split(info.Text(), "|")
		h.Community = infoSlice[0]
		h.Structure = infoSlice[1]
		h.Direction = infoSlice[3]
		h.Status = infoSlice[4]
		fmt.Println(h)
		hs = append(hs, &h)
	})
	return nil, true
}

var rxOk = regexp.MustCompile(`http://sh\.lianjia\.com/ershoufang/(.*)?$`)

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

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://sh.lianjia.com/ershoufang/pudong/")
}
