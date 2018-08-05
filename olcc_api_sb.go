package main

import (
	"net/http"
	"log"
	"net/http/cookiejar"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"io"
	"os"
		"strconv"
	)

const BASE_URL = "http://www.oregonliquorsearch.com"
const S_URL = "http://www.oregonliquorsearch.com/servlet/FrontController?view=productlist&action=display&productSearchParam=MEDOYEFF+STARKA+VODKA&column=Description&pageSize=100"

type Location struct {
	city     string
	street   string
	zip      string
	quantity int
}

type Alcohol struct {
	name      string
	price     string
	category  string
	size      string
	age       string
	proof     string
	locations []Location
}

type Result struct {
	queryStatus bool // true if product found and false otherwise
	alcohol     Alcohol
}

func main() {
	// Define custom client with a timeout to avoid hangs
	// (source: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779)

	//CONSIDER: New cookie for EACH request, or...

	result , _ := Search("")

	log.Println(result.alcohol)
}

func Search(query string) (Result, error) {

	var err = new(error)

	// Build the url with params here maybe use url object
	// ...
	// ...

	// Query the OLCC website and get the resulting html
	//html := QueryOLCCWebsite(query)
	html := loadHtmlFromFile()

	// Parse the html and format the data
	result := ParseHTML(html)

	return result, err
}

// Sends the query to olcc website and returns the resulting html
func QueryOLCCWebsite(query string) io.Reader {

	client := GetHTTPClient(S_URL)

	resp, err := client.Get(S_URL)
	check(err)

	defer resp.Body.Close()
	// TODO: Handle more response codes here

	checkResponseStatus(resp)

	return resp.Body
}

/*
	Returns an http.Client with which to make requests
*/
func GetHTTPClient(reqUrl string) *http.Client {

	// Grab session cookie
	cookies, err := GetCookies()
	check(err)

	// Instantiate a cookiejar
	cookieJar, err := cookiejar.New(nil)
	check(err)

	// Create a url instance
	u, err := url.Parse(reqUrl)
	check(err)

	// Put the cookie in the cookie jar
	cookieJar.SetCookies(u, cookies)

	client := &http.Client{
		//Timeout: time.Second * 10,
		Jar: cookieJar,
	}

	return client
}

/*
	Returns the cookies from the OLCC Search Site
*/
func GetCookies() ([]*http.Cookie, error) {

	resp, err := http.Get(BASE_URL)
	check(err)
	checkResponseStatus(resp)

	defer resp.Body.Close()

	cookies := resp.Cookies()

	log.Println(cookies)

	return cookies, err
}

func ParseHTML(html io.Reader) Result {
	var location Location
	var alcohol Alcohol


	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(html)
	check(err)



	doc.Find("div .list tr").Each(func(i int, item *goquery.Selection) {
		if i > 0 {
			entry := item.Find("td")
			//storeNum := entry.Eq(STORE_NUM).Find("span").Text();
			location.city = entry.Eq(1).Text()
			location.street = entry.Eq(2).Text()
			location.zip = entry.Eq(3).Text()
			location.quantity, err = strconv.Atoi(entry.Eq(6).Text())
			check(err)

			alcohol.locations[i - 1] = location
			//fmt.Println(entry.Eq(1).Text())
		}
	})

	return Result{nil, alcohol}
}

/*****REMOVE WHEN FINISHED TESTING ONLY*****/
func saveHtmlToFile(html io.Reader) {

}

func loadHtmlFromFile() io.Reader {
	htmlFilePath := "C:/Users/taylo/PycharmProjects/olcc_api/olcc_html"
	html, err := os.Open(htmlFilePath)
	check(err)
	return html
}

func checkResponseStatus(resp *http.Response )  {
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
