build:
	@ printf "Building application... "
	@ go build -trimpath -o ./bin/engine ./app/
	@ echo "Done!"

build-migrate:
	@ printf "Building migrate... "
	@ go build -trimpath -o ./bin/migrate ./cmd/migrate.go

up:	dev-env

up-local: dev-env-local dev-air

logs:
	@ docker logs hewpao-backend-app-1 --follow

dev-env:
	@ docker compose up

dev-env-local:
	@ docker compose up -d --build db db-manager

dev-air:
	@ air

migrate-dev: migrate-inside-container

migrate-dev-local:
	@ go run ./cmd/migrate.go

migrate-inside-container:
	@ docker exec -it hewpao-backend-app-1 go run ./cmd/migrate.go 

build-image:
	@ echo "Building docker image..."
	@ docker build \
		--file ./docker/prod.Dockerfile \
		--tag hewpao/hewpao-backend \
		.

test:
	@ go test -v -cover ./usecase/

test-logs:
	@ go test -coverprofile=test-coverage/coverage.out ./usecase/
	@ go tool cover -html=test-coverage/coverage.out -o test-coverage/coverage.html
	@ open test-coverage/coverage.html
