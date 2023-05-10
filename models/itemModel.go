package models

type Item struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"not null" json:"name"`
	Price    float64 `gorm:"not null" json:"price"`
	SellerID uint    `gorm:"not null" json:"seller_id"`
	Seller   User    `json:"seller" gorm:"foreignKey:SellerID"`
}

type ItemList struct {
	Item
	AvgRating float64 `json:"avg_rating"`
}

type ItemChange struct {
	Name  string  `gorm:"not null" json:"name"`
	Price float64 `gorm:"not null" json:"price"`
}

type ItemRating struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	UserID uint    `gorm:"not null" json:"user_id"`
	ItemID uint    `gorm:"not null" json:"item_id"`
	Rating float64 `gorm:"not null" json:"rating"`
	User   User    `gorm:"foreignKey:UserID" json:"user"`
	Item   Item    `gorm:"foreignKey:ItemID" json:"item"`
}

type ItemRatingChange struct {
	Rating float64 `gorm:"not null" json:"rating"`
}

type ItemRatingResponse struct {
	ID     uint    `json:"id"`
	Rating float64 `json:"rating"`
}

type ItemComment struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	ItemID  uint   `gorm:"not null" json:"item_id"`
	Comment string `gorm:"not null" json:"comment"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	Item    Item   `gorm:"foreignKey:ItemID" json:"item"`
}

type ItemCommentChange struct {
	Comment string `gorm:"not null" json:"comment"`
}

type ItemCommentResponse struct {
	ID      uint   `json:"id"`
	Comment string `json:"comment"`
}

type Order struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	ItemID  uint   `gorm:"not null" json:"item_id"`
	Comment string `gorm:"not null" json:"comment"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	Status  string `gorm:"default:'pending';noy null" json:"status"`
}

type OrderResponse struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}

type OrderChange struct {
	Status string `json:"status"`
}
