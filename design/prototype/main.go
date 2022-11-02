package main

import (
	"encoding/json"
	"time"
)

/*原型模式:
通过复制、拷贝或者叫克隆已有对象的方式来创建新对象的设计模式叫做原型模式，被拷贝的对象也被称作原型对象

原型对象按照惯例，会暴露出一个 Clone 方法，给外部调用者一个机会来从自己这里“零成本”的克隆出一个新对象。
这里的“零成本”说的是，调用者啥都不用干，干等着，原型对象在 Clone 方法里自己克隆出自己，给到调用者，所以按照这个约定所有原型对象都要实现一个 Clone 方法

至于原型对象克隆自己的时候用的是深拷贝还是浅拷贝？可以先理解成是都用深拷贝，等完全掌握这种思想后，可以再根据实际情况，
比如为了节省空间、以及减少编写克隆方法的复杂度时可以两者综合使用。

比如全局配置对象这种也可以当成原型对象(典型的如gorm.DB)，如果不想让程序在运行时修改初始化好的原型对象，导致影响其他线程的程序执行的时候，
也可以用原型模式快速拷贝出一份，再在副本上做运行时自定义修改。

使用场景:
	当对象的创建成本比较大，并且同一个类的不同对象间差别不大时（大部分属性值相同），如果对象的属性值需要经过复杂的计算、排序，或者需要从网络、DB等这些慢IO中获取、亦或者或者属性值拥有很深的层级，
这时就是原型模式发挥作用的地方了。因为对象在内存中复制自己远比每次创建对象时重走一遍上面说的操作要来高效的多

*/

// Keyword 搜索关键字
type Keyword struct {
	word      string
	visit     int
	UpdatedAt *time.Time
}

// Clone 这里使用序列化与反序列化的方式深拷贝
func (k *Keyword) Clone() *Keyword {
	var newKeyword Keyword
	b, _ := json.Marshal(k)
	json.Unmarshal(b, &newKeyword)
	return &newKeyword
}

// Keywords 关键字 map
type Keywords map[string]*Keyword

// Clone 复制一个新的 keywords
// updatedWords: 需要更新的关键词列表，由于从数据库中获取数据常常是数组的方式
func (words Keywords) Clone(updatedWords []*Keyword) Keywords {
	newKeywords := Keywords{}

	for k, v := range words {
		// 这里是浅拷贝，直接拷贝了地址
		newKeywords[k] = v
	}

	// 替换掉需要更新的字段，这里用的是深拷贝
	for _, word := range updatedWords {
		newKeywords[word.word] = word.Clone() //深拷贝
	}

	return newKeywords
}

func main() {

}
