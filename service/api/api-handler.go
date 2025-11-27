package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHome)
	rt.router.GET("/fridge", rt.wrap(rt.getFridge))
	rt.router.GET("/fridge/details", rt.wrap(rt.getFridgeDetails))
	rt.router.DELETE("/fridge/item", rt.wrap(rt.deleteItem))
	rt.router.GET("/fridge/item/edit", rt.wrap(rt.getEditForm))
	rt.router.PUT("/fridge/items", rt.wrap(rt.updateItem))

	rt.router.POST("/fridge/items", rt.wrap(rt.addItem))
	rt.router.GET("/fridge/items/form", rt.wrap(rt.getExpirationForm))
	rt.router.GET("/fridge/items/manual-form", rt.wrap(rt.getManualForm))
	rt.router.GET("/fridge/home-items", rt.wrap(rt.getHomeItems))

	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
