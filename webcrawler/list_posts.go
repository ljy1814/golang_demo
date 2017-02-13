package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func postScrape() {
	doc, err := goquery.NewDocument("http://jonathanmh.com")
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("#main article .entry-title").Each(func(
		index int, item *goquery.Selection) {
		title := item.Text()
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		fmt.Printf("Post #%d: %s - %s\n", index, title, link)
	})
}

//登录知乎,
/**
_xsrf 0cba0379060915fdffeb67dc45d76c31
password
email

*/

func main() {
	postScrape()
	zhihu()
}

type MyTransport struct {
	Transport RoundTripper
}

func (t *MyTransport) transport() http.RoundTripper {
	if nil != t.Transport {
		return t.Transport
	}

	return http.DefaultTransport
}

func (t *MyTransport) RoundTrip(req *http.Request) (http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.111 Safari/537.36")
	return t.transport().RoundTrip(req)
}

type Client struct {
	http.Client
}

var c Client

func NewClient() *Client {
	t := &MyTransport{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{Transport: t, Jar: jar}
}

func zhihu() {
	c = NewClient()

	sUrl := "https://www.zhihu.com/#signin"

	v := url.Values{
		"username": "420581a01bv.cdb@sina.cn",
		"password": "ljy#zhihu",
		"_xsrf":    "0cba0379060915fdffeb67dc45d76c31",
	}

	req, err := http.NewRequest("POST", sUrl, v)

	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Body)
}
