package cmap

import (
	"fmt"
	"sync"
)

//https://github.com/orcaman/concurrent-map/blob/master/concurrent_map.go

var SHARD_COUNT = 32

type Stringer interface { //一般接口,仅能用于类型约束,因为同时包含了方法和类型
	fmt.Stringer
	comparable
}

// ConcurrentMap 底层是一个map结构的切片,每一个切片有自己的锁,使锁的粒度最小化
type ConcurrentMap[K comparable, V any] struct {
	shards   []*ConcurrentMapShared[K, V]
	sharding func(key K) uint32
}

// ConcurrentMapShared 代表着最小的一个map单元
type ConcurrentMapShared[K comparable, V any] struct {
	items map[K]V //底层数据用map
	sync.RWMutex
}

// 泛型函数
func create[K comparable, V any](sharding func(key K) uint32) ConcurrentMap[K, V] {
	m := ConcurrentMap[K, V]{
		sharding: sharding,
		shards:   make([]*ConcurrentMapShared[K, V], SHARD_COUNT), //默认分片数量32
	}
	//初始化每一个分片
	for i := 0; i < SHARD_COUNT; i++ {
		m.shards[i] = &ConcurrentMapShared[K, V]{items: make(map[K]V)}
	}
	return m
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func strfnv32[K fmt.Stringer](key K) uint32 {
	return fnv32(key.String())
}

// New 创建一个map,使用默认的分片哈希方法
func New[V any]() ConcurrentMap[string, V] {
	return create[string, V](fnv32) //以string为key
}

func NewStringer[K Stringer, V any]() ConcurrentMap[K, V] {
	return create[K, V](strfnv32[K])
}

func NewWithCustomShardingFunction[K comparable, V any](sharding func(key K) uint32) ConcurrentMap[K, V] {
	return create[K, V](sharding)
}

// GetShard 将key哈希后取余,返回此key落在的分片
func (m ConcurrentMap[K, V]) GetShard(key K) *ConcurrentMapShared[K, V] {
	return m.shards[uint(m.sharding(key))%uint(SHARD_COUNT)]
}

// MSet 批量存入数据
func (m ConcurrentMap[K, V]) MSet(data map[K]V) {
	for k, v := range data {
		shard := m.GetShard(k) //获取分片
		shard.Lock()           //写入该分片当中,加锁保护
		shard.items[k] = v
		shard.Unlock()
	}
}

func (m ConcurrentMap[K, V]) Set(key K, value V) {
	shard := m.GetShard(key) //先hash查找分片
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}
