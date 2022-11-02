package main

import "fmt"

/*一.简单工厂:缺点就是生产的产品如果要扩充,则必须修改工厂函数,增加case
产品过多就会臃肿*/

// Printer 代表产品接口
type Printer interface {
	Print() string
}

func NewPrinter(lan string) Printer {
	switch lan {
	case "ch":
		return new(chPrinter)
	case "en":
		return new(enPrinter)
	default:
		return new(chPrinter)
	}
}

// 具体的产品实现
type chPrinter struct {
}

func (c chPrinter) Print() string {
	return "你好"
}

type enPrinter struct {
}

func (c enPrinter) Print() string {
	return "hello"
}

/*二.工厂方法: 多态性工厂,指的是定义一个创建对象的接口，但由实现这个接口的工厂类来决定实例化哪个产品类，
工厂方法把类的实例化推迟到子类中进行。不再由单一的工厂类生产产品，而是由工厂类的子类实现具体产品的创建。因此，当增加一个产品时，只需增加一个相应的工厂类的子类,
以解决简单工厂生产太多产品时导致其内部代码臃肿（switch … case分支过多）的问题
每个工厂的实现,只能创建一个具体的产品,实现多种的产品,需要拓展对应的工厂实现
*/
//Go中没有继承，所以这里说的工厂子类，其实是直接实现工厂接口的具体工厂类

// OperatorFactory 工厂接口,即能够生产指定类别的产品,才能实现此接口
type OperatorFactory interface {
	Create() MathOperator
}

// MathOperator 代表着产品接口,产品即是具有同类行为的
type MathOperator interface {
	SetOperandA(int)
	SetOperandB(int)
	ComputeResult() int
}

//假设要生产2类计算器,加法和乘法计算器,实现两个工厂类

// PlusOperatorFactory 加法计算器的生产工厂
type PlusOperatorFactory struct {
}

func (p PlusOperatorFactory) Create() MathOperator {
	return &PlusOperator{BaseOperator: &BaseOperator{}} //此处实际开发中会填充复杂的初始化逻辑
}

// PlusOperator 加法计算器具体的产品
type PlusOperator struct {
	*BaseOperator //匿名嵌套基类,公用公共方法
}

func (p PlusOperator) ComputeResult() int {
	return p.operandA + p.operandB
}

// BaseOperator 因为有些公共的方法在不同的产品中很可能是重复的,因此将公共方法抽象为基类,多种实现直接继承即可
type BaseOperator struct {
	operandA, operandB int
}

func (o *BaseOperator) SetOperandA(operand int) {
	o.operandA = operand
}

func (o *BaseOperator) SetOperandB(operand int) {
	o.operandB = operand
}

// MultiOperatorFactory 乘法计算器工厂
type MultiOperatorFactory struct{}

func (mf *MultiOperatorFactory) Create() MathOperator {
	return &MultiOperator{
		BaseOperator: &BaseOperator{}, //注意 实际是复杂初始化逻辑
	}
}

// MultiOperator 实际的产品类--乘法运算器
type MultiOperator struct {
	*BaseOperator
}

func (m *MultiOperator) ComputeResult() int {
	return m.operandA * m.operandB
}

// 需要再拓展生产的产品类别时,只需要再增加新的工厂方法实现即可(增加生产线)
// 灵活性增强，对于新产品的创建，只需多写一个相应的工厂类
// 典型的解耦框架。高层模块只需要知道产品的抽象类，无须关心其他实现类，满足迪米特法则、依赖倒置原则和里氏替换原则。
// 缺点: 只能生产一种具有相同行为的产品，此弊端可使用抽象工厂模式解决
// 工厂方法可用于分层架构,如控制层依赖服务层的工厂方法,而服务层又依赖仓库层的工厂方法,每一层都依赖于下一层的接口
func main() {
	//一.工厂方法使用
	//1.初始化工厂的实现(依赖注入)
	var factory OperatorFactory
	factory = &PlusOperatorFactory{}
	plus := factory.Create()
	plus.SetOperandA(100)
	plus.SetOperandB(1000)
	fmt.Println("plus:", plus.ComputeResult())

	//乘法计算器生产
	factory = &MultiOperatorFactory{}
	m := factory.Create()
	m.SetOperandB(100)
	m.SetOperandA(1000)

	fmt.Println("mu:", m.ComputeResult())

	//二.抽象工厂使用

}

/*三.抽象工厂:用于创建一系列相关的或者相互依赖的对象
注意理解和工厂方法的区别:
抽象工厂模式与工厂方法模式最大的区别在于，工厂方法模式针对的是一个产品等级结构，
而抽象工厂模式则需要面对多个产品等级结构，一个工厂等级结构可以负责多个不同产品等级结构中的产品对象的创建 。
工厂方法模式只有一个抽象产品类，而抽象工厂模式有多个。工厂方法模式的具体工厂类只能创建一个具体产品类的实例，而抽象工厂模式可以创建多个

缺点是规定了所有可能被创建的产品集合，产品族中如果扩展新的产品(增加工厂接口的方法列表)，则需要修改抽象工厂的接口
*/
