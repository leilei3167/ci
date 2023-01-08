package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 操作基本类型
	fmt.Println("操作基本类型:")
	var i interface{} = 10
	var i1 interface{} = 11
	it := reflect.TypeOf(i)
	iv := reflect.ValueOf(i)

	it1 := reflect.TypeOf(&i1)
	iv1 := reflect.ValueOf(&i1)
	fmt.Printf("直接比较Type: %v it.Name:[%v] it.Kind:[%v] it.String:[%v]\n",
		it == it1, it.Name(), it.Kind(), it.String())

	// 如果要修改对象值,必须是可设置的,即,其必须为指针
	fmt.Println("iv: ", iv.CanSet())
	fmt.Printf("iv1 Kind:%v CanSet:%v\n", iv1.Kind(), iv1.CanSet()) // Kind:ptr CanSet:false
	// 注意,此处是interface这个空接口的指针,Kind是interface
	fmt.Printf("iv1.Elem Kind:%v CanSet:%v\n", iv1.Elem().Kind(), iv1.Elem().CanSet()) // Kind:ptr CanSet:false

	s := true
	iv1.Elem().Set(reflect.ValueOf(s))        // 取指针指向的地址
	fmt.Printf("after set bool i1: %v\n", i1) // 空接口由int被设置为了bool

	// 操作复合类型
	// 1.结构体
	s1 := Person{
		Name: "leilei",
		Age:  10,
	}
	s2 := Person{
		Name: "dsa",
		Age:  19,
	}
	s1t := reflect.TypeOf(s1)
	s1v := reflect.ValueOf(s1)
	n := s1t.NumField() // 遍历字段
	for i := 0; i < n; i++ {
		fieldt := s1t.Field(i)
		fieldv := s1v.Field(i)
		if tag := fieldt.Tag.Get("mytag"); tag != "" { // 获取标签
			fmt.Printf("got tag:%v\n", tag)
		}
		// 获取字段的类型或值
		fmt.Printf("字段:%v 值:%v\n", fieldt.Name, fieldv.Interface())
	}
	// 调用方法
	fmt.Printf("\n调用方法:\n")
	m := s1v.MethodByName("Print")
	// 设置参数
	ins := []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(10)}
	result := m.Call(ins) // 获取返回值(注意不要数组越界)
	fmt.Println(result[0].Int())

	fmt.Printf("\n反射的可设置性:\n")
	s2v := reflect.ValueOf(&s2)
	fmt.Printf("s1 CanSet:%v\n", s1v.CanSet())
	fmt.Printf("&s2 CanSet:%v\n", s2v.CanSet())
	// 如果要修改值,传入的必须是指针类型,并且使用Elme获取指针指向的值
	fmt.Printf("&s2.Elem CanSet:%v\n", s2v.Elem().CanSet())

	// 如果本身是引用类型,如切片或map,则获取到其元素就可直接修改元素
	sl := []int{1, 2, 3}
	slv := reflect.ValueOf(sl)
	sl0 := slv.Index(0)
	fmt.Printf("sl.Index(0) CanSet:%v\n", sl0.CanSet())
	fmt.Printf("sl CanSet:%v\n", slv.CanSet()) // 修改切片本身是不可以的,需要传入指针

	slv2 := reflect.ValueOf(&sl)
	sl0 = slv2.Elem().Index(0)
	fmt.Printf("&sl.Index(0) CanSet:%v\n", sl0.CanSet())
	fmt.Printf("&sl CanSet:%v\n", slv2.Elem().CanSet())
	slv2.Elem().SetLen(0) // 清空切片
	// slv2.Elem().Set(reflect.MakeSlice(slv2.Elem().Type(), 0, 0)) //创建新空切片也可
	fmt.Printf("清空切片后:%+v\n", sl)
	slv2.Elem().Set(reflect.Append(slv2.Elem(), reflect.ValueOf(1))) // 追加
	fmt.Printf("追加切片后:%+v\n", sl)
	// 下标截取切片
	// newS := slv2.Elem().Slice(0, slv2.Elem().Len())
}

type Person struct {
	Name string
	Age  int `mytag:"agelalala"`
}

func (p Person) Print(a, b int) int {
	fmt.Printf("%+v:%v\n", p, a+b)
	return a + b
}
