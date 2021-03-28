package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const (
	TextSeparator = " "
)

var (
	RegexUppercaseCharacter         = regexp.MustCompile(`[A-ZА-Я]`)
	RegexPunctuationCharactersStart = regexp.MustCompile(`^[[:punct:]]`)
	RegexPunctuationCharactersEnd   = regexp.MustCompile(`[[:punct:]]$`)
	RegexClearText                  = regexp.MustCompile(`(\n|\r)+|\t+`)
)

type (
	WordMutation  func(string) string
	WordCountMap  map[string]int
	WordCountPair struct {
		Word  string
		Count int
	}
	WordCountList []WordCountPair
)

func Top10(text string, wordMutations ...WordMutation) []string {
	if text == "" {
		return make([]string, 0)
	}

	text = clearText(text)

	words := strings.Split(text, TextSeparator)

	m := make(WordCountMap)

	for _, w := range words {
		for _, mutation := range wordMutations {
			w = mutation(w)
		}

		if w == "" {
			continue
		}

		m[w]++
	}

	l := toPairList(m)

	l = orderByCountThenByWord(l)

	l = l[:minInt(len(l), 10)]

	top10 := make([]string, 0, len(l))

	for _, p := range l {
		top10 = append(top10, p.Word)
	}

	return top10
}

func IgnoreCharCase(w string) string {
	return RegexUppercaseCharacter.ReplaceAllStringFunc(w, strings.ToLower)
}

func IgnorePunctuationCharacters(w string) string {
	w = RegexPunctuationCharactersStart.ReplaceAllStringFunc(w, func(m string) string {
		return ""
	})

	return RegexPunctuationCharactersEnd.ReplaceAllStringFunc(w, func(m string) string {
		return ""
	})
}

func orderByCountThenByWord(l WordCountList) WordCountList {
	sort.Slice(l, func(i, j int) bool {
		switch {
		case l[i].Count > l[j].Count:
			return true
		case l[i].Count < l[j].Count:
			return false
		default:
			return l[i].Word < l[j].Word
		}
	})

	return l
}

func toPairList(m WordCountMap) WordCountList {
	l := make(WordCountList, len(m))

	i := 0
	for k, v := range m {
		l[i] = WordCountPair{k, v}
		i++
	}

	return l
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clearText(text string) string {
	return RegexClearText.ReplaceAllStringFunc(text, func(m string) string {
		if strings.Contains(m, "\t") {
			return " "
		}

		return ""
	})
}
