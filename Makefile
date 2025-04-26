run:
	@go run cmd/api/main.go
migrate-up:
	@go run cmd/api/main.go -m up
