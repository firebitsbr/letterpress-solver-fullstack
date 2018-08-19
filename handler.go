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
	timeInterval := 200 * time.Millisecond

	for i, k := range clickList {
		x := left + width*(k%5)
		y := top + width*(k/5)
		go func(x int, y int, k int, i int) {
			err := exec.Command("adb", "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y)).Run() // tap the tile
			log.Println("touch", k)
			if err != nil {
				log.Println("error: check adb connection.", err)
			}
			if i == len(clickList)-1 {
				time.Sleep(timeInterval * 2)
				exec.Command("adb", "shell", "input", "tap", "1000", "50").Run() // click SUBMIT
			}
		}(x, y, k, i) //pass loop local vars to goroutine!
		time.Sleep(timeInterval)
		if i < 8 {
			timeInterval += 10 * time.Millisecond
		}
	}
}
