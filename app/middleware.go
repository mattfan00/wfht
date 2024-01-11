package app

import (
	"fmt"
	"net/http"
)

func (a *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *App) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		good := a.sessionManager.GetBool(r.Context(), "good")
		if good {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("HX-Request") != "" {
			w.Header().Add("HX-Redirect", "/login")
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})
}
