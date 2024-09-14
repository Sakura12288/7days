package lru

import "container/list"

type Cache struct {
	maxbyte   int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	Onevicted func(key string, value Value)
}
type entry struct {
	key   string //键
	value Value
}

type Value interface {
	Len() int
}

func New(maxbyte int64, onevicted func(key string, value Value)) *Cache {
	return &Cache{
		maxbyte:   maxbyte,
		ll:        list.New(),
		Onevicted: onevicted,
		cache:     make(map[string]*list.Element), //Element 以后面定义的&entry做值
	}
}
func (this *Cache) RemoveOldest() {
	ele := this.ll.Back()
	if ele != nil {
		en := ele.Value.(*entry)
		key, value := en.key, en.value
		this.nbytes -= int64(value.Len()) + int64(len(key))
		delete(this.cache, key)
		this.ll.Remove(ele)
		if this.Onevicted != nil {
			this.Onevicted(key, value)
		}

	}
}

func (this *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := this.cache[key]; ok {
		kv := ele.Value.(*entry)
		this.ll.MoveToFront(ele)
		return kv.value, true
	}
	return
}

func (this *Cache) Add(key string, value Value) {
	ele := &entry{key, value}
	if e, ok := this.cache[key]; ok {
		this.ll.MoveToFront(e)
		en := e.Value.(*entry)
		this.nbytes -= int64(en.value.Len() - value.Len())
		this.cache[key] = this.ll.PushFront(ele)
	} else {
		this.cache[key] = this.ll.PushFront(ele)
		this.nbytes += int64(len(key) + value.Len())
	}
	for this.maxbyte > 0 && this.nbytes > this.maxbyte {
		this.RemoveOldest()
	}
}

func (this *Cache) Len() int {
	return this.ll.Len()
}
