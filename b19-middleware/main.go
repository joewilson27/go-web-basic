package main

/*
Middleware http.Handler

Pada chapter ini, kita akan belajar penggunaan interface http.Handler untuk implementasi custom middleware.
Kita akan menggunakan sample proyek pada chapter sebelumnya B.18. HTTP Basic Auth sebagai dasar bahan pembahasan
chapter ini.

	*Apa itu middleware? Istilah middleware berbeda-beda di tiap bahasa/framework. NodeJS dan Rails ada
	 istilah middleware. Pada pemrograman Java Enterprise, istilah filters digunakan.
	 Pada C# istilahnya adalah delegate handlers. Definisi dari middleware sendiri versi penulis,
	 sebuah BLOK KODE yang DIPANGGIL SEBELUM ataupun SESUDAH http request di proses.

Pada chapter sebelumnya, kalau dilihat, ADA EEBERAPA PROSES YANG DI JALANKAN dalam handler rute /student,
yaitu pengecekan otentikasi dan pengecekan method. Misalnya TERDAPAT RUTE LAGI, MAKA dua validasi tersebut
juga HARUS DIPANGGIL LAGI dalam handlernya.


func ActionStudent(w http.ResponseWriter, r *http.Request) {
    if !Auth(w, r) { return }
    if !AllowOnlyGET(w, r) { return }

    // ...
}

Jika ada BANYAK RUTE, apa yang harus kita lakukan? salah satu solusi yang bisa digunakan adalah dengan memanggil
fungsi Auth() dan AllowOnlyGet() DI SEMUA handler rute yang ada. Namun jelasnya INI BUKAN best practice.
Dan juga belum tentu di tiap rute hanya ada dua validasi ini, bisa saja ADA LEBIH banyak proses,
misalnya pengecekan csrf, authorization, dan lainnya.
*/

// 1. Interface http.Handler
/*
Interface http.Handler MERUPAKAN tipe data PALING POPULER di Go UNTUK KEPERLUAN manajemen middleware.
Struct yang mengimplementasikan interface ini DIWAJIBKAN MEMILIKI method dengan skema ServeHTTP(ResponseWriter, *Request).

Di Go sendiri objek utama untuk KEPERLUAN ROUTING yaitu mux atau multiplexer, adalah mengimplementasikan
interface http.Handler ini.

Dengan memanfaatkan interface ini, kita akan membuat beberapa middleware. Fungsi pengecekan otentikasi dan pengecekan method
akan kita ubah menjadi middleware terpisah.
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func ActionStudent(w http.ResponseWriter, r *http.Request) {
	/*
		PERUBAHAN pada kode ActionStudent() adalah, pengecekan basic auth dan pengecekan method DIHAPUS.
	*/
	if id := r.URL.Query().Get("id"); id != "" {
		OutputJSON(w, SelectStudent(id))
		return
	}

	OutputJSON(w, GetStudents())
}

func main() {
	// 3. Mux / Multiplexer
	/*
		Di Go, mux (kependekan dari multiplexer) adalah router. Semua routing PASTI DILAKUKAN lewat objek mux.

		Apa benar? routing http.HandleFunc() sepertinya tidak menggunakan mux? Begini, SEBENARNYA routing tersebut JUGA MENGGUNAKAN
		mux. Go memiliki default objek mux yaitu http.DefaultServeMux. Routing yang LANGSUNG DILAKUKAN dilakukan dari
		fungsi HandleFunc() milik package net/http sebenarnya mengarah ke method default mux http.DefaultServeMux.HandleFunc().
		Agar lebih jelas, silakan perhatikan dua kode berikut.

				http.HandleFunc("/student", ActionStudent)

				// vs

				mux := http.DefaultServeMux
				mux.HandleFunc("/student", ActionStudent)

		Dua kode di atas menghasilkan hasil yang SAMA PERSIS.
		Mux sendiri adalah BENTUK NYATA struct yang MENGIMPLEMENTASIKAN interface http.Handler. Untuk lebih
		jelasnya silakan baca dokumentasi package net/http di https://golang.org/pkg/net/http/#Handle.

	*/

	mux := http.DefaultServeMux

	mux.HandleFunc("/student", ActionStudent)

	var handler http.Handler = mux
	/*
		Di kode setelah routing di atas, bisa dilihat objek mux ditampung ke variabel baru (handler) bertipe http.Handler
		Seperti ini adalah valid karena memang struct multiplexer MEMENUHI KRITERIA interface http.Handler,
		yaitu MEMILIKI method ServeHTTP().
	*/

	/*
		Lalu dari objek handler tersebut, ke-dua middleware DIPANGGIL dengan parameter adalah
		objek handler ITU SENDIRI dan nilai baliknya ditampung pada objek yang sama.
	*/
	handler = MiddlewareAuth(handler) // param nya adl objek handler itu sendiri dan nilai baliknya ditampung pada objek yang sama
	handler = MiddlewareAllowOnlyGet(handler)

	/*
		Fungsi MiddlewareAuth() dan MiddlewareAllowOnlyGet() adalah middleware yang akan kita buat setelah ini.
		Cara registrasi middleware yang paling populer adalah dengan memanggilnya secara SEKUENSIAL atau BERURUTAN,
		seperti pada kode di atas.
			- MiddlewareAuth() bertugas untuk melakukan pengencekan credentials, basic auth.
			- MiddlewareAllowOnlyGet() bertugas untuk melakukan pengecekan method.


			*Silakan lihat source code beberapa library middleware yang sudah terkenal seperti gorilla,
				gin-contrib, echo middleware, dan lainnya; Semua metode implementasi middleware-nya adalah sama,
				atau paling tidak mirip. Point plus nya, beberapa di antara library tersebut mudah diintegrasikan
				dan compatible satu sama lain.

		Kedua middleware yang akan kita buat tersebut MENGEMBALIKAN fungsi bertipe http.Handler. Eksekusi middleware
		sendiri terjadi pada saat ada http request masuk.
	*/

	/*
		Setelah semua middleware DIREGISTRASI. Masukan objek handler ke property .Handler milik server.
	*/
	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = handler // Masukan objek handler ke property .Handler milik server.

	fmt.Println("server started at localhost:9000")
	server.ListenAndServe()
}
