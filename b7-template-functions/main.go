package main

/*
Template: Functions

Go MENYEDIAKAN BEBERAPA predefiend function yang bisa digunakan dalam file template.
Pada chapter ini kita akan membahas beberapa di antaranya beserta cara penggunaannya.
Cara pemanggilan fungsi atau method sebuah objek pada file template sedikit berbeda dibanding seperti pada chapter sebelumnya.
*/

import (
	"fmt"
	"html/template"
	"net/http"
)

type Superhero struct {
	Name    string
	Alias   string
	Friends []string
}

func (s Superhero) SayHello(from string, message string) string {
	return fmt.Sprintf("%s said: \"%s\"", from, message)
}

/*
Struct Superhero di atas nantinya digunakan untuk membuat objek yang kemudian
disisipkan ke template view.
*/

func main() {
	// 1. Persiapan
	port := 9000
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var person = Superhero{
			Name:    "Bruce Wayne",
			Alias:   "Batman",
			Friends: []string{"Superman", "Flash", "Green Lantern"},
		}

		var tmpl = template.Must(template.ParseFiles("view.html"))
		if err := tmpl.Execute(w, person); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)

}

// 2. Fungsi Escape String
/*
Fungsi pertama yang akan kita bahas adalah html. Fungsi ini digunakan untuk meng-ESCAPE
string. Agar lebih mudah dipahami silakan praktekan kode di bawah ini.
Selain fungsi html, ada juga beberapa fungsi lain yang sudah disediakan oleh Go.
- Fungsi js digunakan untuk meng-escape string javascript
- Fungsi urlquery digunakan untuk meng-escape string url query
*/

// 3. Fungsi Operator Perbandingan
/*
Pada chapter sebelumnya telah dibahas bagaimana penggunaan operator ne pada actions if.
eq dan ne adalah contoh dari fungsi operator perbandingan. Jika digunakan pada seleksi
kondisi yang nilai kondisinya bertipe bool, maka cukup dengan menuliskannya setelah
operator, contohnya.

Jika nilai kondisinya merupakan perbandingan, maka nilai yang dibandingkan harus
dituliskan, sebagai contoh di bawah ini adalah seleksi kondisi memanfaatkan
operator gt untuk deteksi apakah nilai di atas 60.
gt merupakan kependekan dari greater than

Operator	Penjelasan																								Analogi
eq				equal, sama dengan																				a == b
ne				not equal, tidak sama dengan															a != b
lt				lower than, lebih kecil																		a < b
le				lower than or equal, lebih kecil atau sama dengan					a <= b
gt				greater than, lebih besar																	a > b
ge				greater than or equal, lebih besar atau sama dengan				a >= b
*/

// 4. Pemanggilan Method
/*
Cara memanggil method yang disisipkan ke view sama dengan cara pemanggilan fungsi, hanya saja perlu ditambahkan tanda titik . (menyesuaikan scope variabelnya).
*/

// 6. Fungsi len dan index
/*
Kegunaan dari fungsi len seperti yang sudah diketahui adalah untuk MENGHITUNG
jumlah elemen. Sedangkan fungsi index digunakan jika ELEMEN TERTENTU INGIN DIAKSES.
*/

// 7. Fungsi Operator Logika
/*
Selain fungsi operator perbandingan, terdapat juga operator logika or, and, dan not.
Cara penggunaannya adalah dengan dituliskan setelah actions if atau elseif, sebagai
fungsi dengan parameter adalah nilai yang ingin dibandingkan.
	*Fungsi not ekuivalen dengan ne

*/
