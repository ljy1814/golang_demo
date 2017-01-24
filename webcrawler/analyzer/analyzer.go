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

//创建分析器,专注于内部计算
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
	//将类型转换为请求,若是请求类型则转换成功
	req, ok := data.(*base.Request)
	//转换不成功,data是个条目,对于条目,不需要任何处理
	if !ok {
		return append(dataList, data)
	}
	//对请求类型做处理,检验深度是否正确,否则重新生成请求
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		//此处会复制数据,可能会造成垃圾回收的压力,但是只是在深度不正确的情况下才出现
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
