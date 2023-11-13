# setting developemnt flag 
export APP_ENV = development
$ENV:APP_ENV = "development"
dev:
	nodemon --exec go run . --signal SIGTERM
test:
	go test ./...
build:
	go build .