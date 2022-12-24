package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) error {
	sqlStr := "select count(user_id) from users where username= ?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// QueryUserByUsername 根据用户名查询用户
func QueryUserByUsername(username string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select * from users where username= ?"
	err = db.Get(user, sqlStr, username)
	return
}

const secret = "test.com"

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	sqlStr := "insert into users(user_id, username, password, email) values (?, ?, ?, ?)"

	user.Password = encryptPassword(user.Password)
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Email)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPwd string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPwd)))
}

// Login登录逻辑校验
func Login(user *models.User) (err error) {
	password := encryptPassword(user.Password)
	sqlStr := "select username, password, user_id from users where username = ?"
	err = db.Get(user, sqlStr, user.Username)

	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		fmt.Printf("query user failed, err:%v \n", err)
		return err
	}

	if password != user.Password {
		return ErrorInvalidPassword
	}
	return nil
}

// GetUserByID 根据id查询用户信息
func GetUserByID(userID uint64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select * from users where user_id= ?"
	err = db.Get(user, sqlStr, userID)
	return
}
