package models

import "time"

// CollectionStatus 作品狀態
type CollectionStatus string

const (
	CollectionStatusPublished CollectionStatus = "published"
	CollectionStatusDraft     CollectionStatus = "draft"
	CollectionStatusDeleted   CollectionStatus = "deleted"
)

// Collection 資料
type Collection struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Title       string           `json:"title" gorm:"not null"`
	Description string           `json:"description" gorm:"not null"`
	ThumbnailID uint             `json:"thumbnail_id" gorm:"not null"`
	Thumbnail   Image            `json:"thumbnail" gorm:"foreignKey:ThumbnailID"`
	Images      []Image          `json:"images" gorm:"many2many:collection_images"`
	Status      CollectionStatus `json:"status" gorm:"not null"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
