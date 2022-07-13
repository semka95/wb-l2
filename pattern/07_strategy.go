package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

type Context struct {
	strategy func()
}

func (c *Context) Execute() {
	c.strategy()
}

func (c *Context) SetStrategy(strategy func()) {
	c.strategy = strategy
}

// func main() {
// 	concreteStrategyA := func() {
// 		fmt.Println("concreteStrategyA()")
// 	}
// 	concreteStrategyB := func() {
// 		fmt.Println("concreteStrategyB()")
// 	}
// 	context := Context{concreteStrategyA}
// 	context.Execute()
// 	context.SetStrategy(concreteStrategyB)
// 	context.Execute()
// }
