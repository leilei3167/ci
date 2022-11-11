package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	d := &DepositBusinessHandler{userVip: false}

	e.ExecuteBankBusiness()         // 执行模板
}

// 当要做一件事儿的时候，这件事儿的流程和步骤是固定好的，
// 但是每一个步骤的具体实现方式是不一定的。这个时候就可以使用模板模式.

// BankBusinessHandler 代表去银行办理业务的流程.
type BankBusinessHandler interface {
	// 排队拿号
	TakeRowNumber()
	// 等位
	WaitInHead()
	// 处理具体业务
	HandleBusiness()
	// 对服务作出评价
	Commentate()
	// 钩子方法，
	// 用于在流程里判断是不是VIP， 实现类似VIP不用等的需求
	CheckVipIdentity() bool
}

// BankBusinessExecutor 子类继承,具体的模板执行交由此结构.
type BankBusinessExecutor struct {
	handler BankBusinessHandler
}

func NewBankBusinessExecutor(businessHandler BankBusinessHandler) *BankBusinessExecutor {
	return &BankBusinessExecutor{handler: businessHandler}
}

// 模板方法，处理银行业务的模板,如判断是否是vip权限,做出不同的业务流程,子类实现流程的逻辑.
func (b *BankBusinessExecutor) ExecuteBankBusiness() {
	// 适用于与客户端单次交互的流程
	// 如果需要与客户端多次交互才能完成整个流程，
	// 每次交互的操作去调对应模板里定义的方法就好，并不需要一个调用所有方法的模板方法
	b.handler.TakeRowNumber()
	if !b.handler.CheckVipIdentity() {
		b.handler.WaitInHead()
	}
	b.handler.HandleBusiness()
	b.handler.Commentate()
}



// DepositBusinessHandler 实现一个存款的业务.
type DepositBusinessHandler struct {
	*DefaultBusinessHandler
	userVip bool
}

func (*DepositBusinessHandler) TakeRowNumber() {
	fmt.Println("请拿好您的取件码：" + strconv.Itoa(rand.Intn(100)) +
		" ，注意排队情况，过号后顺延三个安排")
}

func (dh *DepositBusinessHandler) WaitInHead() {
	fmt.Println("排队等号中...")
	time.Sleep(5 * time.Second)
	fmt.Println("请去窗口xxx...")
}

func (*DepositBusinessHandler) HandleBusiness() {
	fmt.Println("账户存储很多万人民币...")
}

func (dh *DepositBusinessHandler) CheckVipIdentity() bool {
	return dh.userVip
}

func (*DepositBusinessHandler) Commentate() {
	fmt.Println("请对我的服务作出评价，满意请按0，满意请按0，(～￣▽￣)～")
}


//所以就可以把它们放在抽象类中可以进一步减少代码的重复率(Default)。

// DefaultBusinessHandler 并没有完整实现BankBusinessHandler接口,这么做是为了这个类型只能用于被实现类包装，让 Go 语言的类型检查能够帮我们强制要求，
// 必须用存款或者取款这样子类去实现HandleBusiness方法，整个银行办理业务的流程的程序才能运行起来。
type DefaultBusinessHandler struct{}

func (*DefaultBusinessHandler) TakeRowNumber() {
	fmt.Println("请拿好您的取件码：" + strconv.Itoa(rand.Intn(100)) +
		" ，注意排队情况，过号后顺延三个安排")
}

func (dbh *DefaultBusinessHandler) WaitInHead() {
	fmt.Println("排队等号中...")
	time.Sleep(5 * time.Second)
	fmt.Println("请去窗口xxx...")
}

func (*DefaultBusinessHandler) Commentate() {
	fmt.Println("请对我的服务作出评价，满意请按0，满意请按0，(～￣▽￣)～")
}

func (*DefaultBusinessHandler) CheckVipIdentity() bool {
	// 留给具体实现类实现
	return false
}
