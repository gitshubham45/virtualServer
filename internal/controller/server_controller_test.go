package controller

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/db"
	"github.com/gitshubham45/virtualServer/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// open an in memory  SQlite database a cleanup function
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test databse : %v", err)
	}
	// Migrate models to the test databse
	err = testDB.AutoMigrate(&models.Server{}, &models.ServerLog{})
	if err != nil {
		t.Fatalf("Failed to auto migrate test schema: %v", err)
	}

	originalDB := db.DB
	db.DB = testDB

	cleanup := func(){
		sqlDB , _ := testDB.DB()
		sqlDB.Close() // close the in-memory DB connection
		db.DB = originalDB
	}
	
	return testDB , cleanup
}

func TestGetServersData(t *testing.T){
	// set Gin to test mode to suppress debug output
	gin.SetMode(gin.TestMode)

	// setup the in memory test databse and get cleanup function
	testDB , cleanup := setupTestDB(t)
	defer cleanup()

	testServer := models.Server{
		ID : uuid.New().String(),
		ServerNumber: 101,
		B
	}
}