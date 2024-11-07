all: run

run: build
	docker-compose up

build:
	docker-compose build

clean:
	docker-compose down