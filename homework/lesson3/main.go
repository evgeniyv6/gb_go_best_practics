package main

import "fmt"

func main() {
	fmt.Println("hello from le3")


}

func foo() (map[string]int, error) {
	return map[string]int{"a":1, "b":2}, nil
}