package main

/*
Template: Custom Functions

Pada chapter sebelumnya kita telah mengenal beberapa predefined function (built-in atau bawaan)
yang disediakan oleh Go. Kali ini kita akan belajar tentang fungsi custom,
bagaimana cara membuat dan menggunakannya dalam template.
*/

// 1. Front End
/*
Persiapkan view.html. pada view tersebut ada 2 hal yg perlu diperhatikan
- Fungsi unescape(), digunakan untuk menampilkan string tanpa di-escape
- Fungsi avg(), digunakan untuk mencari rata-rata dari angka-angka yang
	disisipkan sebagai parameter
Kedua fungsi tersebut adalah fungsi kustom yang akan kita buat.

Hal ke-2, terdapat 1 baris statement yang penulisannya agak unik, yaitu {{"</h2>" | unescape}}.
Statement tersebut maknanya adalah string "</h2>" DIGUNAKAN SEBAGAI PARAMETER dalam
PEMANGGILAN fungsi unescape. Tanda pipe atau | adalah penanda bahwa parameter
dituliskan terlebih dahulu sebelum nama fungsi nya.
*/

// 2. Back End

import (
	"fmt"
	"html/template"
	"net/http"
)

/*
Selanjutnya beberapa fungsi akan dibuat, lalu disimpan dalam template.FuncMap.
Pembuatan fungsi dituliskan dalam bentuk key-value atau hash map. Nama fungsi
sebagai key, dan body fungsi sebagai value.
*/

var funcMap = template.FuncMap{
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"avg": func(n ...int) int {
		var total = 0
		for _, each := range n {
			total += each
		}
		return total / len(n)
	},
}

/*
		*Tipe template.FuncMap SEBENARNYA merupakan alias dari map[string]interface{}
Dalam funcMap di atas, dua buah fungsi disiapkan, unescape() dan avg(). Nantinya
fungsi ini kita gunakan di view.
*/

// buat fungsi main dan kemudian disisipkan fungsi yang telah dibuat di atas ke dalamnya
func main() {
	port := 9000
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var tmpl = template.Must(template.New("view.html").
			Funcs(funcMap).
			ParseFiles("view.html")) // ParseFiles disini berbeda dengan ParseFiles() pada materi sebelumnya

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	/*
		Penjelasan:
		- Sebuah template disipakan dengan nama view.html. Pembuatan instance template
			dilakukan melalui fungsi template.New().
		- Fungsi custom yang telah kita buat, diregistrasikan agar dikenali oleh
			template tersebut. Bisa dilihat pada pemanggilan method Funcs().
		- Setelah itu, lewat method ParseFiles(), view view.html di-parsing. Akan dicari
			dalam file tersebut apakah ada template yang didefinisikan dengan nama view.html.
			Karena di dalam template view TIDAK ADA DEKLARASI template sama sekali
			({{template "namatemplate"}}), maka akan dicari view yang namanya adalah
			view.html. Keseluruhan isi view.html akan dianggap sebagai sebuah template
			dengan nama template adalah nama file itu sendiri.
	*/

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)

}

// 4. Perbadaan Fungsi template.ParseFiles() & Method ParseFiles() Milik *template.Template
/*
Pada kode di atas, pemanggilan template.New() MENGHASILKAN OBJEK bertipe
*template.Template.
Pada chapter B.5. Template: Render Partial HTML Template kita telah belajar mengenai
fungsi template.ParseFiles(), yang fungsi tersebut JUGA MENGEMBALIKAN objek bertipe
*template.Template.

Pada kode di atas, method ParseFiles() yang dipanggil BUKANLAH FUNGSI template.ParseFiles()
yang kita telah pelajari sebelumnya. Meskipun namanya sama, kedua fungsi/method ini
berbeda.

- Fungsi template.ParseFiles(), adalah milik package template. Fungsi ini digunakan
	untuk mem-parsing semua view yang disisipkan sebagai parameter.
- Method ParseFiles(), milik *template.Template, digunakan untuk MEMPARSING
	SEMUA VIEW YANG DISISIPKAN sebagai parameter, lalu diambil HANYA BAGIAN yang
	nama template-nya adalah sama dengan nama template yang sudah di-alokasikan menggunakan
	template.New(). Jika template yang dicari tidak ada, maka akan mencari yang nama
	file-nya sama dengan nama template yang sudah ter-alokasi.
*/
