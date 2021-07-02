start docker instance command:
docker-compose up -d

stop docker instance command:
docker stop pg-docker

inspect docker instance command:
docker inspect pg-docker

docker postgres login:
psql -h localhost -p 5432 -U postgres postgres

posgres docker container IP address: 172.18.0.2