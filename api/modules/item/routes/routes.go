package routes

import (
	"github.com/gin-gonic/gin"

	"go-api/modules/item/controller"
)

func Register(router *gin.RouterGroup, ctrl *controller.ItemController) {
	router.GET("/items", ctrl.GetAll)
	router.GET("/items/:id", ctrl.GetByID)
	router.POST("/items", ctrl.Create)
	router.PUT("/items/:id", ctrl.Update)
	router.DELETE("/items/:id", ctrl.Delete)
}
