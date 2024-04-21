package dal

import "testing"

func TestDel(t *testing.T) {
	Del("speak.txt")
	if f := Get("speak.txt"); f.Hash != "" {
		t.Error("del error ")
	}
}
