package scheduler

import (
	"errors"
	"fmt"
	"demo/mylog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
	"demo/webcrawler/analyzer"
	"demo/webcrawler/base"
	"demo/webcrawler/downloader"
	"demo/webcrawler/itempipeline"
	"demo/webcrawler/middleware"
)

//组件唯一编号
const (
	DOWNLOADER_CODE = "downloader"
	ANALYZER_CODE = "analyzer"
	ITEMPIPELINE_CODE = "item_pipeline"
	SCHEDULER_CODE = "scheduler"
)

var logger mylog.Logger = base.NewLogger()
//生成HTTP客户端的函数类型
type GenHttpClient func() *http.Client
//调度器接口类型
type Scheduler interface {
	//开启调度器
	Start(
		channelArgs base.ChannelArgs, //通道参数
		poolBaseArgs base.PoolBaseArgs, //池基本参数容器
		crawlDepth uint32,	//需要爬取的网页的最大深度值,深度大于此值的网页会被忽略
		httpClientGeneratoor GenHttpClient, //生成HTTP客户端的函数
		respParsers []analyzer.ParseResponse, //分析器所需的被用来解析HTTP响应的函数序列
		itemProcessors []itempipeline.ProcessItem, //需要被置入条目处理器管道中的条目处理器的序列
		firstHttpReq *http.Request, //首次请求
	) (err error)
	//停止调度器的运行,所有处理模块执行的流程都会被中止
	Stop() bool
	//判断是否在运行
	Running() bool
	//获得通道错误,调度器以及各个处理器运行模块运行过程中出现的所有错误都会被发送到该通道
	ErrorChan() <-chan error
	Idle() bool //判断所有处理模块是否处于空闲状态
	Summary(prefix string) SchedSummary //摘要信息 
}

//创建调度器
func NewScheduler() Scheduler {
	return &myScheduler{}
}

type myScheduler struct {
	channelArgs base.ChannelArgs //通道参数的容器
	poolBaseArgs base.PoolBaseArgs //池基本参数通道
	crawlDepth uint32 //爬取最大深度
	primaryDomain string //主域名
	chanman middleware.ChannelManager //通道管理器
	stopSign middleware.StopSign //停止信号
	dlpool downloader.PageDownloaderPool //网页下载器池
	analyzerPool analyzer.AnalyzerPool //分析器池
	itemPipeline itempipeline.ItemPipeline //条目处理器通道
	reqCache requestCache //请求缓存
	urlMap map[string]bool  //已请求的URL字典
	running uint32 //运行标记 0-未运行,1-已运行,2-已停止
}

func (sched *myScheduler) Start(channelArgs base.ChannelArgs,
		poolBaseArgs base.PoolBaseArgs,
		crawlDepth uint32,
		httpClientGenerator GenHttpClient,
		respParsers []analyzer.ParseResponse,
		itemProcessors []itempipeline.ProcessItem,
		firstHttpReq *http.Request) (err error) {
	defer func() {   //func1
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Scheduler Error: %s\n", p)
			logger.Fatal(errMsg)
			err = errors.New(errMsg)
		}
	}()
	if atomic.LoadUint32(&sched.running) == 1 {
		return errors.New("The scheduler has been started!\n")
	}
	atomic.StoreUint32(&sched.running, 1)
	if err := poolBaseArgs.Check(); err != nil {
		return nil
	}
	sched.poolBaseArgs = poolBaseArgs
	sched.crawlDepth = crawlDepth
	sched.chanman = generateChannelManager(sched.channelArgs)
	if httpClientGenerator == nil {
		return errors.New("The HTTP Client generator list is invalid!")
	}
	dlpool, err := generatePageDownloaderPool(sched.poolBaseArgs.PageDownloaderPoolSize(), httpClientGenerator)
	if err != nil {
		errMsg := fmt.Sprintf("Occur error when get page downloader pool: %s\n", err)
		return errors.New(errMsg)
	}
	sched.dlpool = dlpool
	analyzerPool, err := generateAnalyzerPool(sched.poolBaseArgs.AnalyzerPoolSize())
	if err != nil {
		errMsg := fmt.Sprintf("Occur error when get analyzer pool: %s\n", err)
		return errors.New(errMsg)
	}
	sched.analyzerPool = analyzerPool
	if itemProcessors == nil {
		return errors.New("The item processor list is invalid!")
	}
	for i, ip := range itemProcessors {
		if ip == nil {
			return errors.New(fmt.Sprintf("The %dth item processor is invalid!", i))
		}
	}
	sched.itemPipeline = generateItemPipeline(itemProcessors)
	if sched.stopSign == nil {
		sched.stopSign = middleware.NewStopSign()
	} else {
		sched.stopSign.Reset()
	}
	sched.reqCache = newRequestCache()
	sched.urlMap = make(map[string]bool)

	sched.startDownloading()
	sched.activateAnalyzers(respParsers)
	sched.openItemPipeline()
	sched.schedule(10 * time.Millisecond)

	if firstHttpReq == nil {
		return errors.New("The first HTTP request is invalid!")
	}
	pd, err := getPrimaryDomain(firstHttpReq.Host)
	if err != nil {
		return err
	}
	sched.primaryDomain = pd
	firstReq := base.NewRequest(firstHttpReq, 0)
	sched.reqCache.put(firstReq)
	return nil
}

func (sched *myScheduler) Stop() bool {
	if atomic.LoadUint32(&sched.running) != 1 {
		return false
	}
	sched.stopSign.Sign()
	sched.chanman.Close()
	sched.reqCache.close()
	atomic.StoreUint32(&sched.running, 2)
	return true
}

func (sched *myScheduler) Running() bool {
	return atomic.LoadUint32(&sched.running) == 1
}

func (sched *myScheduler) ErrorChan() <-chan error {
	if sched.chanman.Status() != middleware.CHANNEL_MANAGER_STATUS_INITIALIZED {
		return nil
	}
	return sched.getErrorChan()
}

func (sched *myScheduler) Idle() bool {
	idleDlPool := sched.dlpool.Used() == 0
	idleAnalyzerPool := sched.analyzerPool.Used() == 0
	idleItemPipeline := sched.itemPipeline.ProcessingNumber() == 0
	if idleDlPool && idleAnalyzerPool && idleItemPipeline {
		return true
	}
	return false
}

func (sched *myScheduler) Summary(prefix string) SchedSummary {
	return NewSchedSummary(sched, prefix)
}

//开始下载
func (sched *myScheduler) startDownloading() {
	go func() {
		for {
			req, ok := <-sched.getReqChan()
			if !ok {
				break
			}
			go sched.download(req)
		}
	}()
}

//下载
func (sched *myScheduler) download(req base.Request) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Download Error: %s\n", p)
			logger.Fatal(errMsg)
		}
	}()
	downloader, err := sched.dlpool.Take()
	if err != nil {
		errMsg := fmt.Sprintf("Downloader pool error: %s\n", err)
		sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		return
	}
	defer func() {
		if err := sched.dlpool.Return(downloader); err != nil {
			errMsg := fmt.Sprintf("Downloader pool error: %s\n", err)
			sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		}
	}()
	code := generateCode(DOWNLOADER_CODE, downloader.Id())
	respp, err := downloader.Download(req)
	if respp != nil {
		sched.sendResp(*respp, code)
	}
	if err != nil {
		sched.sendError(err, code)
	}
}

//激活分析器
func (sched *myScheduler) activateAnalyzers(respParsers []analyzer.ParseResponse) {
	go func() {
		for {
			resp, ok := <-sched.getRespChan()
			if !ok {
				break
			}
			go sched.analyze(respParsers, resp)
		}
	}()
}

//分析
func (sched *myScheduler) analyze(respParsers []analyzer.ParseResponse, resp base.Response) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Analysis Error: %s\n", p)
			logger.Fatal(errMsg)
		}
	}()
	analyzer, err := sched.analyzerPool.Take()
	if err != nil {
		errMsg := fmt.Sprintf("Analyzer pool error: %s", err)
		sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		return
	}
	defer func() {
		err := sched.analyzerPool.Return(analyzer)
		if err != nil {
			errMsg := fmt.Sprintf("Analyzer pool error: %s", err)
			sched.sendError(errors.New(errMsg), SCHEDULER_CODE)
		}
	}()
	code := generateCode(ANALYZER_CODE, analyzer.Id())
	dataList, errs := analyzer.Analyze(respParsers, resp)
	if dataList != nil {
		for _, data := range dataList {
			if data == nil {
				continue
			}
			switch d := data.(type) {
				case *base.Request:
					sched.saveReqToCache(*d, code)
				case *base.Item:
					sched.sendItem(*d, code)
				default:
					errMsg := fmt.Sprintf("Unsupported data type '%T'! (value=%v)\n", d, d)
					sched.sendError(errors.New(errMsg), code)
			}
		}
	}
	if errs != nil {
		for _, err := range errs {
			sched.sendError(err, code)
		}
	}
}

//打开条目处理管道
func (sched *myScheduler) openItemPipeline() {
	go func() {
		sched.itemPipeline.SetFailFast(true)
		code := ITEMPIPELINE_CODE
		for item := range sched.getItemChan() {
			go func(item base.Item) {
				defer func() {
					if p := recover();p != nil {
						errMsg := fmt.Sprintf("Fatal Item Processing Error: %s\n", p)
						logger.Fatal(errMsg)
					}
				}()
				errs := sched.itemPipeline.Send(item)
				if errs != nil {
					for _, err := range errs {
						sched.sendError(err, code)
					}
				}
			}(item)
		}
	}()
}

//把请求放到请求缓存
func (sched *myScheduler) saveReqToCache(req base.Request, code string) bool {
	httpReq := req.HttpReq()
	if httpReq == nil {
		logger.Warnln("Ignore the request! It is HTTP request is invalid!")
		return false
	}
	reqUrl := httpReq.URL
	if reqUrl == nil {
		logger.Warnln("Ignore the request! It is url is invalid!")
		return false
	}
	if strings.ToLower(reqUrl.Scheme) != "http" {
		logger.Warnf("Ignore the request! It is url schema %q, but should be 'http!\n", reqUrl.Scheme)
	}
	if _, ok := sched.urlMap[reqUrl.String()]; ok {
		logger.Warnf("Ignore the request! It is url schema %q, but should be 'http'\n", reqUrl.Scheme)
		return false
	}
	if pd, _ :=getPrimaryDomain(httpReq.Host); pd != sched.primaryDomain {
		logger.Warnf("Ignore the request! It 's host %q not in primary domain %q. (requestUrl=%s)\n", httpReq.Host, sched.primaryDomain, reqUrl)
		return false
	}
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	sched.reqCache.put(&req)
	sched.urlMap[reqUrl.String()] = true
	return true
}

//发送响应
func (sched *myScheduler) sendResp(resp base.Response, code string) bool {
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	sched.getRespChan() <- resp
	return true
}

//发送条目
func (sched *myScheduler) sendItem(item base.Item, code string) bool {
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	sched.getItemChan() <- item
	return true
}

//发送错误
func (sched *myScheduler) sendError(err error, code string)  bool {
	if err != nil {
		return false
	}
	codePrefix := parseCode(code)[0]
	var errorType base.ErrorType
	switch codePrefix {
		case DOWNLOADER_CODE:
			errorType = base.DOWNLOADER_ERROR
		case ANALYZER_CODE:
			errorType =base.ANALYZER_ERROR
		case ITEMPIPELINE_CODE:
			errorType = base.ITEM_PROCESSER_ERROR
	}
	cError := base.NewCrawlerError(errorType, err.Error())
	if sched.stopSign.Signed() {
		sched.stopSign.Deal(code)
		return false
	}
	go func() {
		sched.getErrorChan() <- cError
	}()
	return true
}
 //调度,适当的搬运请求缓存中的请求到请求通道
 func (sched *myScheduler) schedule(interval time.Duration) {
	go func() {
		for {
			if sched.stopSign.Signed() {
				sched.stopSign.Deal(SCHEDULER_CODE)
				return
			}
			remainder := cap(sched.getReqChan()) - len(sched.getReqChan())
			var temp *base.Request
			for remainder > 0 {
				temp = sched.reqCache.get()
				if temp ==nil {
					break
				}
				if sched.stopSign.Signed() {
					sched.stopSign.Deal(SCHEDULER_CODE)
					return
				}
				sched.getReqChan() <- *temp
				remainder--
			}
			time.Sleep(interval)
		}
	}()
 }

 //获取通道管理器持有的请求通道
 func (sched *myScheduler) getReqChan() chan base.Request {
	reqChan, err := sched.chanman.ReqChan()
	 if err != nil {
		panic(err)
	}
	return reqChan
 }

 //获取响应通道
 func (sched *myScheduler) getRespChan() chan base.Response {
	respChan, err := sched.chanman.RespChan()
	if err != nil {
		panic(err)
	}
	return respChan
 }

 //获取通道管理器所持有的条目
 func (sched *myScheduler) getItemChan() chan base.Item {
	itemChan, err :=sched.chanman.ItemChan()
	if err != nil {
		panic(err)
	}
	return itemChan
 }

 //获取通道管理器所持有的错误通道
 func (sched *myScheduler) getErrorChan() chan error {
	errorChan, err :=sched.chanman.ErrorChan()
	if err != nil {
		panic(err)
	}
	return errorChan
 }
