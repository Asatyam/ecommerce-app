package main

import (
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"net/http"
	"strings"
)

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.GuestUser)
			next.ServeHTTP(w, r)
			return
		}

		if len(authorizationHeader) < 8 || !strings.HasPrefix(authorizationHeader, "Bearer ") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]
		v := validator.New()
		if data.ValidateToken(v, tokenString); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, tokenString)

		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.invalidCredentialsResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
func (app *application) requireAuthenticated(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsGuestUser() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
}
func (app *application) requireActivated(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	return app.requireAuthenticated(fn)
}
