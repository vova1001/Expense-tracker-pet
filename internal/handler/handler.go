package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	m "github.com/vova1001/Expense-tracker-pet/internal/model"
	d "github.com/vova1001/Expense-tracker-pet/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func GetTask(userId int) ([]m.Task, error) {
	rows, err := d.DB.Query("SELECT id, name, move, proc, time FROM tasks WHERE user_id= $1 ORDER BY id ASC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []m.Task
	for rows.Next() {
		var t m.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Move, &t.Proc, &t.Time)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil

}

func PostTask(Newtask m.Task, userId int) ([]m.Task, error) {
	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		move TEXT NOT NULL,
		proc BOOLEAN DEFAULT FALSE,
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		duedate TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := d.DB.Exec(createTable)
	if err != nil {
		return nil, err
	}
	fmt.Println("Table created")
	query := `INSERT INTO tasks (user_id, name, move, proc, time, duedate) VALUES ($1, $2, $3, $4, $5, $6)`
	Newtask.Time = time.Now()
	Newtask.UserID = userId
	_, err = d.DB.Exec(query, Newtask.UserID, Newtask.Name, Newtask.Move, Newtask.Proc, Newtask.Time, Newtask.TaskDueDate)
	if err != nil {
		return nil, err
	}

	rows, err := d.DB.Query("SELECT id, name, move, proc, time, duedate FROM tasks WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []m.Task
	for rows.Next() {
		var t m.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Move, &t.Proc, &t.Time, &t.TaskDueDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func DeleteTask(id int, user_id int) error {
	queryDelete := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	res, err := d.DB.Exec(queryDelete, id, user_id)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return fmt.Errorf("no task id")
	}
	return nil
}

func PutTask(updated m.Task, id int, user_id int) (m.Task, error) {
	UpdatedTask := "UPDATE tasks SET name = $1, move = $2, proc = $3, duedate = $4 WHERE id =$5 AND user_id = $6 RETURNING id, name, move, proc, time, duedate"
	var NewUpdated m.Task
	err := d.DB.QueryRow(UpdatedTask, updated.Name, updated.Move, updated.Proc, updated.TaskDueDate, id, user_id).Scan(&NewUpdated.ID, &NewUpdated.Name, &NewUpdated.Move, &NewUpdated.Proc, &NewUpdated.Time, &NewUpdated.TaskDueDate)
	if err != nil {
		return m.Task{}, err
	}
	return NewUpdated, nil

}

func PatchTask(updatedPatch map[string]interface{}, id int) (m.Task, error) {
	var task m.Task
	QueryPatch := "SELECT id, name, move, proc, time FROM tasks WHERE id =$1"
	err := d.DB.QueryRow(QueryPatch, id).Scan(&task.ID, &task.Name, &task.Move, task.Proc, task.Time)
	if err != nil {
		return m.Task{}, err
	}
	if Name, ok := updatedPatch["name"].(string); ok {
		task.Name = Name
	}
	if Move, ok := updatedPatch["move"].(string); ok {
		task.Move = Move
	}
	if Proc, ok := updatedPatch["proc"].(bool); ok {
		task.Proc = Proc
	}

	_, err = d.DB.Exec("UPDATE tasks SET name=$1, move=$2, proc=$3, time=$4 WHERE id=$5", task.Name, task.Move, task.Proc, task.Time, task.ID)
	if err != nil {
		return m.Task{}, err
	}
	return task, nil
}

func ClearAll(user_id int) error {
	Clear := "DELETE FROM tasks WHERE user_id = $1"
	_, err := d.DB.Exec(Clear, user_id)
	if err != nil {
		return err
	}
	return nil
}

func ChekDone(ChekTask m.ChekBox, id int, user_id int) error {
	if ChekTask.Check {
		_, err := d.DB.Exec("UPDATE tasks SET proc=$1 WHERE id=$2 AND user_id = $3", true, id, user_id)
		if err != nil {
			return err
		}
	}
	return nil

}

func DueDateFunc(duedate m.DueDate, id int, user_id int) error {
	_, err := d.DB.Exec("UPDATE tasks SET duedate=$1 WHERE id=$2 AND user_id =$3", duedate.Due, id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func RegisterUser(user m.User) error {
	CreateTableUser := `
	CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	password TEXT NOT NULL,
	email TEXT NOT NULL
	);`
	_, err := d.DB.Exec(CreateTableUser)
	if err != nil {
		return err
	}
	fmt.Println("TableUsers created")
	ok := EmailCheck(user.Email)
	if !ok {
		return err
	}
	var exists bool
	err = d.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return err
	}
	HshPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec("INSERT INTO users(password,email) VALUES($1, $2)", HshPass, user.Email)
	if err != nil {
		return err
	}
	log.Println("User created")
	return nil
}

func Login(user m.User) (error, m.Token) {
	var userID int
	var userPass string
	var userEmail string
	var resultToken m.Token
	err := d.DB.QueryRow("SELECT id, password, email FROM users WHERE email = $1", user.Email).Scan(&userID, &userPass, &userEmail)
	if err != nil {
		return err, m.Token{}
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPass), []byte(user.Password))
	if err != nil {
		return err, m.Token{}
	}
	sk := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"user_email": userEmail,
		"time_life":  time.Now().Add(24 * time.Hour).Unix(),
	})
	SignedToken, err := token.SignedString([]byte(sk))
	if err != nil {
		return err, m.Token{}
	}
	resultToken.Token = SignedToken
	return nil, resultToken
}

func EmailCheck(email string) bool {
	k := 0
	for _, v := range email {
		switch {
		case v == '@':
			k++
			if k > 1 {
				return false
			}
		case !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') || v == '.' || v == '-' || v == '_'):
			return false
		}
	}
	return k == 1
}

func TaskStatus(user_id int) (m.TaskStatus, error) {
	rows, err := d.DB.Query("SELECT proc, duedate FROM tasks WHERE user_id =$1", user_id)
	if err != nil {
		return m.TaskStatus{}, err
	}
	defer rows.Close()
	var StatusTasks m.TaskStatus
	for rows.Next() {
		StatusTasks.AllTasks++
		var proc bool
		var duedate time.Time
		t := time.Now()
		err := rows.Scan(&proc, &duedate)
		if err != nil {
			return m.TaskStatus{}, err
		}
		if proc {
			StatusTasks.CompletedTasks++
		} else {
			StatusTasks.ActiveTasks++
			if t.After(duedate) {
				StatusTasks.OverdueTasks++
			}
		}
	}
	return StatusTasks, err
}
