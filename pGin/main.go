package main

import (
	"engine"
	"fmt"
	"net/http"
)

func main() {
	e := engine.New()
	e.LoadHTMLGlob("./templates/*")
	e.Default().
		Get("/", func(c *engine.Context) {
			c.WriteHTML(http.StatusOK, "<p>Welcome to Naive WebServer</p>\n<p>Author: LZH</p>")
		}).
		Get("/header", func(c *engine.Context) {
			var content string
			for k, v := range c.Request.Header {
				content += fmt.Sprintf("%q => %q\n", k, v)
			}
			c.WriteContent(http.StatusOK, content)
		}).
		Get("/json", func(c *engine.Context) {
			c.WriteJson(http.StatusOK, engine.JSON{
				"code": 1,
				"msg":  "hello!",
			})
		})

	e.Group("/echo").
		Get("/:name/:country", func(c *engine.Context) {
			c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s from %s!", c.Params["name"], c.Params["country"]))
		}).
		AddInterceptor(func(c *engine.Context) bool {
			if c.Param("country") == "Mars" {
				return false
			}
			return true
		})

	//e.Group("/assets").
	//	Get("/*filePath", func(c *engine.Context) {
	//		c.WriteContent(http.StatusOK, fmt.Sprintf("Hey, you are at %s\n", c.Param("filePath")))
	//	})

	e.Group("/assets").BindStatic("/", "./static")

	e.Group("/login").
		Post("/", func(c *engine.Context) {
			if c.Form("username") == "123" && c.Form("password") == "123" {
				c.WriteContent(http.StatusOK, "Login success!\n")
			} else {
				c.WriteContent(http.StatusOK, "Wrong username or password!\n")
			}
		})

	e.Group("/hello").
		Get("/", func(c *engine.Context) {
			c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s!\n", c.Query("name")))
		})

	e.Group("/template").
		Get("/:name", func(c *engine.Context) {
			c.WriteHTMLTemplate(http.StatusOK, "example.tmpl", engine.JSON{"name": c.Param("name")})
		})

	e.Group("/panic").
		Get("/", func(c *engine.Context) {
			var a []int
			fmt.Println(a[5])
		})

	e.Run(8000)
}
