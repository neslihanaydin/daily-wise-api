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

//go:embed data/advices.json
var advicesFS embed.FS

type Advice struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

var advices []Advice

func loadAdvices() {
	// 1) Dosyayı oku
	data, err := advicesFS.ReadFile("data/advices.json")
	if err != nil {
		log.Fatalf("⚠️ ReadFile error: %v", err)
	}
	log.Printf("🔍 ReadFile succeeded, %d bytes read", len(data))

	// 2) Eğer dosya boşsa boylece görebiliriz
	if len(data) == 0 {
		log.Fatal("⚠️ advices.json içeriği boş")
	}

	// 3) JSON’u parse et
	if err := json.Unmarshal(data, &advices); err != nil {
		log.Fatalf("⚠️ JSON unmarshal error: %v", err)
	}
	log.Printf("✅ Loaded %d advice entries", len(advices))
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
