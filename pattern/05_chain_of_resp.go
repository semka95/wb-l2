package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

import (
	"fmt"
)

type Handler interface {
	Request(flag bool)
}

type ConcreteHandlerA struct {
	next Handler
}

func (h *ConcreteHandlerA) Request(flag bool) {
	fmt.Println("ConcreteHandlerA.Request()")
	if flag {
		h.next.Request(flag)
	}
}

type ConcreteHandlerB struct {
	next Handler
}

func (h *ConcreteHandlerB) Request(flag bool) {
	fmt.Println("ConcreteHandlerB.Request()")
}

// func main() {
// 	handlerA := &ConcreteHandlerA{new(ConcreteHandlerB)}
// 	handlerA.Request(true)
// }
