package main

/*
AJAX JSON Payload
Teknik request Form Data digunakan salah satu nya pada request submit lewat <form />.
Pada chapter ini, kita TIDAK AKAN MENGGUNAKAN cara submit lewat form, melainkan menggunakan teknik AJAX
(Asynchronous JavaScript And XML), dengan payload ber-tipe JSON.

Perbedaan antara kedua jenis request tersebut adalah pada ISI HEADER Content-Type, dan bentuk informasi
dikirimkan. Secara default, request lewat <form />, CONTENT TYPE-nya adalah application/x-www-form-urlencoded.
Data dikirimkan dalam BENTUK QUERY STRING (key-value) seperti id=n001&nama=bruce.

	*Ketika di form DITAMBAHKAN ATRIBUT enctype="multipart/form-data", maka content type BERUBAH MENJADI multipart/form-data.

Request Payload JSON SEDIKIT BERBEDA, Content-Type berisikan application/json, dan data disisipkan dalam Body DALAM BENTUK JSON string.
*/

// 2. Front End - HTML
/*
persiapkan view.html
*/

// 3. Front End - HTML
/*
siapkan payload utk kirim data via jquery ajax

Value semua inputan diambil lalu dimasukkan dalam sebuah objek lalu di stringify (AGAR MENJADI JSON string)
*/

// 4. Back End
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		/*
			Isi payload didapatkan dengan cara men-decode body request (r.Body).
			Proses decoding TIDAK DILAKUKAN menggunakan json.Unmarshal() MELAINKAN LEWAT
			json decoder, karena akan lebih efisien untuk jenis kasus seperti ini.

				*Gunakan json.Decoder jika data adalah STREAM io.Reader.
				Gunakan json.Unmarshal() untuk decode data sumbernya sudah ada di memory.

		*/
		payload := struct {
			Name   string `json:"namex"`  // namex dari payload
			Age    int    `json:"age"`    // age dari payload
			Gender string `json:"gender"` // gender dari payload
		}{}
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf(
			"hello, my name is %s. I'm %d year old %s",
			payload.Name,
			payload.Age,
			payload.Gender,
		)
		w.Write([]byte(message))
		return
	}

	http.Error(w, "Only accept POST request", http.StatusBadRequest)
}

func main() {
	port := 9000
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/save", handleSave)

	http.Handle("/static/", // maka nantinya SEMUA REQUEST yang di awali dengan /static/ akan diarahkan ke sini
		http.StripPrefix("/static/", // http.StripPrefix() hanya digunakan untuk membungkus actual handler.
			http.FileServer(http.Dir("assets")))) // Fungsi http.FileServer() mengembalikan objek ber-tipe http.Handler
	// Semua konten, entah file ataupun folder, yang ada di dalam folder assets akan di proses dalam handler.

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)

}
