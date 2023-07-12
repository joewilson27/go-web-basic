package main

/*
AJAX Multiple File Upload

Pada chapter ini, kita akan belajar 3 hal dalam satu waktu, yaitu:
- Bagaiamana cara untuk upload file via AJAX.
- Cara untuk handle upload banyak file sekaligus.
- Cara handle upload file yang lebih hemat memori.

Sebelumnya pada chapter B.13. Form Upload File, pemrosesan file upload dilakukan lewat
ParseMultipartForm, sedangkan pada chapter ini metode yang dipakai BERBEDA, yaitu
menggunakan MultipartReader.

Kelebihan dari MultipartReader adalah, file yang di upload TIDAK DI SIMPAN
sebagai file temporary di lokal terlebih dahulu (tidak seperti ParseMultipartForm),
melainkan langsung diambil dari stream io.Reader.

Di bagian front end, upload file secara asynchronous bisa dilakukan menggunakan
objek FormData. Semua file dimasukkan dalam objek FormData,
setelah itu objek tersebut dijadikan sebagai payload AJAX request.
*/

// 2. Front End
/*
siapkan file view.html

AJAX dilakukan lewat jQuery.ajax. Berikut adalah penjelasan mengenai konfigurasi processData dan contentType
dalam AJAX yang sudah dibuat.
- Konfigurasi contentType PERLU DI SET ke false AGAR HEADER Content-Type yang
	dikirim bisa MENYESUAIKAN data yang disisipkan.
- Konfigurasi processData juga PERLU DI SET ke false, agar data yang akan di kirim TIDAK OTOMATIS DIKONVERSI ke query
	string atau json string (tergantung contentType). Pada konteks ini kita memerlukan payload TETAP dalam tipe FormData.
*/

// 3. Back End
import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
Sebelumnya, pada chapter B.13. Form Upload File, metode yang digunakan untuk handle file upload adalah MENGGUNAKAN
ParseMultipartForm, file DI PROSES DALAM MEMORY dengan alokasi tertentu, dan JIKA MELEBIHI ALOKASI maka akan
disimpan pada temporary file.

Metode tersebut KURANG TEPAT GUNA jika digunakan UNTUK MEMPROSES file yang UKURANNYA BESAR (file size melebihi maxMemory)
atau jumlah file-nya SANGAT BANYAK (memakan waktu, karena isi dari masing-masing file akan ditampung pada file temporary
sebelum benar-benar di-copy ke file tujuan).

SOLUSINYA dari dua masalah di atas adalah MENGGUNAKAN MultipartReader untuk handling file upload. Dengan metode ini,
file destinasi isinya akan di-COPY LANGSUNG DARI STREAM io.Reader, TANPA BUTUH file temporary untuk perantara.
*/

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
		return
	}

	basePath, _ := os.Getwd()
	reader, err := r.MultipartReader()
	/*
		Bisa dilihat, method .MultipartReader() dipanggil dari objek request milik handler.
		Mengembalikan dua objek, pertama *multipart.Reader dan error (jika ada).
	*/
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		Selanjutnya LAKUKAN PERULANGAN terhadap objek reader. Setiap file yang di-upload di proses di masing-masing perulangan.
		Setelah looping berakhir. idealnya semua file sudah terproses dengan benar.
	*/
	for {
		part, err := reader.NextPart()
		/*
			Method .NextPart() mengembalikan 2 informasi, yaitu objek stream io.Reader (dari file yg di upload), dan error.
		*/
		if err == io.EOF {
			/*
				Jika reader.NextPart() mengembalikan error io.EOF, MENANDAKAN BAHWA SEMUA FILE SUDAH DI PROSES,
				maka hentikan perulangan.
			*/
			break
		}

		// File destinasi dipersiapkan,
		fileLocation := filepath.Join(basePath, "files", part.FileName())
		/*
			Fungsi filepath.Join berguna untuk pembentukan path.

			contoh hasil pembentukan path print fileLocation
			C:\go\go-web-basic\b16-ajax-multiple-upload\files\saruman-witch.jpg


		*/
		dst, err := os.Create(fileLocation)
		if dst != nil {
			defer dst.Close()
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// kemudian diisi dengan data dari stream file, menggunakan io.Copy()
		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(`all files uploaded`))
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/upload", handleUpload)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}
