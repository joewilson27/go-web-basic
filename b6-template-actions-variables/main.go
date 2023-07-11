package main

/*
Template: Actions & Variables
Actions adalah predefined keyword YANG SUDAH DISEDIAKAN Go, biasa dimanfaatkan dalam
pembuatan template.

Sebenarnya pada dua chapter sebelumnya, secara tidak sadar KITA TELAH MENGGUNAKAN
beberapa jenis actions, di antaranya:
- Penggunaan pipeline output. Nilai yang diapit tanda {{ }}, yang nantinya akan
	dimunculkan di layar sebagai output, contohnya: {{"hello world"}}.
- Include template lain menggunakan keyword template, contohnya: {{template "name"}}.

Pada chapter ini, kita akan belajar lebih banyak lagi tentang actions lain yang disediakan Go, juga cara pembuatan dan pemanfaatan variabel pada template view.
*/

import (
	"fmt"
	"html/template"
	"net/http"
)

type Info struct {
	Affiliation string
	Address     string
}

type Person struct {
	Name    string
	Gender  string
	Hobbies []string
	Info    Info
}

/*
Pada kode di atas, dua buah struct disiapkan, Info dan Person (YANG DI MANA struct
Info di-embed ke dalam struct Person). Kedua struct tersebut nantinya akan digunakan
untuk pembuatan objek, yang kemudian object tersebut disisipkan ke dalam view.
*/

// 5. Pengaksesan Property Variabel Objek
func (t Info) GetAffiliationDetailInfo() string { // this method belongs to struct Info
	return "have 31 divisions"
}

func main() {
	// 1. Persiapan
	port := 9000
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var person = Person{
			Name:    "Bruce Wayne",
			Gender:  "male",
			Hobbies: []string{"Reading Books", "Travelling", "Shopping"},
			Info:    Info{"Wayne Enterprises", "Gotham City"},
		}

		var tmpl = template.Must(template.ParseFiles("view.html")) // parsing template
		/*
			INGAT, KALAU FILE HTML ADA DI DALAM FOLDER, HARUS DI JOIN DULU DENGAN FOLDER DIMANA TEMPLATE HTML NYA
			jika menggunakan path.Join ~ nevermind.
		*/
		if err := tmpl.Execute(w, person); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	/*
		Pada route handler / di atas, variabel objek person dibuat, lalu disisipkan
		sebagai data pada view view.html yang sebelumya sudah diparsing.

		PERLU DIKETAHUI, ketika data yang disisipkan ke view BERBENTUK map,
		maka key (yang nantinya akan menjadi nama variabel) BOLEH DITULISKAN dalam huruf kecil.
		Sedangkan jika berupa variabel objek struct, maka property HARUS DITULISKAN PUBLIC
		(huruf pertama kapital).

			*Data yang disisipkan ke view, jika tipe nya adalah struct, maka HANYA PROPERTIES
			ber-modifier public (ditandai dengan huruf kapital di awal nama property)
			YANG BISA DIAKSES dari view.

	*/

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)

	// 2. Pipeline Output & Komentar
	/*
		Actions pertama yang akan kita coba terapkan adalah pipeline output, MENAMPILKAN output ke layar. Caranya cukup mudah,
		cukup dengan menuliskan apa yang ingin ditampilkan di layar dengan diapit tanda {{ }}
		(bisa berupa variabel yang dilempar dari back end, bisa juga literal string).
	*/

	// 3. Membuat & Menampilkan Isi Variabel
	/*
		Cara membuat variabel dalam template adalah dengan mendeklarasikannya menggunakan operator :=,
		DENGAN KETENTUAN nama variabel HARUS DIAWALI dengan tanda dollar $.
	*/

	// 4. Perulangan
	/*
		Actions range DIGUNAKAN UNTUK MELAKUKAN PERULANGAN pada template view. Keyword ini bisa diterapkan
		pada tipe data map atau array. Cara penggunaannya sedikit berbeda dibanding penggunaan range pada Go
	*/

	// 5. Pengaksesan Property Variabel Objek
	/*
		Cara mengakses property sebuah variabel objek BISA DILAKUKAN lewat notasi titik ., dengan ketentuan property
		tersebut BERMODIFIER public.

		Sedangkan untuk pengaksesan method, caranya juga sama, hanya saja TIDAK PERLU DITULISKAN tanda kurung method-nya.
		jadi nanti cukup panggil GetAffiliationDetailInfo bukan GetAffiliationDetailInfo()
	*/

	// 6. Penggunaan Keyword with Untuk Mengganti Scope Variabel Pada Suatu Blok
	/*
		SECARA DEFAULT current scope di template view adalah DATA YANG DILEMPAR back end.
		Scope current objek BSIA DIGANTI dengan menggunakan keyword with, sehingga nantinya untuk mengakses sub-property
		variabel objek (seperti .Info.Affiliation), bisa tidak dilakukan dari objek terluar.
			*Current scope yg dimaksud di sini adalah seperti object this ibarat bahasa pemrograman lain.
		Sebagai contoh property Info yang merupakan variabel objek. Kita bisa menentukan scope suatu block adalah mengikuti variabel objek
		tersebut.
		Pada contoh di view, sebuah blok ditentukan scope-nya adalah Info. Maka di dalam blok kode tersebut,
		untuk mengakses sub property-nya (Address, Affiliation, dan GetAffiliationDetailInfo), TIDAK PERLU dituliskan
		dari objek terluar, cukup langsung nama property-nya. Sebagai contoh .Address di atas merujuk ke variabel .Info.
	*/

	// 7. Seleksi Kondisi
	/*
		Seleksi kondisi juga bisa dilakukan pada template view. Keyword actions yang digunakan adalah if dan eq (equal atau sama dengan).

		Untuk seleksi kondisi dengan jumlah kondisi lebih dari satu, bisa gunakan else if.

			{{if pipeline}}
					a
			{{else if pipeline}}
					b
			{{else}}
					c
			{{end}}

		Untuk seleksi kondisi yang kondisinya adalah bersumber dari variabel bertipe bool, maka langsung saja tulis TANPA MENGGUNAKAN eq.
		Jika kondisi yang DIINGINKAN adalah KEBALIKAN dari nilai variabel,
		maka gunakan ne. Contohnya bisa dilihat pada kode berikut.


		misal .IsTrue value bool, maka tulis if tanpa eq
			{{if .IsTrue}}
					<p>true</p>
			{{end}}

			{{isTrue := true}}

			{{if isTrue}}
					<p>true</p>
			{{end}}

		ini dengan menggunakan eq
			{{if eq isTrue}}
					<p>true</p>
			{{end}}

		Jika kita ingin kondisi dibalikkan (negasi), gunakan ne
			{{if ne isTrue}}
					<p>not true (false)</p>
			{{end}}


	*/
}
