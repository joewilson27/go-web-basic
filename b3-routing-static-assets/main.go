package main

/*
Routing Static Assets
Pada bagian ini kita akan belajar bagaimana cara routing static assets atau
static contents. Seperti file css, js, gambar, umumnya dikategorikan
sebagai static assets.
*/

import (
	"fmt"
	"net/http"
)

func main() {
	// 2. Routing
	http.Handle("/static/",
		http.StripPrefix("/static/",
			// ACTUAL HANDLER FOR THIS ROUTE
			http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}

/*
Syarat yang DIBUTUHKAN untuk routing static assets masih sama dengan routing handler,
yaitu perlu DIDEFINISIKAN RUTE-nya dan handler-nya. Hanya saja pembedanya,
dalam routing static assets yang digunakan adalah http.Handle(), BUKAN
http.HandleFunc().
1. Rute terpilih adalah /static/, maka nantinya SEMUA REQUEST yang DI AWALI
	 dengan /static/ AKAN DIARAHKAN ke sini. Registrasi rute menggunakan http.Handle()
	 adalah BERBEDA DENGAN routing menggunakan http.HandleFunc(), lebih jelasnya akan
	 ada sedikit penjelasan pada chapter lain.
2. Sedang untuk handler-nya bisa di-lihat, ada pada parameter ke-2 yang isinya
	 statement http.StripPrefix(). SEBENARNYA ACTUAL handler nya berada pada
	 http.FileServer(). Fungsi http.StripPrefix() HANYA DIGUNAKAN UNTUK MEMBUNGKUS
	 actual handler.

Fungsi http.FileServer() MENGEMBALIKAN objek ber-tipe http.Handler.
Fungsi ini berguna untuk men-serve semua http request, dengan konten
yang didefinisikan pada parameter. Pada konteks ini yang di-maksud
adalah http.Dir("assets"). Semua konten, entah file ataupun folder,
yang ada di dalam folder assets akan di proses dalam handler.

Jalankan main.go, lalu test hasilnya di browser http://localhost:9000/static/.

// 3. Penjelasan
http.Handle("/", http.FileServer(http.Dir("assets")))

Jika dilihat pada struktur folder yang sudah di-buat, di dalam folder assets terdapat
file bernama site.css. Maka dengan bentuk routing pada contoh sederhana di atas,
request ke /site.css akan diarahkan ke path ./site.css (relatif dari folder assets).
Permisalan contoh lainnya:
- Request ke /site.css mengarah path ./site.css relatif dari folder assets
- Request ke /script.js mengarah path ./script.js relatif dari folder assets
- Request ke /some/folder/test.png mengarah path ./some/folder/test.png relatif dari folder assets
- ... dan seterusnya

*Fungsi http.Dir() BERGUNA UNTUK ADJUSTMENT path parameter. Separator dari path yang
di-definisikan akan otomatis di-konversi ke path separator sesuai sistem operasi.

Contoh selanjutnya, silakan perhatikan kode berikut.

http.Handle("/static", http.FileServer(http.Dir("assets")))
Hasil dari routing:
- Request ke /static/site.css mengarah ke ./static/site.css relatif dari folder assets
- Request ke /static/script.js mengarah ke ./static/script.js relatif dari folder assets
- Request ke /static/some/folder/test.png mengarah ke ./static/some/folder/test.png relatif dari folder assets
- ... dan seterusnya

Terlihat bahwa rute yang didaftarkan juga akan digabung dengan
path destinasi file yang dicari, dan ini menjadikan path tidak valid.
File site.css berada pada path assets/site.css, sedangkan dari routing
di atas pencarian file mengarah ke path assets/static/site.css.
Di sinilah kegunaan dari fungsi http.StripPrefix().

Fungsi http.StripPrefix() ini BERGUNA UNTUK MENGHAPUS prefix dari endpoint
yang di-request. Pada contoh paling atas, request ke url yang di awali dengan
/static/ hanya akan di ambil url setelahnya.

- Request ke /static/site.css menjadi /site.css
- Request ke /static/script.js menjadi /script.js
- Request ke /static/some/folder/test.png menjadi /some/folder/test.png
- ... dan seterusnya

Routing static assets menjadi valid, karena file yang di-request akan cocok dengan
path folder dari file yang di request.

*/
