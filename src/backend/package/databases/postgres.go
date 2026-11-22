package databases

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)

	for attempt := 1; attempt <= 5; attempt++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			TranslateError: true,
		})
		if err == nil {
			return db, nil
		}

		fmt.Printf("Failed to connect to database: %v attempt: %d\n", err, attempt)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database")
}
