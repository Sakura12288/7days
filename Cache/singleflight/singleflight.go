package singleflight

import "sync"

//为了避免瞬间多个客户端同时访问一个节点，导致缓存击穿
//可以试试用管道实现

type Group struct {
	mu sync.Mutex
	m  map[string]*call //这里保留对节点的访问记录，当有节点正在访问时，其他也想访问的节点阻塞，直到第一个节点释放，从而得到它的结果
}

type call struct {
	w     sync.WaitGroup //用来添加子线程，进行等待
	value interface{}
	err   error
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//如果ok代表第一个已取到数据,第一个执行的会跳过
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.w.Wait() //等待第一个取数据
		return c.value, c.err
	}
	c := new(call)
	c.w.Add(1)
	g.m[key] = c
	g.mu.Unlock() //尽管已经解锁，但是之后的节点还是会卡在wait，除非第一个执行完
	c.value, c.err = fn()
	c.w.Done()
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.value, c.err

}
