package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project-gotham/config"
	"project-gotham/controllers"
	"project-gotham/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // "_" correspond à l'erreur qu'on ne traite pas.
	db.AutoMigrate(&models.Project{}, &models.Comment{})

	project := models.Project{Name: "Projet test", Description: "Description test"}
	db.Create(&project)

	comment := models.Comment{ProjectID: project.ID, Content: "Commentaire test"}
	db.Create(&comment)

	return db
}

func TestGetProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.DB = setupTestDB()

	r := gin.Default()
	r.GET("/projects", controllers.GetProjects)

	req, _ := http.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	assert.Contains(t, body, "Projet test")
	assert.Contains(t, body, "Commentaire test")

}

func TestPostProject(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.DB = setupTestDB()

	r := gin.Default()
	r.POST("/projects", controllers.PostProject)

	project := map[string]interface{}{
		"name":        "Test project",
		"description": "Test project description",
		"skills":      []string{"Golang", "Testing", "SQLite"},
	}

	data, _ := json.Marshal(project)

	req, _ := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.Contains(t, w.Body.String(), "Test project")
	assert.Contains(t, w.Body.String(), "Test project description")
	assert.Contains(t, w.Body.String(), "Testing")

}
