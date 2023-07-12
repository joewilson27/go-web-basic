package main

/*
Custom Multiplexer

Silakan salin project sebelumnya
*/

// 1. Pembuatan Custom Mux
/*
Pada chapter sebelumnya, default mux milik Go DIGUNAKAN untuk routing dan implementasi middleware.
Kali ini default mux tersebut TIDAK DIGUNAKAN, kita akan buat mux baru.

Namun pembuatan mux baru TIDAKLAH CUKUP, KARENA NATURALLY mux baru tersebut TIDAK AKAN ADA BEDA dengan default mux.
Oleh karena itu agar lebih berguna, kita akan buat tipe mux baru, meng-embed http.ServeMux ke dalamnya,
lalu membuat beberapa hal dalam struct tersebut.
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
	mux := new(CustomMux)
	/*
		Objek mux dicetak dari struct CustomMux yang jelasnya akan di buat.
		Struct ini di dalamnya meng-embed http.ServeMux.
	*/

	mux.HandleFunc("/student", ActionStudent)

	/*
		Registrasi middleware JUGA DIUBAH, sekarang MENGGUNAKAN method .RegisterMiddleware() milik mux.
	*/
	mux.RegisterMiddleware(MiddlewareAuth)
	mux.RegisterMiddleware(MiddlewareAllowOnlyGet)

	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = mux // Masukan objek mux ke property .Handler milik server.

	fmt.Println("server started at localhost:9000")
	server.ListenAndServe()
}
