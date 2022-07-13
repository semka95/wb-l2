package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

type Element interface {
	Accept(Visitor)
}

type ConcreteElementA struct{}

func (e *ConcreteElementA) Accept(visitor Visitor) {
	fmt.Println("ConcreteElementA.Accept()")
	visitor.VisitA(e)
}

type ConcreteElementB struct{}

func (e *ConcreteElementB) Accept(visitor Visitor) {
	fmt.Println("ConcreteElementB.Accept()")
	visitor.VisitB(e)
}

type Visitor interface {
	VisitA(*ConcreteElementA)
	VisitB(*ConcreteElementB)
}

type ConcreteVisitor struct{}

func (v *ConcreteVisitor) VisitA(element *ConcreteElementA) {
	fmt.Println("ConcreteVisitor.VisitA()")
}

func (v *ConcreteVisitor) VisitB(element *ConcreteElementB) {
	fmt.Println("ConcreteVisitor.VisitB()")
}
