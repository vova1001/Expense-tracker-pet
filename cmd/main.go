package main

import (
	"time"

	rout "github.com/vova1001/Expense-tracker-pet/internal/routes"
	d "github.com/vova1001/Expense-tracker-pet/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	d.InitDB()
	defer d.DB.Close()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://172.22.0.3:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	rout.RouterRegister(r)
	r.Run("0.0.0.0:8080")

}
