swag:
	swag init -d ./api
local:
	echo "Starting local environment"
	docker-compose -f docker-compose.yml up --build

down:
	docker-compose -f docker-compose.yml down


