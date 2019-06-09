package app

import (
	"github.com/fernetbalboa/arqweb/src/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	config.Scope = config.TestScope
}

func TestRecoverFromPanics(t *testing.T) {
	r := config.ConfiguredRouter()

	r.GET("/panic", panicFun)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected 500 from panic recovery")
}

func panicFun(_ *gin.Context) {
	panic("Panic!")
}
