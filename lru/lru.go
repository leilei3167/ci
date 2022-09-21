package lru

import (
	"container/list"
)

//每发生一次数据库的读取操作,首先检查待读取的数据是否存在于缓存中,若是,则缓存命中,返回数据
//并将此缓存移动至最新
//否则才从数据库中读取数据,并将读取成功的数据添加到缓存中;在向缓存中添加数据时,如果缓存已满
//则需要删除最早的那条记录,这种更新缓存的方法就是LRU(Least Recently Used),即最近最少使用

//读写的时间复杂度都应该在O(1),map可以实现查询时O(1),但是更新时由于需要先找到最早的记录,需要
//遍历,无法满足要求;因此可以通过map+链表的形式实现
//将缓存项按照访问时间顺序链接起来

type Cache struct {
	MaxEntries int

	ll    *list.List               //链表
	cache map[string]*list.Element //记录链表的元素
}

type entry struct {
	key   string
	value any
}

func New(max int) *Cache {
	return &Cache{
		MaxEntries: max,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (c *Cache) Set(key string, value any) {
	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.ll = list.New()
	}

	//如果其已在缓存中,则将其移动至最前端,并更新值
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)

		ee.Value.(*entry).value = value

		return
	}

	//不存在缓存中,则在最前端插入,并更新map,注意map的key对应的是链表的元素
	ele := c.ll.PushFront(&entry{key: key, value: value})
	c.cache[key] = ele

	//判断是否超出容量,超出的话删除最末尾的
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value any, ok bool) {
	if c.cache == nil {
		return
	}

	if ele, hit := c.cache[key]; hit { //查找到,则更新缓存
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}

	return
}

func (c *Cache) Remove(key string) { //删除指定的key对应的值
	if c.cache == nil {
		return
	}

	if ele, hit := c.cache[key]; hit {
		c.removeEle(ele)
	}
}

func (c *Cache) Len() int {
	if c.cache == nil || c.ll == nil {
		return 0
	}

	return c.ll.Len()
}

func (c *Cache) Clear() {
	if c.cache == nil {
		return
	}

	c.ll = list.New()
	c.cache = make(map[string]*list.Element)
}

func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}

	ele := c.ll.Back()

	if ele != nil {
		c.removeEle(ele)
	}
}

func (c *Cache) removeEle(e *list.Element) {
	c.ll.Remove(e) //从链表删除,也要从map中删除
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}
