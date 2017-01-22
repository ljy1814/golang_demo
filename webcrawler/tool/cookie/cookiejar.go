package cookie

import (
	"code.google.com/p/go.net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

//创建http.CookieJar类型的值
func NewCookieJar() http.CookieJar {
	options := &cookiejar.Options{PublicSuffixList: &myPublicSuffixList{}}
	cj, _ cookiejar.New(options)
	return cj
}

//cookiejar.PublicSuffixList接口的实现
typemyPublicSuffixList struct{}

func (psl *myPublicSuffixList) PublicSuffix(domain string) string {
	suffix, _ := publicsuffix.PublicSuffix(domain)
	return suffix
}

func (psl *myPublicSuffixList) String() string {
	return "Web crawler - public suffix list (rev 1.0) power by 'code.google.com/p/go.net/publicsuffix'"
}
