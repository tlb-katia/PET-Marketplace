serv = server
db = psql_db
cache = redis_db

all: ##@APP application in docker container
	docker-compose-api

docker-compose-api: ##@APP runs application in docker container
	docker build --no-cache -t $(serv) .
	docker-compose up

clean-data: ##@DB clean a database saved data
	rm -rf pkg/repository/db/pgdata
	rm -rf pkg/repository/redis/data
	rm -rf pkg/repository/redis/redis.conf

docker-stop-api: ##@SERVER stops containers
	docker stop $(db)
	docker stop $(serv)
	docker stop $(cache)

docker-clean-api: docker-stop-api ##@SERVER delete server, database and cache containers
	docker rm $(db)
	docker rm $(serv)
	docker rm $(cache)

server-logs: ##@SERVER show logs from server container
	docker logs $(serv)

database-logs:  ##@DB show logs from database container
	docker logs $(db)

cache-logs: ##@CACHE show logs from cache container
	docker logs $(cache)

all-logs: database-logs server-logs cache-logs ##@APP show logs from server and db containers together
