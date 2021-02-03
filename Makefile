build: build-acapy build-issuer build-tails-server

build-tails-server:
	docker-compose --env-file .env.prod build tails-server

build-acapy:
	docker-compose --env-file .env.prod build acapy

build-issuer:
	docker-compose --env-file .env.prod build issuer

up:
	docker-compose -f docker-compose.local.yml --env-file .env.local up -d

logs:
	docker-compose -f docker-compose.local.yml --env-file .env.local logs -f

down:
	docker-compose -f docker-compose.local.yml --env-file .env.local down

up-prod:
	docker-compose -f docker-compose.prod.yml --env-file .env.prod up --force-recreate --no-deps -d acapy issuer nginx
