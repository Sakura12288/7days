package RPC

import (
	"7days/RPC/codec"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber       int           //标记这是geerpc的请求
	CodecType         codec.Type    //编码方式
	ConnectionTimeout time.Duration //0 代表无限制
	HandleTimeout     time.Duration
}

var DefaultOption = &Option{
	MagicNumber:       MagicNumber,
	CodecType:         codec.GobType,
	ConnectionTimeout: 10 * time.Second,
}

//Option固定用json编码，其决定的是后面的头部以及内容

//接下来是服务端，注意解析顺序，Option >- header >- body

type Server struct {
	services sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("连接失败")
			return
		}
		//开启协程，并发处理
		go s.ServeConn(conn)
	}
}

// 改为支持HTTP协议
const (
	connected        = "200 Connected to Gee RPC"
	defaultRPCPath   = "/_geeprc_"
	defaultDebugPath = "/debug/geerpc"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Println("Hijack 失败", req.RemoteAddr, ":", err.Error())
		return
	}
	_, _ = io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	fmt.Println("11")
	s.ServeConn(conn)
}

// HandleHTTP registers an HTTP handler for RPC messages on rpcPath.
// It is still necessary to invoke http.Serve(), typically in a go statement.

func (s *Server) HandleHTTP() {
	http.Handle(defaultRPCPath, s)
	http.Handle(defaultDebugPath, DebugServer{s})
	log.Println("rpc server debug path:", defaultDebugPath)
}

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}

// 默认处理
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.

func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var option = &Option{}
	if err := json.NewDecoder(conn).Decode(option); err != nil {
		log.Println("解析option出问题", err)
		return
	}
	if option.MagicNumber != MagicNumber {
		log.Println("MagicNumber 错误", option.MagicNumber)
		return
	}
	//获取编码方式
	f := codec.NewCodecFuncMap[option.CodecType]

	if f == nil {
		log.Println("编码方式不存在", option.CodecType)
		return
	}
	s.serveCodec(f(conn), option)
}

// invalidRequest is a placeholder for response argv when error occurs

var invalidRequest = struct{}{}

func (s *Server) serveCodec(cc codec.Codec, opt *Option) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	//同时处理多个请求
	for {
		req, err := s.ReadRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handlerRequest(cc, req, sending, wg, opt.HandleTimeout)

	}
	wg.Wait()
	_ = cc.Close()
}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	mtype        *methodType
	svc          *service
}

func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("read header error", err)
		}
		return nil, err
	}
	return &h, nil
}

func (s *Server) ReadRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.svc, req.mtype, err = s.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()
	// make sure that argvi is a pointer, ReadBody need a pointer as parameter
	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("server:读取内容出现错误", err)
		return req, err
	}
	return req, nil
}

func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("回复应答出错", err)
		return
	}
	return
}

func (s *Server) handlerRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	called := make(chan struct{})
	sent := make(chan struct{})
	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()
	if timeout == 0 {
		<-called
		<-sent
		return
	}
	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("服务端：处理请求超时， 希望时间：%s", timeout)
		s.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

func (s *Server) Register(rcvr interface{}) error {
	ser := newService(rcvr)
	if _, ok := s.services.LoadOrStore(ser.name, ser); ok {
		return errors.New("该服务已被注册：" + ser.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

//找service， 调用一般是"structName.methodName"

func (s *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	index := strings.LastIndex(serviceMethod, ".")
	if index == -1 || index == 0 || index == len(serviceMethod)-1 {
		err = errors.New("调用名字不合理" + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:index], serviceMethod[index+1:]
	svci, ok := s.services.Load(serviceName)
	if !ok {
		err = errors.New("该服务不存在" + serviceName)
		return
	}
	svc = svci.(*service)
	mtype, ok = svc.method[methodName]
	if !ok {
		err = errors.New("该方法不存在" + methodName)
		return
	}
	return
}
