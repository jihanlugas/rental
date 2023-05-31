server:
	swag init
	go run main.go server

up:
	go run main.go db up

down:
	go run main.go db down

seed:
	go run main.go db seed

reset-db:
	go run main.go db down
	go run main.go db up
	go run main.go db seed