package main

/*
Template: Render HTML String
Output HTML yang muncul, SELAIN BERSUMBER dari template view, BISA JUGA BERSUMBER
dari sebuah string. Dengan menggunakan method Parse() milik *template.Template
kita bisa menjadikan string html sebagai output.
*/

// 1. Praktek
import (
	"fmt"
	"html/template"
	"net/http"
)

const view string = `<html>
		<head>
				<title>Template</title>
		</head>
		<body>
				<h1>Hello</h1>
		</body>
</html>`

/*
Konstanta bernama view bertipe string disiapkan, dengan isi adalah STRING HTML
yang akan kita jadikan sebagai output nantinya.
*/

/*
buat fungsi main yang isinya route handler /index. Dalam handler tersebut, string html
view diparsing lalu dirender sebagai output.

Tambahkan juga rute /, yang isinya adalah me-redirect request SECARA PAKSA
ke /index menggunakan fungsi http.Redirect().
*/

func main() {
	port := 9000
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		var tmpl = template.Must(template.New("main-template").Parse(view))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// route handler / paksa ke /index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusTemporaryRedirect)
	})
	/*
		Pada kode di atas bisa dilihat, sebuah template bernama main-template disiapkan.
		Template tersebut diisi dengan hasil parsing string html view lewat method Parse().
	*/

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)
}
