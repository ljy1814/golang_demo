package downloader

import (
	"demo/mylog"
	"net/http"
	"demo/webcrawler/base"
	"demo/webcrawler/middleware"
)

//日志记录器
var logger mylog.Logger = base.NewLogger()

var downloaderIdGenerator middleware.IdGenerator = middleware.NewIdGenerator()

//获取downloader的ID
func genDownloaderId() uint32 {
	return downloaderIdGenerator.GetUint32()
}

//网页下载器的接口类型
type PageDownloader interface {
	Id() uint32
	Download(req base.Request) (*base.Response, error)	//根据请求下载网页
}

//创建网页下载器
func NewPageDownloader(client *http.Client) PageDownloader {
	id := genDownloaderId()
	if client == nil {
		client = &http.Client{}
	}
	return &myPageDownloader{
		id: id,
		httpClient: *client,
	}
}

//网页下载器实际类型
type myPageDownloader struct {
	id uint32
	httpClient http.Client
}

func (dl *myPageDownloader) Id() uint32 {
	return dl.id
}

func (dl *myPageDownloader) Download(req base.Request) (*base.Response, error) {
	httpReq := req.HttpReq()
	logger.Infof("Do the request (url=%s)...\n", httpReq.URL)
	httpResp, err := dl.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return base.NewResponse(httpResp, req.Depth()), nil
}
