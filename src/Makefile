build: rolctl rol

rolctl: webapi/swagger/swagger.json
	swagger generate cli -f webapi/swagger/swagger.json -A rol-cmd #--skip-validation
	go mod tidy
	go build -o rolctl cmd/cli/main.go

rol:
	go mod tidy
	go build

clean:
	rm -rf rol rolctl cli/ client/ cmd/ models/
