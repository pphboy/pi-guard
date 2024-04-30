package admin

import (
	"testing"

	"github.com/brianvoe/sjwt"
)

func TestJwt(t *testing.T) {
	c := sjwt.New()
	c.Set("a", "b")
	c.Generate()
}
