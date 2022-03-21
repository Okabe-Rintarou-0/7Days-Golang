# Cash

Cash, is with the same pronunciation as **"Cache"**([kæʃ]). It's a Groupcache-like cache server.

### Get Started

Create a new http pool that running at your given address:

```go
pool := cash.NewHTTPPool("localhost:8000")
```

Create some groups with different namespaces. Groups are independent from each other.

```go
pool.NewGroup(logLevel, maxVolume, "country", cash.GetterFunc(naiveGetter()))
pool.NewGroup(logLevel, maxVolume, "food", cash.GetterFunc(naiveGetter()))
```

Here, you are free to customize your getter function. Once the cash does not hit, it will try get the k-v pair through getter function. Getter function can fetch k-v pair from a database or another k-v store, and it's up to you. 

naiveGetter is a simple example just for display.

```go
func naiveGetter() func(key string) (cash.ByteView, error) {
	kvs := map[string]cash.ByteView{}
	return func(key string) (cash.ByteView, error) {
		if value, ok := kvs[key]; ok {
			return value, nil
		} else {
			return value, fmt.Errorf("no such key")
		}
	}
}
```

### Api

The cash is based on HTTP protocol. If you want to do some basic operations on a specific group, here are some examples.

#### Get the value of key "cake" from group "food"

by using **GET** method:

http://localhost:8000/__cash__/food?key=cake

#### Put "cake" -> "tasty" to group "food"

by using **PUT**  method:

http://localhost:8000/__cash__/food?key=cake&value=tasty

#### Del key "cake" from group "food"

by using **Delete** method:

http://localhost:8000/__cash__/food?key=cake

#### Get info of a given group

by using **GET** method:

http://localhost:8000/__cash__/food/info

Would return message like: 

```json
{"maxVolume":"2000","percent":"0.00","used":"0"}
```

#### Get info of groups

by using **GET** method:

http://localhost:8000/__cash__/__groups__

Would return message like:

```json
["country", "food"]
```

### Console

A very naive one, made by Typescript + React + Ant-Design + Echarts.

Simply monitor the usage of cash. 

![console](README.assets/console-16478359331921.gif)
