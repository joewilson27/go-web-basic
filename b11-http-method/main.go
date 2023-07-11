package main

/*
HTTP Method: POST & GET
Sebuah route handler pada dasarnya BISA MENERIMA SEGALA JENIS request, dalam artian: apapun jenis HTTP method-nya maka akan tetap masuk ke satu handler (seperti POST, GET, dan atau lainnya). Untuk memisah request berdasarkan method-nya, bisa menggunakan seleksi kondisi.
*/

// 1. Praktek
import (
	"fmt"
	"net/http"
)

func main() {
	port := 9000
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		/*
			Struct *http.Request memiliki property bernama Method yang BISA DIGUNAKAN untuk
			MENGECEK METHOD daripada request yang sedang berjalan.
		*/
		switch r.Method {
		case "POST":
			w.Write([]byte("post"))
		case "GET":
			w.Write([]byte("get"))
		default:
			http.Error(w, "", http.StatusBadRequest)
		}
	})
	/*
		Pada contoh di atas, request ke rute / dengan method POST akan menghasilkan output text post,
		sedangkan method GET menghasilkan output text get.

		*test di postman
	*/

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d", port)
	http.ListenAndServe(addr, nil)
}
