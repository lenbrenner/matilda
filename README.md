# Description
The purpose of this project is to explore paths through floor. I am using
Location to represent a square on a grid. The square have transitions from
which we can build a graph. On initialization, we will build all paths from
one location to another. By calculating these ahead of time we should be able
to reduce overall traffic.

# Notes and important permalinks

## Database

### This is where our schema is defined. Yes it should be moved to a file.
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/applications/databaseFactory.go#L10

### Insert - 
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/daos/LocationDao.go#L16

### GetAll -
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/daos/LocationDao.go#L26

## Services

### Uses LocationDao and TransitionDao with a transaction -
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/services/LocationService.go#L102

### Transactions begin and end here
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/services/LocationService.go#L103
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/services/LocationService.go#L117

### This is how I do dependency injection
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/services/LocationService.go#L21

### This is how I bind services
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/applications/Application.go#L26

### This is how I mock a service or DAO testing
https://github.com/lenbrenner/matilda/blob/d1dc267acb41792752cfe958451cb2668c815997/services/LocationService_test.go#L65

## gRPC Apis using protocol buffers

### The idl lives here
https://github.com/lenbrenner/matilda/blob/main/api/matilda.proto

### Client and server stubs are generated using this command and are generated into the api directory as well
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative apis/matilda.proto

### Server code 
https://github.com/lenbrenner/matilda/blob/be5a150930d8e85b7a49dd03dffdaefa4bc2e3c0/server/server.go#L57

### The conversion of entities to resources are a but lengthy. Go 1.18(similar structs map) is one fix
https://github.com/lenbrenner/matilda/blob/853fa813189f2f20427514d826f0acd6070d6ef7/server/server.go#L46

### Client code
https://github.com/lenbrenner/matilda/blob/main/client/client.go

## For web we could use grpc-web or but I believe the decision is to use RESTful APIs so I am investigating:
https://github.com/improbable-eng/grpc-web

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
                              
The stack is composed of:
    gRPC
    sqlx
    Postgres
