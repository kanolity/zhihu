# 仿知乎网页后端
Based on Golang and go-zero.

Before running the project, please make sure to start both **etcd** and **Redis**.
## Notice
To run user service and article service at one time,you should make sure their port should be different

## User Service
### User API Service
The configuration files are located in the `application/user/api/config` directory.  

**Run user rpc first** 
```bash
go run ./application/user/api/user.go
```
### User RPC
The configuration files are located in the `application/user/rpc/config` directory.
```bash
go run ./application/user/rpc/user.go
```

## Article Service
### Article API
The configuration files are located in the `application/article/api/config` directory.

**Run article rpc first**
```bash
go run ./application/user/api/article.go
```
### Article RPC
The configuration files are located in the `application/article/rpc/config` directory.
```bash
go run ./application/user/rpc/article.go
```