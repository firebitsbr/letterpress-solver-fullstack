package solver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

type Conf struct {
	Database struct {
		Drive  string   `json:"drive"`
		User   string   `json:"user"`
		Pass   string   `json:"pass"`
		Host   string   `json:"host"`
		Port   string   `json:"port"`
		Scheme []string `json:"scheme"`
	} `json:"database"`
}

type Word struct {
	id     int    `json:"id"`
	Word   string `json:"word"`
	Length int    `json:"length"`
	A      int    `json:"A"`
	B      int    `json:"B"`
	C      int    `json:"C"`
	D      int    `json:"D"`
	E      int    `json:"E"`
	F      int    `json:"F"`
	G      int    `json:"G"`
	H      int    `json:"H"`
	I      int    `json:"I"`
	J      int    `json:"J"`
	K      int    `json:"K"`
	L      int    `json:"L"`
	M      int    `json:"M"`
	N      int    `json:"N"`
	O      int    `json:"O"`
	P      int    `json:"P"`
	Q      int    `json:"Q"`
	R      int    `json:"R"`
	S      int    `json:"S"`
	T      int    `json:"T"`
	U      int    `json:"U"`
	V      int    `json:"V"`
	W      int    `json:"W"`
	X      int    `json:"X"`
	Y      int    `json:"Y"`
	Z      int    `json:"Z"`
	Valid  int    `json:"valid"`
}

func init() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Conf{}
	err := decoder.Decode(&conf)
	if err != nil {
		panic(err.Error())
	}

	dbconf := conf.Database
	db, err = sql.Open(dbconf.Drive, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbconf.User, dbconf.Pass, dbconf.Host, dbconf.Port, dbconf.Scheme[0]))
	if err != nil {
		log.Panic("open db error", err.Error())
	}

	log.Println("successfully connected to mysql")
}

func selectWordsDb(minLetters string, maxLetters string) []string {
	sqlclause, args := prepareSelectWordsClause(minLetters, maxLetters)

	sql := `SELECT word FROM db_english_all_words WHERE valid = 1 ` + sqlclause + `ORDER BY length ASC LIMIT 999`

	//unpack array as args
	result, err := db.Query(sql, args...)
	if err != nil {
		panic(err.Error())
	}

	res := make([]string, 0, 200)
	for result.Next() {
		var word Word
		err = result.Scan(&word.Word)
		res = append(res, word.Word)
		if err != nil {
			panic(err.Error())
		}
	}
	return res
}

func selectWordsCountDb(minLetters string, maxLetters string) (res int) {
	sqlclause, args := prepareSelectWordsClause(minLetters, maxLetters)

	sql := `SELECT COUNT(*) FROM db_english_all_words WHERE valid = 1 ` + sqlclause

	//unpack array as args
	err := db.QueryRow(sql, args...).Scan(&res)
	if err != nil {
		panic(err.Error())
	}

	return
}

func prepareSelectWordsClause(minLetters string, maxLetters string) (sqlclause string, args []interface{}) {
	loBound := make(map[rune]int)
	hiBound := make(map[rune]int)
	for _, c := range minLetters {
		_, ok := loBound[c]
		if ok {
			loBound[c]++
		} else {
			loBound[c] = 1
		}
	}
	for _, c := range maxLetters {
		_, ok := hiBound[c]
		if ok {
			hiBound[c]++
		} else {
			hiBound[c] = 1
		}
	}

	for _, v := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		l, ok := loBound[v]
		if ok {
			args = append(args, l)
		} else {
			args = append(args, 0)
		}
		h, ok := hiBound[v]
		if ok {
			args = append(args, h)
		} else {
			args = append(args, 0)
		}
		sqlclause = sqlclause + "AND " + string(v) + " >= (?) AND " + string(v) + " <= (?) "
	}
	return
}

func deleteWordDb(word string) {
	sql := `UPDATE db_english_all_words SET valid = 0 WHERE word = (?) `
	_, err := db.Exec(sql, strings.ToLower(word))
	if err != nil {
		panic(err.Error())
	} else {
		log.Println("deleted :", word)
	}
}
