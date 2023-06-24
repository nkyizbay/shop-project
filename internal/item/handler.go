package item

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	itemService Service
}

func Handler(rout *gin.Engine, itemService Service) *handler {
	h := handler{
		itemService: itemService,
	}

	rout.POST("/items", h.CreateItem)
	rout.DELETE("/items/:id", h.CancelItem)

	return &h
}

func (t *handler) CreateItem(c *gin.Context) {
	item := new(Item)
	if err := c.BindJSON(&item); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	if item.CheckFieldsEmpty() {
		c.String(http.StatusBadRequest, "Empty fields error, please fill them")
	}

	if item.IsInvalidPrice() {
		c.String(http.StatusBadRequest, "Invalid price error, please check if price > 0")
	}

	if err := t.itemService.CreateItem(item); err != nil {
		c.String(http.StatusInternalServerError, "Error, something went wrong")
	}

	c.String(http.StatusCreated, "Item created")
}

func (t *handler) CancelItem(c *gin.Context) {
	itemIDStr := c.Param("id")
	itemID, _ := strconv.Atoi(itemIDStr)

	if IsInvalidID(itemID) {
		c.String(http.StatusBadRequest, "Invalid id")
	}

	if err := t.itemService.CancelItem(itemID); err != nil {
		c.String(http.StatusInternalServerError, "Error, something went wrong")
	}

	c.String(http.StatusNoContent, "There is no record")
}
