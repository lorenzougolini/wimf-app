package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lorenzougolini/wimf-app/service/api/reqcontext"
	"github.com/lorenzougolini/wimf-app/service/templates"
)

func (rt *_router) getFridge(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.WriteHeader(http.StatusOK)
	fridgeTemplate := templates.Fridge()
	err := templates.Layout(fridgeTemplate, "Fridge", "/fridge").Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering fridge: %v", err)
	}
}