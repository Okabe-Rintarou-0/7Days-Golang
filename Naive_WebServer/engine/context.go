package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Method  string
	Path    string
	Params  map[string]string
}

type JSON map[string]interface{}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Method:  r.Method,
		Path:    r.URL.Path,
	}
}

func (c *Context) WriteStatus(statusCode int) {
	c.Writer.WriteHeader(statusCode)
}

func (c *Context) Write(content []byte) {
	_, err := c.Writer.Write(content)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) WriteContent(statusCode int, content string) {
	c.Writer.WriteHeader(statusCode)
	c.Write([]byte(content))
}

func (c *Context) WriteJson(statusCode int, obj interface{}) {
	c.WriteStatus(statusCode)
	c.SetHeader(HeaderKeyContentType, HeaderValueJson)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) WriteHTML(statusCode int, html string) {
	c.WriteStatus(statusCode)
	c.SetHeader(HeaderKeyContentType, HeaderValueHTML)
	c.Write([]byte(html))
}

func (c *Context) GetQuery(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) GetForm(key string) string {
	return c.Request.Form.Get(key)
}
