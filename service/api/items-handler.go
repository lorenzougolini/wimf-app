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
	addDate := strings.TrimSpace(r.FormValue("addition_date"))
	manual := strings.TrimSpace(r.FormValue("isManual"))
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

	additionDate := time.Now()
	if manual == "true" {
		additionDate, _ = time.Parse("2006-01-02", addDate)
	}

	itemtToAdd := models.ProductInfo{
		Barcode: barcode,
		Name:    name,
		Brand:   brand,
	}
	err = rt.db.AddItem(itemtToAdd, expirationDate, additionDate)
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
	exists, localItem, err := rt.db.GetItemsByBarcode(barcode)
	var itemtToAdd models.ProductInfo
	if exists {
		itemtToAdd = models.ProductInfo{
			Barcode: barcode,
			Name:    localItem[0].Name,
			Brand:   localItem[0].Brand,
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
	err = templates.ExpirationModal(itemtToAdd, false, "", time.Time{}).Render(r.Context(), w)
	if err != nil {
		ctx.Logger.Errorf("Error rendering modal: %v", err)
		http.Error(w, "Error rendering modal", http.StatusInternalServerError)
	}
}

func (rt *_router) getManualForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	emptyProduct := models.ProductInfo{}

	err := templates.ExpirationModal(emptyProduct, true, "", time.Time{}).Render(r.Context(), w)
	if err != nil {
		ctx.Logger.Errorf("Error rendering manual modal: %v", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
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

func (rt *_router) updateItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	r.ParseForm()
	id := r.FormValue("id")
	name := r.FormValue("name")
	brand := r.FormValue("brand")
	dateStr := r.FormValue("expiration_date")
	date, _ := time.Parse("2006-01-02", dateStr)

	err := rt.db.UpdateItem(id, name, brand, date)
	if err != nil {
		http.Error(w, "Error while updating item", http.StatusInternalServerError)
		message := fmt.Sprintf("Error: %s", err)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	w.Header().Set("HX-Trigger", `{"update-fridge": true}`)
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) deleteItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	id := r.URL.Query().Get("id")
	err := rt.db.DeleteItem(id)
	if err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Trigger", `{"update-fridge": true}`)
	w.WriteHeader(http.StatusOK)
}
