package main

import "fmt"

func main() {
	ids := []int{6}
	newIds := []int{1,2,3,4}

	ids = append(ids, newIds...)

	fmt.Println(ids[len(ids)-1])
	fmt.Println(ids)
}
