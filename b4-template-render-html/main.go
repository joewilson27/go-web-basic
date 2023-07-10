package main

/*
Template: Render HTML Template
Pada bagian ini kita akan belajar BAGAIMANA CARA RENDER file template ber-tipe HTML,
untuk ditampilkan pada browser.

Terdapat banyak jenis template pada Go, yang AKAN KITA PAKAI adalah template HTML.
Package html/template berisi banyak sekali fungsi untuk kebutuhan rendering dan parsing file template jenis ini.
*/

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func main() {
	// 2. Backend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var filepath = path.Join("views", "index.html") // path.Join menggabungkan views dan index.html menghasilkan views/index.html
		var tmpl, err = template.ParseFiles(filepath)
		/*
			Sedangkan template.ParseFiles(), digunakan untuk PARSING FILE template,
			dalam contoh ini file view/index.html. Fungsi ini mengembalikan 2 data,
			yaitu hasil dari proses parsing yang bertipe *template.Template,
			dan informasi error jika ada.
		*/
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			/*
				Fungsi http.Error() digunakan untuk MENANDAI response (http.ResponseWriter)
				bahwa terjadi error, dengan kode error dan pesan error bisa ditentukan. Pada
				contoh di atas yang digunakan adalah 500 - internal server error yang
				direpresentasikan oleh variabel http.StatusInternalServerError.
			*/
			return
		}

		var data = map[string]interface{}{
			"title": "Learning Golang Web",
			"name":  "Ironman",
		}

		err = tmpl.Execute(w, data)
		/*
			Method Execute() milik *template.Template, digunakan untuk MENYISIPKAN DATA
			pada template, untuk kemudian ditampilkan ke browser. Data bisa disipkan DALAM BENTUK
			struct, map, atau interface{}
			- Jika dituliskan dalam bentuk map, maka key akan menjadi nama variabel dan
				value menjadi nilainya
			- Jika dituliskan dalam bentuk variabel objek cetakan struct, nama property
				akan menjadi nama variabel
			Pada contoh di atas, data map yang berisikan key title dan name disisipkan
			ke dalam template yang sudah di parsing.
		*/

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		/*
			Package path (yg kita import) BERISIKAN BANYAK fungsi yang berhubungan dengan
			lokasi folder atau path, yang salah satu di antaranya adalah fungsi path.Join().
			Fungsi ini digunakan untuk MENGGABUNGKAN folder atau file atau keduanya menjadi
			sebuah path, dengan separator relatif terhadap OS yang digunakan.

			*Separator yang digunakan oleh path.Join() adalah \ untuk wind*ws dan / untuk un*x.

			Contoh penerapan path.Join() bisa dilihat di kode di atas, views di-join dengan
			index.html, menghasilkan views/index.html.
		*/

	})

	// routing static assets
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)

	// 3. Front End
	/*
		Untuk menampilkan variabel yang disisipkan ke dalam template, gunakan notasi {{.namaVariabel}}.
		Pada contoh di atas, data title dan name yang dikirim dari back end ditampilkan.

		Tanda titik "." pada {{.namaVariabel}} menerangkan bahwa variabel tersebut diakses dari current scope.
		Dan current scope default adalah data map atau objek yang dilempar back end.
	*/

	// 5. Static File CSS
	/*
		Pada views/index.html, include-kan file css.

			<link rel="stylesheet" href="/static/site.css" />

			perhatikan, langsung tulis /static/site.css, tidak perlu dengan nama directory 'assets' karena otomatis
			dia mengarah ke assets yg sudah di set pada http.Handle

	*/
}
