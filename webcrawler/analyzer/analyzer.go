package analyzer

import (
	"errors"
	"fmt"
	"demo/mylog"
	"net/url"
	"demo/webcrawler/base"
	"demo/webcrawler/middleware"
)

var logger mylog.Logger = base.NewLogger()
var analyzerIdGenerator middleware.IdGenerator = middleware.NewIdGenerator()

func genAnalyzerId() uint32 {
	return analyzerIdGenerator.GetUint32()
}

type Analyzer interface {
	Id() uint32 
	Analyze(
		respParsers []ParseResponse,
		resp base.Response) ([]base.Data, []error)//根据规则分析响应并返回请求条目
}

//创建分析器
func NewAnalyzer() Analyzer {
	return &myAnalyzer{id: genAnalyzerId()}
}

//分析器,只有一个ID成员
type myAnalyzer struct {
	id uint32
}

func (analyzer *myAnalyzer) Id() uint32 {
	return analyzer.id
}

func (analyzer *myAnalyzer) Analyze(respParsers []ParseResponse, resp base.Response) (dataList []base.Data, errorList []error) {
	if respParsers == nil {
		err := errors.New("The response parser list is invalid!")
		return nil, []error{err}
	}
	httpResp := resp.HttpResp()
	if httpResp == nil {
		err := errors.New("The http response is invalid!")
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	logger.Infof("Parse the response (reqUrl=%s)...", reqUrl)
	respDepth := resp.Depth()

	//解析HTTP响应
	dataList = make([]base.Data, 0)
	errorList = make([]error, 0)
	for i, respPaser := range respParsers {
		if respPaser == nil {
			err := errors.New(fmt.Sprintf("The document parser [%d] is invalid!", i))
			errorList = append(errorList, err)
			continue
		}
		pDataList, pErrorList := respPaser(httpResp, respDepth)
		if pDataList != nil {
			for _, pData := range pDataList {
				dataList = appendDataList(dataList, pData, respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				errorList = appendErrorList(errorList, pError)
			}
		}
	}
	return dataList, errorList
}

func appendDataList(dataList []base.Data, data base.Data, respDepth uint32) []base.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*base.Request)
	if !ok {
		return append(dataList, data)
	}
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		req = base.NewRequest(req.HttpReq(), newDepth)
	}
	return append(dataList, req)
}

func appendErrorList(errorList []error, err error) []error {
	if err == nil {
		return errorList
	}
	return append(errorList, err)
}
