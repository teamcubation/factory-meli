package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// Injetando InMemoryDatabase no Usecase
	repo := inMemoryDatabase{
		items: make(map[int]*item),
	}
	u := newItemUsecase(repo)
	h := newHandler(u)

	router := gin.Default()

	router.GET("/", h.home)

	router.POST("/item", h.create)
	router.GET("item/:id", h.findById)
	router.PUT("/item/:id", h.update)
	router.DELETE("/item/:id", h.delete)
	router.GET("/item", h.findAll)

	log.Println("Server started at http://localhost:8080/")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Global error
// /////////////////////////////////////////////////////////////////////////////
type Error struct {
	message string
	code    int
}

func (e Error) Error() string {
	return e.message
}

var errNotFound = Error{message: "resource not found", code: http.StatusNotFound}

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

func (h *handler) create(c *gin.Context) {
	var it createItemRequestDto

	if err := c.BindJSON(&it); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.usecase.create(it)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *handler) update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var it updateItemRequestDto

	if err := c.BindJSON(&it); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it.ID = id

	response, err := h.usecase.update(it)
	if err != nil {
		if errType, ok := err.(Error); ok {
			c.JSON(errType.code, gin.H{"error": errType.message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *handler) delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it := deleteRequestDto{id}

	err = h.usecase.delete(it)
	if err != nil {
		if errType, ok := err.(Error); ok {
			c.JSON(errType.code, gin.H{"error": errType.message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Item deleted successfully"})
}

func (h *handler) findById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it := findByIdRequestDto{id}

	its, err := h.usecase.findById(it)
	if err != nil {
		if errType, ok := err.(Error); ok {
			c.JSON(errType.code, gin.H{"error": errType.message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, its)
}

func (h *handler) findAll(c *gin.Context) {
	its, err := h.usecase.findAll()
	if err != nil {
		if err == errNotFound {
			c.JSON(errNotFound.code, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, its)
}

// ///////////////////////////////////////////////////////////////////////////
// Usecases
// ///////////////////////////////////////////////////////////////////////////
type itemUsecase struct {
	repo ItemsDatabaseRepository
}

func newItemUsecase(repo ItemsDatabaseRepository) *itemUsecase {
	return &itemUsecase{
		repo: repo,
	}
}

func (u *itemUsecase) create(it createItemRequestDto) (createItemResponseDto, error) {
	result, err := u.repo.Create(it)
	if err != nil {
		return createItemResponseDto{}, err
	}
	return result, nil
}

func (u *itemUsecase) update(it updateItemRequestDto) (updateItemResponseDto, error) {
	response, err := u.repo.Update(it)
	if err != nil {
		return updateItemResponseDto{}, err
	}
	return response, nil
}

func (u *itemUsecase) delete(it deleteRequestDto) error {
	err := u.repo.Delete(it)
	if err != nil {
		return err
	}
	return nil
}

func (u *itemUsecase) findById(it findByIdRequestDto) (findByIdResponseDto, error) {
	response, err := u.repo.FindById(it)
	if err != nil {
		return findByIdResponseDto{}, err
	}
	return response, nil
}

func (u *itemUsecase) findAll() ([]findByIdResponseDto, error) {
	its, err := u.repo.FindAll()
	if err != nil {
		return []findByIdResponseDto{}, err
	}
	return its, nil
}

// ///////////////////////////////////////////////////////////////////////////
// DTOs
// ///////////////////////////////////////////////////////////////////////////
type createItemRequestDto struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
}

type createItemResponseDto struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type updateItemRequestDto struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
}

type updateItemResponseDto struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type deleteRequestDto struct {
	ID int
}

type findByIdRequestDto struct {
	ID int `json:"id"`
}

type findByIdResponseDto struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// ///////////////////////////////////////////////////////////////////////////
// Domain
// ///////////////////////////////////////////////////////////////////////////
// Item entity
type item struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func itemFactory(id int, code string, title string, description string, price float64, stock int, status string, createdAt time.Time) *item {
	return &item{id, code, title, description, price, stock, status, createdAt, nil, nil}
}
func (it *item) IsDeleted() bool {
	return it.DeletedAt != nil
}

// ////////////////////////////////////////////////////////////////////////////
// Repositories Interface
// ////////////////////////////////////////////////////////////////////////////
type ItemsDatabaseRepository interface {
	Create(createItemRequestDto createItemRequestDto) (createItemResponseDto, error)
	Update(updateItemRequestDto updateItemRequestDto) (updateItemResponseDto, error)
	Delete(deleteRequestDto deleteRequestDto) error
	FindById(findByIdRequestDto) (findByIdResponseDto, error)
	FindAll() ([]findByIdResponseDto, error)
}

// ////////////////////////////////////////////////////////////////////////////
// InMemoryImplementation
// ////////////////////////////////////////////////////////////////////////////
type inMemoryDatabase struct {
	items map[int]*item
}

func (i inMemoryDatabase) Create(createItemRequestDto createItemRequestDto) (createItemResponseDto, error) {
	id := len(i.items)
	newItem := itemFactory(
		id,
		createItemRequestDto.Code,
		createItemRequestDto.Title,
		createItemRequestDto.Description,
		createItemRequestDto.Price,
		createItemRequestDto.Stock,
		createItemRequestDto.Status,
		time.Now(),
	)
	i.items[id] = newItem

	return createItemResponseDto{
		ID:          newItem.ID,
		Code:        newItem.Code,
		Title:       newItem.Title,
		Description: newItem.Description,
		Price:       newItem.Price,
		Stock:       newItem.Stock,
		Status:      newItem.Status,
		CreatedAt:   newItem.CreatedAt,
	}, nil
}

func (i inMemoryDatabase) Update(updateItemRequestDto updateItemRequestDto) (updateItemResponseDto, error) {
	item, ok := i.items[updateItemRequestDto.ID]
	if !ok || item.IsDeleted() {
		return updateItemResponseDto{}, errNotFound
	}
	item.Code = updateItemRequestDto.Code
	item.Title = updateItemRequestDto.Title
	item.Description = updateItemRequestDto.Description
	item.Price = updateItemRequestDto.Price
	item.Stock = updateItemRequestDto.Stock
	item.Status = updateItemRequestDto.Status

	updatedAt := time.Now()
	item.UpdatedAt = &updatedAt

	return updateItemResponseDto{
		ID:          item.ID,
		Code:        item.Code,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   &updatedAt,
	}, nil
}

func (i inMemoryDatabase) Delete(deleteRequestDto deleteRequestDto) error {
	item, ok := i.items[deleteRequestDto.ID]
	if !ok {
		return errNotFound
	}

	deletedAt := time.Now()
	item.DeletedAt = &deletedAt
	return nil
}

func (i inMemoryDatabase) FindById(request findByIdRequestDto) (findByIdResponseDto, error) {
	item, ok := i.items[request.ID]
	if !ok {
		return findByIdResponseDto{}, errNotFound
	}

	if item.DeletedAt != nil {
		return findByIdResponseDto{}, errNotFound
	}

	return findByIdResponseDto{
		ID:          item.ID,
		Code:        item.Code,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Stock:       item.Stock,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (i inMemoryDatabase) FindAll() ([]findByIdResponseDto, error) {
	response := []findByIdResponseDto{}
	for _, item := range i.items {
		if !item.IsDeleted() {
			itemResponse := findByIdResponseDto{
				ID:          item.ID,
				Code:        item.Code,
				Title:       item.Title,
				Description: item.Description,
				Price:       item.Price,
				Stock:       item.Stock,
				Status:      item.Status,
				CreatedAt:   item.CreatedAt,
				UpdatedAt:   item.UpdatedAt,
			}
			response = append(response, itemResponse)
		}
	}
	return response, nil
}
