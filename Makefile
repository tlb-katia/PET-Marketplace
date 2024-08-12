serv = server
db = psql_db
cache = redis_db

all: docker-compose-api ##@APP application in docker container


docker-compose-api: permission ##@APP runs application in docker container
	docker-compose up --build $(serv)

clean-data: permission ##@DB clean a database saved data
	rm -rf internal/storage/postgres/pgdata
	rm -rf internal/storage/redis/data

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

permission:
	sudo chmod -R 777 /home/katia/GolandProjects/Marketplace/internal/storage/postgres/pgdata

all-logs: database-logs server-logs cache-logs ##@APP show logs from server and db containers together
