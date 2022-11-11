package main

import "fmt"

//参考来源;shttps://segmentfault.com/a/1190000041634906

// 如果要写一个2数相加的函数,如果希望他支持浮点,或者uint,那么就需要重复的写多个函数,如何解决这个问题?
// 如果你经常要分别为不同的类型写完全相同逻辑的代码，那么使用泛型将是最合适的选择
func Add(a int, b int) int {
	return a + b
}

// 逻辑差不多但是需要重复编写函数
func AddFloat32(a float32, b float32) float32 {
	return a + b
}

// 此处的`a string, b string` 被称之为形参列表,a b 形参的类型为int,不具备任何的值,定义函数时并不确定其值,只有我们调用函数传入实参(argument) 之后才有具体的值。
func AddString(a string, b string) string {
	return a + b
}

// 使用泛型写一个泛型函数,使用时需要指定类型实参 AddGeneric[int](2,3)
// Go中是支持自动类型推导的
func AddGeneric[T int | float32 | string](a, b T) T {
	return a + b
}

//将形参的类型也引入 形参和实参的概念,称之为 类型形参和类型实参,在定义函数时,并不指定具体的值,而是在调用时传入类型实参后被指定
//通过引入 类型形参 和 类型实参 这两个概念，我们让一个函数获得了处理多种不同类型数据的能力，这种编程方式被称为 泛型编程

// 如定义一个切片:
type aSlice []int //此切面只能容纳int类型的元素,那如果想要能够存float的类型呢?
//以前只能再定义一个:
//type bSlice []float32

// 引入泛型,在声明时将原先写死的int类型,写作类型形参,也就是类型实参是int float32 uint类型都可以作为元素
// T是类型形参
// int | float32 | uint 是类型约束,多种约束用 | 隔开,代表T的类型实参不能是约束之外的
// [T int | float32 | uint] 整体构成了类型参数列表
// 这一种类型定义中带有类型形参的类型,被称为泛型类型
type cSlice[T int | float32 | uint] []T

// 类型形参的数量可以不止T这么一个,如下,定义了2个类型形参以及他们对应的约束
type MyMap[K int | string, V float32 | float64] map[K]V

// 所有的类型定义都可以使用类型形参,包括接口或结构体
type MyStruct[T int | string] struct {
	Name string
	Data T
}

// 泛型接口
type IPrinter[T int | float32 | string] interface {
	Print(data T)
}

// 泛型channel
type MyChan[T int | float32 | string] chan T

// 类型形参的互相套用
type WowStruct[T int | float32, S []T] struct { // var c WowStruct[int, []int] 实例化
	Data     S
	MaxValue T
	MinValue T
}

/*几种错误*/
// 1.基础类型不能只有类型形参
// type CommonType[T int | float32] T
type Wow[T int | string] int

//var c Wow[string] = "hello" // 编译错误，因为"hello"不能赋值给底层类型int
//这里虽然使用了类型形参，但因为类型定义是 type Wow[T int|string] int ，
//所以无论传入什么类型实参，实例化后的新类型的底层类型都是 int

// 2.类型约束的语法
// type NewType [T * int][]T //会被认为是T乘int,而不是指针
type NewType[T interface{ *int | *uint | float32 }] []T //使用空接口将指针类型包裹科技解决

/*泛型类型的嵌套*/
//和普通类型一样,泛型类型也可以通过组合嵌套组成更为复杂的结构
type Slice[T int | string | float32 | float64] []T

//type UintSlice[T uint | uint8]Slice[T] //错误,泛型类型Slice[T]的类型约束中不包含uint, uint8

type FloatSlice[T float32 | float64] Slice[T] //基于一个泛型类型创建新的泛型类型,新的泛型类型的类型约束必须在旧类型约束之中

type WowMap[T int | string] map[string]Slice[T]

//注意:匿名结构体和匿名函数无法使用泛型(因此在表驱动的单元测试中会比较麻烦)
//匿名函数不能在函数签名中定义自己的类型形参,但是可以使用外部已经定义好的类型形参

/*单独的泛型类型用处不大,和泛型接收者结合使用*/

type MySlice[T int | float32] []T

func (s MySlice[T]) Sum() T { //此方法就能够同时作用于int切片和float32切片的求和,泛型将特别有利于通用数据结构的实现
	var sum T
	for _, value := range s {
		sum += value
	}
	return sum
}

//Go中也不支持泛型方法
//但是因为receiver支持泛型， 所以如果想在方法中使用泛型的话，目前唯一的办法就是曲线救国，
//迂回地通过receiver使用类型形参,其实方法就是特殊的函数

//type A struct {
//}
//
//// 不支持泛型方法
//func (receiver A) Add[T int | float32 | float64](a T, b T) T {
//	return a + b
//}

type A[T int | float32 | float64] struct {
}

// 方法可以使用类型定义中的形参 T
func (receiver A[T]) Add(a T, b T) T {
	return a + b
}

func main() {
	//泛型类型不能直接拿来使用，必须传入类型实参(Type argument) 将其确定为具体的类型之后才可使用。
	//而传入类型实参确定具体类型的操作被称为 实例化(Instantiations)
	c := make(cSlice[int], 0)  //必须通过中括号传入 类型实参
	d := make(cSlice[uint], 0) //可以视为 type cSlice[uint] []uint
	//e:=make(cSlice[string],0)//不能传入在类型约束之外的类型
	fmt.Printf("%T\n%T\n", c, d)

	mMap := make(MyMap[int, float32], 0)
	fmt.Printf("%T\n", mMap)
	//mMap[2.1] = 2.2 //报错,K的类型实参传入的是int

	//调用泛型函数
	fmt.Println(AddGeneric[int](2, 3))
	fmt.Println(AddGeneric("dsa", "3")) //省略类型实参,可以自动推导类型

}
