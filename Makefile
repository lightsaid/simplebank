## postgresql DSN
DATABASE_URL=postgresql://postgres:abc123@localhost:5432/simple_bank?sslmode=disable

## migrate: 生成迁移sql文件, exp: make migrate NAME=init_db
migrate:
	migrate create -seq -ext=.sql -dir=./db/migrations $$NAME

## migrate_up: 向上迁移
migrate_up:
	migrate -database ${DATABASE_URL} -path ./db/migrations -verbose up 1

## migrate_down: 向下迁移
migrate_down:
	migrate -database ${DATABASE_URL} -path ./db/migrations -verbose down 1

## 例如： make migrate_force V=1
migrate_force:
	migrate -database ${DATABASE_URL} -path ./db/migrations force $$V
