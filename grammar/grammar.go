package grammar

// Pronouns for grammar
type Pronouns struct {
	subject string
	object  string
}

// NeutralPronouns for grammar
func NeutralPronouns() Pronouns {
	return Pronouns{"they", "them"}
}

// MasculinPronouns for grammar
func MasculinPronouns() Pronouns {
	return Pronouns{"he", "him"}
}

// FemininePronouns for grammar
func FemininePronouns() Pronouns {
	return Pronouns{"she", "her"}
}
