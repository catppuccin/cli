package utils

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestHandleDir(t *testing.T) {
	assert := is.New(t)

	os.Setenv("CATPPUCCIN", "pog")
	os.Setenv("NEOVIM", "chad/")

	r1 := HandleDir("this/shouldnt/change")
	assert.True(r1 == "this/shouldnt/change") // Unchanging directory should not change

	r2 := HandleDir("$CATPPUCCIN/is/pog")
	assert.True(r2 == "pog/is/pog") // Catppuccin should always be pog

	r3 := HandleDir("catppuccin/is/$CATPPUCCIN/nice")
	assert.True(r3 == "catppuccin/is/pog/nice") // Catppuccin should always be pog, even in the middle

	r4 := HandleDir("neovim/is/a/$NEOVIM/text/editor")
	assert.True(r4 == "neovim/is/a/chad/text/editor") // Neovim should always be a chad editor
}
