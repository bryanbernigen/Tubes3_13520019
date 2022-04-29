package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "root"
	DB_NAME     = "dnadb"
)

type prediksi struct {
	Tanggalprediksi string `json:"tanggalprediksi"`
	Namapasien      string `json:"namapasien"`
	Namapenyakit    string `json:"namapenyakit"`
	Statuspenyakit  bool   `json:"statuspenyakit"`
}

type Penyakit struct {
	Namapenyakit string `json:"namapenyakit"`
	Rantaidna    string `json:"rantaidna"`
}

type Input struct {
	Input string `json:"input"`
}

func main() {
	fmt.Println("Server started on port 8080")
	r := mux.NewRouter()

	r.HandleFunc("/api/submitdisease/{namapenyakit}/{rantaidna}", addpenyakit).Methods("POST")
	r.HandleFunc("/api/getpredictionKMP/{namapasien}/{rantaidna}/{namapenyakit}", addprediksiKMP).Methods("POST")
	r.HandleFunc("/api/getpredictionBM/{namapasien}/{rantaidna}/{namapenyakit}", addprediksiBM).Methods("POST")
	r.HandleFunc("/api/searchdisease/{input}", searchpenyakit).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

func readDNAFromFile(rantaidna string) string {
	var str string = ""
	str = string(rantaidna)
	str = strings.Replace(str, "\n", "", -1)
	// fmt.Println(str)
	return str
}

func validateDNA(DNA string) bool {
	regex := regexp.MustCompile("^[ATCG]+$")
	if regex.MatchString(DNA) {
		fmt.Println("Valid DNA")
		return true
	} else {
		fmt.Println("Invalid DNA")
		return false
	}
}

func KMPMatch(pattern string, text string) bool {
	result := false
	pattern_lenght := len(pattern)
	string_lenght := len(text)

	var lps = make([]int, pattern_lenght)

	len := 0
	lps[0] = 0

	i := 1
	for i < pattern_lenght {
		if pattern[i] == pattern[len] {
			len++
			lps[i] = len
			i++
		} else {
			if len != 0 {
				len = lps[len-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	i = 0
	j := 0
	for i < string_lenght {
		if pattern[j] == text[i] {
			j++
			i++
		}
		if j == pattern_lenght {
			result = true
			fmt.Println("Pattern found at index ", i-j)
			j = lps[j-1]
		} else if i < string_lenght && pattern[j] != text[i] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i = i + 1
			}
		}
	}
	return result
}

func BoyerMooreMatch(pattern string, text string) bool {
	pattern_len := len(pattern)
	text_len := len(text)
	var appeared_char = [256]int{}
	var result bool = false
	var i int = 0
	for i < 256 {
		appeared_char[i] = -1
		i++
	}

	for i = 0; i < pattern_len; i++ {
		appeared_char[int(pattern[i])] = i
	}

	i = 0
	for i <= (text_len - pattern_len) {
		j := pattern_len - 1
		for j >= 0 && pattern[j] == text[i+j] {
			j--
		}
		if j < 0 {
			fmt.Printf("Pattern found at index %d\n", i)
			result = true
			if i+pattern_len < text_len {
				i += pattern_len - appeared_char[int(text[i+pattern_len])]
			} else {
				i += 1
			}

		} else {
			if j-appeared_char[int(text[i+j])] > 0 {
				i += j - appeared_char[int(text[i+j])]
			} else {
				i += 1
			}
		}
	}
	return result
}

func rowExists(query string) bool {
	var exists bool
	db := setupDB()
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("No rows: %s", query)
	}
	return exists
}

func addpenyakit(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)

	namapenyakit := params["namapenyakit"]
	rantaidna := readDNAFromFile(params["rantaidna"])

	db := setupDB()

	if validateDNA(rantaidna) {
		res, err := db.Query("INSERT INTO penyakit (namapenyakit, rantaidna) VALUES ('" + namapenyakit + "','" + rantaidna + "')")
		if err != nil {
			log.Fatal(err)
		}

		defer res.Close()
		defer db.Close()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func shownamapenyakit() {
	db := setupDB()

	res, err := db.Query("SELECT namapenyakit FROM penyakit")
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {
		var namapenyakit string
		err := res.Scan(&namapenyakit)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", namapenyakit)
	}

	defer res.Close()
	defer db.Close()
}

func addprediksiBM(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)

	namapenyakit := params["namapenyakit"]
	rantaidna := readDNAFromFile(params["rantaidna"])
	namapasien := params["namapasien"]

	if validateDNA(rantaidna) {

		db := setupDB()
		if rowExists("SELECT rantaidna FROM penyakit WHERE namapenyakit = '" + namapenyakit + "'") {
			res, err := db.Query("SELECT rantaidna FROM penyakit WHERE namapenyakit = '" + namapenyakit + "'")
			res.Next()

			var pattern string
			res.Scan(&pattern)
			checkErr(err)

			hasil := BoyerMooreMatch(pattern, rantaidna)
			tm := time.Now()
			if hasil {
				fmt.Println("Prediksi berhasil " + tm.Format("2006-01-02") + " " + namapasien + " " + namapenyakit + " " + "true")

				ret := prediksi{Tanggalprediksi: tm.Format("2006-01-02"), Namapasien: namapasien, Namapenyakit: namapenyakit, Statuspenyakit: true}
				jsonResponse, jsonError := json.Marshal(ret)
				if jsonError != nil {
					fmt.Println("Unable to encode JSON")
				}

				fmt.Println(string(jsonResponse))

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)

				res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Format("2006-01-02") + "','" + namapasien + "','" + namapenyakit + "','1')")
				if err != nil {
					log.Fatal(err)
				}

				defer res.Close()
				defer db.Close()
			} else {
				fmt.Println("Prediksi berhasil " + tm.Format("2006-01-02") + " " + namapasien + " " + namapenyakit + " " + "false")

				ret := prediksi{Tanggalprediksi: tm.Format("2006-01-02"), Namapasien: namapasien, Namapenyakit: namapenyakit, Statuspenyakit: false}
				fmt.Println(ret)
				jsonResponse, jsonError := json.Marshal(ret)
				if jsonError != nil {
					log.Fatal(jsonError)
				}

				fmt.Println(string(jsonResponse))

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)

				res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Format("2006-01-02") + "','" + namapasien + "','" + namapenyakit + "','0')")
				if err != nil {
					log.Fatal(err)
				}

				defer res.Close()
				defer db.Close()
			}
		} else {
			fmt.Println("Penyakit tidak ditemukan")
			ret := prediksi{Tanggalprediksi: "", Namapasien: namapasien, Namapenyakit: namapenyakit, Statuspenyakit: false}
			fmt.Println(ret)
			jsonResponse, jsonError := json.Marshal(ret)
			if jsonError != nil {
				log.Fatal(jsonError)
			}

			fmt.Println(string(jsonResponse))

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		}
	}
}

func addprediksiKMP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)

	namapenyakit := params["namapenyakit"]
	rantaidna := readDNAFromFile(params["rantaidna"])
	namapasien := params["namapasien"]

	if validateDNA(rantaidna) {

		db := setupDB()

		if rowExists("SELECT rantaidna FROM penyakit WHERE namapenyakit = '" + namapenyakit + "'") {
			res, err := db.Query("SELECT rantaidna FROM penyakit WHERE namapenyakit = '" + namapenyakit + "'")
			res.Next()

			var pattern string
			res.Scan(&pattern)
			checkErr(err)

			hasil := KMPMatch(pattern, rantaidna)
			tm := time.Now()
			if hasil {
				fmt.Println("Prediksi berhasil " + tm.Format("2006-01-02") + " " + namapasien + " " + namapenyakit + " " + "true")

				ret := prediksi{Tanggalprediksi: tm.Format("2006-01-02"), Namapasien: namapasien, Namapenyakit: namapenyakit, Statuspenyakit: true}
				jsonResponse, jsonError := json.Marshal(ret)
				if jsonError != nil {
					fmt.Println("Unable to encode JSON")
				}

				fmt.Println(string(jsonResponse))

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)

				res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Format("2006-01-02") + "','" + namapasien + "','" + namapenyakit + "','1')")
				checkErr(err)

				defer res.Close()
				defer db.Close()
			} else {
				fmt.Println("Prediksi berhasil " + tm.Format("2006-01-02") + " " + namapasien + " " + namapenyakit + " " + "false")

				ret := prediksi{Tanggalprediksi: tm.Format("2006-01-02"), Namapasien: namapasien, Namapenyakit: namapenyakit, Statuspenyakit: false}
				fmt.Println(ret)
				jsonResponse, jsonError := json.Marshal(ret)
				if jsonError != nil {
					log.Fatal(jsonError)
				}

				fmt.Println(string(jsonResponse))

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)

				res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Format("2006-01-02") + "','" + namapasien + "','" + namapenyakit + "','0')")
				checkErr(err)

				defer res.Close()
				defer db.Close()
			}
		} else {
			fmt.Println("Penyakit tidak ditemukan")
			ret := prediksi{Tanggalprediksi: "", Namapasien: namapasien, Namapenyakit: namapenyakit, Statuspenyakit: false}
			fmt.Println(ret)
			jsonResponse, jsonError := json.Marshal(ret)
			if jsonError != nil {
				log.Fatal(jsonError)
			}

			fmt.Println(string(jsonResponse))

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
		}
	}
}

func searchpenyakit(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)

	inp := params["input"]
	fmt.Println("ini masukannya : ", inp)

	datemonthyear := ""
	namapenyakit := ""
	//Regex dengan Format: tanggal bulan tahun namapenyakit
	re, _ := regexp.Compile("^[^\\d][\\w|-]*$")
	if re.FindString((inp)) != "" {
		namapenyakit = re.FindString((inp))
		fmt.Println("Penyakit Saja")
	} else {
		re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]*(0?[1-9]|1[012]|[\\w]{3,9})[^\\d]*([1-9][0-9]{3})$")
		if re.FindString((inp)) != "" {
			re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]")
			date := re.FindString(inp)
			if len(date) == 2 {
				date = "0" + date
			}
			re, _ = regexp.Compile("^([0][1-9]|[12][0-9]|3[01])")
			date = re.FindString(date)
			//Cari bulan dengan nilai 01-12 atau 1-12
			re, _ = regexp.Compile("[^\\d](0?[1-9]|1[012]|[\\w]{3,9})[^\\d]")
			month := re.FindString(inp)
			fmt.Println("a" + month + "a")
			if len(month) == 3 {
				re, _ = regexp.Compile("[1-9]")
				month = re.FindString(month)
				month = "0" + month
			}
			monthInInt := monthInInt(month)
			re, _ = regexp.Compile("([0][1-9]|1[012])")
			month = re.FindString(monthInInt)
			fmt.Println(month)
			//Cari tahun dengan nilai 1000-9999
			re, _ = regexp.Compile("([1-9][0-9]{3})")
			year := re.FindString(inp)
			datemonthyear = year + "-" + month + "-" + date
			fmt.Println("Tanggal Saja")
		} else {
			re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]*(0?[1-9]|1[012]|[\\w]{3,9})[^\\d]*([1-9][0-9]{3})[^\\d]*[\\w|-]*$")
			if re.FindString((inp)) != "" {
				//Cari tanggal dengan nilai 01-31 atai 1-31
				re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]")
				date := re.FindString(inp)
				if len(date) == 2 {
					date = "0" + date
				}
				re, _ = regexp.Compile("^([0][1-9]|[12][0-9]|3[01])")
				date = re.FindString(date)
				//Cari bulan dengan nilai 01-12 atau 1-12
				re, _ = regexp.Compile("[^\\d](0?[1-9]|1[012]|[\\w]{3,9})[^\\d]")
				month := re.FindString(inp)
				if len(month) == 3 {
					re, _ = regexp.Compile("[1-9]")
					month = re.FindString(month)
					month = "0" + month
				}
				monthInInt := monthInInt(month)
				re, _ = regexp.Compile("([0][1-9]|1[012])")
				month = re.FindString(monthInInt)
				fmt.Println(month)
				//Cari tahun dengan nilai 1000-9999
				re, _ = regexp.Compile("([1-9][0-9]{3})")
				year := re.FindString(inp)
				//Cari nama penyakit dengan whitespace kecuali '-' sebagai pemisah dengan tahun
				re, _ = regexp.Compile("[\\w|-]*$")
				namapenyakit = re.FindString(inp)
				datemonthyear = year + "-" + month + "-" + date
				fmt.Println("Tanggal dan Penyakit Saja")
			}
		}
	}

	db := setupDB()

	var semuaprediksi []prediksi
	if namapenyakit == "" {
		fmt.Println(datemonthyear)

		// Check existence of data searched
		if rowExists("SELECT * FROM prediksi WHERE tanggalprediksi like '" + datemonthyear + "'") {
			res, err := db.Query("SELECT * FROM prediksi WHERE tanggalprediksi like '" + datemonthyear + "'")
			checkErr(err)
			for res.Next() {
				var varprediksi prediksi
				err := res.Scan(&varprediksi.Tanggalprediksi, &varprediksi.Namapasien, &varprediksi.Namapenyakit, &varprediksi.Statuspenyakit)

				checkErr(err)
				semuaprediksi = append(semuaprediksi, prediksi{Tanggalprediksi: varprediksi.Tanggalprediksi, Namapasien: varprediksi.Namapasien, Namapenyakit: varprediksi.Namapenyakit, Statuspenyakit: varprediksi.Statuspenyakit})
			}
			checkErr(err)

			defer res.Close()
		}
	} else {
		if datemonthyear == "" {
			fmt.Println(namapenyakit)
			if rowExists("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "'") {
				res, err := db.Query("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "'")
				if err != nil {
					log.Fatal(err)
				}
				for res.Next() {
					var varprediksi prediksi
					err := res.Scan(&varprediksi.Tanggalprediksi, &varprediksi.Namapasien, &varprediksi.Namapenyakit, &varprediksi.Statuspenyakit)

					checkErr(err)
					semuaprediksi = append(semuaprediksi, prediksi{Tanggalprediksi: varprediksi.Tanggalprediksi, Namapasien: varprediksi.Namapasien, Namapenyakit: varprediksi.Namapenyakit, Statuspenyakit: varprediksi.Statuspenyakit})
				}

				defer res.Close()
			}
		} else {
			fmt.Println(datemonthyear + " " + namapenyakit)
			if rowExists("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "' AND tanggalprediksi like '" + datemonthyear + "'") {
				res, err := db.Query("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "' AND tanggalprediksi like '" + datemonthyear + "'")
				checkErr(err)

				for res.Next() {
					var varprediksi prediksi
					err := res.Scan(&varprediksi.Tanggalprediksi, &varprediksi.Namapasien, &varprediksi.Namapenyakit, &varprediksi.Statuspenyakit)

					checkErr(err)
					semuaprediksi = append(semuaprediksi, prediksi{Tanggalprediksi: varprediksi.Tanggalprediksi, Namapasien: varprediksi.Namapasien, Namapenyakit: varprediksi.Namapenyakit, Statuspenyakit: varprediksi.Statuspenyakit})
				}

				defer res.Close()
			}
		}
	}
	fmt.Println(semuaprediksi)
	// POST
	jsonResponse, jsonError := json.Marshal(semuaprediksi)
	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}
	fmt.Println(string(jsonResponse))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func monthInInt(input string) string {
	re, _ := regexp.Compile("[\\d]*")
	if re.FindString(input) != "" {
		return input
	}
	re, _ = regexp.Compile("[jJ][aA][nN][\\w]*")
	if re.FindString(input) != "" {
		return "01"
	}
	re, _ = regexp.Compile("[fF][eE][bB][\\w]*")
	if re.FindString(input) != "" {
		return "02"
	}
	re, _ = regexp.Compile("[mM][aA][rR][\\w]*")
	if re.FindString(input) != "" {
		return "03"
	}
	re, _ = regexp.Compile("[aA][pP][rR][\\w]*")
	if re.FindString(input) != "" {
		return "04"
	}
	re, _ = regexp.Compile("[mM]([aA][yY]|[eE][iI])[\\w]*")
	if re.FindString(input) != "" {
		return "05"
	}
	re, _ = regexp.Compile("[jJ][uU][nN][\\w]*")
	if re.FindString(input) != "" {
		return "06"
	}
	re, _ = regexp.Compile("[jJ][uU][lL][\\w]*")
	if re.FindString(input) != "" {
		return "07"
	}
	re, _ = regexp.Compile("[aA][uU][gG][\\w]*")
	if re.FindString(input) != "" {
		return "08"
	}
	re, _ = regexp.Compile("[sS][eS][pP][\\w]*")
	if re.FindString(input) != "" {
		return "09"
	}
	re, _ = regexp.Compile("[oO]([cC]|[kK])[tT][\\w]*")
	if re.FindString(input) != "" {
		return "10"
	}
	re, _ = regexp.Compile("[nN][oO][vV][\\w]*")
	if re.FindString(input) != "" {
		return "11"
	}
	re, _ = regexp.Compile("[dD][eE]([cC]|[sS])[\\w]*")
	if re.FindString(input) != "" {
		return "12"
	}
	return input
}
