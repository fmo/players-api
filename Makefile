PLAYERS_API_BINARY=playersApiApp

players_api:
	@echo "Building binary..."
	go build -o ${PLAYERS_API_BINARY} ./cmd/api/
	@echo "Done!"

players_api_amd:
	@echo "Building binary..."
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${PLAYERS_API_BINARY} ./cmd/api/
	@echo "Done!"
