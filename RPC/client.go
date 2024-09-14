package RPC

import (
	"7days/RPC/codec"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

//call 代表一次调用

type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Done          chan *Call //一次调用结束标志
	Error         error
}

func (c *Call) done() {
	c.Done <- c
}

type Client struct {
	cc       codec.Codec
	seq      uint64
	sending  sync.Mutex
	mu       sync.Mutex
	pending  map[uint64]*Call //剩下未处理的调用
	opt      *Option
	closing  bool //用户关闭
	shutdown bool //服务器关闭，或不正常关闭
	h        codec.Header
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("连接被关闭")

//关闭连接

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}
	c.closing = true
	return c.cc.Close()
}

//判断cilent是否可用

func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.closing && !c.shutdown
}

//一次调用发起，存入client中
//call的序号由client分配

func (c *Client) registerCall(call *Call) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing || c.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = c.seq
	c.pending[call.Seq] = call
	c.seq++
	return call.Seq, nil
}

func (c *Client) removeCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	call := c.pending[seq]
	delete(c.pending, seq)
	return call
}

//由于产生了错误，所以终止，错误即为参数err

func (c *Client) terminateCalls(err error) {
	c.sending.Lock()
	defer c.sending.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shutdown = true
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
}

//call 不存在，可能是请求没有发送完整，或者因为其他原因被取消，但是服务端仍旧处理了。
//call 存在，但服务端处理出错，即 h.Error 不为空。
//call 存在，服务端处理正常，那么需要从 body 中读取 Reply 的值。

func (c *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = c.cc.ReadHeader(&h); err != nil {
			break
		}
		call := c.removeCall(h.Seq) //收到回复代表调用结束，可以移除此次调用
		switch {
		case call == nil:
			err = c.cc.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = c.cc.ReadBody(nil)
			call.done()
		default:
			err = c.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New(" 读取内容失败 " + err.Error())
			}
			call.done()
		}
	}
	//跳出循环则证明发生错误，终止剩下的
	c.terminateCalls(err)
}

//创建客户端
//创建 Client 实例时，首先需要完成一开始的协议交换，即发送 Option 信息给服务端。协商好消息的编解码方式之后，
//再创建一个子协程调用 receive() 接收响应。

//实现HTTP链接

func NewHTTPClient(conn net.Conn, opt *Option) (client *Client, err error) {
	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", defaultRPCPath))
	fmt.Println("33")
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("nihao", p)
			fmt.Println(err)
		}
	}()
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	fmt.Println("kk")
	if resp == nil {
		fmt.Println("sb")
	}
	fmt.Println(resp.Status, err)
	if err == nil && resp.Status == connected {
		fmt.Println("22")
		return NewClient(conn, opt)
	}
	if err == nil {
		err = errors.New("unexpected HTTP response: " + resp.Status)
	}
	return nil, err
}

func NewClient(conn net.Conn, opt *Option) (client *Client, err error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err = fmt.Errorf("编码方式不存在 %s", opt.CodecType)
		log.Println("客户端 ： 编码错误", err)
		return nil, err
	}
	if err = json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("option 出错", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq:     1,
		pending: make(map[uint64]*Call),
		opt:     opt,
		cc:      cc,
	}
	go client.receive() //注意开始启动客户端
	return client
}

//解析协议

func parseOptions(opt ...*Option) (*Option, error) {
	if len(opt) == 0 || opt[0] == nil {
		return DefaultOption, nil
	}
	if len(opt) != 1 {
		err := fmt.Errorf("opts 个数超过一个了")
		return nil, err
	}
	op := opt[0]
	op.MagicNumber = DefaultOption.MagicNumber
	if op.CodecType == "" {
		op.CodecType = DefaultOption.CodecType
	}
	return op, nil
}

//超时处理机制

type newClientFunc func(conn net.Conn, opt *Option) (client *Client, err error)

type clientResult struct {
	client *Client
	err    error
}

func dialTimeout(f newClientFunc, network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	fmt.Println("1")
	if err != nil {
		log.Println("opts错误", err)
		return nil, err
	}
	fmt.Println("2")
	conn, err := net.DialTimeout(network, address, opt.ConnectionTimeout)
	if err != nil {
		log.Println("无法连接服务端", err)
		return nil, err
	}
	fmt.Println("3")
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()
	ch := make(chan clientResult)
	fmt.Println("4")
	go func() {
		client, err = f(conn, opt)
		ch <- clientResult{client: client, err: err}
	}()
	fmt.Println("5")
	if opt.ConnectionTimeout == 0 {
		result := <-ch
		fmt.Println("6")
		return result.client, result.err
	}
	select {
	case <-time.After(opt.ConnectionTimeout):
		return nil, fmt.Errorf("客户端: 连接超时，希望 %s 之内", opt.ConnectionTimeout)
	case result := <-ch:
		fmt.Println("nini")
		return result.client, result.err
	}
}
func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	return dialTimeout(NewClient, network, address, opts...)
}

func DialHTTP(network, addr string, opts ...*Option) (client *Client, err error) {
	return dialTimeout(NewHTTPClient, network, addr, opts...)
}

//通用接口
// XDial calls different functions to connect to a RPC server
// according the first parameter rpcAddr.
// rpcAddr is a general format (protocol@addr) to represent a rpc server
// eg, http@10.0.0.1:7001, tcp@10.0.0.1:9999, unix@/tmp/geerpc.sock

func XDial(rpcAddr string, opts ...*Option) (client *Client, err error) {
	parts := strings.Split(rpcAddr, "@")
	if len(parts) != 2 {
		log.Println("地址有误 :", rpcAddr)
		return nil, fmt.Errorf("rpc client err: wrong format '%s', expect protocol@addr", rpcAddr)
	}
	switch parts[0] {
	case "http":
		return DialHTTP("tcp", parts[1], opts...)
	default:
		return Dial(parts[0], parts[1], opts...)
	}
}

func (c *Client) send(call *Call) {
	c.sending.Lock()
	defer c.sending.Unlock()
	seq, err := c.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	c.h.Seq = seq
	c.h.ServiceMethod = call.ServiceMethod
	c.h.Error = ""
	if err = c.cc.Write(&c.h, call.Args); err != nil {
		call := c.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
			return
		}
	}
}

func (c *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	c.send(call)
	return call
}

func (c *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	call := c.Go(serviceMethod, args, reply, make(chan *Call, 1))
	//ctx 可设置withtimeout判断是否超时
	select {
	case <-ctx.Done():
		c.removeCall(call.Seq)
		return errors.New("rpc client : call failed" + ctx.Err().Error())
	case call := <-call.Done:
		return call.Error
	}
}

//例如
//ctx, _ := context.WithTimeout(context.Background(), time.Second)
//var reply int
//err := client.Call(ctx, "Foo.Sum", &Args{1, 2}, &reply)
//...
