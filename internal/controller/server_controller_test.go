package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitshubham45/virtualServer/internal/db"
	"github.com/gitshubham45/virtualServer/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	Basic = "basic"
	Plus  = "plus"
	Prime = "prime"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// open an in memory(created only for testing period)
	//  SQlite database
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

	// cleanup function
	cleanup := func() {
		sqlDB, _ := testDB.DB()
		sqlDB.Close() // close the in-memory DB connection
		db.DB = originalDB
	}

	return testDB, cleanup
}

type mockErrorDB struct{}

// A helper method to simulate a GORM Create() returning an error
func (m *mockErrorDB) Create(value interface{}) *gorm.DB {
	return &gorm.DB{Error: errors.New("mock database connection failed")}
}

func TestCreateServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testDB , cleanup := setupTestDB(t)
	defer cleanup()

	serverConfig := map[string]interface{}{
		"region": "India",
		"type":   Prime,
	}

	jsonData, err := json.Marshal(serverConfig)
	if err != nil {
		t.Fatalf("Failed to encode json : %v", err)
		panic(err)
	}

	router := gin.Default()
	router.POST("/api/server", CreateServer)

	// --- Test case 1 : Server Creaetd ----
	t.Run("Server Created", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/server", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		//A Asserions
		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d for created , got %d", http.StatusCreated, rec.Code)
		}

		var response struct {
			Message string `json:"message"`
			ID      string `json:"id"`
			Status  string `json:"status"`
		}

		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Faield to unmarshal the response body: %v", err)
		}

		if response.Message != "success" {
			t.Errorf("Expectd message = 'success' , got '%s' ", response.Message)
		}
		if response.Status != "running" {
			t.Errorf("Expected status = 'running' , got '%s'", response.Status)
		}
	})

	// -- Server not created --
	t.Run("Invalid request", func(t *testing.T) {
		// check what is byte.Buffer()
		invalidJSON := []byte(`{"region": "India", "type":}`)
		req, err := http.NewRequest(http.MethodPost, "/api/server", bytes.NewBuffer(invalidJSON))
		if err != nil {
			t.Fatalf("Failed to create req:%v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Assertion
		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status as %d for internal server error , got %d", http.StatusBadRequest, rec.Code)
		}

		// var response struct{
		// 	Error error `json:"error"`
		// }

		// if err := json.Unmarshal(rec.Body.Bytes() , &response) ; err != nil {
		// 	t.Fatalf("Faield to unmarshal the response body: %v", err)
		// }

		// if err == nil {
		// 	t.Errorf("Expected error but not got any")
		// }
	})

	t.Run("Database Error", func(t *testing.T) {
		
		err := testDB.Migrator().DropTable(&models.Server{})
		if err != nil {
			t.Fatalf("Failed to drop table : %v", err)
		}

		validJSON := []byte(`{"region":"India" ,            "type" : "Prime"}`)
		fmt.Println(validJSON)
		req , err := http.NewRequest(http.MethodPost , "/api/server" , bytes.NewBuffer(validJSON)) // converted to stream of byte
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type" , "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec,req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d , got %d" , http.StatusInternalServerError , rec.Code)
		}
	})

}
func TestGetServersData(t *testing.T) {
	// set Gin to test mode to suppress debug output
	gin.SetMode(gin.TestMode)

	// setup the in memory test databse and get cleanup function
	testDB, cleanup := setupTestDB(t)
	defer cleanup()

	testServer := models.Server{
		ID:           uuid.New().String(),
		ServerNumber: 101,
		BillingRate:  5.0,
		Status:       "running",
		Region:       "US East",
		Type:         "basic",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert the tese server into in-memory db
	if err := testDB.Create(&testServer).Error; err != nil {
		t.Fatalf("Failed to create test server in DB: %v", err)
	}

	// create a new Gin router
	router := gin.Default()
	router.GET("/api/servers/:id", GetServersData)

	// --- Test Case 1 : Server  Found case(Success) ---
	t.Run("Server Found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/servers/"+testServer.ID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req) // call the handler

		// Assertions
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d for found seevr, got %d", http.StatusOK, rec.Code)
		}

		var response struct {
			Message string        `json:"message"`
			Server  models.Server `json:"server"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		if err != nil {
			t.Fatalf("Failed to unmarshal response body : %v", err)
		}

		// Basic check fot ID , more comprehensive checks would compare all fields
		if response.Server.ID != testServer.ID {
			t.Errorf("Expected server ID %s , got %s", testServer.ID, response.Server.ID)
		}
		if response.Server.Status != testServer.Status {
			t.Errorf("Expected server status %s , got %s", testServer.Status, response.Server.Status)
		}
	})

	// --- Test Case 2: Server Not Found ---
	t.Run("Server Not Found", func(t *testing.T) {
		nonExistentID := uuid.New().String() // A random ID that wont be in DB
		req, err := http.NewRequest(http.MethodGet, "/api/servers/"+nonExistentID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Assertions
		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d for not found server , got %d", http.StatusNotFound, rec.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		if response["message"] != fmt.Sprintf("Server with ID '%s' not found.", nonExistentID) {
			t.Errorf("Expected message 'Server not found' , got %s", response["message"])
		}
	})

	// --- Test Case 3: Invalid ID Format (Gin's param binding might catch this or it might proceed to DB) ---
	t.Run("Invalid ID Format", func(t *testing.T) {
		// If Gin allows non-UUID strings to pass as params, it will hit DB.
		// If you have validation middleware for UUID format, it would catch here.
		// For simplicity, we'll assume it hits DB and gets Not Found or Internal Error.
		invalidID := "not-a-uuid"
		req, err := http.NewRequest(http.MethodGet, "/api/servers/"+invalidID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Assertions (depends on your actual handler's behavior for invalid UUIDs before DB query)
		// Assuming GORM's First will return ErrRecordNotFound or a type error if UUID conversion fails internally
		if rec.Code != http.StatusNotFound && rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d or %d for invalid ID, got %d", http.StatusNotFound, http.StatusInternalServerError, rec.Code)
		}
	})

}
