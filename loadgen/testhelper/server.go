package testhelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"
//	"demo/display"
)

type ServerReq struct {
	Id int64
	Operands []int
	Operator string
}

type ServerResp struct {
	Id int64
	Formula string
	Result int
	Err error
}

func op(operands []int, operator string) int {
	var result int
	switch {
	case operator == "+":
		for _, v := range operands {
			if result == 0 {
				result  = v
			} else {
				result += v
			}
		}
	case operator == "-":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result -= v
			}
		}
	case operator == "*":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result *= v
			}
		}
	case operator == "/":
		for _, v := range operands {
			if result == 0 {
				result = v
			} else {
				result /= v
			}
		}
	}
	return result
}

func genFormula(operands []int, operator string, result int, equal bool) string {
	var buff bytes.Buffer
	n := len(operands)
	for i := 0; i < n; i++ {
		if i > 0 {
			buff.WriteString(" ")
			buff.WriteString(operator)
			buff.WriteString(" ")
		}
		buff.WriteString(strconv.Itoa(operands[i]))
	}
	if equal {
		buff.WriteString(" = ")
	} else {
		buff.WriteString(" != ")
	}
	buff.WriteString(strconv.Itoa(result))
	return buff.String()
}

func reqHandler(conn net.Conn) {
	var errMsg  string
	var sresp ServerResp
	req, err := read(conn, DELIM)
	if err != nil {
		errMsg = fmt.Sprintf("Server: Req Read Error : %s\n", err)
	} else {
		var sreq ServerReq
		if err := json.Unmarshal(req, &sreq); err != nil {
			errMsg = fmt.Sprintf("Server: Req Unmarshal Error: %s", err)
		} else {
			sresp.Id = sreq.Id
			sresp.Result = op(sreq.Operands, sreq.Operator)
			sresp.Formula = genFormula(sreq.Operands, sreq.Operator, sresp.Result, true)
		}
	}
	if errMsg != "" {
		sresp.Err = errors.New(errMsg)
	}
	bytes, err := json.Marshal(sresp)
	if err != nil {
		fmt.Errorf("Server: Resp Marshal Error : %s", err)
	}
	_, err = write(conn, bytes, DELIM)
	if err != nil {
		fmt.Errorf("Server: Resp Write error: %s", err)
	}
}

type TcpServer struct {
	listener net.Listener
	active bool
	lock *sync.Mutex
}

func (self *TcpServer) init(addr string) error {
//	fmt.Println("enter TcpServer init...")
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		return nil
	}
	ln, err := net.Listen("tcp", addr)
//	fmt.Println(err)
	if err != nil {
//		display.Display("listener err", err)
		return err
	}
	self.listener = ln
	self.active = true
//	display.Display("TcpServer init", self)
	return nil
}

func (self *TcpServer) Listen(addr string) error {
//	display.Display("TcpServer", self)
//	fmt.Println("------------")
	if err := self.init(addr); err != nil {
//		display.Display("listener err", err)
		return err
	}
	go func(active *bool) {
		for {
			conn, err := self.listener.Accept()
//			display.Display("conn", conn)
			if err != nil {
				fmt.Errorf("Server: Request Acception Error: %s\n", err)
				continue
			}
			//处理请求
			go reqHandler(conn)
			runtime.Gosched()
		}
	}(&self.active)
	return nil
}

func (self *TcpServer) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.listener.Close()
		self.active = false
		return true
	} else {
		return false
	}
}

func NewTcpServer() *TcpServer {
	return &TcpServer{lock: new(sync.Mutex)}
}
