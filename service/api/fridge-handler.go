package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lorenzougolini/wimf-app/service/api/reqcontext"
	"github.com/lorenzougolini/wimf-app/service/models"
	"github.com/lorenzougolini/wimf-app/service/templates"
)

func (rt *_router) getFridge(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	items, err := rt.db.GetFridge()
	if err != nil {
		http.Error(w, "Error retrieving the fridge", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		// Just refresh the table part (No Header, No Footer)
		templates.FridgeTable(items).Render(r.Context(), w)
	} else {
		// Full Page Load (Includes Header, Footer, CSS)
		templates.Fridge(items).Render(r.Context(), w)
	}

	// fridgeTemplate := templates.Fridge(items)
	// err = templates.Layout(fridgeTemplate, "Fridge", "/fridge").Render(r.Context(), w)
	// if err != nil {
	// 	log.Printf("Error rendering fridge: %v", err)
	// }
}

func (rt *_router) getFridgeDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	barcode := r.URL.Query().Get("barcode")
	_, items, err := rt.db.GetItemsByBarcode(barcode)
	if err != nil {
		http.Error(w, "Error retrieving fridge details", http.StatusInternalServerError)
		return
	}

	templates.FridgeDetailModal(items).Render(r.Context(), w)
}

func (rt *_router) getEditForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.URL.Query().Get("id")

	item, _ := rt.db.GetItemById(id)

	info := models.ProductInfo{
		Barcode: item.Barcode,
		Name:    item.Name,
		Brand:   item.Brand,
	}
	templates.ExpirationModal(info, true, id, item.ExpirationDate).Render(r.Context(), w)
}
