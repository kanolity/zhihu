# 仿知乎网页后端
## Run program
Based on Golang and go-zero.  

Before running the project, please make sure to start both **etcd** and **Redis**.
### User API
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