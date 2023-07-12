package main

/*
Download File
Sebenarnya download file bisa dengan mudah di-implementasikan menggunakan teknik routing static file,
langsung akses url dari public assets di browser. Namun outcome dari teknik ini SANGAT TERGANTUNG pada browser.
Tiap browser MEMILIKI BEHAVIOR berbeda, ada yang file tidak di-download melainkan dibuka di tab, ada yang ter-download.

Dengan menggunakan teknik berikut, file pasti akan ter-download.
*/

// 2. Front End
/*
as usual siapin view.html

Tag <ul /> nantinya AKAN BERISIKAN list semua file yang ada dalam folder files.
Data list file didapat dari back end. Diperlukan sebuah AJAX untuk pengambilan data tersebut.

Siapkan sebuah fungsi dengan nama Xo atau bisa lainnya, fungsi ini berisikan closure renderData(),
getAllListFiles(), dan method init(). Buat instance object baru dari Xo, lalu akses method init(),
tempatkan dalam event window.onload.
*/

// 3. Back End
import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type M map[string]interface{}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
handleListFiles() untuk handler /
Isi dari handler ini adalah membaca semua file yang ada pada folder files untuk kemudian dikembalikan sebagai output berupa JSON. Endpoint ini akan diakses oleh AJAX dari front end.
*/
func handleListFiles(w http.ResponseWriter, r *http.Request) {
	files := []M{}
	basePath, _ := os.Getwd()
	/*
		Fungsi os.Getwd() MENGEMBALIKAN INFORMASI ABSOLUTE PATH di mana aplikasi di-eksekusi.
		Path tersebut kemudian di gabung dengan folder bernama files lewat fungsi filepath.Join dibawah
	*/
	filesLocation := filepath.Join(basePath, "files")
	/*
				*Fungsi filepath.Join AKAN MENGGABUNGKAN item-item dengan path separator
				sesuai dengan sistem operasi di mana program dijalankan. \ untuk Windows dan / untuk Linux/Unix.

		contoh jika di windows seperti ini:
		C:\go\go-web-basic\b16-ajax-multiple-upload\files\saruman-witch.jpg
	*/

	err := filepath.Walk(filesLocation, func(path string, info os.FileInfo, err error) error {
		/*
			Fungsi filepath.Walk berguna untuk MEMBACA ISI DARI SEBUAH DIREKTORI,
			apa yang ada di dalamnya (file maupun folder) akan di-loop. Dengan memanfaatkan
			callback parameter kedua fungsi ini (yang bertipe filepath.WalkFunc),
			kita bisa mengambil informasi tiap item satu-per satu.
		*/

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, M{"filename": info.Name(), "path": path})
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

/*
handleDownload() untuk handler /download

	Implementasi teknik download PADA DASARNYA SAMA pada semua bahasa pemrograman,
	yaitu dengan memainkan header Content-Disposition pada response.
*/
func handleDownload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.FormValue("path")
	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", f.Name())
	w.Header().Set("Content-Disposition", contentDisposition)
	/*
		Content-Disposition adalah salah satu ekstensi MIME protocol, berguna untuk MENGINFORMASIKAN BROWSER
		bagaimana dia harus berinteraksi dengan output.
		ADA BANYAK JENIS value content-disposition, salah satunya adalah attachment.
		Pada kode di atas, header Content-Disposition: attachment; filename=filename.json
		menghasilkan output response berupa attachment atau file,
		yang kemudian akan di-download oleh browser.
	*/

	/*
		Objek file yang direpresentasikan variabel f, isinya di-copy ke objek response
		lewat statement io.Copy(w, f).
	*/
	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/list-files", handleListFiles)
	http.HandleFunc("/download", handleDownload)

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", nil)
}
