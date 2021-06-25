package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type response struct {
	Name string `json:"file_name"`
	Size int64  `json:"file_size"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", index)
	//router.GET("/display", displayFile)
	router.POST("/upload", uploadFile)
	router.DELETE("/delete", deleteFile)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("cannot run http,err: [%v]", err)
	}
}

func index(c *gin.Context) {
	files, err := ioutil.ReadDir("fileList")
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	var data []response
	for _, file := range files {
		data = append(data, response{file.Name(), file.Size()})
	}
	c.JSON(http.StatusOK, data)
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	err = os.MkdirAll("./fileList", os.ModePerm)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	if err := c.SaveUploadedFile(file, "./fileList/"+file.Filename); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully uploaded",
	})
}

func deleteFile(c *gin.Context) {
	fn := c.Query("name")
	str := "fileList/" + fn
	err := os.Remove(str)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}
