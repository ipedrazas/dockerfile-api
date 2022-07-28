package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type DckrFileRes struct {
	Name   string
	Src    string
	Data   string
	Commit string
}

type INFO struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Golang  string `json:"go_version"`
}

func main() {

	router := gin.Default()
	router.GET("/", info)
	router.GET("/repo", getDockerfiles)
	router.Use(cors.Default())
	router.Run("0.0.0.0:8060")
}

func getDockerfiles(c *gin.Context) {

	gitUrl := strings.ToLower(c.Query("url"))
	dfUrl := parseRepoUrl(gitUrl)

	_, err := url.ParseRequestURI(dfUrl)
	if err != nil {
		c.IndentedJSON(http.StatusPreconditionFailed, err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, dfUrl, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	dckf := &DckrFileRes{
		Name: "Dockerfile",
		Src:  gitUrl,
		Data: string(resBody),
	}
	c.IndentedJSON(http.StatusOK, dckf)

}

func parseRepoUrl(repo string) string {

	if strings.Contains(repo, "https://github.com") {
		return repo + "/raw/main/Dockerfile"
	}
	if strings.Contains(repo, "https://gitea.") {
		return repo + "/raw/branch/main/Dockerfile"
	}
	return ""
}

func info(c *gin.Context) {
	info := &INFO{
		Name:    "Dockerfiles API",
		Version: "v0.1.3",
		Golang:  "go1.18",
	}
	c.IndentedJSON(http.StatusOK, info)
}
