package service

import (
	"chat/pkg/e"
	"chat/serializer"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	code := e.SUCCESS
	//w := c.Writer
	req := c.Request
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		code = e.ErrorGetFile
		c.JSON(500, serializer.Response{
			Status: e.ErrorGetFile,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		})
		return
	}
	suffix := ".png"
	filename := head.Filename
	tem := strings.Split(filename, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31, suffix)
	dstFile, err := os.Create("./pkg/upload/data/" + fileName)
	if err != nil {
		code = e.ErrorCreateFile
		c.JSON(500, serializer.Response{
			Status: e.ErrorCreateFile,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		})
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		code = e.ErrorPlaceFile
		c.JSON(500, serializer.Response{
			Status: e.ErrorPlaceFile,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		})
		return
	}
	url := "./pkg/upload/data/" + fileName
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(code),
		Data:   map[string]string{"url": url},
	})
	return
}
