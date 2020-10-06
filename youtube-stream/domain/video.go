package domain

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model
	VideoId    string `json:"id" gorm:"column:video_id"`
	CategoryID string `json:"categoryId" gorm:"column:category_id"`
	ETag       string `json:"etag" gorm:"column:etag"`
}

func (v Video) TableName() string {
	// double check here, make sure the table does exist!!
	return "video" // default table name
}
