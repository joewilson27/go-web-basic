package main

/*
Template: Render Partial HTML Template
Satu buah halaman yang berisikan html, BISA TERBENTUK DARI BANYAK TEMPLATE html (parsial).
Pada chapter ini kita akan belajar bagaimana membuat, mem-parsing, dan me-render semua file tersebut.

Ada beberapa metode yang bisa digunakan, dari ke semuanya akan kita bahas 2 di antaranya, yaitu:
- Menggunakan fungsi template.ParseGlob().
- Menggunakan fungsi template.ParseFiles().
*/

import (
	"fmt"
	"html/template"
	"net/http"
)

type M map[string]interface{}

/*
Tipe M merupakan ALIAS dari map[string]interface{}, DISIAPKAN untuk MEMPERSINGKAT
penulisan tipe map tersebut. Pada pembahasan-pembahasan selanjutnya kita AKAN BANYAK
MENGGUNAKAN tipe ini.
*/

func main() {
	// 2. Back End
	//var tmpl, err = template.ParseGlob("views/*") // parsing file html di awal
	/*
		fungsi template.ParseGlob() dipanggil, dengan parameter adalah pattern path "views/*".
		Fungsi ini digunakan untuk MEMPARSING SEMUA FILE yang MATCH dengan pattern yang
		ditentukan, dan fungsi ini mengembalikan 2 objek: *template.Template & error

			*Pattern path pada fungsi template.ParseGlob() nantinya akan di proses oleh filepath.Glob()

		Proses PARSING SEMUA FILE html dalam folder views DILAKUKAN DI AWAL,
		AGAR KETIKA MENGAKSES rute-tertentu-yang-menampilkan-html,
		tidak terjadi proses parsing lagi (INTINYA DI LOAD DI AWAL BIAR FILE HTML KE
		LOAD DARI AWAL).

			*Parsing semua file menggunakan template.ParseGlob() yang dilakukan di luar handler,
			 TIDAK DIREKOMENDASIKAN dalam fase development. Karena akan mempersulit testing html.
			 Lebih detailnya akan dibahas di bagian bawah.

	*/
	// if err != nil {
	// 	panic(err.Error())
	// 	return
	// }

	// http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
	// 	var data = M{"name": "Batman"} // bikin map dengan alias
	// 	err = tmpl.ExecuteTemplate(w, "index", data)
	/*
		Karena SEMUA FILE HTML SUDAH DI PARSING DI AWAL, maka untuk render template tertentu
		CUKUP DENGAN MEMANGGIL method ExecuteTemplate(), dengan menyisipkan
		3 parameter berikut:
		- Parameter ke-1, objek http.ResponseWriter --> w
		- Parameter ke-2, nama template --> "index" => {{define "index"}}
		- Parameter ke-3, data --> data

		*Nama template BUKANLAH nama file. Setelah masuk ke bagian front-end,
		 akan diketahui apa yang dimaksud dengan nama template.
	*/

	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })

	// http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
	// 	var data = map[string]interface{}{
	// 		"name": "Batman",
	// 	}
	// ingat, ketika define variable map dengan value interface kosong, harus di tutup koma dalam valuenya

	// 	err = tmpl.ExecuteTemplate(w, "about", data)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })
	/*
		Kedua rute tersebut sama, PEMBEDANYA adalah TEMPLATE yang di-render.
		Rute /index me-render template bernama index, dan rute /about me-render template
		bernama about.

		Karena SEMUA file html SUDAH DI PARSING DI AWAL, maka untuk render template tertentu
		cukup dengan memanggil method ExecuteTemplate(), dengan menyisipkan
		3 parameter berikut:
			- Parameter ke-1, objek http.ResponseWriter
			- Parameter ke-2, nama template
			-	Parameter ke-3, data

		*Nama template BUKANLAH nama file. Setelah masuk ke bagian front-end,
			 akan diketahui apa yang dimaksud dengan nama template.
	*/

	// tester purpose
	// http.HandleFunc("/tester-page", func(w http.ResponseWriter, r *http.Request) {
	// 	err = tmpl.ExecuteTemplate(w, "_message", nil)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })

	// 3. Front End
	/*
		Ada 4 buah template yang harus kita siapkan satu per satu.

		index.html
		Berikut adalah penjelasannya kode index.html.
		- Statement {{define "index"}}, digunakan untuk MENDEFINISIKAN nama template.
			Semua blok kode setelah statement tersebut (batasnya adalah hingga statement {{end}})
			adalah milik template dengan nama index. keyword define digunakan dalam penentuan nama template.
		- Statement {{template "_header"}} artinya adalah template bernama _header
			di-include ke bagian itu. keyword template DIGUNAKAN untuk INCLUDE template lain.
		- Statement {{template "_message"}}, sama seperti sebelumnya,
			template bernama _message akan di-include.
		- Statement {{.name}} akan memunculkan data, name, yang data ini sudah
			disisipkan oleh back end pada saat rendering.
		- Statement {{end}} adalah penanda batas akhir pendefinisian template.

		Template _header.html
			*Nama file BISA DITULIS dengan diawali karakter underscore atau _.
			Pada chapter ini, nama file yang diawali _ kita ASUMSIKAN SEBAGAI
			template parsial, template yang nantinya di-include-kan ke template utama.

	*/

	// 3.6. Parsing Banyak File HTML Menggunakan template.ParseFiles()
	/*
		Metode parsing menggunakan template.ParseGlob() MEMILIKI KEKURANGAN yaitu SANGAT TERGANTUNG
		terhadap pattern path yang digunakan. Jika dalam suatu proyek terdapat sangat banyak
		file html dan folder, sedangkan hanya beberapa yang digunakan, pemilihan pattern path
		yang kurang tepat akan menjadikan file lain ikut ter-parsing dengan sia-sia.

		Dan juga, karena statement template.ParseGlob() dieksekusi diluar handler, maka KETIKA ADA PERUBAHAN
		pada salah satu view, lalu halaman di refresh, OUTPUT yang dihasilkan AKAN TETAP SAMA.
		SOLUSI dari masalah ini adalah dengan MEMANGGIL template.ParseGlob() DI TIAP handler rute-rute
		yang diregistrasikan.

			*Best practices yang bisa diterapkan, ketika environment adalah production, maka tempatkan
			 template.ParseGlob() DI LUAR (sebelum) handler. Sedangkan pada environment development,
			 taruh template.ParseGlob() DI DALAM masing-masing handler. Gunakan seleksi kondisi
			 untuk mengakomodir skenario ini.

		Alternatif METODE LAIN yang BISA DIGUNAKAN, yang lebih efisien, adalah dengan memanfaatkan
		fungsi template.ParseFiles(). Fungsi ini selain bisa digunakan untuk PARSING SATU BUAH FILE
		saja (seperti yang sudah dicontohkan pada chapter sebelumnya), bisa digunakan untuk parsing
		BANYAK FILE

		Mari kita praktekan. Ubah handler rute /index dan /about. Gunakan template.ParseFiles() dengan
		isi parameter (variadic-banyak parameter) adalah path dari file-file html yang akan
		dipergunakan di masing-masing rute. Lalu hapus statement template.ParseGlob() (yg sebelumnya
		untuk memparsing template di awal)
	*/
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		var data = M{"name": "Batman"}
		/*
			Wrap fungsi template.ParseFiles() DALAM template.Must(). Fungsi ini berguna untuk
			DETEKSI ERROR pada saat membuat instance *template.Template baru atau
			ketika sedang mengolahnya. Ketika ada error, panic dimunculkan.
		*/
		var tmpl = template.Must(template.ParseFiles(
			"views/index.html",
			"views/_header.html",
			"views/_message.html",
		))

		var err = tmpl.ExecuteTemplate(w, "index", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]interface{}{
			"name": "Batman",
		}
		var tmpl = template.Must(template.ParseFiles(
			"views/about.html",
			"views/_header.html",
			"views/_message.html",
		))

		var err = tmpl.ExecuteTemplate(w, "about", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// tester
	http.HandleFunc("/tester-page", func(w http.ResponseWriter, r *http.Request) {
		var tmpl = template.Must(template.ParseFiles(
			"views/_message.html",
		))

		var err = tmpl.ExecuteTemplate(w, "_message", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}
