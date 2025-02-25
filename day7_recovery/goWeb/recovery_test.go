package goWeb

import (
	"net/http"
	"testing"
)

func TestRecovery(t *testing.T) {
	r := Default()
	r.GET("/", func(c *Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
