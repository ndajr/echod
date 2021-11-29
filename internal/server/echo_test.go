package server

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestEcho(t *testing.T) {
	g := Goblin(t)

	g.Describe("Echo", func() {
		g.It("valid", func() {
			e, _ := Echo([]byte(`{"username":"xyz","upload":"xyz"}`))
			g.Assert(e).Equal(echoBody{
				"username": "xyz",
				"upload":   "xyz",
				"echoed":   true,
			})
		})

		g.It("echoed as false", func() {
			e, _ := Echo([]byte(`{"username":"xyz","upload":"xyz","echoed":false}`))
			g.Assert(e).Equal(echoBody{
				"username": "xyz",
				"upload":   "xyz",
				"echoed":   true,
			})
		})

		g.It("echoed as true", func() {
			_, err := Echo([]byte(`{"username":"xyz","upload":"xyz","echoed":true}`))
			g.Assert(err.Error()).Equal("echoed already set")
		})

		g.It("invalid json", func() {
			_, err := Echo([]byte(""))
			g.Assert(err.Error()).Equal("invalid json")
		})

		g.It("empty valid json", func() {
			e, _ := Echo([]byte(`{}`))
			g.Assert(e).Equal(echoBody{"echoed": true})
		})
	})

}

func ExampleEcho() {
	e, _ := Echo([]byte(`{"foo": "bar"}`))
	fmt.Println(e)
	// Output: map[echoed:true foo:bar]
}

func ExampleEcho_echoedAlreadySet() {
	_, err := Echo([]byte(`{"foo": "bar", "echoed": true}`))
	fmt.Println(err)
	// Output: echoed already set
}

func ExampleEcho_invalidJSON() {
	_, err := Echo([]byte(""))
	fmt.Println(err)
	// Output: invalid json
}
