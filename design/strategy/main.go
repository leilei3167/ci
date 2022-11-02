package main

func main() {

}

type scanStrategy interface { //定义扫描策略
	Start() string
}

type operator struct {
	s scanStrategy
}

func (o operator) Start() string {
	return o.s.Start()
}

func (o operator) setStrategy(s scanStrategy) {
	o.s = s
}
