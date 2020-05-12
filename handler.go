package solver

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func setMatch(jsonBytes []byte) {
	matchInfo = jsonBytes
}

func checkVowalsMatch(jsonBytes []byte) {
	match := &MatchInfoSingle{}
	json.Unmarshal(jsonBytes, match)
	tiles := match.Match.ServerData.Tiles
	letters := ""
	for _, tile := range tiles {
		letters += tile.T
	}

	vowals := ","
	for _, v := range []rune("AEIOU") {
		if strings.ContainsRune(letters, v) {
			vowals += string(v) + ","
		}
	}

	if len(letters) > 0 {
		res := selectWordsCountDb("", letters)
		exec.Command("say", strconv.Itoa(res/1000), "k").Run()
	}
}

func clickTiles(clickList []int) {
	width := 216
	left := 108
	top := 803
	timeInterval := 220 * time.Millisecond

	for i, k := range clickList {
		x := left + width*(k%5)
		y := top + width*(k/5)
		go func(x int, y int, k int, i int) {
			err := exec.Command("adb", "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y)).Run() // tap the tile
			if err != nil {
				log.Println("error: check adb connection.", err)
			}
			// if i == len(clickList)-1 {
			// 	time.Sleep(timeInterval * 2)
			// 	exec.Command("adb", "shell", "input", "tap", "1000", "50").Run() // click SUBMIT
			// }
		}(x, y, k, i) //pass loop local vars to goroutine!
		time.Sleep(timeInterval)
		if i < 6 {
			timeInterval += 15 * time.Millisecond
		}
	}
}

func markPlayedWord(playedWords []string) {
	tagPlayedWordDb(playedWords[len(playedWords)-1])
}

func addLastPlayedWords(bs []byte) {
	matches := &MatchInfo{}
	json.Unmarshal(bs, matches)
	lastPlayedWords := make([]string, 0)
	for _, match := range matches.Matches {
		uw := match.ServerData.UsedWords
		status := match.MatchStatus
		currentID := match.Participants[match.CurrentPlayerIndex].UserID
		if len(uw) > 0 && (status != 1 || contains(conf.Letterpress.UserIDs, currentID)) {
			lastPlayedWords = append(lastPlayedWords, uw[len(uw)-1])
		}
	}
	addWordsDB(lastPlayedWords)
}
func addLastPlayedWord(bs []byte) {
	match := &MatchInfoSingle{}
	json.Unmarshal(bs, match)
	lastPlayedWords := make([]string, 0)
	uw := match.Match.ServerData.UsedWords
	if len(uw) > 0 {
		lastPlayedWords = append(lastPlayedWords, uw[len(uw)-1])
	}

	addWordsDB(lastPlayedWords)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
