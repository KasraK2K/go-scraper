package scraper

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Products struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	URL   string `json:"url"`
}

type Result struct {
	Result []Products `json:"result"`
}

func JackJones() {
	c := colly.NewCollector(colly.AllowedDomains("www.jackjones.com"))
	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	var allProducts []Products

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Scraping:", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status:", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("div.product-tile__content-wrapper", func(h *colly.HTMLElement) {
		products := Products{
			Name:  h.ChildText("a.product-tile__name__link.js-product-tile-link"),
			Price: h.ChildText("em.value__price"),
			URL:   h.ChildAttr("a.product-tile__name__link.js-product-tile-link", "href"),
		}

		allProducts = append(allProducts, products)

		content, err := json.Marshal(Result{allProducts})
		if err != nil {
			fmt.Println(err.Error())
		}
		os.WriteFile("jack-shoes.json", content, 0644)
	})

	c.OnHTML("a.paging-controls__next.js-page-control", func(p *colly.HTMLElement) {
		nextPage := p.Request.AbsoluteURL(p.Attr("data-href"))
		c.Visit(nextPage)
	})

	c.Visit("https://www.jackjones.com/nl/en/jj/shoes/")
}
