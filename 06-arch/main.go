package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	u := newItemUsecase()
	h := newHandler(u)

	router := gin.Default()

	router.GET("/", h.home)
	router.GET("/hello", h.hello)
	router.POST("/bye", h.bye)

	router.POST("/items", h.saveItem)
	router.GET("/items", h.listItems)
	router.GET("/items/:id", h.getItemByID)
	router.PUT("/items/:id", h.updateItemByID)
	router.DELETE("/items/:id", h.deleteItemByID)

	log.Println("Server started at http://localhost:8080/")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Global error
// /////////////////////////////////////////////////////////////////////////////
var errNotFound = errors.New("not found")

// /////////////////////////////////////////////////////////////////////////////
// In-memory storage
// /////////////////////////////////////////////////////////////////////////////
var itemsDB = make(map[int]item)

// /////////////////////////////////////////////////////////////////////////////
// Handler
// /////////////////////////////////////////////////////////////////////////////
type handler struct {
	usecase ItemUsecase
}

// Handler constructor
func newHandler(u ItemUsecase) *handler {
	return &handler{
		usecase: u,
	}
}

func (h *handler) home(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the home page!")
}

func (h *handler) hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!")
}

func (h *handler) bye(c *gin.Context) {
	var msg map[string]string
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}
	message, exists := msg["message"]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message field is missing"})
		return
	}
	c.String(http.StatusOK, "Received POST request with message: %s", message)
}

func (h *handler) saveItem(c *gin.Context) {
	var it item

	if err := c.BindJSON(&it); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.saveItem(it); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "item saved successfully")
}

func (h *handler) listItems(c *gin.Context) {
	its, err := h.usecase.listItems()
	if err != nil {
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, its)
}

func (h *handler) getItemByID(c *gin.Context) {
	id, ok := parseIdParam(c, "id")
	if !ok {
		return
	}

	item, err := h.usecase.getItemByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *handler) updateItemByID(c *gin.Context) {
	id, ok := parseIdParam(c, "id")
	if !ok {
		return
	}

	var updatedItem item
	if err := c.BindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	updatedItem, err := h.usecase.updateItemByID(id, updatedItem)
	if err != nil {
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully", "item": updatedItem})
}

func (h *handler) deleteItemByID(c *gin.Context) {
	id, ok := parseIdParam(c, "id")
	if !ok {
		return
	}

	err := h.usecase.deleteItemByID(id)
	if err != nil {
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Item deleted successfully"})
}

// /////////////////////////////////////////////////////////////////////////////
// Usecases
// /////////////////////////////////////////////////////////////////////////////
type ItemUsecase interface {
	saveItem(it item) error
	listItems() (map[int]item, error)
	getItemByID(id int) (item, error)
	updateItemByID(id int, updatedItem item) (item, error)
	deleteItemByID(id int) error
}

var _ ItemUsecase = &itemUsecase{}

type itemUsecase struct{}

func newItemUsecase() *itemUsecase {
	return &itemUsecase{}
}

func (u *itemUsecase) saveItem(it item) error {
	if _, exists := itemsDB[it.ID]; exists {
		return errors.New("item with this ID already exists")
	}

	itemsDB[it.ID] = it
	return nil
}

func (u *itemUsecase) listItems() (map[int]item, error) {
	return itemsDB, nil
}

func (u *itemUsecase) getItemByID(id int) (item, error) {
	res, exists := itemsDB[id]

	if !exists {
		return item{}, errNotFound
	}

	return res, nil
}

func (u *itemUsecase) updateItemByID(id int, updatedItem item) (item, error) {
	if _, exists := itemsDB[id]; !exists {
		return item{}, errNotFound
	}

	updatedItem.ID = id
	itemsDB[id] = updatedItem
	return updatedItem, nil
}

func (u *itemUsecase) deleteItemByID(id int) error {
	if _, exists := itemsDB[id]; !exists {
		return errNotFound
	}
	delete(itemsDB, id)
	return nil
}

// ////////////////////////////////////////////////////////////////////////////
// Domain
// ///////////////////////////////////////////////////////////////////////////
// Item entity
type item struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ///////////////////////////////////////////////////////////////////////////
// Utils
// ///////////////////////////////////////////////////////////////////////////
func parseIdParam(c *gin.Context, param string) (int, bool) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, false
	}

	return id, true
}
