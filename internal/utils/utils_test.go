package utils

import (
	"os"
	"testing"
)

func TestHandleDir(t *testing.T) {
	os.Setenv("CATPPUCCIN", "pog")
	os.Setenv("NEOVIM", "chad/")
	r1 := HandleDir("this/shouldnt/change")
	r2 := HandleDir("$CATPPUCCIN/is/pog")
	r3 := HandleDir("catppuccin/is/$CATPPUCCIN/nice")
	r4 := HandleDir("neovim/is/a/$NEOVIM/text/editor")
	if r1 != "this/shouldnt/change" {
		t.Log("Unchanging directory changed.")
		t.Fail()
	}
	if r2 != "pog/is/pog" {
		t.Log("Catppuccin wasn't pog, which is obviously untrue.")
		t.Log(r2)
		t.Fail()
	}
	if r3 != "catppuccin/is/pog/nice" {
		t.Log("Catppuccin in the middle? Still pog.")
		t.Log(r3)
		t.Fail()
	}
	if r4 != "neovim/is/a/chad/text/editor" {
		t.Log("Catppuccin in the middle? Still pog.")
		t.Log(r4)
		t.Fail()
	}
}
