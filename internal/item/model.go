package item

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID          int     `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Price       float64 `gorm:"not null;check:price>0" json:"price"`
	Description string  `gorm:"not null" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (t *Item) IsNameEmpty() bool {
	return t.Name == ""
}

func (t *Item) IsDescriptionEmpty() bool {
	return t.Description == ""
}

func (t *Item) CheckFieldsEmpty() bool {
	return t.IsNameEmpty() || t.IsDescriptionEmpty()
}

func (t *Item) IsValidPrice() bool {
	return t.Price >= 0
}

func (t *Item) IsInvalidPrice() bool {
	return !t.IsValidPrice()
}

func IsValidID(id int) bool {
	return id > 0
}

func IsInvalidID(id int) bool {
	return !IsValidID(id)
}
