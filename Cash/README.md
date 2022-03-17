# Cash

Cash, is with the same pronunciation as **"Cache"**([kæʃ]). It's a Groupcache-like cache server.

### Get Started

```go
c := cash.New(maxVolume, logLevel)
```

### Three Basic KV Method

```go
c.Put("Beijing", []byte("China"))
c.Put("Tokyo", []byte("Japan"))

c.Get("Beijing")
c.Del("Beijing")
```

