package pkg

import (
	"chat/conf"
	"chat/pkg/e"
	"chat/serializer"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenToken(id string) (aToken, rToken string, err error) {
	claims := MyClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(aTokenTime).Unix(),
			Issuer:    "waterSystem",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	aToken, err = token.SignedString(MySecret)
	// rToken 不需要存储任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(rTokenTime).Unix(), // 过期时间
		Issuer:    "my-project",                      // 签发人
	}).SignedString(MySecret)
	return aToken, rToken, nil
}

type MyClaims struct {
	ID string `json:"student_id"`
	jwt.StandardClaims
}

var MySecret = []byte("5cq2up0t")

const aTokenTime = 8 * time.Hour
const rTokenTime = 14 * 24 * time.Hour

// ParasToken 解析 access_token
func ParasToken(aToken string) (claims *MyClaims, err error) {
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(aToken, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid { // token 是否有效
		err = errors.New("invalidToken")
		return nil, err
	}
	return claims, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return MySecret, nil
}

// NewToken 刷新token
func NewToken(aToken, rToken string) (newToken, newRToken string, err error) {
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return "", "", err
	}
	var claims MyClaims
	token, err := jwt.ParseWithClaims(aToken, &claims, keyFunc)
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid access_token")
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return GenToken(claims.ID)
	}
	return "", "", nil
}

// AuthMiddleware 中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		aToken := c.GetHeader("Authorization")
		if aToken == "" {
			code := e.ErrorHeaderData
			c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			})
			return
		}
		claims, err := ParasToken(aToken)
		if err != nil {
			code := e.ErrorValidToken
			c.AbortWithStatusJSON(http.StatusUnauthorized, serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			})
			return
		}
		c.Set("ID", claims.ID)
		c.Next()
	}
}

func ParseSet(c *gin.Context) (uint, error) {
	ID, exists := c.Get("ID")
	if !exists {
		return 0, errors.New("not exists")
	}
	id, err := strconv.Atoi(ID.(string))
	if err != nil {
		return 0, err
	}
	return uint(id), nil

}

// SaveImg 保存图片到指定文件目录下，自定义文件后缀
func SaveImg(path, suffix string, srcFile multipart.File, head *multipart.FileHeader) (string, error) {
	fmt.Println(path)
	filename := head.Filename
	tem := strings.Split(filename, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31, suffix)
	dstFile, err := os.Create(path + fileName)
	if err != nil {
		fmt.Println("at os.Create err happen")
		return "", err
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("io.Copy err happen")
		return "", err
	}
	url := path + fileName
	return url, nil
}
func SendCheckCode(email string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000)) // 生成随机的六位验证码
	fmt.Println(conf.EmailUser, conf.EmailPassword, conf.EmailPort, conf.EmailHost, "over")
	msg := fmt.Sprintf("您的验证码为：%s", code)

	m := gomail.NewMessage()
	m.SetHeader("From", conf.EmailUser) // 发件人邮箱，替换为你的QQ邮箱
	m.SetHeader("To", email)            // 收件人邮箱
	m.SetHeader("Subject", "验证码")
	m.SetBody("text/html", msg)
	port, _ := strconv.Atoi(conf.EmailPort)
	d := gomail.NewDialer(conf.EmailHost, port, conf.EmailUser, conf.EmailPassword) // 设置QQ邮箱的SMTP服务器和端口号
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}                             // 跳过证书验证，可在测试阶段使用，生产环境请勿使用此设置

	err := d.DialAndSend(m)
	if err != nil {
		log.Println("发送邮件时出现错误:", err)
		return "", err
	}

	log.Println("验证码已成功发送到邮箱:", email)
	return code, nil
}
