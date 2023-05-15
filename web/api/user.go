package api

import (
	"Lesson_4/web/api/middleware"
	"Lesson_4/web/dao"
	"Lesson_4/web/model"
	"Lesson_4/web/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func register(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": "verification failed",
		})
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否重复
	flag := dao.SelectUser(username)
	fmt.Println(flag)
	if flag {
		// 以 JSON 格式返回信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "user already exists",
		})
		return
	}

	dao.AddUser(username, password)
	// 以 JSON 格式返回信息
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "add user successful",
	})
}

// 仅有登录部分有改动
func login(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否存在
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}

	// 正确则登录成功
	// 创建一个我们自己的声明
	claim := model.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "Yxh",                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	utils.RespSuccess(c, tokenString)
}

// 新增以下代码
func getUsernameFromToken(c *gin.Context) {
	username, _ := c.Get("username")
	utils.RespSuccess(c, username.(string))
}

func changePassword(c *gin.Context) {
	form := model.User{}
	if err := c.ShouldBind(&form); err != nil {
		fmt.Println(err)
		utils.RespSuccess(c, "verification failed")
		return
	}

	// 获取用户信息
	username := c.GetString("username")
	oldPassword := form.Password
	newPassword := form.NewPassword

	// 验证旧密码是否正确
	selectPassword := dao.SelectPasswordFromUsername(username)
	if selectPassword != oldPassword {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}

	// 更新密码
	dao.AddUser(username, newPassword)
	utils.RespSuccess(c, "change password successful")
}

// 找回密码
func retrievePassword(c *gin.Context) {
	form := model.User{}
	if err := c.ShouldBind(&form); err != nil {
		fmt.Println(err)
		utils.RespSuccess(c, "verification failed")
		return
	}

	// 获取用户信息
	username := form.Username
	//email := form.Password

	// 检查该用户名是否存在
	flag := dao.SelectUser(username)
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}

	// 发送密码重置链接
	// ... （省略具体实现）
	utils.RespSuccess(c, "retrieve password successful")
}

/*func register(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": "verification failed",
		})
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否重复
	flag := dao.SelectUser(username)
	fmt.Println(flag)
	if flag {
		// 以 JSON 格式返回信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "user already exists",
		})
		return
	}

	dao.AddUser(username, password)
	// 以 JSON 格式返回信息
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "add user successful",
	})
}

func login(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  500,
			"message": "verification failed",
		})
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否存在
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "user doesn't exists",
		})
		return
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "wrong password",
		})
		return
	}

	// 正确则登录成功
	c.SetCookie("gin_demo_cookie", "test", 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "login successful",
	})
}*/
