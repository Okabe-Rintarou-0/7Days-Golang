package main

import (
	"engine"
	"fmt"
	"net/http"
)

func main() {
	e := engine.NewEngine()

	e.NewGroup("").
		Bind(http.MethodGet, "/", func(c *engine.Context) {
			c.WriteHTML(http.StatusOK, "<p>Welcome to Naive WebServer</p>\n<p>Author: LZH</p>")
		}).
		Bind(http.MethodGet, "/header", func(c *engine.Context) {
			var content string
			for k, v := range c.Request.Header {
				content += fmt.Sprintf("%q => %q\n", k, v)
			}
			c.WriteContent(http.StatusOK, content)
		}).
		Bind(http.MethodGet, "/json", func(c *engine.Context) {
			c.WriteJson(http.StatusOK, engine.JSON{
				"code": 1,
				"msg":  "hello!",
			})
		})

	e.NewGroup("/echo").
		Bind(http.MethodGet, "/:name", func(c *engine.Context) {
			c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s!\n", c.Params["name"]))
		}).
		Bind(http.MethodGet, "/:name/:country", func(c *engine.Context) {
			c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s from %s!", c.Params["name"], c.Params["country"]))
		})

	e.NewGroup("/assets").
		Bind(http.MethodGet, "/*addr", func(c *engine.Context) {
			c.WriteContent(http.StatusOK, fmt.Sprintf("Hey, you are at %s\n", c.Params["addr"]))
		})

	e.NewGroup("/login").
		Bind(http.MethodPost, "/", func(c *engine.Context) {
			if c.GetForm("username") == "123" && c.GetForm("password") == "123" {
				c.WriteContent(http.StatusOK, "Login success!\n")
			} else {
				c.WriteContent(http.StatusOK, "Wrong username or password!\n")
			}
		})

	e.NewGroup("/hello").
		Bind(http.MethodGet, "/", func(c *engine.Context) {
			c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s!\n", c.GetQuery("name")))
		})

	e.Run(8000)
}
