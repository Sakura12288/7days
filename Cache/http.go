package Cache

import (
	"7days/Cache/consistent"
	"7days/Cache/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//类似于中间点，用来节点之间相互通信，基于http实现

const (
	defaultpath     = "/_geecache/" //默认前缀
	defaultreplices = 50
)

type HTTPPool struct {
	selfpath    string //IP地址
	basepath    string //前缀
	mu          sync.Mutex
	peers       *consistent.Map
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(selfpath string) *HTTPPool {
	return &HTTPPool{
		selfpath: selfpath,
		basepath: defaultpath,
	}
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//用于节点间通讯，检查前缀

	if !strings.HasPrefix(r.URL.Path, h.basepath) {
		panic("can found " + r.URL.Path)
	}
	h.Log("%s %s", r.Method, r.URL.Path)
	//提取缓存名字以及键
	path := strings.SplitN(r.URL.Path[len(h.basepath):], "/", 2)
	if len(path) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	name := path[0] //缓存名
	key := path[1]
	g := GetGroup(name)
	if g == nil {
		http.Error(w, "不存在"+name, http.StatusNotFound)
		return
	}
	v, err := g.Get(key)
	body, err := proto.Marshal(&pb.Response{Value: v.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	//转义保证安全
	pa := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	res, err := http.Get(pa)
	if err != nil {
		return fmt.Errorf("sbbba %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}
	w, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response err %v", err)
	}
	if err = proto.Unmarshal(w, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

var _ PeerGetter = (*httpGetter)(nil)

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistent.NewMap(defaultreplices, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basepath} //ip地址与端口不一样，文件地址保持一样
	}
}

func (h *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if pe := h.peers.Get(key); pe != "" && pe != h.selfpath {
		//h.Log
		h.Log("Pick peer rrr %s", pe)
		return h.httpGetters[pe], true
	}
	return nil, false
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.selfpath, fmt.Sprintf(format, v...))
}

var _ PeerPicker = (*HTTPPool)(nil)
