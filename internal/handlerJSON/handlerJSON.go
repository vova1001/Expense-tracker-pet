package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	h "github.com/vova1001/Expense-tracker-pet/internal/handler"
	m "github.com/vova1001/Expense-tracker-pet/internal/model"
)

func GET_JSON(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)

	procStr := ctx.Query("proc")
	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")
	periodStr := ctx.Query("period")
	searchStr := ctx.Query("search")

	ok := ChekParam(procStr, limitStr, pageStr, periodStr, searchStr)
	if ok {
		TasksParam, err := h.GetTaskParam(userIdInt, pageStr, limitStr, procStr, periodStr, searchStr)
		if err != nil {
			log.Println("Error in GetTaskParam:", err)
			ctx.JSON(400, gin.H{"err": "err TaskParam"})
			return
		}
		ctx.JSON(200, TasksParam)
		return
	}

	task, msg := h.GetTask(userIdInt)
	if len(task) == 0 {
		ctx.JSON(http.StatusOK, msg)
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func POST_JSON(ctx *gin.Context) {
	var task m.Task
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err})
		return
	}
	ResultTask, err := h.PostTask(task, userIdInt)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err post"})
		return
	}
	ctx.JSON(http.StatusOK, ResultTask)

}

func DELETE_JSON(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"err": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	IdStr := ctx.Param("id")
	id, err := strconv.Atoi(IdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid ID": err})
		return
	}
	ResultTask_err := h.DeleteTask(id, userIdInt)
	if ResultTask_err != nil {
		ctx.JSON(http.StatusNotFound, ResultTask_err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func PUT_JSON(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"err": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	var UpdatedTask m.Task
	err := ctx.ShouldBindJSON(&UpdatedTask)
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
	ResultTask, err := h.PutTask(UpdatedTask, Id, userIdInt)
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
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"err": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	h.ClearAll(userIdInt)
	ctx.JSON(200, gin.H{"Result": "Cleared successful"})
}

func CheckTask_JSON(ctx *gin.Context) {
	var ChekTask m.ChekBox
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"err": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	err := ctx.ShouldBindJSON(&ChekTask)
	if err != nil {
		ctx.JSON(400, gin.H{"Invalid JSON": err})
		return
	}
	StrId := ctx.Param("id")
	id, err := strconv.Atoi(StrId)
	if err != nil {
		ctx.JSON(400, gin.H{"Invalid id": err})
		return
	}
	err = h.ChekDone(ChekTask, id, userIdInt)
	if err != nil {
		ctx.JSON(400, gin.H{"err ResultCheckFunc": err})
		return
	}
	ctx.JSON(200, "Check succssesful")
}

func DueDate_JSON(ctx *gin.Context) {
	var DueDate m.DueDate
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"err": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	err := ctx.ShouldBindJSON(&DueDate)
	if err != nil {
		ctx.JSON(400, gin.H{"Invalid JSON": err})
		return
	}
	StrId := ctx.Param("id")
	id, err := strconv.Atoi(StrId)
	if err != nil {
		ctx.JSON(400, gin.H{"err ID": err})
		return
	}
	err = h.DueDateFunc(DueDate, id, userIdInt)
	if err != nil {
		ctx.JSON(400, gin.H{"err DueDateFunc": err})
		return
	}
	ctx.JSON(200, "DueDate created")
}

func Register_JSON(ctx *gin.Context) {
	var User m.User
	err := ctx.ShouldBindJSON(&User)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err JSON"})
		return
	}
	err = h.RegisterUser(User)
	if err != nil {
		ctx.JSON(404, gin.H{"err": "err handler func"})
	}
	ctx.JSON(200, "User is registered")

}

func Login_JSON(ctx *gin.Context) {
	var user m.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "Invalid JSON"})
	}
	token, err := h.Login(user)
	if err != nil {
		ctx.JSON(401, gin.H{"err": err.Error()})
	}
	ctx.JSON(200, token)

}

func validateTokenSignature(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("invalid singature")
	}
	secret := []byte(os.Getenv("JWT_SECRET"))
	return secret, nil
}

func JWTMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Header := ctx.GetHeader("Authorization")
		if Header == "" {
			ctx.JSON(401, gin.H{"err": "Authorization missing"})
			ctx.Abort()
			return
		}
		partsHeader := strings.Split(Header, " ")
		if len(partsHeader) != 2 || partsHeader[0] != "Bearer" {
			ctx.JSON(401, gin.H{"err": "Invalid header"})
			ctx.Abort()
			return
		}
		tokenStr := partsHeader[1]

		token, err := jwt.Parse(tokenStr, validateTokenSignature)
		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"err": "Invalid or expired token"})
			ctx.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(401, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		iuserID, ok := claims["user_id"].(float64)
		if !ok {
			ctx.JSON(401, gin.H{"error": "user_id not found in token"})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", int(iuserID))
		ctx.Next()
	}
}

func TaskStatus_JSON(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(401, gin.H{"err": "User not authenticated"})
		return
	}
	userIdInt := userId.(int)
	Result, err := h.TaskStatus(userIdInt)
	if err != nil {
		ctx.JSON(400, gin.H{"err": err})
		return
	}
	ctx.JSON(200, Result)
}

func ChekParam(proc, limit, page, periodStr, searchStr string) bool {
	if proc != "" || page != "" || limit != "" || periodStr != "" || searchStr != "" {
		return true
	}
	return false
}
