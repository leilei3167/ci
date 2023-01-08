package main

import (
	"fmt"
	"reflect"

	"github.com/leilei3167/ci/reflection/orm"
)

func main() {
	p := new(orm.Person)
	sql, err := orm.ConstructSQL(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)

	orm.Simple()
	orm.Compound()
	orm.Call()
	orm.Settable()

	fmt.Println("---------------------------------------------------------")
	m := NewManager()
	m.Put(A{})
	m.Put(B{})
	m.Put(B{})
	m.Put(B{})
	m.Put(B{})
	m.Put(A{})
	m.Put(struct{}{}) // 未知类型
	m.Range()
	fmt.Println("--------------------Send----------------")
	m.Send()
	fmt.Println("--------------reset----------------")
	m.Reset()
	m.Range()
}

type Flusher interface {
	Flush()
}

type (
	A  struct{}
	As []A
)

func (a As) Flush() {
	fmt.Println("As")
}

type (
	B  struct{}
	Bs []B
)

func (b Bs) Flush() {
	fmt.Println("Bs")
}

type Manager struct {
	queue []Flusher
	m     map[string]Flusher
	len   int
}

func NewManager() *Manager {
	c := &Manager{}
	q := []Flusher{
		&Bs{}, &As{},
	}
	m := make(map[string]Flusher)

	for _, v := range q {
		// 获取类型(指针需要使用Elem)
		t := reflect.TypeOf(v).Elem()
		// 以其元素名为key
		fmt.Printf("元素为:%v\n", t.Elem().Name())
		m[t.Elem().Name()] = v
	}
	c.queue = q
	c.m = m
	c.len = 0
	return c
}

func (m *Manager) Put(toPut interface{}) {
	// 获取类型名
	t := reflect.TypeOf(toPut)
	k := t.Name() // 不能是指针哦
	tv := reflect.ValueOf(toPut)

	q, ok := m.m[k]
	if ok {
		// 取出存放的底层值
		qv := reflect.ValueOf(q).Elem() // 可设置的切片
		qv.Set(reflect.Append(qv, tv))
		m.len++
	} else {
		fmt.Printf("未知类型:%v %T\n", k, toPut)
	}
}

func (m *Manager) Reset() {
	m.len = 0
	for _, q := range m.m {
		qv := reflect.ValueOf(q).Elem() // 可设置的切片
		qv.Set(reflect.MakeSlice(qv.Type(), 0, 0))
	}
}

func (m *Manager) Range() {
	for _, q := range m.m {
		qv := reflect.ValueOf(q).Elem() // 可设置的切片
		l := qv.Len()
		for i := 0; i < l; i++ {
			fmt.Printf("[%s slice] idex:%d v:%v\n", qv.Type().Name(), i, qv.Index(i))
		}
	}
}

func (m *Manager) Send() {
	for _, q := range m.m {
		qv := reflect.ValueOf(q).Elem() // 可设置的切片
		l := qv.Len()
		// 相当于[0:len()],还原成接口,强制转换为Flusher
		newSlice := qv.Slice(0, l).Interface().(Flusher)
		fmt.Printf("newSlice Type:%T v:%v\n", newSlice, newSlice)
	}
}
