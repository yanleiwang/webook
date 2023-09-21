package webook

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMy(t *testing.T) {
	type Foo struct {
		Name  string
		age   int
		Phone int
	}

	foo := Foo{}
	s, _ := json.Marshal(foo)

	fmt.Println(string(s))

}
