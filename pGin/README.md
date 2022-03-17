# pGin

pGin, namely **pseudo-Gin**, is a Gin-like webserver framework.

### Start an engine

Import the basic package **engine**, and create a new instance of engine by calling **engine.New()**

```go
e := engine.New()
```

### Run it

You can specify the port the webserver works on, for example, 8000

```go
e.run(8000)
```

### Group and bindings

Routes are grouped by its prefix in pGin, so you need to get a group first.

 **e.Default()** refers to the group with prefix "/"

You are also free to create new groups by calling **e.Group(prefix)**.

Here are are some examples:

```go
e.Default().
   Get("/header", func(c *engine.Context) {
      var content string
      for k, v := range c.Request.Header {
         content += fmt.Sprintf("%q => %q\n", k, v)
      }
      c.WriteContent(http.StatusOK, content)
   })
```

Notice that you can chain-call the method, which is very convenient. 

You are free to bind any valid http methods. The example above binds a handler only for **GE**T method.

Here is another example of binding handler.

```go
e.Group("/login").
   Post("/", func(c *engine.Context) {
      if c.Form("username") == "123" && c.Form("password") == "123" {
         c.WriteContent(http.StatusOK, "Login success!\n")
      } else {
         c.WriteContent(http.StatusOK, "Wrong username or password!\n")
      }
   })
```

The handler is only for **POST** method, and will check the form.

### Dynamic route

Support dynamic routing.

```go
e.Group("/echo").
   Get("/:name/:country", func(c *engine.Context) {
      c.WriteContent(http.StatusOK, fmt.Sprintf("Hello, %s from %s!", c.Params["name"], c.Params["country"]))
   })
```

Try it on "http://localhost:8000/echo/yourname/yourcountry".

```go
e.Group("/assets").
    Get("/*filePath", func(c *engine.Context) {
        c.WriteContent(http.StatusOK, fmt.Sprintf("Hey, you are at %s\n", c.Param("filePath")))
    })
```

Try it on "http://localhost:8000/assets/China/Shanghai/SJTU"

### Interceptors

You are free to customize your interceptors, here is a simple example.

```go
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
```

### Static resources

```go
e.Group("/assets").BindStatic("/", "./static")
```

Run this project and try to access "http://localhost:8000/assets/Ganyu.png"

**I love her, yes, absolutely.**

### Templates

```go
e.Group("/template").
   Get("/:name", func(c *engine.Context) {
      c.WriteHTMLTemplate(http.StatusOK, "example.tmpl", engine.JSON{"name": c.Param("name")})
   })
```

### Panic Recovery

Try to access "http://localhost:8000/panic"	