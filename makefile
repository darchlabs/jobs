install:
	@echo "[install] Installing dependencies":
	@npm i -g nodemon
	@go get

dev:
	@echo "[dev] Running":
	@nodemon --exec go run cmd/api/main.go --signal SIGTERM 