package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lorenzougolini/wimf-app/service/api/reqcontext"
	"github.com/lorenzougolini/wimf-app/service/models"
	"github.com/lorenzougolini/wimf-app/service/templates"
)

func (rt *_router) addItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid request body")
		return
	}
	barcode := strings.TrimSpace(r.FormValue("barcode"))
	name := strings.TrimSpace(r.FormValue("name"))
	brand := strings.TrimSpace(r.FormValue("brand"))
	expDate := strings.TrimSpace(r.FormValue("expiration_date"))
	var message string

	// check valid barcode and parse date
	if barcode == "" || len(barcode) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided barcode '%s' is not valid", barcode)
		_ = json.NewEncoder(w).Encode(message)
		return
	}
	expirationDate, err := time.Parse("2006-01-02", expDate)
	if err != nil {
		expirationDate = time.Now().AddDate(0, 0, 14)
	}

	itemtToAdd := models.ProductInfo{
		Barcode: barcode,
		Name:    name,
		Brand:   brand,
	}
	err = rt.db.AddItem(itemtToAdd, expirationDate)
	if err != nil {
		ctx.Logger.Errorf("Error while adding item: adding new item", err)
		http.Error(w, "Error while adding item: adding new item", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info("Item added succesfully")
	w.Header().Set("HX-Trigger", `{"item-added": true}`)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Item{Barcode: barcode})
}

func (rt *_router) getExpirationForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid request body")
		return
	}
	barcode := strings.TrimSpace(r.URL.Query().Get("barcode"))
	var message string

	// check valid barcode
	if barcode == "" || len(barcode) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided barcode '%s' is not valid", barcode)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// check if item already exists, if yes add, else create new
	exists, localItem, err := rt.db.GetItemByBarcode(barcode)
	var itemtToAdd models.ProductInfo
	if exists {
		itemtToAdd = models.ProductInfo{
			Barcode: barcode,
			Name:    localItem.Name,
			Brand:   localItem.Brand,
		}
	} else {
		apiInfo, err := rt.foodApi.GetProductByBarcode(barcode)
		if err != nil {
			ctx.Logger.Errorf("Failed to fetch product info", err)
			http.Error(w, "Failed to fetch product info", http.StatusInternalServerError)
			return
		}
		itemtToAdd = apiInfo
	}

	// render the expiration modal
	err = templates.ExpirationModal(itemtToAdd).Render(r.Context(), w)
	if err != nil {
		ctx.Logger.Errorf("Error rendering modal: %v", err)
		http.Error(w, "Error rendering modal", http.StatusInternalServerError)
	}
}

func (rt *_router) getHomeItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	homeItems := models.HomeItems{}

	recentItems, err := rt.db.GetNItemsBy(10, "latest")
	if err != nil {
		ctx.Logger.Errorf("Failed to get home items", err)
		http.Error(w, "Failed loading home items", http.StatusInternalServerError)
		return
	}
	homeItems.RecentItems = recentItems

	expiringItems, err := rt.db.GetNItemsBy(10, "expiring")
	if err != nil {
		ctx.Logger.Errorf("Failed to get home items", err)
		http.Error(w, "Failed loading home items", http.StatusInternalServerError)
		return
	}
	homeItems.ExpiringItems = expiringItems

	err = templates.HomeItems(homeItems).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "HomeItems render error", http.StatusInternalServerError)
	}
}
