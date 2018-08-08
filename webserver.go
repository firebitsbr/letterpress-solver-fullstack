package solver

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	matchInfo []byte
)

type Words []struct {
	Word string `json:"word"`
}

//RunWeb run a webserver
func RunWeb(port string) {

	r := mux.NewRouter()
	// 1. LP solver
	r.HandleFunc("/match", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(matchInfo)
	}).Methods("GET")
	r.HandleFunc("/words", findWords).Methods("GET")
	r.HandleFunc("/word", deleteWord).Methods("DELETE")

	r.PathPrefix("/solver/").Handler(http.StripPrefix("/solver/", http.FileServer(http.Dir("./lpsolver/dist"))))

	// Use default options
	handler := cors.AllowAll().Handler(r)

	log.Println("web server at port", port)
	http.ListenAndServe(":"+port, handler)

}

func findWords(w http.ResponseWriter, r *http.Request) {

	minLetters, _ := r.URL.Query()["selected"]
	maxLetters, _ := r.URL.Query()["letters"]

	res := selectWordsDb(minLetters[0], maxLetters[0])
	ws, _ := json.Marshal(res)
	log.Println("Found words: ", len(res))
	w.Header().Set("Content-Type", "application/json")
	w.Write(ws)
}

func deleteWord(w http.ResponseWriter, r *http.Request) {
	word, _ := r.URL.Query()["delete"]
	log.Println(word[0])
	deleteWordDb(word[0])
}

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

	res := selectWordsDb("", letters)
	exec.Command("say", strconv.Itoa(len(res)), vowals).Run()
}

//MatchInfo ...
type MatchInfo struct {
	Success bool `json:"success"`
	Matches []struct {
		MatchID            string `json:"matchId"`
		MatchIDNumber      int    `json:"matchIdNumber"`
		MatchURL           string `json:"matchURL"`
		CreateDate         string `json:"createDate"`
		UpdateDate         string `json:"updateDate"`
		MatchStatus        int    `json:"matchStatus"`
		CurrentPlayerIndex int    `json:"currentPlayerIndex"`
		Letters            string `json:"letters"`
		RowCount           int    `json:"rowCount"`
		ColumnCount        int    `json:"columnCount"`
		TurnCount          int    `json:"turnCount"`
		MatchData          string `json:"matchData"`
		ServerData         struct {
			Language  int   `json:"language"`
			UsedTiles []int `json:"usedTiles"`
			Tiles     []struct {
				T string `json:"t"`
				O int    `json:"o"`
			} `json:"tiles"`
			UsedWords  []string `json:"usedWords"`
			MinVersion int      `json:"minVersion"`
		} `json:"serverData"`
		Participants []struct {
			UserID                string      `json:"userId"`
			UserName              string      `json:"userName"`
			PlayerIndex           int         `json:"playerIndex"`
			PlayerStatus          string      `json:"playerStatus"`
			LastTurnStatus        string      `json:"lastTurnStatus"`
			MatchOutcome          string      `json:"matchOutcome"`
			TurnDate              string      `json:"turnDate"`
			TimeoutDate           interface{} `json:"timeoutDate"`
			AvatarURL             string      `json:"avatarURL"`
			IsFavorite            bool        `json:"isFavorite"`
			UseBadWords           bool        `json:"useBadWords"`
			BlockChat             bool        `json:"blockChat"`
			DeletedFromPlayerList bool        `json:"deletedFromPlayerList"`
			Online                bool        `json:"online"`
			ChatsUnread           int         `json:"chatsUnread"`
			MuteChat              bool        `json:"muteChat"`
			AbandonedMatch        bool        `json:"abandonedMatch"`
			IsBot                 bool        `json:"isBot"`
			BannedChat            bool        `json:"bannedChat"`
		} `json:"participants"`
	} `json:"matches"`
}

//MatchInfo ...
type MatchInfoSingle struct {
	Success bool `json:"success"`
	Match   struct {
		MatchID            string `json:"matchId"`
		MatchIDNumber      int    `json:"matchIdNumber"`
		MatchURL           string `json:"matchURL"`
		CreateDate         string `json:"createDate"`
		UpdateDate         string `json:"updateDate"`
		MatchStatus        int    `json:"matchStatus"`
		CurrentPlayerIndex int    `json:"currentPlayerIndex"`
		Letters            string `json:"letters"`
		RowCount           int    `json:"rowCount"`
		ColumnCount        int    `json:"columnCount"`
		TurnCount          int    `json:"turnCount"`
		MatchData          string `json:"matchData"`
		ServerData         struct {
			Language  int   `json:"language"`
			UsedTiles []int `json:"usedTiles"`
			Tiles     []struct {
				T string `json:"t"`
				O int    `json:"o"`
			} `json:"tiles"`
			UsedWords  []string `json:"usedWords"`
			MinVersion int      `json:"minVersion"`
		} `json:"serverData"`
		Participants []struct {
			UserID                string      `json:"userId"`
			UserName              string      `json:"userName"`
			PlayerIndex           int         `json:"playerIndex"`
			PlayerStatus          string      `json:"playerStatus"`
			LastTurnStatus        string      `json:"lastTurnStatus"`
			MatchOutcome          string      `json:"matchOutcome"`
			TurnDate              string      `json:"turnDate"`
			TimeoutDate           interface{} `json:"timeoutDate"`
			AvatarURL             string      `json:"avatarURL"`
			IsFavorite            bool        `json:"isFavorite"`
			UseBadWords           bool        `json:"useBadWords"`
			BlockChat             bool        `json:"blockChat"`
			DeletedFromPlayerList bool        `json:"deletedFromPlayerList"`
			Online                bool        `json:"online"`
			ChatsUnread           int         `json:"chatsUnread"`
			MuteChat              bool        `json:"muteChat"`
			AbandonedMatch        bool        `json:"abandonedMatch"`
			IsBot                 bool        `json:"isBot"`
			BannedChat            bool        `json:"bannedChat"`
		} `json:"participants"`
	} `json:"match"`
}
