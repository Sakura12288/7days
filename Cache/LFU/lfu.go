package LFU

import "container/list"

type entry struct {
	key, val, freq int
}

type LFUCache struct {
	keyToValue map[int]*list.Element
	freqToList map[int]*list.List
	capacity   int
	minFreq    int
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		capacity:   capacity,
		minFreq:    1,
		keyToValue: make(map[int]*list.Element),
		freqToList: make(map[int]*list.List),
	}
}

func (c *LFUCache) Remove() {
	ele := c.freqToList[c.minFreq].Back()
	entr := ele.Value.(*entry)
	delete(c.keyToValue, entr.key)
	c.freqToList[c.minFreq].Remove(ele)
	if c.freqToList[c.minFreq].Len() == 0 {
		c.minFreq++
	}
}

func (c *LFUCache) GetEntry(key int) *entry {
	ele := c.keyToValue[key]
	if ele == nil {
		return nil
	}
	entr := ele.Value.(*entry)
	entr.freq++
	c.freqToList[entr.freq-1].Remove(ele)
	if c.freqToList[entr.freq] == nil {
		c.freqToList[entr.freq] = list.New()
	}
	el := c.freqToList[entr.freq].PushFront(entr)
	c.keyToValue[entr.key] = el
	if c.freqToList[c.minFreq].Len() == 0 {
		c.minFreq++
	}
	return entr
}

func (c *LFUCache) PutEntry(entr *entry) {
	if e := c.GetEntry(entr.key); e != nil {
		e.val = entr.val
		return
	}
	entr.freq = 1
	if len(c.keyToValue) >= c.capacity {
		c.Remove()
	}
	c.minFreq = 1
	if c.freqToList[c.minFreq] == nil {
		c.freqToList[c.minFreq] = list.New()
	}
	e := c.freqToList[c.minFreq].PushFront(entr)
	c.keyToValue[entr.key] = e
}

func (this *LFUCache) Get(key int) int {
	en := this.GetEntry(key)
	if en == nil {
		return -1
	}
	return en.val
}

func (this *LFUCache) Put(key int, value int) {
	en := &entry{key: key, val: value}
	this.PutEntry(en)
}
