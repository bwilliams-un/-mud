package character

import (
	grammar "github.com/bwilliams-un/mud/grammar"
)

// Character object
type Character struct {
	name     string
	pronouns grammar.Pronouns
}
