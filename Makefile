
.PHONY: help run dev start build docker start down logs clean

help:
	@echo "Comandos disponíveis:"
	@echo "  make run     	 - Roda o client em Go"
	@echo "  make dev     	 - Roda o client em Go"
	@echo "  make start      - Roda o client em Go"
	@echo "  make build      - Compila o client em Go"

run dev start:
	go mod tidy & go run src/main.go

build:
	go mod tidy && go build -o dist/devpipe src/main.go

docker start:
	docker build -t devpipe-cli .
	docker run -it --rm devpipe-cli 

clean:
	rm -f dist
