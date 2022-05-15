package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	//Create csv file
	file, err := os.Create("tripadvisor.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	//Create csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	//Create csv headers
	headers := []string{"Hotel", "location", "Url"}
	writer.Write(headers)

	//Create colly collector and allow only tripadvisor url
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
	)

	// collector.Limit(&colly.LimitRule{
	// 	// Filter domains affected by this rule
	// 	DomainGlob: "tripadvisor.com/*",
	// 	// Set a delay between requests to these domains
	// 	Delay: 1 * time.Second,
	// 	// Add an additional random delay
	// 	RandomDelay: 1 * time.Second,
	// })

	//Keep track of visited urls
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting...", r.URL)
	})

	//Get Hotel description Url
	collector.OnHTML(`div.relWrap`,
		func(h *colly.HTMLElement) {
			h.ForEach("div[id=taplc_hsx_hotel_list_lite_dusty_hotels_combined_sponsored_ad_density_control_0]", func(i int, h *colly.HTMLElement) {
				reviewHotelUrl := h.Request.AbsoluteURL(h.ChildAttrs(".respListingPhoto", "href")[0])
				//Here visiting specific hotel detail
				collector.Visit(reviewHotelUrl)
			})
		})

	//Get next page url
	collector.OnHTML("div[data-trackingstring=pagination_h]", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
		collector.Visit(nextPage)
	})

	//ALL ABOUT DETAIL HOTEL PAGE

	collector.OnHTML(".page", func(h *colly.HTMLElement) {
		hotelName := h.ChildText("h1")
		fmt.Println(hotelName)
	})

	//END

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	collector.Visit("https://www.tripadvisor.com/Hotels-g482884-Zanzibar_Island_Zanzibar_Archipelago-Hotels.html")
}
