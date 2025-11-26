package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lorenzougolini/wimf-app/service/templates"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	homeTemplate := templates.Home()
	err := templates.Layout(homeTemplate, "Home", "/").Render(r.Context(), w)
	if err != nil {
		// log.Printf("error rendering the home: %v", err)
		http.Error(w, "Failed to render home page", http.StatusInternalServerError)
	}
}
