package xclient

import (
	"7days/RPC"
	"context"
	"reflect"
	"sync"
)

type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *RPC.Option
	mu      sync.Mutex
	clients map[string]*RPC.Client //以地址映射服务端，不存在就创建
}

func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.clients {
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

func NewXClient(d Discovery, mode SelectMode, opt *RPC.Option) *XClient {
	return &XClient{
		opt:     opt,
		d:       d,
		mode:    mode,
		clients: make(map[string]*RPC.Client),
	}
}

func (xc *XClient) dial(rpcAddr string) (*RPC.Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	client, ok := xc.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(xc.clients, rpcAddr)
		client = nil
	}
	client, err := RPC.XDial(rpcAddr, xc.opt)
	if err == nil {
		xc.clients[rpcAddr] = client
		return client, err
	}
	return nil, err
}

func (xc *XClient) call(rpcAddr string, serviceMethod string, ctx context.Context, args, reply interface{}) error {
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	err = client.Call(ctx, serviceMethod, args, reply)
	return err
}

func (xc *XClient) Call(serviceMethod string, ctx context.Context, args, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, serviceMethod, ctx, args, reply)
}

//调用所有的服务处理，得到结果则返回

func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	services, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error
	replyDone := (reply == nil)
	ctx, cancel := context.WithCancel(ctx)
	for _, rpcAddr := range services {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()
			var clonedReply interface{}
			if reply != nil {
				clonedReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			err := xc.call(rpcAddr, serviceMethod, ctx, args, clonedReply)
			mu.Lock()
			if err != nil && e == nil {
				e = err
				cancel()
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clonedReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
