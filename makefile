swag:
	swag init -d ./api
cserv:
	echo "Starting customer service"
	docker-compose -f docker-compose.customer.yml up --build
aserv:
	echo "Starting account service"
	docker-compose -f docker-compose.account.yml up --build
serv:
	echo "Starting services"
	docker-compose -f docker-compose.customer.yml up --build
	docker-compose -f docker-compose.account.yml up --build
down:
	docker-compose -f docker-compose.yml down


