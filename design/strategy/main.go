package main

import "fmt"

/*
策略模式经常会和模板模式配合使用,注意区分他们的区别,模板模式是针对一个流程共性梳理出的固定的执行步骤,具体的
执行方式由子类来实现
而策略模式是侧重于让完成某个任务的具体方式可以相互切换
两个模式的解耦维度不同,策略模式在抽象方法的实现里,经常会用到模板模式
如,多种支付策略实际他们的流程都是类似的,都是查询余额 确认金额,扣款等,不同的策略实现中,就可以将这个支付过程
抽象为一个模板,即使是不同的策略,也会出现流程一致的情况


*/

func main() {
	wxPay := &WxPay{}
	px := NewPayCtx(wxPay) //设置微信支付策略,此处设置哪一种策略可以增加逻辑判断从而实现动态的策略选择
	px.Pay()

}

// PayBehavior 策略接口
type PayBehavior interface {
	OrderPay(ps *PayCtx)
}

type WxPay struct {
}

func (w WxPay) OrderPay(px *PayCtx) {
	fmt.Printf("Wx支付加工支付请求 %v\n", px.payParams)
	fmt.Println("正在使用Wx支付进行支付")
}

// 三方支付
type ThirdPay struct{}

func (*ThirdPay) OrderPay(px *PayCtx) {
	fmt.Printf("三方支付加工支付请求 %v\n", px.payParams)
	fmt.Println("正在使用三方支付进行支付")
}

// PayCtx 代表策略执行所需的上下文,由其维护策略及其参数
type PayCtx struct {
	// 提供支付能力的接口实现
	payBehavior PayBehavior
	// 支付参数
	payParams map[string]interface{}
}

func (px *PayCtx) setPayBehavior(p PayBehavior) {
	px.payBehavior = p //设置策略
}

func (px *PayCtx) Pay() {
	px.payBehavior.OrderPay(px)
}

func NewPayCtx(p PayBehavior) *PayCtx {
	params := map[string]interface{}{
		"appId": "234fdfdngj4",
		"mchId": 123456,
	}
	return &PayCtx{
		payBehavior: p,
		payParams:   params,
	}
}
