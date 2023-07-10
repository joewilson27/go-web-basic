package main

/*
Routing http.HandleFunc
Dalam Go, routing bisa dilakukan dengan beberapa cara, di antaranya:
- Dengan memanfaatkan fungsi http.HandleFunc()
- Mengimplementasikan interface http.Handler pada suatu struct,
	untuk kemudian digunakan pada fungsi http.Handle()
- Membuat multiplexer sendiri dengan memanfaatkan struct http.ServeMux
- Dan lainnya

untuk penjelasan pada chapter ini hanya http.HandleFunc() yang kita pelajari.

*Metode routing cara pertama dan cara kedua MEMILIKI KESAMAAN
yaitu sama-sama menggunakan DefaultServeMux untuk pencocokan rute/endpoint
yang diregistrasikan. Mengenai apa itu DefaultServeMux akan kita bahas lebih mendetail pada chapter lain.
*/

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Penggunaan http.HandleFunc()
	handlerIndex := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/index", handlerIndex)

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello again"))
	})
	/*
		handler bisa anonymous function atau closure (closure -->  fungsi yang bisa disimpan dalam variabel)
	*/

	// contoh dengan closure
	var withClosure = func(w http.ResponseWriter, r *http.Request) {
		content := "Hello from closure!"
		w.Write([]byte(content))
	}

	http.HandleFunc("/closure", withClosure)

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}
