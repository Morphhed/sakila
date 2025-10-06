.PHONY: mysql createall dropdb stop start rm migrate-up migrate-down updb sqlc

# Jalankan MySQL container
mysql:
	docker run --name sakila-mysql \
		-e MYSQL_ROOT_PASSWORD=123 \
		-p 3306:3306 \
		-d mysql:8.0

# Buat database sakila (import schema + data)
updb:
	docker exec -i sakila-mysql mysql -uroot -p123 -e "CREATE DATABASE sakila;"

# Hapus database sakila
dropdb:
	docker exec -i sakila-mysql \
		mysql -uroot -p123 -e "DROP DATABASE IF EXISTS sakila;"




# Stop container
stop:
	docker stop sakila-mysql

# start container
start:
	docker start sakila-mysql

# Remove container
rm:
	docker rm -f sakila-mysql




# Buat database sakila (import schema + data)
createall:
	docker exec -i sakila-mysql \
		mysql -uroot -p123 < /home/ciruno/proyek/sakila/rizz/db/query/sakila-schema.sql
	docker exec -i sakila-mysql \
		mysql -uroot -p123 < /home/ciruno/proyek/sakila/rizz/db/query/sakila-data.sql

# Jalankan migration UP
migrate-up:
	migrate -path /home/ciruno/proyek/sakila/rizz/db/migration \
		-database "mysql://root:123@tcp(localhost:3306)/sakila" -verbose up

# Jalankan migration DOWN
migrate-down:
	migrate -path /home/ciruno/proyek/sakila/rizz/db/migration \
		-database "mysql://root:123@tcp(localhost:3306)/sakila" -verbose down

# Jalankan SQLC
sqlc:
	sqlc generate 
