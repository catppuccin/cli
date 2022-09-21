package utils

import (
  "testing"
  "os"
)

func TestHandleDir(t *testing.T) {
  os.Setenv("CATPPUCCIN", "pog")
  r1 := HandleDir("this/shouldnt/change")
  r2 := HandleDir("$CATPPUCCIN/is/pog")
  r3 := HandleDir("catppuccin/is/$CATPPUCCIN/nice")
  if r1 != "this/shouldnt/change" {
    t.Log("Unchanging directory changed.")
    t.Fail()
  }
  if r2 != "pog/is/pog" {
    t.Log("Catppuccin wasn't pog, which is obviously untrue.")
  }
  if r3 != "catppuccin/is/pog/nice" {
    t.Log("Catppuccin in the middle? Still pog.")
  }
}
