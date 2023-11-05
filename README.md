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
$ /join {"cmd":"join", "id":"XVlBzgba"}
```

### guess
```shell
$ cd client
$ /guess {"cmd": "guess", "id":"XVlBzgba", "room":"room1699210612327731000", "data":4}
```


### status

- [ ] http stats
- [x] ws /guess
- [x] ws /gameOver
- [x] in memory usage
- [x] concurrency 
- [x] waiting room
- [ ] matchmaking algorithm
- [ ] unit test