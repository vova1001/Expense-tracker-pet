package routes

import (
	"github.com/gin-gonic/gin"

	hj "github.com/vova1001/Expense-tracker-pet/internal/handlerJSON"
)

func RouterTaskRegister(r *gin.Engine) {
	task := r.Group("/task")
	{
		task.GET("/", hj.GET_JSON)
		task.POST("/", hj.POST_JSON)
		task.DELETE("/:id", hj.DELETE_JSON)
		task.PUT("/:id", hj.PUT_JSON)
		task.PATCH("/:id", hj.PATCH_JSON)
	}
	r.DELETE("/clear", hj.ClearAll_JSON)
}
