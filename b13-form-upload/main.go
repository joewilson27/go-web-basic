package main

/*
Form Upload File
Yo saatnya belajar bagaimana cara meng-handle upload file lewat form.
Di beberapa bagian caranya MIRIP SEPERTI pada chapter sebelumnya, hanya
PERLU DITAMBAHKAN proses untuk handling file yang di-upload.
File tersebut disimpan ke dalam path/folder tertentu.
*/

//
/*
Program sederhana yang akan kita buat, memiliki satu form dengan 2 inputan,
alias dan file. Data file nantinya disimpan pada folder files yang telah dibuat,
dengan nama sesuai nama file aslinya. Kecuali ketika user mengisi
inputan alias, maka nama tersebut yang akan digunakan sebagai nama file tersimpan.
*/

// 2. Front End
/*
Di bagian front end, isi file view.html dengan kode berikut.
Template file ini nantinya yang dimunculkan sebagai landing page.

Perlu diperhatikan, pada tag <form> perlu ditambahkan atribut
enctype="multipart/form-data", agar http request mendukung upload file.
*/

// 3. Back End
/*
siapkan 2 buah route handler, satu untuk landing page, dan satunya lagi
digunakan ketika proses upload selesai
*/

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

/*
Handler route / isinya proses untuk menampilkan landing page (file view.html).
Method yang diperbolehkan mengakses rute ini hanya GET.
*/
func routeIndexGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var tmpl = template.Must(template.ParseFiles("view.html"))
	var err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
Selanjutnya siapkan handler untuk rute /proccess, yaitu fungsi routeSubmitPost.
Gunakan statement r.ParseMultipartForm(1024) untuk parsing form data yang dikirim.
*/
func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil { // argumen 1024 adl maxMemory
		/*
			Pemanggilan method tersebut membuat file yang terupload DISIMPAN SEMENTARA pada memory dengan
			alokasi adalah SESUAI DENGAN maxMemory. Jika ternyata kapasitas yang sudah dialokasikan tersebut TIDAK CUKUP,
			maka file akan disimpan dalam temporary file.
		*/
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		Method ParseMultipartForm() digunakan untuk mem-parsing form data YANG ADA
		data file nya. Argumen 1024 pada method TERSEBUT ADALAH maxMemory. Pemanggilan
		method tersebut membuat file yang terupload disimpan sementara pada
		memory dengan alokasi adalah sesuai dengan maxMemory. Jika ternyata kapasitas
		yang sudah dialokasikan tersebut tidak cukup, maka file akan disimpan dalam
		temporary file.
	*/

	alias := r.FormValue("alias") // jika alias di isi, maka nantinya nama file yg di upload menggunakan alias ini

	uploadedFile, handler, err := r.FormFile("file")
	/*
		Statement r.FormFile("file") digunakan untuk mengambil file yg di upload, mengembalikan 3 objek:

		- Objek bertipe multipart.File (yang merupakan turunan dari *os.File)
		- Informasi header file (bertipe *multipart.FileHeader)
		- Dan error jika ada
	*/

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		Tahap selanjutnya adalah, MENAMBAHKAN KODE membuat file baru, yang nantinya file ini akan diisi dengan
		isi dari file yang ter-upload. Jika inputan alias di-isi, maka nama nilai inputan tersebut dijadikan sebagai nama file.
	*/
	filename := handler.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
		/*
			Fungsi filepath.Ext digunakan untuk MENGAMBIL EKSTENSI dari sebuah file

			Pada kode di atas, handler.Filename yang berisi nama file terupload diambil ekstensinya,
			LALU DIGABUNG dengan alias yang sudah terisi.
		*/
	}

	fileLocation := filepath.Join(dir, "files", filename)
	/*
		Fungsi filepath.Join berguna untuk pembentukan path.

		contoh hasil pembentukan path print fileLocation
		C:\go\go-web-basic\b13-form-upload\files\main-qimg-2b414536020cd00309f1dc5b4d31e8fe.webp

		bila rename dengan alias
		C:\go\go-web-basic\b13-form-upload\files\saruman-witch.jpg

	*/

	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	/*
		Fungsi os.OpenFile digunakan untuk MEMBUKA FILE. Fungsi ini membutuhkan 3 buah parameter:
		- Parameter pertama merupakan path atau lokasi dari file yang ingin di buka
		- Parameter kedua adalah flag mode, apakah read only, write only, atau keduanya, atau lainnya.
			- os.O_WRONLY|os.O_CREATE maknanya, file yang dibuka HANYA akan bisa di TULIS SAJA (write only konsantanya adalah os.O_WRONLY),
				dan file tersebut akan dibuat jika belum ada (konstantanya os.O_CREATE).
		- Sedangkan parameter terakhir adalah permission dari file, yang digunakan dalam pembuatan file itu sendiri.
	*/
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		/*
			Fungsi io.Copy AKAN MENGISI konten file parameter pertama (targetFile) dengan isi parameter kedua
			(uploadedFile). File kosong yang telah kita buat tadi akan diisi dengan data file yang tersimpan di memory.
		*/
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("done"))

}

func main() {
	http.HandleFunc("/", routeIndexGet)
	http.HandleFunc("/process", routeSubmitPost)

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}
