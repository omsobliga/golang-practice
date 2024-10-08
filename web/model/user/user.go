package user

import (
	"strings"
	"web/model"
)

type User struct {
	ID   int64
	Name string
	Age  int32
}

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (userDAO *UserDAO) Insert(name string, age int32) int64 {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := model.DB().Exec(sqlStr, name, age)
	if err != nil {
		panic(err)
	}
	theID, err := ret.LastInsertId()
	if err != nil {
		panic(err)
	}
	return theID
}

func (userDAO *UserDAO) GetByID(id int64) (User, error) {
	sqlStr := "select id, name, age from user where id=?"
	var u User
	err := model.DB().QueryRow(sqlStr, id).Scan(&u.ID, &u.Name, &u.Age)
	return u, err
}

func (userDAO *UserDAO) AsyncGetByID(id int64, chUser chan<- *User) {
	user, err := userDAO.GetByID(id)
	if err != nil {
		chUser <- &user
	} else {
		chUser <- nil
	}
}

func (userDAO *UserDAO) GetByIDs(ids []int64) []User {
	sqlStr := "select id, name, age from user where id in (?" + strings.Repeat(",?", len(ids)-1) + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := model.DB().Query(sqlStr, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]User, len(ids))
	user_map := make(map[int64]User, len(ids))

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Age)
		if err != nil {
			panic(err)
		}
		user_map[u.ID] = u
	}

	for i, id := range ids {
		users[i] = user_map[id]
	}
	return users
}
