package main

import (
	"engine"
	"fmt"
	"net/http"
)

func main() {
	e := engine.NewEngine()
	e.Bind(engine.GET, "/", func(c *engine.Context) {
		c.WriteHTML(http.StatusOK, "<p>Welcome to Naive WebServer</p>\n"+
			"<p>Author: LZH</p>")
	})

	e.Bind(engine.GET, "/header", func(c *engine.Context) {
		var content string
		for k, v := range c.Request.Header {
			content += fmt.Sprintf("%q => %q\n", k, v)
		}
		c.WriteContent(http.StatusOK, content)
	})

	e.Bind(engine.GET, "/json", func(c *engine.Context) {
		c.WriteJson(http.StatusOK, engine.JSON{
			"code": 1,
			"msg":  "hello!",
		})
	})

	e.Bind(engine.GET, "/echo/:name", func(c *engine.Context) {
		c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s!\n", c.Params["name"]))
	})

	e.Bind(engine.GET, "/assets/*addr", func(c *engine.Context) {
		c.WriteContent(http.StatusOK, fmt.Sprintf("Hey, you are at %s\n", c.Params["addr"]))
	})

	e.Run(8000)
}
