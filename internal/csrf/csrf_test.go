package csrf

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerate(t *testing.T) {
	conn := &gin.Context{}
	Generate()(conn)
}
