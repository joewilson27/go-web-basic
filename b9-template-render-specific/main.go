package main

/*
Template: Render Specific HTML Template
pada bagian ini kita akan belajar bagaimana cara untuk render template html TERTENTU.
Sebuah file view BISA BERISIKAN BANYAK template. Template mana yang ingin di-render
bisa ditentukan.
*/

// 1. Front End
/*
Buat template view.html
Pada file view.html, terlihat TERDAPAT 2 template didefinisikan dalam 1 file,
template index dan test. Rencananya template index akan ditampilkan ketika rute / diakses,
dan template test ketika rute /test diakses.
*/

// 2. Back End
import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	port := 9000
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var tmpl = template.Must(template.New("index").ParseFiles("view.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var tmpl = template.Must(template.New("test").ParseFiles("view.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	/*
		Pada kode di atas bisa dilihat, terdapat 2 rute yang masing-masing memparsing file
		YANG SAMA, tapi spesifik template yang dipilih untuk di-render BERBEDA.

		Contoh di rute /, sebuah template dialokasikan dengan nama index, kemudian
		di-parsing-lah view bernama view.html menggunakan method ParseFiles().
		Golang secara cerdas akan melakukan mencari dalam file view tersebut,
		APAKAH ADA TEMPLATE yang namanya adalah index atau tidak.
		Jika ada akan ditampilkan. Hal ini juga berlaku pada rute /test,
		jika isi dari template bernama test akan ditampilkan tiap kali
		rute tersebut diakses.
	*/

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)
}
