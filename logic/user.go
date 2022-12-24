package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"fmt"
)

// 注册逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 使用sonyflake生成UID
	userID, _ := snowflake.GenID2()
	fmt.Println("userID=", userID)
	// 构造user实例
	u := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
		Email:    p.Email,
	}
	// 保存到数据库
	return mysql.InsertUser(&u)
}

// 登录逻辑
func Login(p *models.ParamLogin) (user *models.User, err error) {

	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	// 生成JWT token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return
}
