include .env

stop_containers:
	@echo "Stopping other docker containers"
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers"; \
		docker stop $$(docker ps -q); \
	else \
		echo "no containers running"; \
	fi

create_container:
	docker run --name ${DB_DOCKER_CONTAINER} -p ${DB_PORT}:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:12-alpine

create_db:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}