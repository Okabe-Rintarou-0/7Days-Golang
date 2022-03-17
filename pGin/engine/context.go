package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Context struct {
	engine       *Engine
	Writer       http.ResponseWriter
	Request      *http.Request
	Method       string
	Pattern      string
	Params       map[string]string
	Body         []byte
	Middlewares  []FuncHandler
	Interceptors []Interceptor
}

type JSON map[string]interface{}

func (json *JSON) Get(key string) interface{} {
	if value, ok := (*json)[key]; ok {
		return value
	}
	return nil
}

func (json *JSON) Put(key string, value interface{}) {
	(*json)[key] = value
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	body, _ := ioutil.ReadAll(r.Body)
	return &Context{
		Writer:  w,
		Request: r,
		Method:  r.Method,
		Pattern: r.URL.Path,
		Body:    body,
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

func (c *Context) NotFound() {
	http.Error(c.Writer, "Not Found", http.StatusNotFound)
}

func (c *Context) Forbidden() {
	http.Error(c.Writer, "Forbidden", http.StatusForbidden)
}

func (c *Context) InternalServerError(errMsg string) {
	http.Error(c.Writer, errMsg, http.StatusInternalServerError)
}

func (c *Context) WriteJson(statusCode int, obj interface{}) {
	c.WriteStatus(statusCode)
	c.SetHeader(HeaderKeyContentType, HeaderValueJson)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		c.InternalServerError(err.Error())
	}
}

func (c *Context) WriteHTML(statusCode int, html string) {
	c.WriteStatus(statusCode)
	c.SetHeader(HeaderKeyContentType, HeaderValueHTML)
	c.Write([]byte(html))
}

func (c *Context) WriteHTMLTemplate(statusCode int, template string, data interface{}) {
	c.WriteStatus(statusCode)
	c.SetHeader(HeaderKeyContentType, HeaderValueHTML)
	err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, template, data)
	if err != nil {
		c.InternalServerError(err.Error())
	}
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Form(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) ParseBody(v interface{}) interface{} {
	err := json.Unmarshal(c.Body, v)
	if err != nil {
		c.InternalServerError(err.Error())
		return nil
	}
	return v
}

func (c *Context) Body2Text() string {
	return string(c.Body)
}

func (c *Context) Body2Json() JSON {
	body := JSON{}
	err := json.Unmarshal(c.Body, &body)
	if err != nil {
		c.InternalServerError(err.Error())
		return nil
	}
	return body
}
