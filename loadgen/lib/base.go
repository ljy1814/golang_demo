package lib

import (
	"time"
)

type RawReq struct {
	Id int64
	Req []byte
}

type RawResp struct {
	Id int64
	Resp []byte
	Err error
	Elapse time.Duration
}

type ResultCode int

//保留1-1000给载荷承受者使用
const (
//	RESULT_CODE_SUCCESS	= 200 //成功
	RESULT_CODE_SUCCESS	ResultCode = 0 //成功
	RESULT_CODE_WARNING_CALL_TIMEOUT	ResultCode = 1001	//超时警告
	RESULT_CODE_ERROR_CALL	ResultCode = 2001	//调用错误
	RESULT_CODE_ERROR_RESPONSE	ResultCode = 2002	//响应内容错误
	RESULT_CODE_ERROR_CALEE	ResultCode = 2003	//被调用方(被测软件)内部错误
	RESULT_CODE_FATAL_CALL	ResultCode = 3001	//调用过程中发生了致命错误
)


func GetResultCodePlain(code ResultCode) string {
	var codePlain string
	switch code {
		case RESULT_CODE_SUCCESS:
			codePlain = "Success"
		case RESULT_CODE_WARNING_CALL_TIMEOUT:
			codePlain = "Call Timeout Warning"
		case RESULT_CODE_ERROR_CALL:
			codePlain = "Call Error"
		case RESULT_CODE_ERROR_RESPONSE:
			codePlain = "Response Error"
		case RESULT_CODE_ERROR_CALEE:
			codePlain = "Callee Error"
		case RESULT_CODE_FATAL_CALL:
			codePlain = "Call Fatal Error"
		default:
			codePlain = "Unknown result code"
	}
	return codePlain
}

//调用结果的结构体
type CallResult struct {
	Id int64		//ID,唯一标识
	Req RawReq		//原生请求
	Resp RawResp	//原生响应
	Code ResultCode	//响应码
	Msg	string	//结果成因的简述
	Elapse time.Duration	//耗时
}

//荷载发生器的状态类型
type GenStatus int

const (
	STATUS_ORIGINAL GenStatus = 0
	STATUS_STARTED GenStatus = 1
	STATUS_STOPPED GenStatus = 2
)

//荷载发生器的接口

type Generator interface {
	//启动荷载发生器
	Start()
	//停止荷载发生器
	//第一个结果代表已发生的荷载总数,而且在第二个结果值为true时有效
	//第二个结果值表示是否成功将荷载发生器转换为已停止状态
	Stop() (uint64, bool)
	//获取状态
	Status() GenStatus
}
