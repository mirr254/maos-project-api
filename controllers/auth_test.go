package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"maos-cloud-project-api/models"
	"maos-cloud-project-api/utils"

	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {

	r := utils.SetUpRouter()
	r.POST("/signup", Signup)
	// userID := xid.New().String()
	user := models.User{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "test1234",
		Role:     "admin",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

}

func TestDashboard(t *testing.T) {
	// mockSucessResponse := `{
	// 	"success": "customer dashboard",
	// 	"role":    "admin"
	// }`

	mockUnAuthorizedResponse := `{"error":"unauthorized"}`

	r := utils.SetUpRouter()
	r.GET("/dashboard", Dashboard)
	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockUnAuthorizedResponse, string(responseData))
	assert.Equal(t, 401, w.Code)

}
