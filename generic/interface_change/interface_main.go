package main

/*Go引入泛型后对接口也作出了较大的变化*/

// 有时候使用泛型编程时，我们会书写长长的类型约束，如下：
// 一个可以容纳所有int,uint以及浮点类型的泛型切片
type Slice[T ~int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64] []T

// Go支持将类型约束单独拿出来定义到接口中，从而让代码更容易维护
type intUintFloat interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Slice1[T intUintFloat] []T //需要指定类型约束的时候直接使用接口 IntUintFloat 即可

// 接口和普通类型作为类型约束是可以通过 | 来组合,接口中也能够组合其他的接口
type Slice2[T intUintFloat | string] []T

/*指定底层的类型*/
var s1 Slice[int]

type MyInt int

var s2 Slice[MyInt] //错误,因为MyInt并非是int类型,需要使用 ~ 符合,在类型约束中带有~符合的指的是其底层的数据类型
// ~ 后面不能是接口,只能是基础数据类型

/*从方法集到类型集
我们可以把 ReaderWriter 接口看成代表了一个 类型的集合，
所有实现了 Read() Writer() 这两个方法的类型都在接口代表的类型集合当中

当满足以下条件时，我们可以说 类型 T 实现了接口 I ( type T implements interface I)：
	T 不是接口时：类型 T 是接口 I 代表的类型集中的一个成员 (T is an element of the type set of I)
	T 是接口时： T 接口代表的类型集是 I 代表的类型集的子集(Type set of T is a subset of the type set of I)

*/

// 在接口中,使用|连接代表的是并集,而一个接口有多行类型定义,代表的是交集
type A interface {
	~int | ~uint | ~float64
}

type B interface {
	~float64 | ~uint
}

type C interface { //代表的是取交集,等同于~float64|~uint
	A
	B
}

// 如果多行的写法没有交集时,将被视为空集(没有交叉部分),没有一种类型是属于空集的
// 这样写没有任何意义
type D interface {
	~int
	~float64
}

/*空接口和any*/
//空接口代表所有的类型的集合
//虽然空接口内没有写入任何的类型，但它代表的是所有类型的集合，而非一个 空集
//类型约束中指定 空接口 的意思是指定了一个包含所有类型的类型集，并不是类型约束限定了只能使用 空接口 来做类型形参
//any是空接口的别名,方便使用

/*comparable和ordered*/
//对于一些数据类型，我们需要在类型约束中限制只接受能 != 和 == 对比的类型，如map：
//type MyMap[KEY any, VALUE any] map[KEY]VALUE // 错误。因为 map 中键的类型必须是可进行 != 和 == 比较的类型

// 提供了comparable 的接口，它代表了所有可用 != 以及 == 对比的类型
type MyMap[KEY comparable, VALUE any] map[KEY]VALUE // 正确
//可比较不代表着可排序,可比较指的是能够使用==和!=进行操作的类型,并不代表该类型能够使用<><=>=,如结构体不可比大小

// 要使用类型集的概念去看待1.18之后的接口
// 接口类型 ReadWriter 代表了一个类型集合，所有以 string 或 []rune 为底层类型，
// 并且实现了 Read() Write() 这两个方法的类型都在 ReadWriter 代表的类型集当中
type ReadWriter interface {
	~string | ~[]rune

	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

// 要实现此接口,该类型的底层类型必须是string或者[]byte类型,且必须有Read和Write方法
type myString string

func (m myString) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (m myString) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

/*基本接口和一般接口
考虑到定义一个 ReadWriter 类型的接口变量，然后接口变量赋值的时候不光要考虑到方法的实现，
还必须考虑到具体底层类型,心智负担也太大了,同时也为了保证兼容性,1.18开始将接口分为了两种类型
*/

//基本接口:即接口定义中只有方法的接口,也就是1.18版本之前的接口,用法和之前完全一致,但也可以用于类型约束中

//一般接口:接口定义中不光有方法,还同时具有类型的接口,就像上述的ReadWriter
//一般接口类型不能用来定义变量，只能用于泛型的类型约束中

/*泛型接口
所有的类型定义中都可以使用类型形参,接口也不例外
*/

type DataProcessor[T any] interface { //因为引入了类型形参,所以这个接口属于泛型类型,在使用此接口时也必须指定类型实参
	Process(o T) (new T)
	Save(data T) error
}

// var s DataProcessor[string]

func main() {

}
