package ai

import (
	"strings"
)

// Common rhyming endings for basic rhyme matching
var rhymeEndings = map[string][]string{
	"ove":  {"love", "dove", "above", "shove", "glove"},
	"ay":   {"day", "way", "say", "stay", "play", "away", "today", "okay"},
	"ight": {"night", "light", "right", "fight", "sight", "bright", "flight", "might", "tight", "white"},
	"ound": {"sound", "found", "ground", "around", "bound", "round", "pound"},
	"ine":  {"line", "mine", "fine", "shine", "wine", "sign", "divine", "define"},
	"art":  {"heart", "start", "part", "apart", "smart", "chart", "dart"},
	"eel":  {"feel", "real", "deal", "heal", "steal", "reveal", "wheel"},
	"ain":  {"rain", "pain", "chain", "brain", "train", "gain", "main", "remain", "explain"},
	"ong":  {"song", "long", "wrong", "strong", "along", "belong"},
	"ame":  {"name", "same", "game", "flame", "blame", "shame", "frame", "came"},
	"one":  {"alone", "stone", "phone", "bone", "tone", "zone", "known", "own", "grown"},
	"ear":  {"hear", "near", "fear", "clear", "dear", "year", "tear", "appear"},
	"oo":   {"you", "true", "blue", "through", "new", "knew", "view"},
	"ow":   {"know", "show", "grow", "flow", "slow", "go", "low", "below"},
	"all":  {"call", "fall", "wall", "small", "tall", "ball", "hall"},
	"ing":  {"thing", "sing", "bring", "ring", "spring", "king", "string", "wing"},
	"eak":  {"speak", "weak", "seek", "peak", "leak", "sneak", "creek"},
	"ive":  {"live", "give", "drive", "five", "alive", "survive", "arrive"},
	"ace":  {"place", "face", "space", "grace", "race", "trace", "embrace"},
	"ore":  {"more", "before", "store", "floor", "door", "shore", "explore", "core"},
	"ide":  {"side", "hide", "ride", "guide", "wide", "pride", "inside", "decide"},
	"eam":  {"dream", "stream", "team", "beam", "seem", "cream", "scheme"},
	"end":  {"end", "friend", "send", "spend", "blend", "bend", "trend", "pretend"},
	"ake":  {"break", "make", "take", "wake", "shake", "mistake", "lake"},
	"ife":  {"life", "wife", "knife", "strife"},
	"ope":  {"hope", "rope", "scope", "slope"},
	"ime":  {"time", "crime", "climb", "prime", "rhyme", "sublime"},
	"urn":  {"turn", "burn", "learn", "return", "concern", "yearn"},
	"old":  {"old", "cold", "gold", "hold", "told", "bold", "sold", "fold"},
	"ate":  {"late", "fate", "hate", "wait", "great", "state", "create", "date"},
}

// FindBasicRhymes finds words that rhyme with the given word using basic pattern matching
func FindBasicRhymes(word string) []string {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return []string{}
	}

	var rhymes []string
	
	// Check if the word matches any ending patterns
	for ending, words := range rhymeEndings {
		if strings.HasSuffix(word, ending) {
			for _, w := range words {
				if w != word {
					rhymes = append(rhymes, w)
				}
			}
		}
	}
	
	// If no matches found, try to find rhymes by ending
	if len(rhymes) == 0 {
		// Extract the ending (last 2-3 characters)
		for endLen := 3; endLen >= 2 && len(rhymes) == 0; endLen-- {
			if len(word) >= endLen {
				ending := word[len(word)-endLen:]
				for _, words := range rhymeEndings {
					for _, w := range words {
						if strings.HasSuffix(w, ending) && w != word {
							rhymes = append(rhymes, w)
						}
					}
				}
			}
		}
	}
	
	// Remove duplicates
	seen := make(map[string]bool)
	unique := []string{}
	for _, r := range rhymes {
		if !seen[r] {
			seen[r] = true
			unique = append(unique, r)
		}
	}
	
	return unique
}
