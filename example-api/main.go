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

	router.GET("/items", h.listItems)
	router.POST("/items", h.saveItem)

	router.GET("/items/:id", h.getItem)
	router.PUT("/items/:id", h.updateItem)
	router.DELETE("/items/:id", h.deleteItem)

	log.Println("Server started at http://localhost:8080/")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Global error
// /////////////////////////////////////////////////////////////////////////////
var errNotFound = errors.New("not found")
var errItemAlreadyExists = errors.New("item ID already exists")

// ////////////////////////////////////////////////////////////////////////////
// Handler
// ////////////////////////////////////////////////////////////////////////////
type handler struct {
	usecase *itemUsecase
}

// Handler constructor
func newHandler(u *itemUsecase) *handler {
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
		if err == errItemAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, "item saved successfully")
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

func (h *handler) getItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	it, err := h.usecase.getItem(id)

	if err != nil {
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, it)
}

func (h *handler) updateItem(c *gin.Context) {
	var it item
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err = c.BindJSON(&it); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	it.ID = id

	if err = h.usecase.updateItem(it); err != nil {
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, "item updated successfully")
}

func (h *handler) deleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.deleteItem(id); err != nil {
		if err == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, "item deleted successfully")
}

// ///////////////////////////////////////////////////////////////////////////
// Usecases
// ///////////////////////////////////////////////////////////////////////////
type itemUsecase struct {
	storage map[int]item
}

func newItemUsecase() *itemUsecase {
	return &itemUsecase{}
}

func (u *itemUsecase) saveItem(it item) error {
	if u.storage == nil {
		u.storage = make(map[int]item)
	}

	if _, exists := u.storage[it.ID]; exists {
		return errItemAlreadyExists
	}

	now := time.Now()

	it.CreatedAt = now
	it.UpdatedAt = now

	u.storage[it.ID] = it
	return nil
}

func (u *itemUsecase) listItems() (map[int]item, error) {
	return u.storage, nil
}

func (u *itemUsecase) getItem(id int) (item, error) {
	it, exists := u.storage[id]

	if !exists {
		return item{}, errNotFound
	}

	return it, nil
}

func (u *itemUsecase) updateItem(item item) error {
	_, exists := u.storage[item.ID]

	if !exists {
		return errNotFound
	}

	item.UpdatedAt = time.Now()

	u.storage[item.ID] = item

	return nil
}

func (u *itemUsecase) deleteItem(id int) error {
	_, exists := u.storage[id]

	if !exists {
		return errNotFound
	}

	delete(u.storage, id)

	return nil
}

// ///////////////////////////////////////////////////////////////////////////
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
