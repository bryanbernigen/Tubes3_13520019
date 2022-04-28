package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	// "os"

	_ "github.com/lib/pq"
)

const (
	DB_USER 	= "postgres"
	DB_PASSWORD = "root"
	DB_NAME 	= "dnadb"
)

type Prediksi struct {
	tanggalprediksi string	`json:"tanggalprediksi"`
	namapasien      string	`json:"namapasien"`
	namapenyakit    string	`json:"namapenyakit"`
	statuspenyakit  bool	`json:"statuspenyakit"`
}

type Penyakit struct {
	namapenyakit string		`json:"namapenyakit"`
	rantaidna    string		`json:"rantaidna"`
}

type Input struct {
	input string	`json:"input"`
}

// var dummy []Penyakit

func main() {
	fmt.Println("Server started on port 8080")
	r := mux.NewRouter()

	type person struct {
		Name     string
		LastName string
		Age      uint8
	}

	//CONTOH HELLO WORLD
	r.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)

		fmt.Println("HELLO WORLD is called!")
		person := person{Name: "Shashank", LastName: "Tiwari", Age: 30}

		jsonResponse, jsonError := json.Marshal(person)
		
		if jsonError != nil {
		fmt.Println("Unable to encode JSON")
		}
		
		fmt.Println(string(jsonResponse))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	}).Methods("GET")


	r.HandleFunc("/{penyakit}",
		func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)
		// get data from body
		var data map[string]interface{}
		json.NewDecoder(r.Body).Decode(&data)
		fmt.Println("data")
		fmt.Println(data)
		fmt.Println(data["name"])

		vars := mux.Vars(r)
		penyakit := vars["penyakit"]
		fmt.Println("penyakti yg dikirim adalah", penyakit);

		person := person{Name: "ini yg dari post", LastName: "huhu", Age: 14}

		jsonResponse, jsonError := json.Marshal(person)
		
		if jsonError != nil {
		fmt.Println("Unable to encode JSON")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	}).Methods("POST");
	r.HandleFunc("/api/submitdisease", addpenyakit).Methods("POST")
	r.HandleFunc("/api/getprediction", addprediksi).Methods("POST")
	r.HandleFunc("/api/searchdisease", searchpenyakit).Methods("GET")
	http.ListenAndServe(":8080", r)
}

// func getPort() (string, error) {
// 	// the PORT is supplied by Heroku
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 	  return "", fmt.Errorf("$PORT not set")
// 	}
// 	return ":" + port, nil
// }
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func setupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)

    checkErr(err)

    return db
}

func readDNAFromFile(fileName string) string {
	var str string = ""
	fileName = "../test/" + fileName
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	str = string(b)
	str = strings.Replace(str, "\n", "", -1)
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
	var pattern_len int = len(pattern)
	var text_len int = len(text)
	var kmparray = make([]int, pattern_len)
	var result bool = false
	//Cari Suffix Prefix terbesar
	len := 0
	kmparray[0] = 0
	var j int = 1 //Pointer buat pattern
	for j < pattern_len-1 {
		if pattern[j] == pattern[len] {
			len++
			kmparray[j] = len
			j++
		} else {
			if len != 0 {
				len = kmparray[len-1]
			} else {
				kmparray[j] = 0
				j++
			}
		}
	}

	//String Matching
	var i int = 0 //Pointer buat text
	j = 0         //Pointer buat pattern
	for i < text_len {
		if text[i] == pattern[j] {
			i++
			j++
		} else {
			if j != 0 {
				j = kmparray[j-1]
			} else {
				i++
			}
		}
		if j == pattern_len {
			fmt.Printf("Pattern found at index %d\n", i-j)
			result = true
			j = kmparray[j-1]
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

func addpenyakit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db := setupDB()

	if (validateDNA(params["rantaidna"])) {
		res, err := db.Query("INSERT INTO penyakit (namapenyakit, rantaidna) VALUES ('" + params["namapenyakit"] + "','" + params["rantaidna"] + "')")
		if err != nil {
			log.Fatal(err)
		}

		defer res.Close()
		defer db.Close()
	}
}

func shownamapenyakit() {
	// db, err := sql.Open("mysql", db_username+":"+db_password+"@tcp(127.0.0.1:3306)/"+db_name)
	// if err != nil {
	// 	panic(err.Error())
	// }

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



func addprediksi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	if (validateDNA(params["rantaidna"])) {

		db := setupDB()

		res, err := db.Query("SELECT rantaidna FROM penyakit WHERE namapenyakit = '" + params["namapenyakit"] + "'")
		res.Next()
		var pattern string
		res.Scan(&pattern)
		if err != nil {
			log.Fatal(err)
		}

		hasil := KMPMatch(pattern, params["sequencedna"])
		tm := time.Now()
		if hasil {
			res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Format("2006-01-02") + "','" + params["namapengguna"] + "','" + params["namapenyakit"] + "','1')")
			if err != nil {
				log.Fatal(err)
			}

			defer res.Close()
			defer db.Close()

			fmt.Println(json.NewEncoder(w).Encode(Prediksi{tanggalprediksi: tm.Format("2006-01-02"), namapasien: params["namapasien"], namapenyakit: params["namapenyakit"], statuspenyakit: true}))
		} else {
			res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Format("2006-01-02") + "','" + params["namapengguna"] + "','" + params["namapenyakit"] + "','0')")
			if err != nil {
				log.Fatal(err)
			}

			defer res.Close()
			defer db.Close()

			json.NewEncoder(w).Encode(Prediksi{tanggalprediksi: tm.Format("2006-01-02"), namapasien: params["namapasien"], namapenyakit: params["namapenyakit"], statuspenyakit: false})
		}
	}
}

func searchpenyakit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var input string = params["input"]

	datemonthyear := ""
	namapenyakit := ""
	//Regex dengan Format: tanggal bulan tahun namapenyakit
	re, _ := regexp.Compile("^[\\w|-]*$")
	if re.FindString((input)) != "" {
		namapenyakit = re.FindString((input))
		fmt.Println("Penyakit Saja")
	} else {
		re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]*(0?[1-9]|1[012]|[\\w]{3,9})[^\\d]*([1-9][0-9]{3})$")
		if re.FindString((input)) != "" {
			re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]")
			date := re.FindString(input)
			if len(date) == 2 {
				date = "0" + date
			}
			re, _ = regexp.Compile("^([0][1-9]|[12][0-9]|3[01])")
			date = re.FindString(date)
			//Cari bulan dengan nilai 01-12 atau 1-12
			re, _ = regexp.Compile("[^\\d](0?[1-9]|1[012]|[\\w]{3,9})[^\\d]")
			month := re.FindString(input)
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
			year := re.FindString(input)
			datemonthyear = year + "-" + month + "-" + date
			fmt.Println("Tanggal Saja")
		} else {
			re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]*(0?[1-9]|1[012]|[\\w]{3,9})[^\\d]*([1-9][0-9]{3})[^\\d]*[\\w|-]*$")
			if re.FindString((input)) != "" {
				//Cari tanggal dengan nilai 01-31 atai 1-31
				re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]")
				date := re.FindString(input)
				if len(date) == 2 {
					date = "0" + date
				}
				re, _ = regexp.Compile("^([0][1-9]|[12][0-9]|3[01])")
				date = re.FindString(date)
				//Cari bulan dengan nilai 01-12 atau 1-12
				re, _ = regexp.Compile("[^\\d](0?[1-9]|1[012]|[\\w]{3,9})[^\\d]")
				month := re.FindString(input)
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
				year := re.FindString(input)
				//Cari nama penyakit dengan whitespace kecuali '-' sebagai pemisah dengan tahun
				re, _ = regexp.Compile("[\\w|-]*$")
				namapenyakit = re.FindString(input)
				datemonthyear = year + "-" + month + "-" + date
				fmt.Println("Tanggal dan Penyakit Saja")
			}
		}
	}

	db := setupDB()

	var semuaprediksi []Prediksi
	if namapenyakit == "" {
		fmt.Println(datemonthyear)
		res, err := db.Query("SELECT * FROM prediksi WHERE tanggalprediksi like '" + datemonthyear + "'")
		if err != nil {
			log.Fatal(err)
		}
		for res.Next() {
			var prediksi Prediksi
			err := res.Scan(&prediksi.tanggalprediksi, &prediksi.namapasien, &prediksi.namapenyakit, &prediksi.statuspenyakit)
			
			if err != nil {
				log.Fatal(err)
			}
			semuaprediksi = append(semuaprediksi, Prediksi{tanggalprediksi: prediksi.tanggalprediksi, namapasien: prediksi.namapasien, namapenyakit: prediksi.namapenyakit, statuspenyakit: prediksi.statuspenyakit})
		}
		defer res.Close()
		json.NewEncoder(w).Encode(semuaprediksi)
	} else {
		if datemonthyear == "" {
			fmt.Println(namapenyakit)
			res, err := db.Query("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "'")
			if err != nil {
				log.Fatal(err)
			}
			for res.Next() {
				var prediksi Prediksi
				err := res.Scan(&prediksi.tanggalprediksi, &prediksi.namapasien, &prediksi.namapenyakit, &prediksi.statuspenyakit)

				if err != nil {
					log.Fatal(err)
				}
				semuaprediksi = append(semuaprediksi, Prediksi{tanggalprediksi: prediksi.tanggalprediksi, namapasien: prediksi.namapasien, namapenyakit: prediksi.namapenyakit, statuspenyakit: prediksi.statuspenyakit})
			}
			fmt.Println(semuaprediksi)
			defer res.Close()
		}else {
			fmt.Println(datemonthyear + " " + namapenyakit)
			res, err := db.Query("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "' AND tanggalprediksi like '" + datemonthyear + "'")
				log.Fatal(err)
			
			for res.Next() {
				var prediksi Prediksi
				err := res.Scan(&prediksi.tanggalprediksi, &prediksi.namapasien, &prediksi.namapenyakit, &prediksi.statuspenyakit)

				if err != nil {
					log.Fatal(err)
				}
				semuaprediksi = append(semuaprediksi, Prediksi{tanggalprediksi: prediksi.tanggalprediksi, namapasien: prediksi.namapasien, namapenyakit: prediksi.namapenyakit, statuspenyakit: prediksi.statuspenyakit})
			}
			defer res.Close()
			json.NewEncoder(w).Encode(semuaprediksi)
		}
	}
	fmt.Println(semuaprediksi)
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