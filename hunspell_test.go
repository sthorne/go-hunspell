package hunspell

import (

  "testing"

  )


func TestSuggest(t *testing.T) {

  h := Hunspell("/usr/share/hunspell/en_US.aff", "/usr/share/hunspell/en_US.dic")

  words := []string{
    "aplez",
    "yaho",
  }

  //
  // TEST

  for _, word := range words {

    y := Spell(t, word, h)

    if y {
      t.Error(`"` + word + `" is now a valid word and should not be`)
    }

    s := h.Suggest(word)

    for _, v := range s {
      t.Log("suggestion: " + v)
    }

  }

}

func TestSpell(t *testing.T) {

  h := Hunspell("/usr/share/hunspell/en_US.aff", "/usr/share/hunspell/en_US.dic")

  y := Spell(t, "xerocole", h)

  if y {
    t.Error(`"xerocole is now a valid word and should not be in the dictionary"`)
  }

  t.Log(`Adding "xerocole" to Dictionary`)

  b := h.Add("xerocole")

  if !b {
    t.Error(`Failed to Add "xerocole" to dictionary`)
  }

  y = Spell(t, "xerocole", h)

  if !y {
    t.Error(`"xerocole" is not a valid word after being added to dictionary`)
  }

}

func Spell(t *testing.T, word string, h *Hunhandle) bool {
  
  t.Log(`Testing "` + word + `"`)

  return h.Spell(word)

}