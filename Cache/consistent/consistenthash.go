package consistent

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//一致性哈希算法，实现节点的选择，同时预防缓存雪崩

//可以试着实现删除某个节点

type Hash func([]byte) uint32
type Map struct {
	hash     Hash
	replices int            //虚拟节点的个数
	keys     []int          //存储所有的节点
	search   map[int]string //通过节点编号找节点名字
}

func NewMap(replices int, hash Hash) *Map {
	m := &Map{
		replices: replices,
		hash:     hash,
		search:   make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//添加节点

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 1; i <= m.replices; i++ {
			va := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, va)
			m.search[va] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	v := int(m.hash([]byte(key)))
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= v
	})
	return m.search[m.keys[index%len(m.keys)]]
}
