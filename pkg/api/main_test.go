package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	registerTagName()
	registerValidDObValidator()
	registerValidEducationLevelFieldValidator()
	os.Exit(m.Run())
}
