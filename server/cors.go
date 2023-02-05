package server

import (
	"net/http"
)

// CORS stands for Cross Origin Resource Sharing. For security reasons, browsers and websites restrict cross-origin HTTP requests eminating from scripts. By setting "Access-Control-Allow-Origin" to *, CORS should be enabled for this server. This is necessary because the script accessing our site (in the browser extension front-end) is not actually served by our site. It's essentially running locally on the user computer. The logic for this function was taken from: https://stackoverflow.com/questions/64062803/how-to-enable-cors-in-go
func CORS(handler http.Handler) http.Handler {

	next := handler.ServeHTTP

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}
		next(w, r)
	})
}
