package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type HotelUrl struct {
	url string
}

type HotelInfo struct {
	hotelName                 string
	hotelAddress              string
	hotelOverallRatingLabel   string
	hotelOverallRating        string
	hotelOverallReviewCount   string
	hotelOverallDescription   string
	reviewerName              string
	reviewerLocation          string
	reviewerRating            string
	reviewerCommentTitle      string
	reviewerCommentDescrition string
	reviewerStayTime          string
	tripType                  string
}

func main() {
	scrapHotelUrlFromCSV()
}

func scrapHotelUrlFromCSV() {
	//Open hotel url csv file
	file, err := os.Open("hotelUrl.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	//Read my file
	fileReader := csv.NewReader(file)
	for {
		hotelUrlDataSet, err := fileReader.Read()
		if err != nil || err == io.EOF {
			log.Fatal(err)
			break
		}

		for value := range hotelUrlDataSet {
			getHotelDetails(hotelUrlDataSet[value])
			//fmt.Printf("%s\n", hotelUrlDataSet[value])
		}

	}

	if err != nil {
		fmt.Println(err)
	}

}

func getHotelDetails(url string) {

	//Initiate Hotel Struct
	hotel := HotelInfo{}

	//Create csv file
	file, err := os.OpenFile("hotel.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	//Create csv writer
	writer := csv.NewWriter(file)
	//defer writer.Flush()

	//Create csv headers
	headers := []string{"hotelName", "hotelLocation", "hotelRatingLabel", "hotelRating", "hotelReviewCount", "reviewerName", "reviewerLocation", "reviewerStayTime",
		"reviewerCommentTitle", "reviewerCommentDescrition", "tripType", "reviewerRating"}
	writer.Write(headers)

	//Create colly collector and allow only tripadvisor url
	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
		colly.CacheDir("./cache"),
	)

	collector.Limit(&colly.LimitRule{
		Parallelism: 2,
		// Filter domains affected by this rule
		DomainGlob: "tripadvisor.com/*",
		// Set a delay between requests to these domains
		Delay: 100 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	//Keep track of visited urls
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting...", r.URL)
	})

	collector.OnHTML(".page", func(h *colly.HTMLElement) {
		hotel.hotelName = h.ChildText("h1")
		hotel.hotelAddress = h.ChildText("div.ApqWZ.S4.H3.f.u.eEkxn > span.eWZDY._S.eCdbd.yYjkv > span.ceIOZ.yYjkv")
		hotel.hotelOverallRating = h.ChildText("span.bvcwU.P")
		hotel.hotelOverallRatingLabel = h.ChildText("div.cNJsa")
		hotel.hotelOverallReviewCount = h.ChildText("span.btQSs.q.Wi.z.Wc")
		hotel.hotelOverallDescription = h.ChildText("div.pIRBV._T")

		h.ForEach("div[data-test-target=reviews-tab]", func(i int, h *colly.HTMLElement) {
			// fmt.Print(reviewerCommentTitle)

			h.ForEach("div.cWwQK.MC.R2.Gi.z.Z.BB.dXjiy", func(i int, h *colly.HTMLElement) {
				hotel.reviewerName = h.ChildText("div.xMxrO > div.bJaRP._Z.o > div.bcaHz > span > a.ui_header_link.bPvDb")
				hotel.reviewerCommentTitle = h.ChildText("div.cqoFv._T > div.fpMxB.MC._S.b.S6.H5._a > a.fCitC > span > span")
				hotel.reviewerCommentDescrition = h.ChildText("div.dovOW > div.duhwe._T.bOlcm.dMbup > div.pIRBV._T > q > span")
				hotel.reviewerLocation = h.ChildText("div.xMxrO > div.bJaRP._Z.o > div.BZmsN > span.fSiLz > span.default.ShLyt.small")
				hotel.reviewerRating = h.ChildAttr("div.cqoFv._T > div.elFlG.f.O > div.emWez.F1 > span", "class")[24:][:1] + "." + h.ChildAttr("div.emWez.F1 > span", "class")[24:][1:]
				h.ForEach("div.cqoFv._T > div.dovOW > div.bzjij > span.euPKI._R.Me.S4.H3", func(i int, h *colly.HTMLElement) {
					stayTime := h.Text
					if stayTime != "" {
						hotel.reviewerStayTime = h.Text[14:]
					} else {
						hotel.reviewerStayTime = ""
					}
				})

				tripType := h.ChildText("span.eHSjO._R.Me")
				if tripType != "" {
					hotel.tripType = tripType[11:]
				} else {
					hotel.tripType = ""
				}

				csvRow := []string{hotel.hotelName, hotel.hotelAddress, hotel.hotelOverallRatingLabel, hotel.hotelOverallRating, hotel.hotelOverallReviewCount,
					hotel.reviewerName, hotel.reviewerLocation, hotel.reviewerStayTime, hotel.reviewerCommentTitle,
					hotel.reviewerCommentDescrition, hotel.tripType, hotel.reviewerRating}

				writer.Write(csvRow)

			})

		})
		fmt.Println(hotel.hotelName)

		nextUrl := h.ChildAttr("a.ui_button.nav.next.primary", "href")
		//fmt.Println(hotel)

		h.Request.Visit(nextUrl)

	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	collector.Visit(url)

}

func scrapeHotelUrl() {
	indexPage := 0

	//Create csv file
	file, err := os.Create("hotelUrl.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	//Create csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	//Create csv headers
	headers := []string{"Url"}
	writer.Write(headers)

	//Create colly collector and allow only tripadvisor url
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"),
	)

	//detailCollector := collector.Clone()

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

	//Get Hotel description Url div.relWrap
	collector.OnHTML(".photo-wrapper", func(h *colly.HTMLElement) {
		//initialized Hotel Object
		hotel := HotelUrl{}
		hotel.url = h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
		fmt.Println(hotel)
		csvRow := []string{hotel.url}
		writer.Write(csvRow)
	})

	//Get next page url
	collector.OnHTML("div[data-trackingstring=pagination_h]", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
		//h.Request.Visit(nextPage)
		indexPage++
		fmt.Println("NextPage" + strconv.Itoa(indexPage))
		collector.Visit(nextPage)
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	collector.Visit("https://www.tripadvisor.com/Hotels-g482884-Zanzibar_Island_Zanzibar_Archipelago-Hotels.html")
}
