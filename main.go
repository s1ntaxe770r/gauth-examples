package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"

	"encoding/json"

	"golang.org/x/exp/slog"
)

var (
	SECRET_KEY   = ""
	GITHUB_TOKEN = ""
)

func main() {
	r := gin.Default()
	ghProvider := github.New(SECRET_KEY, GITHUB_TOKEN, "http://localhost:8080/callback", "user:email repo")
	goth.UseProviders(ghProvider)

	htmlFormat := `<html><body>%v</body></html>`
	r.GET("/", func(c *gin.Context) {
		html := fmt.Sprintf(htmlFormat, `<a href="/github">Login through github</a>`)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/github", func(ctx *gin.Context) {
		q := ctx.Request.URL.Query()
		q.Add("provider", "github")
		ctx.Request.URL.RawQuery = q.Encode()
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	})

	r.GET("/callback", func(c *gin.Context) {

		q := c.Request.URL.Query()
		q.Add("provider", "github")
		c.Request.URL.RawQuery = q.Encode()
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		slog.Debug(user.Email)
		res, err := json.Marshal(user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		jsonString := string(res)
		html := fmt.Sprintf(htmlFormat, jsonString)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	r.Run()

}
