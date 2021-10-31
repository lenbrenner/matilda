# Description
Matilda is a waltzer.

The purpose of this project is to explore paths through floor. I am using
Location to represent a square on a grid. The square have transitions from
which we can build a graph. On initialization, we will build all paths from
one location to another. By calculating these ahead of time we should be able
to reduce overall traffic.

# Run the sample code
To compile and run the server, assuming you are in the root where this README.md
file lives:

```sh
$ go run server/server.go
```

Likewise, to run the client:

```sh
$ go run client/client.go
```

# Optional command line flags
The server and client both take optional command line flags. For example, the
client and server run without TLS by default. To enable TLS:

```sh
$ go run server/server.go -tls=true
```

and

```sh
$ go run client/client.go -tls=true
```
