#  How to use

### run the server on ws
```shell
$  go run main.go
```

### post user from http
```shell
$  curl -X POST -H "Content-Type: application/json" -d http://localhost:8081/register {"name":"nickname"}
```

### run the clients
```shell
$ cd client
$ go run main.go
```

### join
```shell
$ cd client
$ /join {“cmd”:”join”, “id”:uuid}
```


### todo

* http stats
* ws /guess
* ws /gameOver
* in memory usage
* concurrency 
* waiting room
* matchmaking algorithm
* unit test