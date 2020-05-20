package solver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db    *sql.DB
	conf  Conf
	table = "en_words" //db_english_all_words en_words
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
	Letterpress struct {
		UserIDs []string `json:"UserIDs"`
	} `json:"Letterpress"`
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

	sql := `SELECT word,valid FROM ` + table + ` WHERE valid IN (1,2) ` + sqlclause + `ORDER BY frequency DESC, length ASC LIMIT 1999`
	result, err := db.Query(sql, args...)
	if err != nil {
		panic(err.Error())
	}

	res := make([]string, 0, 200)
	for result.Next() {
		var word Word
		err = result.Scan(&word.Word, &word.Valid)
		//double the played word when build the search result, so that the frontend can know the played word
		if word.Valid == 2 {
			res = append(res, word.Word+"*")
		} else {
			res = append(res, word.Word)
		}
		if err != nil {
			panic(err.Error())
		}
	}
	return res
}

func selectWordsCountDb(minLetters string, maxLetters string) (res int) {
	sqlclause, args := prepareSelectWordsClause(minLetters, maxLetters)

	sql := `SELECT COUNT(*) FROM ` + table + ` WHERE valid > 0 ` + sqlclause

	//unpack array as args
	err := db.QueryRow(sql, args...).Scan(&res)
	if err != nil {
		panic(err.Error())
	}

	return
}

func selectWordsFreqeuncyDb(minLetters string, maxLetters string) (res []int) {
	res = make([]int, 26)
	arr := make([]interface{}, 26)
	sqlclause, args := prepareSelectWordsClause(minLetters, maxLetters)

	sql := `SELECT `
	for i, l := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		sql += `SUM(` + string(l) + `),`
		arr[i] = &res[i]
	}
	sql = sql[:len(sql)-1]
	sql += ` FROM ` + table + ` WHERE valid > 0 ` + sqlclause

	err := db.QueryRow(sql, args...).Scan(arr...)
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

func deleteWordDb(inValidWord string) {

	sql := `UPDATE ` + table + ` SET valid = 0 WHERE word = (?) `
	_, err := db.Exec(sql, strings.ToLower(inValidWord))
	if err != nil {
		panic(err.Error())
	} else {
		log.Println("deleted :", inValidWord)
	}
}

func tagPlayedWordDb(word string) {
	sql := `UPDATE ` + table + ` SET valid = 2 WHERE word = (?) `
	_, err := db.Exec(sql, strings.ToLower(word))
	if err != nil {
		panic(err.Error())
	} else {
		log.Println("played :", word)
	}
}

func removeByValue(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func addWordsDB(words []string) {
	if len(words) == 0 {
		return
	}

	// Filter out existing words in db
	sql := `SELECT word FROM ` + table + ` WHERE word in ('` + strings.Join(words, "','") + `')`
	result, err := db.Query(sql)
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		var word Word
		err := result.Scan(&word.Word)
		words = removeByValue(words, word.Word)
		if err != nil {
			panic(err.Error())
		}
	}

	sql = `INSERT IGNORE INTO ` + table + `(word,length,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,frequency) 
	VALUES ((?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),(?),123456)`

	for _, word := range words {
		result, err := db.Exec(sql, word,
			len(word),
			strings.Count(word, "a"),
			strings.Count(word, "b"),
			strings.Count(word, "c"),
			strings.Count(word, "d"),
			strings.Count(word, "e"),
			strings.Count(word, "f"),
			strings.Count(word, "g"),
			strings.Count(word, "h"),
			strings.Count(word, "i"),
			strings.Count(word, "j"),
			strings.Count(word, "k"),
			strings.Count(word, "l"),
			strings.Count(word, "m"),
			strings.Count(word, "n"),
			strings.Count(word, "o"),
			strings.Count(word, "p"),
			strings.Count(word, "q"),
			strings.Count(word, "r"),
			strings.Count(word, "s"),
			strings.Count(word, "t"),
			strings.Count(word, "u"),
			strings.Count(word, "v"),
			strings.Count(word, "w"),
			strings.Count(word, "x"),
			strings.Count(word, "y"),
			strings.Count(word, "z"),
		)
		if err != nil {
			panic(err.Error())
		} else {
			lastInsertID, _ := result.LastInsertId()
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected > 0 {
				exec.Command("say", "new word; ;"+word).Run()
				log.Printf("add word: %s; %v; %v", word, lastInsertID, rowsAffected)
			}
		}
	}
}

func allLetterFrequencyDb() {
	as := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	m := make(map[rune]int)
	for _, l := range as {
		sql := `SELECT SUM(` + string(l) + `) FROM ` + table + ` WHERE valid = 1 `
		i := 0
		err := db.QueryRow(sql).Scan(&i)
		if err != nil {
			panic(err.Error())
		}
		m[l] = i
		fmt.Printf("%q\t%v\n", l, i)
	}

	type kv struct {
		Key   rune
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	for _, kv := range ss {
		fmt.Printf("%q", kv.Key)
	}
}
