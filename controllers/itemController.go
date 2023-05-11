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
		Select("items.id, items.name, items.price, AVG(item_ratings.rating) as avg_rating, users.name as seller_name").
		Joins("LEFT JOIN item_ratings ON items.id = item_ratings.item_id").
		Joins("INNER JOIN users ON items.seller_id = users.id").
		Group("items.id, users.name").
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
	var item models.ItemDetail
	itemID := ctx.Param("id")
	ic.DB.Raw(`
		  SELECT 
			items.id, 
			items.name, 
			items.price, 
			AVG(item_ratings.rating) as avg_rating, 
			users.name as seller_name,
			json_agg(item_comments.comment) as comments
		  FROM items
		  LEFT JOIN item_ratings ON items.id = item_ratings.item_id
		  INNER JOIN users ON items.seller_id = users.id
		  LEFT JOIN item_comments ON items.id = item_comments.item_id
		  WHERE items.id = ?
		  GROUP BY items.id, users.name
		`, itemID).Scan(&item)
	ic.DB.Table("item_comments").
		Select("item_comments.comment, users.email as email").
		Joins("INNER JOIN users ON item_comments.user_id = users.id").
		Where("item_comments.item_id = ?", itemID).
		Order("item_comments.ID DESC").Scan(&item.Comments)

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
		Name:     payload.Name,
		Price:    payload.Price,
		SellerID: ctx.MustGet("currentUser").(models.User).ID,
	}

	result := ic.DB.Create(&newItem)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	newItemResponse := models.ItemChange{
		Name:  payload.Name,
		Price: payload.Price,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item": newItemResponse})
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
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "item_comment": NewItemCommentResponse})
}

func (ic *ItemController) PurchaseItem(ctx *gin.Context) {
	itemID := ctx.Param("id")
	var item models.Item
	ic.DB.First(&item, itemID)
	if item.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "item not found"})
		return
	}
	currentUser := ctx.MustGet("currentUser").(models.User)

	newOrder := models.Order{
		ItemID: item.ID,
		UserID: currentUser.ID,
	}
	result := ic.DB.Create(&newOrder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	NewOrderResponce := models.OrderResponse{
		ID:     newOrder.ID,
		Status: newOrder.Status,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "order": NewOrderResponce})
}

func (ic *ItemController) OrderStatus(ctx *gin.Context) {
	orderID := ctx.Param("id")
	var order models.Order
	ic.DB.First(&order, orderID)
	if order.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "order not found"})
		return
	}

	var payload *models.OrderChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	order.Status = payload.Status

	result := ic.DB.Save(order)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}
	NewOrderResponce := models.OrderResponse{
		ID:     order.ID,
		Status: order.Status,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "order": NewOrderResponce})
}
