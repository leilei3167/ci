package orm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

/*
	反射的唯一的两个入口: TypeOf,ValueOf;分别返回reflect.Type和reflect.Value
	其中Value


*/

func ConstructSQL(obj interface{}) (stmt string, err error) {
	// 先判断是否是指针
	typ := reflect.TypeOf(obj)
	// 如果是指针,获取其底层值
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	// 必须是结构体
	if typ.Kind() != reflect.Struct {
		err = errors.New("must be a struct")
		return
	}
	buffer := bytes.NewBufferString("")
	buffer.WriteString("SELECT ")

	// 字段
	if typ.NumField() == 0 {
		err = errors.New("struct must have some fields")
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		// 依次获取字段
		field := typ.Field(i)

		if i != 0 {
			buffer.WriteString(", ")
		}
		column := field.Name
		// 获取tag
		if tag := field.Tag.Get("orm"); tag != "" {
			column = tag
		}
		buffer.WriteString(column)
	}
	stmt = fmt.Sprintf("%s FROM %s", buffer.String(), typ.Name())
	return
}

// 简单原生类型的反射
func Simple() {
	var b bool = true
	val := reflect.ValueOf(b)
	typ := reflect.TypeOf(b)
	fmt.Println(typ.Name(), val.Bool())

	// 如果是指针类型,必须通过Elem()获取其指向的值,才能调用对应的取值方法
	// val := reflect.ValueOf(&b)
	// typ := reflect.TypeOf(&b)
	// fmt.Println(typ.Elem().Name(), val.Elem().Bool())

	// 匿名的函数,或某些没有确定名称的类型,Name()将返回空字符串,可以使用String()返回器描述字符串
	fn := func(a, b int) int {
		return a + b
	}
	val = reflect.ValueOf(fn)
	typ = reflect.TypeOf(fn)
	fmt.Println(typ.Kind(), typ.String())

	// Kind返回类别,如不同变量的指针,类别都是ptr
	p := (*int)(nil)
	p1 := (*float64)(nil)
	fmt.Println(reflect.TypeOf(p).Kind(), reflect.TypeOf(p1).Kind())
}

func Compound() {
	// 切片
	sl := []int{5, 6, 6, 4, 2, 1}
	val := reflect.ValueOf(sl)
	typ := reflect.TypeOf(sl)

	fmt.Printf("\nlen:%d cap:%d\n", val.Len(), val.Cap())
	fmt.Printf("Type:%v String:%v Name:%v\n", typ.Kind(), typ.String(), typ.Name())
	// 反射遍历切片,使用Index()获取切片元素的值
	for i := 0; i < val.Len(); i++ {
		fmt.Printf("index:%d value:%d\n", i, val.Index(i).Int())
	}

	// map
	m := map[string]int{
		"tony": 1,
		"tom":  2,
		"john": 3,
	}
	val = reflect.ValueOf(m)
	typ = reflect.TypeOf(m)
	fmt.Printf("\nlen:%d\n", val.Len())
	fmt.Printf("Type:%v String:%v Name:%v\n", typ.Kind(), typ.String(), typ.Name())
	// 遍历map
	iter := val.MapRange() // 创建map迭代器
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		fmt.Printf("k:%v v:%v\n", k, v)
	}

	// 结构体
	p := Person{
		Name: "leilei",
		Age:  10,
		De:   "school",
	}

	val = reflect.ValueOf(p)
	typ = reflect.TypeOf(p)
	// 使用Field取出对应的字段
	fmt.Printf("\nName:%s Age:%d De:%s\n", val.Field(0).String(), val.Field(1).Int(), val.Field(2).String())
	fmt.Printf("Type:%v String:%v Name:%v\n", typ.Kind(), typ.String(), typ.Name()) // String()会包含包名

	// channel
	ch := make(chan int, 100)
	val = reflect.ValueOf(ch)
	typ = reflect.TypeOf(ch)
	ch <- 17
	fmt.Printf("\nlen:%d cap:%d\n", val.Len(), val.Cap()) // channel中的数量,以及缓冲区大小
	fmt.Printf("Type:%v String:%v Name:%v\n", typ.Kind(), typ.String(), typ.Name())
	// 操作通道
	v, ok := val.Recv()
	if ok {
		fmt.Printf("get:%v\n", v.Int())
	}
}

func Call() {
	// 调用函数
	f := reflect.ValueOf(Add)
	a, b := 1, 4
	// 转换为反射对象
	vals := []reflect.Value{reflect.ValueOf(a), reflect.ValueOf(b)}
	ret := f.Call(vals)
	fmt.Printf("\n调用结果:%v\n", ret[0].Int())

	// 调用方法,需要先获取到他的方法名称,传参数都是一样的
	c := reflect.ValueOf(Person{})
	m := c.MethodByName("Add")
	ret = m.Call(vals)
	fmt.Printf("方法调用结果:%v\n", ret[0].Int())
}

func Add(a, b int) int {
	return a + b
}

type Person struct {
	Name string
	Age  int
	De   string `orm:"depart_ment"`
}

func (p Person) Add(i, j int) int {
	return i + j
}

func Settable() {
	// 如果希望反射值能够被修改,则传入reflect.ValueOf的时候必须传递指针
	i := 17
	//	val := reflect.ValueOf(i)
	//	val.SetInt(1) // panic: reflect: reflect.Value.SetInt using unaddressable value
	val := reflect.ValueOf(&i)
	val.Elem().SetInt(1) // Elem相当于*,能够获取到指针指向的值,并做修改
	fmt.Println(i)

	// Value提供了CanSet,CanAddr,CanInterface,来判断反射对象是否可设置(settable),可寻址,可恢复为interface{}
	fmt.Printf("\n*int:\n")
	fmt.Printf("settable:%v,addr:%v,Interface:%v\n", val.CanSet(), val.CanAddr(), val.CanInterface())
	fmt.Printf("\n*int Elem:\n")
	fmt.Printf("settable:%v,addr:%v,Interface:%v\n", val.Elem().CanSet(), val.Elem().CanAddr(), val.Elem().CanInterface())

	fmt.Printf("\nint:\n")
	val = reflect.ValueOf(i)
	fmt.Printf("settable:%v,addr:%v,Interface:%v\n", val.CanSet(), val.CanAddr(), val.CanInterface())

	// 切片
	s1 := []int{1, 2, 3}
	val = reflect.ValueOf(s1)
	// 切片本身是不可设置的
	fmt.Printf("\nslice:\n")
	fmt.Printf("settable:%v,addr:%v,Interface:%v\n", val.CanSet(), val.CanAddr(), val.CanInterface())
	// 但是其某一个元素却是可以的(相当于通过下标修改切片元素)
	fmt.Printf("\nslice elem:\n")
	val = val.Index(0)
	fmt.Printf("settable:%v,addr:%v,Interface:%v\n", val.CanSet(), val.CanAddr(), val.CanInterface())
}
