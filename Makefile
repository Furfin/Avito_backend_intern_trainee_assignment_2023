all: test start


test:
	docker-compose up -d db
	go test -v ./handlers/user/
	go test -v ./handlers/segment/
	docker-compose down

start:
	docker-compose up --build