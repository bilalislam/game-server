#  How to use

### run the server on ws
```shell
$  go run main.go
```

### post user from http
```shell
$  curl -X POST -H "Content-Type: application/json" -d http://localhost:8080/register {"name":"nickname"}
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
- [x] matchmaking algorithm
- [ ] unit test
- [ ] void to command after timeout
- [ ] handle timeout , guess as zero and set rank is -1 to server
- [ ] start new game when no ones give an answer 
- [ ] start game until room count is 3 players 
