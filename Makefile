build:
	@ printf "Building application... "
	@ go build -trimpath -o ./bin/engine ./app/
	@ echo "Done!"

build-migrate:
	@ printf "Building migrate... "
	@ go build -trimpath -o ./bin/migrate ./cmd/migrate.go

up:	dev-env

logs:
	@ docker logs hewpao-backend-app-1 --follow

dev-env:
	@ docker compose up

migrate-dev:
	@ go run ./cmd/migrate.go

build-image:
	@ echo "Building docker image..."
	@ docker build \
		--file ./docker/prod.Dockerfile \
		--tag hewpao/hewpao-backend \
		.
