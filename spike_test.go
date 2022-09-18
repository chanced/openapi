package openapi_test

import (
	"fmt"
	"testing"

	"github.com/chanced/uri"
)

type Spike1 struct {
	Key string
	Val string
}

func (s *Spike1) K() string { return s.Key }
func (s *Spike1) V() string { return s.Val }

type Spike2 struct {
	Key2 string
	Val2 string
}

func (s *Spike2) K() string { return s.Key2 }
func (s *Spike2) V() string { return s.Val2 }

type S interface {
	K() string
	V() string
}

type C struct {
	One *Spike1
	Two *Spike2
}

type X struct {
	v interface{}
}
type V struct {
	S *Spike2
}

func (c *C) R1() interface{} { return &c.One }
func (c *C) R2() interface{} { return &c.Two }

func TestSpike(t *testing.T) {
	u, err := uri.Parse("#asd")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("u:", u)
	// ones := []S{
	// 	&Spike1{Key: "1one", Val: "1one"},
	// 	&Spike1{Key: "1two", Val: "1two"},
	// 	&Spike1{Key: "1three", Val: "1three"},
	// }
	// twos := []S{
	// 	&Spike2{Key2: "2one", Val2: "2one"},
	// 	&Spike2{Key2: "2two", Val2: "2two"},
	// 	&Spike2{Key2: "2three", Val2: "2three"},
	// }
	// _ = ones
	// _ = twos
	// c := &C{}

	// r1 := reflect.ValueOf(c.R1())
	// v1 := reflect.ValueOf(ones[0])
	// r1.Elem().Set(v1)
	// fmt.Println(c.One)
	// fmt.Println(r1.Type())

	// v2 := reflect.ValueOf(twos[0])

	// fmt.Println("v2 assignable to r1.Type.Elem()", v2.Type().AssignableTo(r1.Type().Elem()))
	// fmt.Println("v1 assignable to r1.Type.Elem()", v1.Type().AssignableTo(r1.Type().Elem()))

	// v := &V{}

	// x := &X{v: &v.S}
	// xr := reflect.ValueOf(x.v)

	// fmt.Println(xr.CanAddr())

	// fmt.Printf("V.S: %+v\n", v.S)
}
