# Tubes3_13520019
Tugas Besar III IF2211 Strategi Algoritma Semester II Tahun 2020/2021 Penerapan String Matching dan Regular Expression dalam DNA Pattern Matching

# Author
> Maharani Putri Ayu Irawan (13520019)
> Bryan Bernigen (13520034)
> Ng Kyle (13520040)

## Table of Contents
* [General Info](#general-information)
* [Technologies Used and Requirements](#technologies-used-and-requirements)
* [Usage](#usage)
* [Setup](#setup)
* [Project Status](#project-status)


## General Information
- Program is a Web based Application for DNA Testing using KMP Pattern Search and Boyer-Moore Pattern Search
- Implementation including: KMP Search, Boyer-Moore Search, Web (Front-End and Back-End)

## Technologies Used and Requirements
- Go Lang Version 1.18 
- Postgressql


## Usage
1.	Run backend :
- cd to path Tubes3_13520019\src
- run : go run backend.go
2.	Run frontend :
- cd to Tubes3_13520019\src\frontend\dnapatternmatching
- run : npm run dev
3.	Web App Usage :
a.	Fitur Add Disease
- Masukkan nama penyakit yang akan diinput
- Gunakan button Choose File untuk memasukkan file txt yang berisi konfigurasi DNA penyakit
- Tekan tombol Submit, penyakit otomatis tersimpan pada database.
b.	Fitur Check Disease
- Masukkan Nama pasien yang akan diperiksa
- Gunakan button Choose File untuk memasukkan file txt yang berisi konfigurasi DNA pasien
- Masukkan Nama Penyakit yang akan diperiksa
- Pilih salah satu opsi (KMP atau Boyer-Moore) dengan cara klik button, lalu hasil akan tertampil.
c.	Fitur History
- Masukkan Query berdasarkan format yang dijelaskan pada subbab sebelumnya. 
- Hasil pencarian akan tertampilkan (lebih dari 1) atau tidak.


## Setup
You need to have go installed, postgressql.
Create a database with name dnadb. Load dump file.


## Project Status
Project is: _completed_
