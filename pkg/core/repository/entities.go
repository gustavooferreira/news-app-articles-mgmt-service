package repository

import "time"

// Article represents the 'articles' table in the database.
type Article struct {
	GUID          string `gorm:"primaryKey;type:varchar(500);not null"`
	Provider      Provider
	ProviderID    uint64 `gorm:"not null"` // Foreign Key
	Category      Category
	CategoryID    uint64    `gorm:"not null"` // Foreign Key
	Title         string    `gorm:"type:varchar(500);not null"`
	Description   string    `gorm:"not null"`
	Link          string    `gorm:"type:varchar(500);not null"`
	PublishedDate time.Time `gorm:"index;not null"`
}

// Provider represents the 'providers' table in the database.
type Provider struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"type:varchar(30);uniqueIndex;not null"`
}

// Category represents the 'categories' table in the database.
type Category struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"type:varchar(30);uniqueIndex;not null"`
}
