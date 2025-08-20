package routes

import (
	"github.com/gin-gonic/gin"

	hj "github.com/vova1001/Expense-tracker-pet/internal/handlerJSON"
)

func RouterRegister(r *gin.Engine) {
	task := r.Group("/task")
	task.Use(hj.JWTMiddleWare())
	{
		task.GET("/", hj.GET_JSON)
		task.POST("/", hj.POST_JSON)
		task.DELETE("/:id", hj.DELETE_JSON)
		task.PUT("/:id", hj.PUT_JSON)
		task.PATCH("/:id", hj.PATCH_JSON)
		task.PUT("/done/:id", hj.CheckTask_JSON)
		task.PUT("/duedate/:id", hj.DueDate_JSON)
		task.DELETE("/clearAll", hj.ClearAll_JSON)
		task.GET("/stats", hj.TaskStatus_JSON)
	}
	r.POST("/register", hj.Register_JSON)
	r.POST("/login", hj.Login_JSON)

}
