server: 
	go run ./cmd/main.go

migratedown:
	go run ./cmd/migration/main.go down

migrateup:
	go run ./cmd/migration/main.go up

