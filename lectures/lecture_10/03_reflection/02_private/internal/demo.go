package internal

import "fmt"

type PublicStruct struct {
	privateField string
}

func (p *PublicStruct) String() string {
	return fmt.Sprintf("%s", p.privateField)
}

func NewPublicStruct() *PublicStruct {
	return &PublicStruct{
		privateField: "123",
	}
}
