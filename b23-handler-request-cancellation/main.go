package main

/*
Server Handler HTTP Request Cancellation
Dalam konteks web application, KADANG KALA sebuah http request BUTUH WAKTU CUKUP LAMA untuk selesai, bisa jadi karena
KODE yang KURANG DI OPTIMASI, atau prosesnya memang lama, atau mungkin faktor lainnya. Dari sisi client,
biasanya ada handler untuk cancel request jika melebihi batas timeout yang sudah didefinisikan,
dan ketika itu terjadi di client akan sangat mudah untuk antisipasinya.

BERBEDA dengan handler di back end-nya, by default request yang sudah di-cancel oleh client
TIDAK TERDETEKSI (proses di back end akan tetap lanjut). UMUMNYA TIDAK ADA MASALAH mengenai ini,
TAPI ada kalanya KITA PERLU men-treat cancelled request dengan baik untuk keperluan lain
(logging, atau lainnya).

	*Chapter ini FOKUS TERHADAP cancellation pada client http request. Untuk cancellation pada proses konkuren
	silakan merujuk ke materi sebelumnya: A.64. Concurrency Pattern: Context Cancellation Pipeline.

*/

// 1. Praktek
/*
Dari objek *http.Request bisa diambil objek context lewat method .Context(), dan dari context
tersebut kita BISA MENDETEKSI apakah sebuah request di-cancel atau tidak oleh client.

Object context MEMILIKI method .Done() yang NILAI BALIKNYA berupa channel. Dari channel
tersebut kita BISA DETEKSI APAKAH request di-cancel atau tidak, caranya dengan cara mengecek
apakah ADA DATA YANG TERKIRIM lewat channel tersebut,
jika ada maka lakukan pengecekan pada error message-nya,
jika ada keterangan "cancelled" maka diasumsikan request tersebut dibatalkan.
*/

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	/*
		Di dalam handleIndex() disimulasikan sebuah PROSES MEMBUTUHKAN WAKTU LAMA untuk selesai (kita gunakan time.Sleep() untuk ini).
		Umumnya kode akan dituliskan langsung dalam handler tersebut, tapi pada kasus ini tidak. Untuk bisa mendeteksi
		sebuah request di-cancel atau tidak, harus di-dispatch sebuah goroutine baru.
		- Cara ke-1: bisa dengan menaruh proses utama di dalam gorutine tersebut, dan menaruh kode untuk deteksi
			di luar (di dalam handler-nya).
		- Cara ke-2: Atau sebaliknya. Menaruh proses utama di dalam handler, dan menempatkan deteksi cancelled
			request dalam goroutine baru.

		Pada contoh berikut, kita gunakan cara pertama. Tulis kode berikut dalam handler.
	*/
	done := make(chan bool)

	go func() {
		// 1. Praktek
		// do the process here
		// simulate a long-time request by putting 10 seconds sleep
		// time.Sleep(10 * time.Second)

		// done <- true // data dikirimkan ke channel done

		// 2. Handle Cancelled Request yang ada Payload-nya
		// do the process here
		// simulate a long-time request by putting 10 seconds sleep
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		time.Sleep(10 * time.Second)

		done <- true
	}()
	/*
		Pada kode di atas terlihat, PROSES UTAMA DIBUNGKUS dalam goroutine.
		Ketika selesai, sebuah data dikirimkan ke channel done.
	*/

	// Channel - Select
	/*
		Lalu diluar, keyword select dipergunakan untuk DETEKSI PENGIRIMAN terhadap dua channel.
		- Channel r.Context().Done(), jika channel ini menerima data MAKA DIASUMSIKAN request selesai.
			Selanjutnya lakukan pengecekan pada objek error milik konteks untuk deteksi apakah selesai-nya
			request ini karena memang selesai, atau di-cancel oleh client, atau faktor lainnya.
		- Channel <-done, jika channel ini menerima data, maka proses utama adalah selesai.
	*/
	select {
	case <-r.Context().Done():
		if err := r.Context().Err(); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "canceled") {
				log.Println("request cancelled")
			} else {
				log.Println("unknown error occured.", err.Error())
			}
		}
	case <-done:
		log.Println("done")
	}

}

func main() {
	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8080", nil)
	/*
		run project seperti biasa, lalu buka cli baru utk test request

		test dengan reques curl pada cli:
		curl -X GET http://localhost:8080/

		ketika kita biarkan proses selesai normal pada cli, maka akan muncul pesan 'done' pada cli tersebut,
		coba sekarang jalankan kembali request curl di atas dan PAKSA cancel dengan ctrl+c
		maka akan muncul 'request cancelled' pada cli yg running-in project

	*/

	// 2. Handle Cancelled Request yang ada Payload-nya
	/*
		KHUSUS UNTUK REQUEST dengan HTTP method yang MEWAJIBKAN UNTUK ADA request body-nya (payload),
		maka channel r.Context().Done() TIDAK AKAN MENERIMA data hingga terjadi proses read pada body payload.

		jalankan pada cli baru, kita tambahkan payload object kosong:
		curl -X POST http://localhost:8080/ -H 'Content-Type: application/json' -d '{}'
	*/
}
