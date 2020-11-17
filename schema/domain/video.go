package domain

type Video struct {
	VideoId    string `json:"id" gorm:"column:video_id"`
	CategoryID string `json:"categoryId" gorm:"column:category_id"`
	ETag       string `json:"etag" gorm:"column:etag"`
	Author     string `json:"author,omitempty" gorm:"column:author"`
}

func (v Video) TableName() string {
	// double check here, make sure the table does exist!!
	return "video" // default table name
}
