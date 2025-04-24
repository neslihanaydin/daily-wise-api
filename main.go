package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var advicesFS embed.FS

type Advice struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

var advices []Advice

func loadAdvices() {
	b, _ := advicesFS.ReadFile("data/advices.json")
	if err := json.Unmarshal(b, &advices); err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadAdvices()
	r := gin.Default()

	r.GET("/advices", func(c *gin.Context) {
		c.JSON(http.StatusOK, advices)
	})

	r.GET("/advices/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for _, a := range advices {
			if a.ID == id {
				c.JSON(http.StatusOK, a)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})

	r.GET("/advices/today", func(c *gin.Context) {
		idx := time.Now().Day() % len(advices)
		c.JSON(http.StatusOK, advices[idx])
	})

	r.Run()
}
