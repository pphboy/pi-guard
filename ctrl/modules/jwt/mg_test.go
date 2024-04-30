package jwt

import (
	"go-ctrl/models"
	"testing"
)

func TestMG(t *testing.T) {
	g := NewMyG[models.PiManager]()
	s := g.GetJwtToken("aaa", models.PiManager{
		ManagerUsername: "username",
	})
	t.Log(s)

	asdf, err := g.ParseToken(s)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(asdf)
	t.Log(asdf["ManagerUsername"])
}
