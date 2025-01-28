package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inv/internal/app"
	"inv/internal/tenant/domain"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"inv/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ResponseBody struct {
	ID        string     `json:"ID"`
	CreatedAt string     `json:"CreatedAt"`
	UpdatedAt string     `json:"UpdatedAt"`
	DeletedAt *time.Time `json:"DeletedAt"`
	Email     string     `json:"Email"`
	Name      string     `json:"Name"`
	Password  string     `json:"Password"`
	Tenants   []any      `json:"Tenants"`
}

type TenantResponse struct {
	ID            string  `json:"ID"`
	CreatedAt     string  `json:"CreatedAt"`
	UpdatedAt     string  `json:"UpdatedAt"`
	DeletedAt     *string `json:"DeletedAt"`
	TenantOwnerID string  `json:"tenant_owner_id"`
	DatabaseName  string  `json:"database_name"`
	DatabaseHost  string  `json:"database_host"`
	DatabasePort  string  `json:"database_port"`
	DatabaseUser  string  `json:"database_user"`
	DatabasePass  string  `json:"database_pass"`
	Status        string  `json:"status"`
	Owner         *string `json:"Owner"`
}

type TestTenantConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var testApp *app.App

func TestMain(m *testing.M) {
	// Setup test environment
	setTestEnv()

	if err := setupTest(); err != nil {
		fmt.Printf("Failed to setup test: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	cleanDatabase()
	os.Exit(code)
}

func setTestEnv() {
	os.Setenv("APP_DB_HOST", "localhost")
	os.Setenv("APP_DB_PORT", "5432")
	os.Setenv("APP_DB_USER", "postgres")
	os.Setenv("APP_DB_PASSWORD", "postgres")
	os.Setenv("APP_DB_NAME", "inventory_test")
	os.Setenv("APP_DB_SSLMODE", "disable")
	os.Setenv("GIN_MODE", "test")
}

func setupTest() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create test app
	testApp, err = app.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	// Clean database before tests
	if err := cleanDatabase(); err != nil {
		return fmt.Errorf("failed to clean database: %w", err)
	}

	return nil
}

func cleanDatabase() error {
	db := testApp.GetDB()

	// Drop all tables
	if err := db.Exec("DROP SCHEMA public CASCADE").Error; err != nil {
		return err
	}

	// Recreate public schema
	if err := db.Exec("CREATE SCHEMA public").Error; err != nil {
		return err
	}

	// Run migrations
	if err := testApp.GetDB().AutoMigrate(
		&domain.Tenant{},
		&domain.TenantOwner{},
	); err != nil {
		return err
	}

	return nil
}

func TestTenantFlow(t *testing.T) {
	var ownerID string
	var tenantID string
	var tenantDatabaseName string
	var tenantConfig TestTenantConfig

	defer func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			tenantConfig.Host,
			tenantConfig.User,
			tenantConfig.Password,
			tenantConfig.DBName,
			tenantConfig.Port,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("Failed to get database instance: %v", err)
		}
		defer sqlDB.Close()

		// Drop the tenant database
		dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", tenantConfig.DBName)
		if err := db.Exec(dropQuery).Error; err != nil {
			t.Fatalf("Failed to drop tenant database: %v", err)
		}
	}()

	// Create Tenant Owner
	t.Run("Create Tenant Owner", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"email":    "test.owner@example.com",
			"name":     "Test Owner",
			"password": "password123",
		}

		w := performRequest(testApp.Router(), "POST", "/api/tenants/owners", requestBody)

		var response ResponseBody
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		ownerID = response.ID

		if w.Code != 201 {
			t.Fatalf("Expected status code 201, got %d", w.Code)
		}

	})

	// Get Created Owner
	t.Run("Get Created Owner", func(t *testing.T) {
		w := performRequest(testApp.Router(), "GET", "/api/tenants/owners/"+ownerID, nil)
		fmt.Println(w.Body.String())
		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
	})

	// Update Owner
	t.Run("Update Owner", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"email":    "test.owner2@example.com",
			"name":     "Owner2",
			"password": "newsecret123",
		}

		w := performRequest(testApp.Router(), "PUT", "/api/tenants/owners/"+ownerID, requestBody)

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
	})

	// Create Tenant
	t.Run("Create Tenant", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"tenant_owner_id": ownerID,
		}

		w := performRequest(testApp.Router(), "POST", "/api/tenants", requestBody)

		var tenantResp TenantResponse
		if err := json.Unmarshal(w.Body.Bytes(), &tenantResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		tenantID = tenantResp.ID
		tenantDatabaseName = tenantResp.DatabaseName

		tenantConfig = TestTenantConfig{
			Host:     tenantResp.DatabaseHost,
			Port:     tenantResp.DatabasePort,
			DBName:   tenantResp.DatabaseName,
			User:     tenantResp.DatabaseUser,
			Password: tenantResp.DatabasePass,
		}

		if w.Code != 201 {
			t.Fatalf("Expected status code 201, got %d", w.Code)
		}

	})

	// Get Created Tenant
	t.Run("Get Created Tenant", func(t *testing.T) {
		w := performRequest(testApp.Router(), "GET", "/api/tenants/"+tenantID, nil)

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
	})

	// List All Tenants
	t.Run("List All Tenants", func(t *testing.T) {
		w := performRequest(testApp.Router(), "GET", "/api/tenants", nil)

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
	})

	// Update Tenant
	t.Run("Update Tenant", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"tenant_owner_id": ownerID,
			"database_name":   tenantDatabaseName,
			"database_host":   "localhost",
			"database_port":   "5435",
		}

		w := performRequest(testApp.Router(), "PUT", "/api/tenants/"+tenantID, requestBody)

		// fmt.Println(w.Body.String())

		if w.Code != 200 {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
	})

	// Delete Tenant
	t.Run("Delete Tenant", func(t *testing.T) {
		w := performRequest(testApp.Router(), "DELETE", "/api/tenants/"+tenantID, nil)

		// fmt.Println(w.Body.String())
		if w.Code != 204 {
			t.Fatalf("Expected status code 204, got %d", w.Code)
		}
	})

	// Delete Owner
	t.Run("Delete Owner", func(t *testing.T) {
		w := performRequest(testApp.Router(), "DELETE", "/api/tenants/owners/"+ownerID, nil)

		fmt.Println(w.Body.String())

		if w.Code != 204 {
			t.Fatalf("Expected status code 204, got %d", w.Code)
		}
	})

	// Verify Tenant Deletion
	t.Run("Verify Tenant Deletion", func(t *testing.T) {
		w := performRequest(testApp.Router(), "GET", "/api/tenants/"+tenantID, nil)

		if w.Code != 404 {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}
	})

	// Verify Owner Deletion
	t.Run("Verify Owner Deletion", func(t *testing.T) {
		w := performRequest(testApp.Router(), "GET", "/api/tenants/owners/"+ownerID, nil)

		if w.Code != 404 {
			t.Fatalf("Expected status code 404, got %d", w.Code)
		}
	})
}

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	w := httptest.NewRecorder()

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Add("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
		req.Header.Add("Content-Type", "application/json")
	}

	r.ServeHTTP(w, req)
	return w
}
