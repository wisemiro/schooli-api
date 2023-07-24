DB_URL=postgresql://schooli:schooli@localhost:5432/schooli?sslmode=disable

migrateup:
	migrate -path assets/migrations -database "$(DB_URL)" -verbose up

migrations:
	 @read -p "Migration Name? : " Name; \
 	 migrate create -ext sql -dir assets/migrations -seq $${Name}

migrateupF:
	 @read -p "Version to force to ? : " version; \
	migrate -path assets/migrations -database "$(DB_URL)" force $${version}

migrateDA:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateDV:
	 @read -p "Version to downgrade to ? : " version; \
	migrate -path internal/postgresql/migrations -database "$(DB_URL)" down $${version}