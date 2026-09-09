package example_hello_world

import "fmt"

func ExampleHelloWorld(name string) {
	if name == "" {
		name = "World"
	}
	fmt.Println("Hello, " + name + "!")
}
