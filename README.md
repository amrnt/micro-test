## How to run

Open 4 terminals

```
# First download deps
go mod download

# In terminal 1 to run rabbitmq
docker-compose up

# in terminal 2
go run srv/sub/*.go

# in terminal 3
go run api/pub/*.go

# in terminal 4
curl http://localhost:3000
```
