# certwatch
CertWatch is a lightweight, locally deployable monitoring system that tracks SSL/TLS certificate expiration.

It allows users to input their domain names, set custom alerts for upcoming expirations, and stores data locally in SQLite. 

The system reads configurations from a file in its first version, providing an easy setup for anyone concerned about certificate management.


## 运行
```
go run cmd/certwatch.go
```

编译后运行
```
go build -o certwatch cmd/certwatch.go
./certwatch
```


中英文介绍