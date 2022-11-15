package main

/*
SingleFlight 的作用是将并发请求合并成一个请求，以减少对下层服务的压力
注意区分和Once的区别,sync.Once 主要是用在单次初始化场景中，而 SingleFlight 主要用在合并并发请求的场景中，尤其是缓存场景


*/

func main() {

}
