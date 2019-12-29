# Golang gRPC-API

## Start Proxy
```
./cloud_sql_proxy -instances=xxx=tcp:3306 -credential_file=credentials.json
```

## Start Server
```
go run cmd/server/main.go -db-password=xxx -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00 -db-host=xxx
```

## Start Client
```
go run cmd/client-grpc/main.go -server=localhost:8080
```

## Run with Docker
```
docker run -d -p 7000:8080 -v credentials.json:/app/credentials.json -e instances='XXX' -e dbPw='XXX' gcr.io/mytests-262609/goserver:1.0
````

## Deploy application on Kubernetes
```
export my_zone=us-central1-a  
export my_cluster=standard-cluster-1  

gcloud container clusters create $my_cluster \  
  --num-nodes 1 --zone $my_zone  
  
gcloud container clusters get-credentials $my_cluster --zone $my_zone  

kubectl create secret generic grpc-project-x-key \  
 --from-file=credentials.json=credentials.json  

kubectl create secret generic cloudsql-pw \  
 --from-literal=rootPw=XXX  
 
kubectl create configmap myconfigmap --from-literal=cloudSqlInstance=XXX  

kubectl applay -f deployment.yaml  
```

