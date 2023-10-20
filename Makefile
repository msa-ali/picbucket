start-db:
	docker-compose up

stop-db:
	docker compose down

db-cmd:
	docker exec -it lenslocked-db-1 /usr/bin/psql -U admin -d picbucket

# Migration Commands
mig-up:
	cd migrations && goose postgres "host=localhost port=5432 user=admin password=admin dbname=picbucket sslmode=disable" up

mig-down:
	cd migrations && goose postgres "host=localhost port=5432 user=admin password=admin dbname=picbucket sslmode=disable" down

mig-status:
	cd migrations && goose postgres "host=localhost port=5432 user=admin password=admin dbname=picbucket sslmode=disable" status

mig-fix:
	cd migrations && goose fix