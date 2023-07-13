package main

/*
Simple Configuration
Dalam development, PASTI BANYAK SEKALI variabel dan konstanta yang DIPERLUKAN.
Mulai dari variabel yang dibutuhkan untuk start server seperti port, timeout, hingga variabel
global dan variabel shared lainnya.

Pada chapter ini, kita akan belajar CARA MEMBUAT config file modular.
*/

import (
	"fmt"
	"log"
	"net/http"
	"simple-configuration/conf"
	"time"
)

type CustomMux struct {
	http.ServeMux
}

func (c CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if conf.Configuration().Log.Verbose {
		log.Println("Incoming request from", r.Host, "accessing", r.URL.String())
	}
	/*
		Bisa dilihat dalam method ServeHTTP() ini, ada pengecekan salah satu konfigurasi,
		yaitu Log.Verbose. Cara pengaksesannya cukup mudah, yaitu lewat fungsi Configuration()
		milik package conf yang telah di-import.
	*/

	c.ServeMux.ServeHTTP(w, r)
}

func main() {
	router := new(CustomMux)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	router.HandleFunc("/howareyou", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("How are you?"))
	})

	server := new(http.Server)
	/*
		Objek baru bernama server telah dibuat dari struct http.Server.
		Untuk start server cukup panggil method ListenAndServe() milik objek tersebut.

		Dengan memanfaatkan struct ini, KITA BISA meng-CUSTOM beberapa konfigurasi default pada Go web server.
		Di antaranya seperti ReadTimeout dan WriteTimeout.
	*/
	server.Handler = router
	server.ReadTimeout = conf.Configuration().Server.ReadTimeout * time.Second
	server.WriteTimeout = conf.Configuration().Server.WriteTimeout * time.Second
	server.Addr = fmt.Sprintf(":%d", conf.Configuration().Server.Port)
	/*
		Pada kode di atas bisa kita lihat, ada 4 buah properti milik server di-isi.
		- server.Handler. Properti ini wajib di isi dengan custom mux yang dibuat.
		- server.ReadTimeout. Adalah timeout ketika memproses sebuah request. Kita isi dengan nilai dari configurasi.
		- server.WriteTimeout. Adalah timeout ketika memproses response.
		- server.Addr. Port yang digunakan web server pada saat start.
	*/

	if conf.Configuration().Log.Verbose {
		log.Printf("Starting server at %s \n", server.Addr)
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	/*
		Coba ubah konfigurasi pada config.json nilai log.verbose menjadi false.
		Lalu restart aplikasi, maka log (di CLI) tidak muncul.
	*/

	// 5. Kekurangan Konfigurasi File
	/*
		• Tidak mendukung komentar
		Komentar sangat penting karena untuk aplikasi besar yang konfigurasi item-nya sangat banyak - akan susah untuk
		dipahami. Sebenarnya perihal ini bisa di-resolve MENGGUNAKAN jenis konfigurasi lain seperti YAML, .env, atau lainnya.

		• Nilai konfigurasi harus diketahui di awal
		Kita HARUS TAHU SEMUA value tiap-tiap konfigurasi terlebih dahulu, dan dituliskan ke file,
		sebelum aplikasi di-up. Dari sini akan sangat susah jika misal ada beberapa konfigurasi yang kita tidak tau
		nilainya tapi tau cara pengambilannya.
		Contohnya pada beberapa kasus, seperti di AWS, database server yang di-setup secara automated akan meng-generate
		connection string yang host-nya BISA BERGANTI-GANTI tiap start-up, dan tidak hanya itu,
		bisa saja username, password dan lainnya juga tidak statis.
		Dengan ini akan sangat susah jika kita harus cari terlebih dahulu value konfigurasi tersebut
		untuk kemudian dituliskan ke file. Memakan waktu dan kurang baik dari banyak sisi.

		• Tidak terpusat
		Dalam pengembangan aplikasi, banyak konfigurasi yang nilai-nya akan didapat lewat jalan lain,
		seperti environment variables atau command arguments.
		Akan lebih mudah jika hanya ada satu sumber konfigurasi saja untuk dijadikan acuan.

		• Statis (tidak dinamis)
		Konfigurasi umumnya dibaca hanya jika diperlukan. Penulisan konfigurasi dalam file membuat proses pembacaan
		file harus dilakukan di awal, harus kemudian kita bisa ambil nilai konfigurasi dari data yang sudah ada di memori.
		Hal tersebut memiliki beberapa konsekuensi, untuk aplikasi yang di-manage secara automated, sangat mungkin
		adanya perubahan nilai konfigurasi. Dari sini berarti pembacaan konfigurasi file tidak boleh hanya dilakukan
		di awal saja. Tapi juga tidak boleh dilakukan di setiap waktu, karena membaca file itu ada cost-nya dari prespektif I/O.

		• Solusi
		Akan di bahas pada chapter berikutnya
	*/
}
