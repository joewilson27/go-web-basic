package main

/*
HTTP Cookie
Cookie adalah data dalam bentuk teks yang disimpan pada komputer (oleh web browser)
ketika pengunjung sedang surfing ke sebuah situs. Cookie dapat dibuat dari sisi front end (javascript)
maupun back end (dalam konteks ini Go).

Cookie merupakan salah satu aspek penting dalam pengembangan aplikasi web. Sangat sering kita
membutuhkan sebuah data bisa disimpan dan diakses untuk keperluan aplikasi web kita,
seperti pengecekan preferensi pengunjung, pengecekan status login tidak nya user.


install package:
 go get github.com/novalagung/gubrak/v2
*/

// 1. Praktek
import (
	"fmt"
	"net/http"
	"time"
)

type M map[string]interface{}

var cookieName = "CookieData"

// Di dalam fungsi ini, data berupa random string disimpan dalam cookie.
func ActionIndex(w http.ResponseWriter, r *http.Request) {
	cookieName := "CookieData"

	c := &http.Cookie{} // cetak objek baru utk cookie baru dengan http.Cookie

	/*
		Cookie BISA DIAKSES lewat method .Cookie() milik objek *http.Request.
		Method ini mengembalikan 2 informasi:
			- Objek cookie
			- Error, jika ada
	*/
	if storedCookie, _ := r.Cookie(cookieName); storedCookie != nil {
		/*
			ketika storedCookie nilainya BUKANLAH nil (berarti cookie dengan nama cookieName
			sudah dibuat), maka objek cookie tersebut disimpan dalam c
		*/
		c = storedCookie
	}

	/*
		Jika c.Value adalah kosong, kita asumsikan bahwa cookie belum pernah dibuat
		(atau expired), maka kita BUAT COOKIE baru dengan data adalah random string.
	*/
	if c.Value == "" {
		c = &http.Cookie{}
		c.Name = cookieName
		c.Value = "abcdefghijklmnopqrstuvwxyzABCDE" // gubrak.RandomString(32)           // gubrak.RandomString(32) akan menghasilkan string 32 karakter.
		c.Expires = time.Now().Add(5 * time.Minute) // Cookie bisa expired. Lama cookie aktif ditentukan lewat property Expires
		http.SetCookie(w, c)                        // Gunakan http.SetCookie() untuk menyimpan cookie yang baru dibuat.
	}

	w.Write([]byte(c.Value))
}

// hapus cookie
func ActionDelete(w http.ResponseWriter, r *http.Request) {
	/*
		Cara menghapus cookie adalah dengan menge-set ulang cookie dengan NAMA YANG SAMA,
		dengan isi property Expires = time.Unix(0, 0) dan MaxAge = -1.
		Tujuannya agar cookie expired.
	*/
	c := &http.Cookie{}
	c.Name = cookieName
	c.Expires = time.Unix(0, 0)
	c.MaxAge = -1
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main() {
	port := 9000
	http.HandleFunc("/", ActionIndex)
	http.HandleFunc("/delete", ActionDelete)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)
}

/*
Variabel cookieName berisikan string, nantinya digunakan sebagai nama cookie.
- Rute / bertugas untuk MEMBUAT cookie baru (jika belum ada atau cookie sudah ada
	namun expired).
- Rute /delete mempunyai tugas untuk MENGHAPUS cookie, lalu redirect ke /
	sehingga cookie baru akan dibuat



	testing

	Coba refresh page beberapa kali, informasi header cookie dan data yang muncul
	adalah TETAP SAMA. Karena ketika cookie sudah pernah dibuat, maka seterusnya
	endpoint ini AKAN MENGGUNAKAN data cookie YANG SUDAH TERSIMPAN tersebut.

	Selanjutnya, buka url /delete, halaman akan di redirect kembali ke /,
	dan random string baru beserta COOKIE BARU TERBUAT.
*/
