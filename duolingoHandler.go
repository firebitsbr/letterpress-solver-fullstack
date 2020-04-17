package solver

import (
	"regexp"
	"strconv"
	"strings"
)

// Trim duolingo challenge choices by 4
func trimDuolingoChoices(s string) string {
	var reps []string
	reChoices := regexp.MustCompile(`"choices":\[([^\]]+)\]`)
	reIndices := regexp.MustCompile(`"correctIndices":\[([^\]]+)\]`)
	matchChoices := reChoices.FindAllString(s, -1)
	matchIndices := reIndices.FindAllString(s, -1)
	for _, ms := range matchChoices {
		if !strings.Contains(ms, "\"text\"") {
			continue
		}
		ts := strings.Split(ms, ",") // tokens
		if len(ts) <= 4 {
			continue
		}

		reps = append(reps, ms) // all choices
		trimNum := 4
		if strings.Contains(ms, "https") {
			trimNum = 8
		}
		for i, t := range ts {
			ts[i] = strings.Replace(t, "\"}", " "+strconv.Itoa(i*4/trimNum)+"\"}", 1)
		}
		rs := strings.Join(ts[:len(ts)-trimNum], ",") + "]"
		reps = append(reps, rs) // all choices - 4
	}

	// Only tap the first!
	for _, cis := range matchIndices {
		if strings.Contains(cis, ",") {
			reps = append(reps, cis)
			reps = append(reps, "\"correctIndices\":[0]")
		}
	}
	replacer := strings.NewReplacer(reps...)
	news := replacer.Replace(s)
	return news
}

// Mark correct single choice
func markDuolingoSingleChoice(s string, duo []Challenge) string {
	for _, challenge := range duo {
		var sentences []string
		var mentences []string
		if challenge.Metadata.Sentences != nil {
			for _, sen := range challenge.Metadata.Sentences {
				sentences = append(sentences, "\""+sen.Sentence+"\"")
				if sen.Correct {
					mentences = append(mentences, "\""+sen.Sentence+" *\"")
				} else {
					mentences = append(mentences, "\""+sen.Sentence+"\"")
				}
			}
			s = strings.Replace(s, strings.Join(sentences, ","), strings.Join(mentences, ","), 1)
		}
	}
	return s
}
