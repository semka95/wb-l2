package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

import (
	"fmt"
)

type Creator struct {
	factory factory
}

func (c *Creator) Operation() {
	product := c.factory.factoryMethod()
	product.method()
}

type factory interface {
	factoryMethod() Product
}

type ConcreteCreator struct{}

func (c *ConcreteCreator) factoryMethod() Product {
	return new(ConcreteProduct)
}

type Product interface {
	method()
}

type ConcreteProduct struct{}

func (p *ConcreteProduct) method() {
	fmt.Println("ConcreteProduct.method()")
}

// func main() {
// 	creator := Creator{new(ConcreteCreator)}
// 	creator.Operation()
// }
