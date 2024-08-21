package logout

import (
	Generic "Amazon_Server/Generic"

	"fmt"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request){
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Unix(0, 0),
			Path:    "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Logged out successfully")
	}
}