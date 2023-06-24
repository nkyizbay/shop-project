package item

import (
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type Repository interface {
	Create(item *Item) error
	Delete(id int) error
}

type defaultRepository struct {
	database *gorm.DB
}

func NewItemRepository(database *gorm.DB) Repository {
	return &defaultRepository{database: database}
}

func (t *defaultRepository) Create(item *Item) error {

	if err := t.database.Model(&Item{}).Create(item).Error; err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (t *defaultRepository) Delete(id int) error {

	if err := t.database.Delete(&Item{}, id).Error; err != nil {
		log.Error(err)
		return err
	}

	return nil
}
