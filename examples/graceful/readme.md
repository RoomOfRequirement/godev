# Server graceful restart with Go

reference: https://github.com/Scalingo/go-graceful-restart-example

## Run the server

```
$ go run ping.go
2014/12/14 20:26:42 [Server - 4301] Listen on [::]:12345
[...]
```

## Connect with the client

```
$ go run client/pong.go
```

## Graceful restart

```
# The server pid is included in its log, in the example: 4301

$ kill -HUP <server pid>
```

## Stop with timeout

Let 10 seconds for the current requests to finish.

```
$ kill -TERM <server pid>
```
