DB_URL=postgresql://schooli:schooli@localhost:5432/schooli?sslmode=disable

migrateup:
	migrate -path assets/migrations -database "$(DB_URL)" -verbose up
