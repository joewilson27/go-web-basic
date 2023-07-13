package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var shared *_Configuration

type _Configuration struct {
	Server struct {
		Port         int           `json:"port"`         // port adl properties pada config.json
		ReadTimeout  time.Duration `json:"read_timeout"` // read_timeout adl properties pada config.json
		WriteTimeout time.Duration `json:"write_timeout"`
	} `json:"server"`

	Log struct {
		Verbose bool `json:"verbose"` // verbose adl properties pada config.json -> penentu apakah log di print atau tidak.
	} `json:"log"`
}

/*
Struct _Configuration dibuat berisikan banyak property yang STRUKTURNYA SAMA PERSIS
dengan isi file config.json. Dengan desain seperti ini, akan SANGAT MEMUDAHKAN developer
dalam pengaksesan konfigurasi.

Dari struct tersebut TERCETAK private objek bernama shared.
Variabel inilah yang nantinya AKAN DIKEMBALIKAN lewat fungsi yang akan kita buat.
*/

/*
Selanjutnya, isi init() dengan beberapa proses: membaca file json, lalu di decode ke
object shared.

Dengan menuliskan proses barusan ke fungsi init(), pada saat package conf ini di
import ke package lain maka file config.json AKAN OTOMATIS di parsing. Dan dengan
menambahkan sedikit validasi, parsing hanya akan terjadi sekali di awal.
*/
func init() {
	if shared != nil {
		return
	}

	basePath, err := os.Getwd()
	if err != nil {
		panic(err)
		return
	}

	bts, err := ioutil.ReadFile(filepath.Join(basePath, "conf", "config.json"))
	/*
		*Fungsi filepath.Join AKAN MENGGABUNGKAN item-item dengan path separator
		sesuai dengan sistem operasi di mana program dijalankan. \ untuk Windows dan / untuk Linux/Unix.
	*/
	if err != nil {
		panic(err)
		return
	}

	shared = new(_Configuration)
	err = json.Unmarshal(bts, &shared)
	/*
		 Unmarshal -> deserializing data in various formats, such as JSON or XML, INTO Go data
		 structures.


			The json.Unmarshal function takes two arguments:
			- the JSON data as a byte slice or string,
			- and a pointer to the variable where the unmarshaled data should be stored.
			It returns an error if the unmarshaling process fails.
	*/
	if err != nil {
		panic(err)
		return
	}
}

// fungsi yang mengembalikan object shared
func Configuration() _Configuration {
	return *shared
}
