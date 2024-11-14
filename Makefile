all: build run

run:
	docker-compose up

build:
	docker-compose build

clean:
	docker-compose down
