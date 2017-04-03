package visit

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func TrimUrl(url string) string {
	url = strings.TrimRight(url, "/")
	url = strings.TrimRight(url, "/?")
	return url
}

func GetHostName(url string) string {
	strs := strings.Split(url, "//")
	if len(strs) == 0 {
		return ""
	}

	host_name := strs[1]
	strs = strings.Split(host_name, "/")
	if len(strs) == 0 {
		return ""
	}
	return strs[0]
}

func GetTitle(url string) string {
	page := findPage(url)
	if page.Id != 0 {
		return page.Title
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return url
	}

	title := doc.Find("title").Text()
	title = strings.Trim(title, " ")
	if len(title) == 0 {
		return url
	}
	return title
}
