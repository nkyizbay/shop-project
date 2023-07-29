package item

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nkyizbay/shop-project/internal/auth"
)

type handler struct {
	itemService Service
}

func Handler(rout *gin.Engine, itemService Service) *handler {
	h := handler{
		itemService: itemService,
	}

	rout.POST("/items", auth.AdminMiddleware(), h.CreateItem)
	rout.DELETE("/items/:id", auth.AdminMiddleware(), h.CancelItem)

	return &h
}

func (t *handler) CreateItem(c *gin.Context) {
	item := new(Item)
	if err := c.BindJSON(&item); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if item.CheckFieldsEmpty() {
		c.String(http.StatusBadRequest, "Empty fields error, please fill them")
		return
	}

	if item.IsInvalidPrice() {
		c.String(http.StatusBadRequest, "Invalid price error, please check if price > 0")
		return
	}

	if err := t.itemService.CreateItem(item); err != nil {
		c.String(http.StatusInternalServerError, "Error, something went wrong")
		return
	}

	c.String(http.StatusCreated, "Item created")
}

func (t *handler) CancelItem(c *gin.Context) {
	itemIDStr := c.Param("id")
	itemID, _ := strconv.Atoi(itemIDStr)

	if IsInvalidID(itemID) {
		c.String(http.StatusBadRequest, "Invalid id")
		return
	}

	if err := t.itemService.CancelItem(itemID); err != nil {
		c.String(http.StatusInternalServerError, "Error, something went wrong")
		return
	}

	c.String(http.StatusNoContent, "There is no record")
}
