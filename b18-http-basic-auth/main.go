package main

/*
HTTP Basic Authentication
HTTP Basic Auth adalah salah satu TEKNIK OTENTIKASI http request.
Metode ini MEMBUTUHKAN INFORMASI username dan password untuk DISISIPKAN
dalam header request (dengan format tertentu), jadi cukup sederhana,
TIDAK MEMERLUKAN COOKIES maupun SESSION. Lebih jelasnya silakan baca https://datatracker.ietf.org/doc/html/rfc7617

Informasi username dan password tidak serta merta disisipkan dalam header,
informasi tersebut HARUS di-encode TERLEBIH DAHULU ke dalam format yg sudah
ditentukan sesuai spesifikasi, sebelum dimasukan ke header.

contoh penulisan basic auth:
// Request header
Authorization: Basic c29tZXVzZXJuYW1lOnNvbWVwYXNzd29yZA==

Informasi DISISIPKAN dalam request header dengan key Authorization, dan value adalah Basic
spasi hasil enkripsi dari data username dan password. Data username dan password
digabung dengan separator tanda titik dua (:), lalu di-encode dalam format
encoding Base 64.

// Username password encryption
base64encode("someusername:somepassword")
// Hasilnya adalah c29tZXVzZXJuYW1lOnNvbWVwYXNzd29yZA==

Golang MENYEDIAKAN FUNGSI untuk meng-handle request basic auth dengan
cukup mudah, jadi TIDAK PERLU untuk memparsing header request terlebih
dahulu untuk mendapatkan informasi username dan password.

*/

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ActionStudent(w http.ResponseWriter, r *http.Request) {
	if !Auth(w, r) {
		return
	}
	/*
		Auth() untuk mengecek apakah request merupakan valid basic auth request atau tidak.
	*/

	if !AllowOnlyGET(w, r) {
		return
	}
	/*
		AllowOnlyGET(), gunanya untuk MEMASTIKAN HANYA request dengan method GET
		yang diperbolehkan masuk.
	*/

	/*
	 jika ada parameter student id, maka hanya user dengan id yg DIINGINKAN
	 yg dijadikan nilai balik, lewat fungsi SelectStudent(id).
	*/
	if id := r.URL.Query().Get("id"); id != "" {
		OutputJSON(w, SelectStudent(id))
		return
	}

	/*
		Ketika TIDAK ADA parameter student id, maka endpoint ini mengembalikan
		semua data user yang ada, lewat pemanggilan fungsi GetStudents().
	*/
	OutputJSON(w, GetStudents())
}

// OutputJSON(). Fungsi ini digunakan untuk mengkonversi data menjadi JSON string.
func OutputJSON(w http.ResponseWriter, o interface{}) {
	res, err := json.Marshal(o)
	/*
		KONVERSI dari objek atau slice ke JSON string BISA dilakukan
		dengan memanfaatkan json.Marshal.
	*/
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	http.HandleFunc("/student", ActionStudent)

	server := new(http.Server)
	server.Addr = ":9000"

	fmt.Println("server started at localhost:9000")
	server.ListenAndServe()

	/*
		run on terminal:
		go run *.go

		JANGAN MENGGUNAKAN go run main.go, dikarenakan dalam package main terdapat beberapa file
		lain yang harus di-ikut-sertakan pada saat runtime.


		pada postman test endpoint berikut dengan basic authentication username dan password
		http://localhost:9000/student
		http://localhost:9000/student?id=s001
	*/
}
