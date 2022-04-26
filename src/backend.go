package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Prediksi struct {
	tanggalprediksi string
	namapasien      string
	namapenyakit    string
	statuspenyakit  bool
}

var db_username string = "root"
var db_password string = ""
var db_name string = "dna"

func main() {
	// a := readDNAFromFile("homo_sapiens.txt")
	// fmt.Print(a)
	searchpenyakit("07-02-2022")
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

func addpenyakit(namapenyakit string, rantaidna string) {
	db, err := sql.Open("mysql", db_username+":"+db_password+"@tcp(127.0.0.1:3306)/"+db_name)
	if err != nil {
		panic(err.Error())
	}

	var DNA string = rantaidna
	res, err := db.Query("INSERT INTO penyakit (namapenyakit, rantaidna) VALUES ('" + namapenyakit + "','" + DNA + "')")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close()
	defer db.Close()
}

func shownamapenyakit() {
	db, err := sql.Open("mysql", db_username+":"+db_password+"@tcp(127.0.0.1:3306)/"+db_name)
	if err != nil {
		panic(err.Error())
	}

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

func addprediksi(namapengguna string, sequencedna string, namapenyakit string) {
	db, err := sql.Open("mysql", db_username+":"+db_password+"@tcp(127.0.0.1:3306)/"+db_name)
	if err != nil {
		panic(err.Error())
	}

	res, err := db.Query("SELECT rantaidna FROM penyakit WHERE namapenyakit = '" + namapenyakit + "'")
	res.Next()
	var pattern string
	res.Scan(&pattern)
	if err != nil {
		log.Fatal(err)
	}

	hasil := KMPMatch(pattern, sequencedna)
	tm := time.Now()
	if hasil {
		res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Local().Local().Format("2006-01-01") + "','" + namapengguna + "','" + namapenyakit + "','1')")
		if err != nil {
			log.Fatal(err)
		}

		defer res.Close()
		defer db.Close()
	} else {
		res, err := db.Query("INSERT INTO prediksi VALUES('" + tm.Local().Format("2006-01-01") + "','" + namapengguna + "','" + namapenyakit + "','0')")
		if err != nil {
			log.Fatal(err)
		}

		defer res.Close()
		defer db.Close()
	}
}

func searchpenyakit(input string) {
	//Regex dengan Format: tanggal bulan tahun namapenyakit
	re, _ := regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])[^\\d]*(0?[1-9]|1[012]|[\\w]{3,9})[^\\d]*([1-9][0-9]{3})[^\\d]*[\\w|-]*$")
	if re.FindString((input)) != "" {
		//Cari tanggal dengan nilai 01-31 atai 1-31
		re, _ = regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])")
		date := re.FindString(input)
		if len(date) == 1 {
			date = "0" + date
		}
		//Cari bulan dengan nilai 01-12 atau 1-12
		re, _ = regexp.Compile("(0?[1-9]|1[012]|[\\w]{3,9})")
		month := re.FindString(input)
		if len(month) == 1 {
			month = "0" + month
		}
		monthInInt := monthInInt(month)
		//Cari tahun dengan nilai 1000-9999
		re, _ = regexp.Compile("([1-9][0-9]{3})")
		year := re.FindString(input)
		//Cari nama penyakit dengan whitespace kecuali '-' sebagai pemisah dengan tahun
		re, _ = regexp.Compile("[\\w|-]*$")
		namapenyakit := re.FindString(input)
		if namapenyakit == year {
			namapenyakit = ""
		}
		datemonthyear := year + "-" + monthInInt + "-" + date
		datemonthyearinput := date + " " + month + " " + year
		//Connect dengan Database
		db, err := sql.Open("mysql", db_username+":"+db_password+"@tcp(127.0.0.1:3306)/"+db_name)
		if err != nil {
			log.Fatal(err)
		}
		var semuaprediksi []Prediksi
		if namapenyakit == "" {
			fmt.Println(datemonthyear)
			res, err := db.Query("SELECT * FROM prediksi WHERE tanggalprediksi = '" + datemonthyear + "'")
			if err != nil {
				log.Fatal(err)
			}
			for res.Next() {
				var prediksi Prediksi
				err := res.Scan(&prediksi.tanggalprediksi, &prediksi.namapasien, &prediksi.namapenyakit, &prediksi.statuspenyakit)

				if err != nil {
					log.Fatal(err)
				}
				prediksi.tanggalprediksi = datemonthyearinput
				semuaprediksi = append(semuaprediksi, prediksi)
			}
			fmt.Println(semuaprediksi)
			defer res.Close()
		} else {
			fmt.Println(datemonthyear + " " + namapenyakit)
			res, err := db.Query("SELECT * FROM prediksi WHERE namapenyakit like '" + namapenyakit + "' AND tanggalprediksi = '" + datemonthyear + "'")
			if err != nil {
				log.Fatal(err)
			}
			for res.Next() {
				var prediksi Prediksi
				err := res.Scan(&prediksi.tanggalprediksi, &prediksi.namapasien, &prediksi.namapenyakit, &prediksi.statuspenyakit)

				if err != nil {
					log.Fatal(err)
				}
				prediksi.tanggalprediksi = datemonthyearinput
				semuaprediksi = append(semuaprediksi, prediksi)
			}
			fmt.Println(semuaprediksi)
			defer res.Close()
		}
	}
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
	re, _ = regexp.Compile("[oO][cC][tT][\\w]*")
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
