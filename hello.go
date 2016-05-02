package main

import (
	"./nodes"
	"./server"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func sum(nums ...int) (total int) {
	fmt.Println(nums, " ")
	for _, num := range nums {
		total += num
	}
	return
}

func main() {

	fmt.Println("Map example")

	m := make(map[string]int)
	m["one"] = 1
	m["two"] = 2

	v1 := m["one"]
	fmt.Println("v1:", v1)

	kvs := map[string]string{"a": "avalue", "b": "bvalue"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	fmt.Println("Sum of the values: ")
	vals := []int{1, 2, 3}

	fmt.Println(sum(vals...))
	n1 := nodes.NewRandomNodeId()
	fmt.Println(n1.String())

	fmt.Println("Trying some RPC stuff")

	if err := server.StartServer(); err != nil {
		panic(err)
	}

}
