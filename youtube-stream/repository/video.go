package repository

import (
	"context"

	"github.com/dungtc/kafka-playground/youtube-stream/domain"
	"github.com/dungtc/kafka-playground/youtube-stream/storage"
)

type VideoRepository interface {
	Save(ctx context.Context, video *domain.Video) error
}

func NewVideoRepository(storage storage.Storage) VideoRepository {

	storage.GetDBConn().AutoMigrate(&domain.Video{})
	return videoRepository{
		storage,
	}
}

type videoRepository struct {
	storage storage.Storage
}

func (y videoRepository) Save(ctx context.Context, video *domain.Video) error {
	return y.storage.WithContext(ctx).Create(video).Error
}
