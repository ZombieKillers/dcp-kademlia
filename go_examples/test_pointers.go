package main

import "fmt"

type Contact struct {
	id int
}

type ContactStruct struct{
	node *Contact
	value int
}

func NewContact(id int) (*Contact){
	return &Contact{id: id}
}


func createContactStructs(l *[]*ContactStruct,  contacts []*Contact){

	for i, v := range contacts {
		*l = append(*l ,&ContactStruct{v, i})
	}

}

func main() {
	contact1 := NewContact(3)
	contact2 := NewContact(4)

	contacts := make([]*Contact, 2)
	contacts[0] = contact1
	contacts[1] = contact2
	contactStructs := make([]*ContactStruct, 0)
	createContactStructs(&contactStructs, contacts)
	for _,v := range contactStructs {
		fmt.Println("Node:", v.node.id)
		fmt.Println("Value:", v.value)
	}
}
