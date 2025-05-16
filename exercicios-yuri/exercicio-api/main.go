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

	router.GET("/items/:id", h.itemById)
	router.PUT("/items", h.changeItem)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (h *handler) itemById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, errUsecase := h.usecase.itemById(id)

	if errUsecase != nil {
		if errUsecase == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": errUsecase.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errUsecase.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *handler) changeItem(c *gin.Context) {

	var it item

	if err := c.BindJSON(&it); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.changeItem(it); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "item updated successfully")
}

func (h *handler) deleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	errUsecase := h.usecase.deleteItem(id)

	if errUsecase != nil {
		if errUsecase == errNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": errUsecase.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errUsecase.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, "item deleted successfully")
}

// ///////////////////////////////////////////////////////////////////////////
// Usecases
// ///////////////////////////////////////////////////////////////////////////
type itemUsecase struct {
	items map[int]item
}

func newItemUsecase() *itemUsecase {
	return &itemUsecase{
		items: make(map[int]item),
	}
}

func (u *itemUsecase) saveItem(it item) error {
	_, err := u.itemExist(it.ID)
	if err == nil {
		return errors.New("saveItem: item already exist")
	}
	now := time.Now()
	it.CreatedAt = now
	it.UpdatedAt = now
	u.items[it.ID] = it
	return nil
}

func (u *itemUsecase) listItems() (map[int]item, error) {
	its := u.items
	if len(its) < 1 {
		return its, errNotFound
	}
	return its, nil
}

func (u *itemUsecase) itemById(id int) (item, error) {
	it, err := u.itemExist(id)
	if it == (item{}) {
		return it, err
	}
	return it, nil
}

func (u *itemUsecase) changeItem(it item) error {
	_, err := u.itemExist(it.ID)
	if err != nil {
		return err
	}
	it.UpdatedAt = time.Now()
	u.items[it.ID] = it
	return nil
}

func (u *itemUsecase) deleteItem(id int) error {
	it, err := u.itemExist(id)
	if it == (item{}) {
		return err
	}
	delete(u.items, id)
	return nil
}

func (u *itemUsecase) itemExist(id int) (item, error) {
	it, exist := u.items[id]
	if !exist {
		return item{}, errNotFound
	}
	return it, nil
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
