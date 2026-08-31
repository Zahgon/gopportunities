package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/arthur404dev/gopportunities/handler"
	"github.com/arthur404dev/gopportunities/router"
	"github.com/arthur404dev/gopportunities/schemas"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbCounter atomic.Int64

func newTestServer(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()

	dsn := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", dbCounter.Add(1))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&schemas.Opening{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	srv := router.SetupRouter()
	handler.SetDB(db)
	return srv, db
}

func do(t *testing.T, srv http.Handler, method, target string, body interface{}) (*httptest.ResponseRecorder, map[string]interface{}) {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, target, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	var parsed map[string]interface{}
	if rec.Body.Len() > 0 {
		_ = json.Unmarshal(rec.Body.Bytes(), &parsed)
	}
	return rec, parsed
}

func seedOpening(t *testing.T, db *gorm.DB) schemas.Opening {
	t.Helper()
	opening := schemas.Opening{
		Role:     "Backend Engineer",
		Company:  "Acme",
		Location: "Remote",
		Remote:   true,
		Link:     "https://acme.example/jobs/1",
		Salary:   120000,
	}
	if err := db.Create(&opening).Error; err != nil {
		t.Fatalf("failed to seed opening: %v", err)
	}
	return opening
}

func TestCreateOpeningSuccess(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodPost, "/api/v1/opening", map[string]interface{}{
		"role":     "Backend Engineer",
		"company":  "Acme",
		"location": "Remote",
		"remote":   true,
		"link":     "https://acme.example/jobs/1",
		"salary":   120000,
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "operation from handler: create-opening successfull" {
		t.Fatalf("unexpected message: %v", got)
	}
	if _, ok := parsed["data"]; !ok {
		t.Fatalf("expected data field in response, got: %v", parsed)
	}
	if ct := rec.Header().Get("Content-Type"); ct == "" {
		t.Fatalf("expected a Content-Type header to be set")
	}
}

func TestCreateOpeningValidationError(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodPost, "/api/v1/opening", map[string]interface{}{
		"company":  "Acme",
		"location": "Remote",
		"remote":   true,
		"link":     "https://acme.example/jobs/1",
		"salary":   120000,
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "param: role (type: string) is required" {
		t.Fatalf("unexpected validation message: %v", got)
	}
	if got := parsed["errorCode"]; got != float64(http.StatusBadRequest) {
		t.Fatalf("expected errorCode 400, got: %v", got)
	}
}

func TestShowOpeningSuccess(t *testing.T) {
	srv, db := newTestServer(t)
	opening := seedOpening(t, db)

	rec, parsed := do(t, srv, http.MethodGet, "/api/v1/opening?id=1", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "operation from handler: show-opening successfull" {
		t.Fatalf("unexpected message: %v", got)
	}
	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got: %v", parsed["data"])
	}
	if data["Role"] != opening.Role {
		t.Fatalf("expected role %q, got %v", opening.Role, data["Role"])
	}
}

func TestShowOpeningMissingID(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodGet, "/api/v1/opening", nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "param: id (type: queryParameter) is required" {
		t.Fatalf("unexpected message: %v", got)
	}
}

func TestShowOpeningNotFound(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodGet, "/api/v1/opening?id=999", nil)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "opening not found" {
		t.Fatalf("unexpected message: %v", got)
	}
}

func TestListOpenings(t *testing.T) {
	srv, db := newTestServer(t)
	seedOpening(t, db)
	seedOpening(t, db)

	rec, parsed := do(t, srv, http.MethodGet, "/api/v1/openings", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "operation from handler: list-openings successfull" {
		t.Fatalf("unexpected message: %v", got)
	}
	data, ok := parsed["data"].([]interface{})
	if !ok {
		t.Fatalf("expected data array, got: %v", parsed["data"])
	}
	if len(data) != 2 {
		t.Fatalf("expected 2 openings, got %d", len(data))
	}
}

func TestUpdateOpeningSuccess(t *testing.T) {
	srv, db := newTestServer(t)
	seedOpening(t, db)

	rec, parsed := do(t, srv, http.MethodPut, "/api/v1/opening?id=1", map[string]interface{}{
		"role": "Staff Engineer",
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "operation from handler: update-opening successfull" {
		t.Fatalf("unexpected message: %v", got)
	}
	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got: %v", parsed["data"])
	}
	if data["Role"] != "Staff Engineer" && data["role"] != "Staff Engineer" {
		t.Fatalf("expected updated role, got: %v", data)
	}
}

func TestUpdateOpeningMissingID(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodPut, "/api/v1/opening", map[string]interface{}{
		"role": "Staff Engineer",
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "param: id (type: queryParameter) is required" {
		t.Fatalf("unexpected message: %v", got)
	}
}

func TestUpdateOpeningNotFound(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodPut, "/api/v1/opening?id=999", map[string]interface{}{
		"role": "Staff Engineer",
	})

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "opening not found" {
		t.Fatalf("unexpected message: %v", got)
	}
}

func TestDeleteOpeningSuccess(t *testing.T) {
	srv, db := newTestServer(t)
	seedOpening(t, db)

	rec, parsed := do(t, srv, http.MethodDelete, "/api/v1/opening?id=1", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "operation from handler: delete-opening successfull" {
		t.Fatalf("unexpected message: %v", got)
	}

	var count int64
	db.Model(&schemas.Opening{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected opening to be deleted, %d remain", count)
	}
}

func TestDeleteOpeningMissingID(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodDelete, "/api/v1/opening", nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got != "param: id (type: queryParameter) is required" {
		t.Fatalf("unexpected message: %v", got)
	}
}

func TestDeleteOpeningNotFound(t *testing.T) {
	srv, _ := newTestServer(t)

	rec, parsed := do(t, srv, http.MethodDelete, "/api/v1/opening?id=999", nil)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if got := parsed["message"]; got == nil {
		t.Fatalf("expected a not-found message, got nil")
	}
}
