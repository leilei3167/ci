package main

import "reflect"

/*基于泛型的队列.*/

// Queue 泛型队列的定义,可储存任意类型的数据.
type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Put(value T) {
	//value.(int) //不能对类型形参使用类型断言或类型选择
	//但是可以通过反射来得知T的类型,进行不同的操作
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int:
	case reflect.String:
	default:
	}
	//不建议这样使用泛型,你为了避免使用反射而选择了泛型，结果到头来又为了一些功能在在泛型中使用反射
	q.elements = append(q.elements, value)
}

func (q *Queue[T]) Pop() (T, bool) {
	var value T
	if len(q.elements) == 0 {
		return value, true
	}
	value = q.elements[0]
	q.elements = q.elements[1:]
	return value, len(q.elements) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

//var q2 Queue[string]
