build: build-acapy build-issuer

build-acapy:
	docker-compose --env-file .env.prod build acapy

build-issuer:
	docker-compose --env-file .env.prod build issuer

start-db:
	docker-compose -f docker-compose.local.yaml up -d db

start-local:
	docker-compose -f docker-compose.local.yaml --env-file .env.local up --force-recreate --no-deps -d acapy issuer

start-prod:
	docker-compose -f docker-compose.prod.yaml --env-file .env.prod up --force-recreate --no-deps -d acapy issuer nginx

logs-local:
	docker-compose -f docker-compose.local.yaml --env-file .env.local logs -f