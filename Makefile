start-db:
	docker-compose up

stop-db:
	docker compose down

db-cmd:
	docker exec -it lenslocked-db-1 /usr/bin/psql -U admin -d picbucket
