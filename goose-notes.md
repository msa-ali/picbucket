# Migration using goose

```zsh
mkdir migrations
cd migrations

goose create widgets sql -> create a new file

goose postgres "host=localhost port=5432 user=admin password=admin dbname=picbucket sslmode=disable" status

goose postgres "host=localhost port=5432 user=admin password=admin dbname=picbucket sslmode=disable" up

 goose postgres "host=localhost port=5432 user=admin password=admin dbname=picbucket sslmode=disable" down

 goose fix
```
