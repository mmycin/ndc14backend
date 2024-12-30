build:
	@go build -o bin/server cmd/main.go
	@echo Build completed successfully

run:
	@cls
	@go run cmd/main.go
	@pause
	@cls

watch:
	@cls
	@air cmd/main.go

format:
	@go fmt ./...
	@echo "Code formatted successfully"

clean:
	@cls	
	@if exist bin rmdir /s /q bin
	@if exist bin/ rm -rf bin
	@echo Cleanup completed successfully

serve:
	@cls
	@echo "Starting server..."
	@go run cmd/main.go

test:
	@cls
	@go test -v ./tests