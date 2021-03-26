package repository

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the database manager connecting to the database.
type Database struct {
	conn *gorm.DB
}

// NewDatabase returns a new Database.
func NewDatabase(host string, port int, username string, password string, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	// TODO: Setup logger for gorm here
	// I should implement GORM's logger interface on core.AppLogger and pass it to gorm config.
	// I should also pass log level to GORM.

	// dbconn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// dbconn = dbconn.Debug()
	dbconn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	dbconn = dbconn.Session(&gorm.Session{})
	db := Database{conn: dbconn}
	return &db, nil
}

// Close closes all database connections.
func (db *Database) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// HealthCheck checks whether the database is still around.
func (db *Database) HealthCheck() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

// FindAllArticleRecords finds all the article records.
func (db *Database) FindAllArticleRecords(provider string, category string) ([]Article, error) {
	var articleResults []Article
	chain := db.conn.Joins("Provider").Joins("Category")

	if provider != "" {
		chain = chain.Where("`Provider`.`name` = ?", provider)
	}

	if category != "" {
		chain = chain.Where("`Category`.`name` = ?", category)
	}

	// chain = chain.Where(&Feed{Enabled: &enabled})
	result := chain.Find(&articleResults)
	return articleResults, result.Error
}

// InsertArticleRecord inserts a new article record in the database.
func (db *Database) InsertArticleRecord(guid string, provider string, category string) error {
	// Add Provider if it doesn't exist
	var providerRecord Provider
	result := db.conn.Where(Provider{Name: provider}).FirstOrCreate(&providerRecord)
	if result.Error != nil {
		return result.Error
	}

	// Add Category if it doesn't exist
	var categoryRecord Category
	result = db.conn.Where(Category{Name: category}).FirstOrCreate(&categoryRecord)
	if result.Error != nil {
		return result.Error
	}

	articleRecord := Article{
		GUID:       guid,
		ProviderID: providerRecord.ID,
		CategoryID: categoryRecord.ID,
	}

	result = db.conn.Create(&articleRecord)
	return result.Error
}
