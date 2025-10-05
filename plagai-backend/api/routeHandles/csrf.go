/*
UNUSED FOR NOW
MAYBE WE'LL NEED IT LATER
*/

package routeHandles

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
)

func GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "OK"}`)
}