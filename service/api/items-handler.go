package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lorenzougolini/wimf-app/service/api/reqcontext"
	"github.com/lorenzougolini/wimf-app/service/models"
	"github.com/lorenzougolini/wimf-app/service/templates"
)

func (rt *_router) addItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	itemid := strings.TrimSpace(r.URL.Query().Get("itemid"))
	var message string

	// check valid itemid
	if itemid == "" || len(itemid) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		message = fmt.Sprintf("The provided itemid '%s' is not valid", itemid)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	// check if item already exists, if yes add, else create new
	// user, err := rt.db.GetByItemid(username)

	// if err != nil {

	// 	// user doesn't exists, create a new one
	// 	generateID, _ := uuid.NewV4()
	// 	newUserID = formatId(generateID.String())
	// 	if exists, err := rt.db.CheckIDExistence(newUserID); err != nil || exists {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}

	// 	err = rt.db.SetUser(newUserID, username)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		_ = json.NewEncoder(w).Encode(err)
	// 		return
	// 	}
	// 	ctx.Logger.Info("User created")
	// 	w.WriteHeader(http.StatusOK)
	// _ = json.NewEncoder(w).Encode(User{UserID: newUserID, Username: username})

	// } else {

	// 	ctx.Logger.Info("User logged in")
	// 	_ = json.NewEncoder(w).Encode(user)
	// }
	err := rt.db.AddItem(itemid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		ctx.Logger.Printf(err.Error())
		return
	}
	ctx.Logger.Info("Item added")
	w.Header().Set("HX-Trigger", `{"item-added": true}`)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.Item{ItemID: itemid})
}

func (rt *_router) getHomeItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	items, err := rt.db.GetLastItems(10)
	if err != nil {
		http.Error(w, "Failed loading items", http.StatusInternalServerError)
		return
	}

	err = templates.HomeItems(items).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "HomeItems render error", http.StatusInternalServerError)
	}
}
