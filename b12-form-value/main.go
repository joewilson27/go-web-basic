package main

/*
Form Value
Pada step ini kita belajar cara submit data, dari form di layer front end, ke back end.
*/

// 1. Front End
/*
siapkan sebuah file template view view.html. Pada file ini PERLU DIDEFINISIKAN 2 buah
template, yaitu form dan result. Template pertama (form) dijadikan landing page
program, isinya beberapa inputan untuk submit data.

Aksi dari form di atas adalah /process, yang di mana url tersebut nantinya akan
MENGEMBALIKAN output berupa html hasil render template result. Silakan tulis template
result berikut dalam view.html (jadi file view ini berisi 2 buah template).
*/

// 2. Back End
import (
	"fmt"
	"html/template"
	"net/http"
)

/*
Handler route / dibungkus dalam fungsi bernama routeIndexGet. Di dalamnya, template
form dalam file template view.html akan di-render ke view. Request dalam handler
ini HANYA DIBATASI untuk method GET saja, request dengan method lain akan MENGHASILKAN
response 400 Bad Request.
*/
func routeIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

/*
Fungsi routeSubmitPost yang merupakan handler route /process, berisikan PROSES YANG MIRIP
seperti handler route /, yaitu parsing view.html untuk di ambil template result-nya.
Selain itu, pada handler ini ADA PROSES PENGAMBILAN DATA yang DIKIRIM DARI FORM
ketika di-submit, untuk kemudian disisipkan ke template view.
*/
func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("view.html"))

		/*
			Method ParseForm() pada statement r.ParseForm() BERGUNA UNTUK PARSING form
			data YANG DIKIRIM DARI VIEW, sebelum akhirnya bisa diambil data-datanya.
			Method tersebut MENGEMBALIKAN data error jika proses PARSING GAGAL
			(kemungkinan karena data yang dikirim ada yang tidak valid).
		*/
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var name = r.FormValue("name")
		/*
			PENGAMBILAN DATA yang dikirim dari view DILAKUKAN LEWAT method FormValue().
			r.FormValue("name") akan mengembalikan data inputan name
			dari inputan <input name="name" />
		*/
		var message = r.Form.Get("message")
		/*
			Selain lewat method FormValue(), PENGAKSESAN DATA juga BISA DILAKUKAN
			dengan cara mengakses property Form terlebih dahulu, kemudian mengakses
			method Get(). Contohnya seperti r.Form.Get("message"), yang akan menghasilkan
			data inputan message. Hasil dari kedua cara di atas adalah sama.
		*/

		// setelah value form ditangkap, kemudian di tampung variable data bertipe map[string]string
		// yang kemudian disisipkan ke view
		var data = map[string]string{"name": name, "message": message}

		// sisipkan data ke view
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func main() {
	port := 9000
	/*
		siapkan 2 route pada main:
		- Route / adalah landing page, menampilkan form input.
		- Route /process sebagai action dari form input, menampilkan text.
	*/
	http.HandleFunc("/", routeIndexGet)
	http.HandleFunc("/process", routeSubmitPost)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)
}
