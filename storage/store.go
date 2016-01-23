package storage

import (
	"github.com/t-yuki/zipkin-go/models"
	"github.com/t-yuki/zipkin-go/storage/mysql"
)

func Open() (Storage, error) {
	stor, err := mysql.Open()
	return stor, err
}

type Storage interface {
	StoreSpans(spans models.ListOfSpans) error
	Close() error
}
