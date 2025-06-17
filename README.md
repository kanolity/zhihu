# 仿知乎网页后端
## run program
based on Golang and go-zero 
Before running the project, please make sure to start both **etcd** and **Redis**.
### user api 
The configuration files are located in the `application/user/api/config` directory.
**run user rpc first** 
```bash
go run ./application/user/api/user.go
```
### user rpc
The configuration files are located in the `application/user/api/config` directory.
```bash
go run ./application/user/rpc/user.go
```
