package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	h "github.com/vova1001/Expense-tracker-pet/internal/handler"
	m "github.com/vova1001/Expense-tracker-pet/internal/model"
)

func GET_JSON(ctx *gin.Context) {
	task, msg := h.GetTask()
	if len(task) == 0 {
		ctx.JSON(http.StatusOK, msg)
		return
	}
	ctx.JSON(http.StatusOK, task)

}

func POST_JSON(ctx *gin.Context) {
	var task m.Task
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err})
		return
	}
	ResultTask := h.PostTask(task)
	ctx.JSON(http.StatusOK, ResultTask)

}

func DELETE_JSON(ctx *gin.Context) {
	var task m.Task
	err := ctx.ShouldBindJSON(&task)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err})
		return
	}
	IdStr := ctx.Param("id")
	id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid ID": err})
		return
	}
	ResultTask_err := h.DeleteTask(id)
	if ResultTask_err != nil {
		ctx.JSON(http.StatusNotFound, ResultTask_err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func PUT_JSON(ctx *gin.Context) {
	var UpdatedTask m.Task
	err := ctx.ShouldBindJSON(&UpdatedTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err})
		return
	}
	IdStr := ctx.Param("id")
	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error strconv": err})
		return
	}
	UpdatedTask.ID = Id
	ResultTask, err := h.PutTask(UpdatedTask)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Task not found": err})
		return
	}
	ctx.JSON(http.StatusOK, ResultTask)
}

func PATCH_JSON(ctx *gin.Context) {
	var UpdatedPatch map[string]interface{}
	err := ctx.ShouldBindJSON(&UpdatedPatch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err})
		return
	}
	IdStr := ctx.Param("id")
	Id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid ID": err})
		return
	}
	ResultTaskPatch, err := h.PatchTask(UpdatedPatch, Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err result patch": err})
		return
	}
	ctx.JSON(http.StatusOK, ResultTaskPatch)
}

func ClearAll_JSON(ctx *gin.Context) {
	h.ClearAll()
	ctx.JSON(200, gin.H{"Result": "Cleared successful"})
}
