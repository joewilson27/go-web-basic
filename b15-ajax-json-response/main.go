package main

/*
AJAX JSON Response
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func ActionIndex(w http.ResponseWriter, r *http.Request) {
	data := []struct {
		Nama string
		Age  int
	}{
		{"Richard Grayson", 24},
		{"Jason Todd", 23},
		{"Tim Drake", 22},
		{"Damian Wayne", 21},
	}

	jsonInBytes, err := json.Marshal(data) // process of converting Go data structures or objects into a serialized format, such as JSON, XML, or binary data.
	/*
		pada contoh ini, kita ingin convert object variable data ke json menggunakan Marshal.
		Fungsi ini mengembalikan dua nilai balik
		- data json (dalam bentuk []byte)
		- dan error jika ada.

		*Untuk MENGAMBIL bentuk STRING dari hasil konversi JSON,
		cukup lakukan casting pada data slice bytes tersebut. Contoh: string(jsonInBytes)

	*/
	var printConvert string = string(jsonInBytes)
	fmt.Println("original data", jsonInBytes)
	fmt.Printf("hasil convert --> %s\n", printConvert)
	/*
		jsonInBytes data original, slice byte []byte
		[91 123 34 78 97 109 97 34 58 34 82 105 99 104 97 114 100 32 71 114 97 121 115 111 110 34 44 34 65 103 101 34 58 50 52 125 44 123 34 78 97 109 97 34 58 34 74 97 115 111 110 32 84 111 100 100 34 44 34 65 103 101 34 58 50 51 125 44 123 34 78 97 109 97 34 58 34 84 105 109 32 68 114 97 107 101 34 44 34 65 103 101 34 58 50 50 125 44 123 34 78 97 109 97 34 58 34 68 97 109 105 97 110 32 87 97 121 110 101 34 44 34 65 103 101 34 58 50 49 125 93]

		jsonInBytes hasil convert ke string :
		[{"Nama":"Richard Grayson","Age":24},{"Nama":"Jason Todd","Age":23},{"Nama":"Tim Drake","Age":22},{"Nama":"Damian Wayne","Age":21}]
	*/

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
}

func ActionWithEncoder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// buat object slice struct anonymous
		data := []struct {
			Nama string
			Age  int
		}{
			{"JRR Tolkien", 87},
			{"Sir Christopher Lee", 102},
			{"George Washington", 85},
			{"Marthin Luther Jr", 79},
		}

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func ActionHardCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		person := struct {
			Name  string
			Age   int
			Email string
		}{
			Name:  "Joe Wilson",
			Age:   20,
			Email: "joewilson@example.com",
		}

		// Create a file to write the JSON data (data fisik)
		file, err := os.Create("person.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Create a new JSON encoder
		encoder := json.NewEncoder(file)

		// Encode the person object into JSON and write it to the file (fisik)
		err = encoder.Encode(person)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		// to print into response
		// jsonInBytes, err := json.Marshal(person)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// w.Header().Set("Content-Type", "application/json")
		// _, err = w.Write(jsonInBytes)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		w.Header().Set("Content-Type", "application/json")
		errs := json.NewEncoder(w).Encode(person)
		if errs != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Only accept GET request", http.StatusBadRequest)
	}
}

func main() {
	port := 9000
	// 1. Praktek
	http.HandleFunc("/", ActionIndex) // dengan JSON Marshal

	// 3. JSON Response menggunakan JSON.Encoder
	/*
		Pada chapter sebelumnya sudah disinggung, bahwa LEBIH BAIK menggunakan json.Decoder jika ingin men-decode
		data yang sumbernya ada di stream io.Reader.

		Package json juga memiliki fungsi lain-nya yaitu json.Encoder, yang sangat cocok digunakan untuk meng-encode
		data menjadi JSON dengan tujuan objek langsung ke stream io.Reader.

		Karena tipe http.ResponseWriter adalah meng-EMBED io.Reader, maka jelasnya BISA KITA TERAPKAN penggunaan encoder di sini.
	*/

	http.HandleFunc("/with-encoder", ActionWithEncoder) // JSON.Encoder
	http.HandleFunc("/with-encoder-hard-copy", ActionHardCopy)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("server started at localhost:%d\n", port)
	http.ListenAndServe(addr, nil)
}
