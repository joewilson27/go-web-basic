package main

import (
	"fmt"
	"net/http"
)

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	var message = "Welcome"
	w.Write([]byte(message)) // method Write() menerima argumen type []byte, maka perlu dilakukan casting dari string ke []byte
	/*
		Method Write() milik parameter pertama (yang bertipe http.ResponseWrite),
		digunakan untuk meng-OUTPUT-kan nilai balik data.
		Argumen method adalah data yang ingin dijadikan output,
		ditulis dalam bentuk []byte.
	*/
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	var message = "Hello world!"
	w.Write([]byte(message))
}

/*
Fungsi dengan struktur di atas DIPERLUKAN oleh http.HandleFunc untuk KEPERLUAN PENANGANAN request ke rute yang ditentukan
*/

func main() {
	// 1. Pembuatan Aplikasi
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/index", handlerIndex)
	http.HandleFunc("/hello", handlerHello)
	/*
		Fungsi http.HandleFunc() digunakan untuk routing. Parameter pertama adalah rute dan parameter ke-2 adalah handler-nya.
	*/

	// var address = "localhost:9000"
	// fmt.Printf("server started at %s\n", address)
	// err := http.ListenAndServe(address, nil) // Fungsi ini mengembalikan nilai balik ber-tipe error. Jika proses pembuatan web server baru gagal, maka kita bisa mengetahui root-cause nya apa.
	/*
		Fungsi http.ListenAndServe() digunakan MEMBUAT SEKALIGUS start server baru, dengan parameter pertama
		adalah alamat web server yang diiginkan (bisa diisi host, host & port, atau port saja).
		Parameter kedua merupakan object mux atau multiplexer.

		*Dalam chapter ini kita menggunakan default mux yang sudah disediakan oleh Go, jadi untuk parameter ke-2 cukup isi dengan nil.
	*/
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	/*
		Penggunaan http.HandleFunc()
		Fungsi ini digunakan untuk routing, menentukan aksi dari sebuah url tertentu ketika diakses
		(di sini url tersebut kita sebut sebagai rute/route).

		Penjelasan Mengenai Handler
		pada contoh disini, fungsi handler adalah handlerIndex dan handlerHello,
		Route handler atau handler atau parameter kedua fungsi http.HandleFunc(), adalah sebuah fungsi
		dengan ber-skema func (ResponseWriter, *Request).

		Penggunaan http.ListenAndServe()
		Fungsi ini digunakan untuk membuat web server baru. Pada contoh yang telah dibuat, web server di-start pada port 9000
		(bisa dituliskan dalam bentuk localhost:9000 atau cukup :9000 saja).
		Fungsi http.ListenAndServe() bersifat blocking, menjadikan semua statement setelahnya tidak akan dieksekusi, sebelum di-stop.
	*/

	// 2. Web Server Menggunakan http.Server
	/*
		Selain menggunakan http.ListenAndServe(), ADA CARA LAIN YANG BISA DITERAPKAN untuk start web server,
		yaitu dengan memanfaatkan struct http.Server.

		Kode di bagian start server yang sudah kita buat sebelumnya, jika diubah ke cara ini,
		kurang lebih menjadi seperti berikut.
	*/
	var address = ":9000"
	fmt.Printf("server started at %s\n", address)

	server := new(http.Server)
	server.Addr = address
	err := server.ListenAndServe()
	/*
		Informasi host/port perlu dimasukan dalam property .Addr MILIK objek server.
		Lalu dari objek tersebut panggil method .ListenAndServe() untuk start web server.
	*/
	if err != nil {
		fmt.Println(err.Error())
	}
	/*
		KELEBIHAN menggunakan http.Server salah satunya adalah KEMAMPUAN UNTUK MENGUBAH
		beberapa konfigurasi default web server Go.

		Contoh, pada kode berikut, timeout untuk read request dan write request di ubah
		menjadi 10 detik.

			server.ReadTimeout = time.Second * 10
			server.WriteTimeout = time.Second * 10

		Ada banyak lagi property dari struct http.Server ini, yang pastinya akan dibahas pada
		pembahasan-pembahasan selanjutnya.
	*/
}
