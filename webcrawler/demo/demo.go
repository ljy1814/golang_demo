package main

import (
	"demo/mylog"
	"demo/webcrawler/analyzer"
	"demo/webcrawler/base"
	"demo/webcrawler/itempipeline"
	"demo/webcrawler/scheduler"
	"demo/webcrawler/tool"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var logger mylog.Logger = mylog.NewSimpleLogger()

//处理条目
func processItem(item base.Item) (result base.Item, err error) {
	if item == nil {
		return nil, errors.New("Invalid item!")
	}
	//生成结果
	result = make(map[string]interface{})
	for k, v := range item {
		result[k] = v
	}
	//	fmt.Println("-----------------------")
	//	fmt.Println(result)
	//	fmt.Println("-----------------------")
	if _, ok := result["number"]; !ok {
		result["number"] = len(result)
	}
	time.Sleep(10 * time.Millisecond)
	return result, nil
}

//响应解析函数
func parseForATag(httpResp *http.Response, respDepth uint32) ([]base.Data, []error) {
	//TODO 支持更多的HTTP响应状态
	if httpResp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Unsupported status code %d. (httpResponse=%v)", httpResp))
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	var httpRespBody io.ReadCloser = httpResp.Body
	defer func() {
		if httpResp != nil {
			httpRespBody.Close()
		}
	}()
	dataList := make([]base.Data, 0)
	errs := make([]error, 0)
	doc, err := goquery.NewDocumentFromReader(httpRespBody)
	if err != nil {
		errs = append(errs, err)
		return dataList, errs
	}
	//查找A标签并提取链接
	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		//前期过滤
		if !exists || href == "" || href == "#" || href == "/" {
			return
		}
		href = strings.TrimSpace(href)
		lowerHref := strings.ToLower(href)
		//暂时不支持对JavaScript代码的解析
		if href != "" && !strings.HasPrefix(lowerHref, "javascript") {
			aUrl, err := url.Parse(href)
			if err != nil {
				errs = append(errs, err)
				return
			}
			if !aUrl.IsAbs() {
				aUrl = reqUrl.ResolveReference(aUrl)
			}
			httpReq, err := http.NewRequest("GET", aUrl.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				req := base.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		}
		text := strings.TrimSpace(sel.Text())
		if text != "" {
			imap := make(map[string]interface{})
			imap["parent_url"] = reqUrl
			imap["a.text"] = text
			item := base.Item(imap)
			dataList = append(dataList, &item)
			imap["a.index"] = index
		}
	})
	return dataList, errs
}

//获取响应解析函数的序列
func getResponseParsers() []analyzer.ParseResponse {
	parsers := []analyzer.ParseResponse{
		parseForATag,
	}

	return parsers
}

//获得条目处理器的序列
func getItemProcessors() []itempipeline.ProcessItem {
	itemProcessors := []itempipeline.ProcessItem{
		processItem,
	}
	return itemProcessors
}

//生成HTTP客户端
func genHttpClient() *http.Client {
	return &http.Client{}
}

func record(level byte, content string) {
	if content == "" {
		return
	}
	switch level {
	case 0:
		logger.Infoln(content)
	case 1:
		logger.Warnln(content)
	case 2:
		logger.Errorln(content)
	}
}

func main() {
	fmt.Println("vim-go")
	//创建调度器
	scheduler1 := scheduler.NewScheduler()
	//准备监控参数
	intervalNs := 10 * time.Millisecond
	maxIdleCount := uint(1000)
	//开始监控
	checkCountChan := tool.Monitoring(scheduler1, intervalNs, maxIdleCount, true, false, record)
	//准备启动参数
	channelArgs := base.NewChannelArgs(10, 10, 10, 10)
	poolBaseArgs := base.NewPoolBaseArgs(3, 3)
	crawlDepth := uint32(1)
	httpClientGenerator := genHttpClient
	respParsers := getResponseParsers()
	itemProcessors := getItemProcessors()
	//	startUrl := "http://www.sogou.com"
	startUrl := "https://www.zhihu.com"
	firstHttpReq, err := http.NewRequest("GET", startUrl, nil)
	if err != nil {
		logger.Errorln(err)
		return
	}
	//开启调度器
	scheduler1.Start(channelArgs, poolBaseArgs, crawlDepth, httpClientGenerator, respParsers, itemProcessors, firstHttpReq)
	//等待监控结束
	<-checkCountChan
}
