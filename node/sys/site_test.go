package sys

import "testing"

func TestKillPro(t *testing.T) {
	err := KillProcessByPort("8081")
	t.Log(err)
}
