package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", index )
	router.POST("/upload", uploadFile)
	router.DELETE("/delete", deleteFile)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("cannot run http,err: [%v]", err)
	}
}

func index(c *gin.Context){
	files, err := ioutil.ReadDir("fileList")
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}
	for _, file := range files{
		c.JSON(http.StatusOK, gin.H{
			"name": file.Name(),
			"size": file.Size(),
		})
	}
}

func uploadFile (c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: [%s]", err.Error()))
		return
	}
	err = os.MkdirAll("./fileList", os.ModePerm)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("creat folder err: [%s]", err.Error()))
		return
	}
	if err := c.SaveUploadedFile(file, "./fileList/"+file.Filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: [%s]", err.Error()))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}

func deleteFile(c *gin.Context){
	fn := c.Query("name")
	str := "fileList/" + fn
	err := os.Remove(str)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("cannot remove, err: [%s]", err.Error()))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("File %s delete successfully", fn))
}