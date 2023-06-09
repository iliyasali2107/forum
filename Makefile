all: build run


server: 
	go run ./cmd/main.go



build:
	docker build -t forum .
run:
	docker run -dp 8080:8080 --rm --name forum_container forum
