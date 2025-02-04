build:
	@ printf "Building application... "
	@ go build -trimpath -o ./bin/engine ./app/
	@ echo "Done!"

up:	dev-env dev-air

dev-env:
	@ docker compose up -d --build db db-manager

dev-air:
	@ air

migrate-dev:
	@ go run ./cmd/migrate.go
