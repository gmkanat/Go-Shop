package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gmkanat/Go-Shop/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type ItemController struct {
	DB *gorm.DB
}

func NewItemController(DB *gorm.DB) ItemController {
	return ItemController{DB}
}

// GetItems [...] Get all items
func (ic *ItemController) GetItems(ctx *gin.Context) {
	var items []models.ItemList
	query := ic.DB.Table("items").
		Select("items.id, items.name, items.price, AVG(item_ratings.rating) as avg_rating").
		Joins("LEFT JOIN item_ratings ON items.id = item_ratings.item_id").
		Group("items.id").
		Order("items.id")

	ratingGTEFilter, err := strconv.Atoi(ctx.Query("rating_gte"))
	if err == nil {
		query.Where("rating >= ?", ratingGTEFilter).Find(&items)
	}
	ratingLTEFilter, err := strconv.Atoi(ctx.Query("rating_lte"))
	if err == nil {
		query.Where("rating <= ?", ratingLTEFilter).Find(&items)
	}
	priceGTEFilter, err := strconv.Atoi(ctx.Query("price_gte"))
	if err == nil {
		query.Where("price >= ?", priceGTEFilter)
	}
	priceLTEFilter, err := strconv.Atoi(ctx.Query("price_lte"))
	if err == nil {
		query.Where("price <= ?", priceLTEFilter)
	}

	searchFilter := strings.ToLower(ctx.Query("search"))
	if searchFilter != "" {
		query.Where("name ILIKE ?", "%"+searchFilter+"%")
	}
	query.Scan(&items)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "items": items})
}

// GetItem [...] Get item by id
func (ic *ItemController) GetItem(ctx *gin.Context) {
	var item models.Item
	itemID := ctx.Param("id")
	ic.DB.First(&item, itemID)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item": item})
}

// CreateItem [...] Create item
func (ic *ItemController) CreateItem(ctx *gin.Context) {
	var payload *models.ItemChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	newItem := models.Item{
		Name:  payload.Name,
		Price: payload.Price,
	}
	result := ic.DB.Create(&newItem)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item": newItem})
}

// UpdateItem [...] Update item
func (ic *ItemController) UpdateItem(ctx *gin.Context) {
	var payload *models.ItemChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	itemID := ctx.Param("id")
	var item models.Item
	ic.DB.First(&item, itemID)
	if payload.Price != 0 {
		item.Price = payload.Price
	}
	if payload.Name != "" {
		item.Name = payload.Name
	}
	result := ic.DB.Save(&item)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item": item})
}

// DeleteItem [...] Delete item
func (ic *ItemController) DeleteItem(ctx *gin.Context) {
	itemID := ctx.Param("id")
	var item models.Item
	ic.DB.First(&item, itemID)
	if item.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "item not found"})
		return
	}
	result := ic.DB.Delete(&item)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ic *ItemController) GiveRatingToItem(ctx *gin.Context) {
	itemID := ctx.Param("id")
	var item models.Item
	ic.DB.First(&item, itemID)
	if item.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "item not found"})
		return
	}
	var payload *models.ItemRatingChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	currentUser := ctx.MustGet("currentUser").(models.User)
	// check if user already give rating to this item
	var itemRating models.ItemRating
	ic.DB.Where("item_id = ? AND user_id = ?", item.ID, currentUser.ID).First(&itemRating)
	if itemRating.ID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "user already give rating to this item"})
		return
	}
	newItemRating := models.ItemRating{
		ItemID: item.ID,
		Rating: payload.Rating,
		User:   currentUser,
	}
	result := ic.DB.Create(&newItemRating)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	NewItemRatingResponse := models.ItemRatingResponse{
		ID:     newItemRating.ID,
		Rating: newItemRating.Rating,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item_rating": NewItemRatingResponse})
}

func (ic *ItemController) CommentItem(ctx *gin.Context) {
	itemID := ctx.Param("id")
	var item models.Item
	ic.DB.First(&item, itemID)
	if item.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "item not found"})
		return
	}
	var payload *models.ItemCommentChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	currentUser := ctx.MustGet("currentUser").(models.User)

	newItemComment := models.ItemComment{
		ItemID:  item.ID,
		Comment: payload.Comment,
		User:    currentUser,
	}
	result := ic.DB.Create(&newItemComment)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	NewItemCommentResponse := models.ItemCommentResponse{
		ID:      newItemComment.ID,
		Comment: newItemComment.Comment,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item_rating": NewItemCommentResponse})
}
