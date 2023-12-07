swag:
	swag init -d ./api
cserv:
	echo "Starting customer service"
	docker-compose -f docker-compose.customer.yml up --build
aserv:
	echo "Starting account service"
	docker-compose -f docker-compose.account.yml up --build
down:
	docker-compose -f docker-compose.customer.yml down
	docker-compose -f docker-compose.account.yml down
serv:
	echo "Starting app"
	docker-compose -f docker-compose.customer.yml -f docker-compose.account.yml up --build
