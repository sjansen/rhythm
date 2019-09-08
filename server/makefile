deploy:
	GOOS=linux GOARCH=amd64 go build -o build/lambda-handler cmd/lambda-handler/main.go
	cd build && zip -9 lambda-handler.zip lambda-handler
	cd ../terraform/ && terraform apply -auto-approve

runserver:
	go run cmd/localdev/main.go -b ../webui/build/
