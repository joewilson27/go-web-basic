package main

import "net/http"

const USERNAME = "ironman"
const PASSWORD = "secret"

type CustomMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

/*
Middleware yang didaftarkan DITAMPUNG oleh slice .middlewares.
*/
func (c *CustomMux) RegisterMiddleware(next func(next http.Handler) http.Handler) { // this method belongs to CustomMux struct object
	c.middlewares = append(c.middlewares, next)
}

/*
Lalu buat method ServeHTTP. Method ini DIPERLUKAN dalam custom mux agar memenuhi kriteria interface http.Handler.
*/
func (c *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &c.ServeMux

	for _, next := range c.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
	/*
		Method ServeHTTP() milik mux adalah method YANG PASTI DIPANGGIL pada web server,
		di setiap request yang masuk.
	*/

	/*
		Dengan perubahan di atas, SETIAP KALI ADA REQUERST MASUK pasti akan melewati
		middleware-middleware terlebih dahulu secara berurutan.
		Jika lolos middleware ke-1, lanjut ke-2; jika lolos middleware ke-2, lanjut ke-3;
		dan seterusnya.
	*/
}

func MiddlewareAuth(next http.Handler) http.Handler {
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

func MiddlewareAllowOnlyGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Write([]byte("Only GET is allowed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
