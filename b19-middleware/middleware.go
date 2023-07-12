package main

import "net/http"

const USERNAME = "ironman"
const PASSWORD = "secret"

// ini adalah fungsi Auth() pada materi sebelumnya
func MiddlewareAuth(next http.Handler) http.Handler {
	// fungsi http.HandlerFunc built-in Go mengembalikan nilai http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Write([]byte(`something went wrong`))
			return
		}

		isValid := (username == USERNAME) && (password == PASSWORD)
		if !isValid {
			w.Write([]byte(`wrong username/password`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
Idealnya fungsi middleware HARUS MENGEMBALIKAN struct yang implements http.Handler.
Beruntungnya, Go SUDAH MENYIAPKAN fungsi ajaib untuk MEMPERSINGKAT pembuatan struct-yang-implemenets-http.Handler.
Fungsi tersebut adalah http.HandlerFunc, cukup bungkus callback func(http.ResponseWriter,*http.Request)
sebagai tipe http.HandlerFunc dan semuanya beres.

Isi dari MiddlewareAuth() sendiri adalah pengecekan basic auth, sama seperti pada chapter sebelumnya.
*/

// di materi sebelumnya ini adalah fungsi AllowOnlyGet()
func MiddlewareAllowOnlyGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Write([]byte("Only GET is allowed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
