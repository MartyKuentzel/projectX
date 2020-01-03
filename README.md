# Golang gRPC-API

## Start Proxy
```
./cloud_sql_proxy -instances=xxx=tcp:3306 -credential_file=credentials.json
```

## Start Server
```
go run cmd/server/main.go -db-password=xxx -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00
```

## Start Client
```
go run cmd/client-grpc/main.go -server=localhost:8080
```

