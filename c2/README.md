# Getting protoc working
Because Google has absolutely no idea how to update docs, you need the following to get this working:
* protoc
* protoc-gen-go
* protoc-gen-go-grpc

# Intalling Protoc
## Ubuntu 
```shell
$ sudo apt install -y protobuf-compiler
```
## Arch
```shell
$ sudo pacman -S protobuf
```
## MacOS 
```shell
$ brew install protobuf
```

# Installing the other packages 
```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

# Generating the code
```shell
$ protoc --go-grpc_out=. --go_out=. spec/innocent_pictures.proto
```