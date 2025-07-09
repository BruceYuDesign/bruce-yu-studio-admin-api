package models

import "time"

// ImageStatus 圖片狀態
type ImageStatus string

const (
	ImageStatusActive  ImageStatus = "active"
	ImageStatusDeleted ImageStatus = "deleted"
)

// ImageFormat 圖片格式
type ImageFormat string

const (
	ImageFormatJpeg ImageFormat = "jpeg"
	ImageFormatJpg  ImageFormat = "jpg"
	ImageFormatPng  ImageFormat = "png"
)

// Image 資料
type Image struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Url       string      `json:"url" gorm:"not null"`
	Alt       string      `json:"alt" gorm:"not null"`
	Width     int         `json:"width" gorm:"not null"`
	Height    int         `json:"height" gorm:"not null"`
	Format    ImageFormat `json:"format" gorm:"not null"`
	Status    ImageStatus `json:"status" gorm:"not null"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
