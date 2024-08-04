package main

import (
	"errors"
	"fmt"

	"github.com/chriswtodd/go-stack-patterns/stack"
)

func main() {
	var se = stack.NewStack[string]()
	var s1 = stack.NewStack[string]().Push("Hello").Push("World").Push("!!!")
	var s2 = stack.NewStack[string]().Push("Hello").Push("World").Push("!!!")

	fmt.Println("Empty stack and two stacks with 3 items")
	fmt.Println(se, s1, s2)

	fmt.Println("First item on the stack")
	fmt.Println(s1.TopOrElse(func() interface{} { return errors.New("should not throw") }))

	fmt.Println("Stack after TopOrElse call")
	fmt.Println(s1)

	fmt.Println("First item off the stack")
	fmt.Println(s1.PopOrElse(func() interface{} { return errors.New("should not throw") }))

	fmt.Println("Pop when empty")
	fmt.Println(stack.NewStack[string]().PopOrElse(func() interface{} { return errors.New("error the stack is empty") }))
}
