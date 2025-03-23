generate:
	mkdir -p pb
	protoc --proto_path=proto proto/*.proto \
	  --go_out=pb --go_opt=paths=source_relative \
	  --go-grpc_out=pb --go-grpc_opt=paths=source_relative

server:
	go run server/main.go 

client: 
	go run client/main.go 
