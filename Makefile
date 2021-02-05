.PHONY: issuer

build: build-acapy build-issuer build-tails-server

build-tails-server:
	docker-compose build tails-server

build-acapy:
	docker-compose build acapy

build-issuer:
	docker-compose build issuer

up:
	docker-compose up --force-recreate -d

issuer:
	docker-compose up --force-recreate --no-deps -d issuer

logs:
	docker-compose logs -f --tail=100

down:
	docker-compose down

up-prod:
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod pull
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up --force-recreate -d

issuer-prod:
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod pull issuer
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up --force-recreate --no-deps -d issuer

acapy-prod:
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod pull acapy
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up --force-recreate --no-deps -d acapy


logs-remote:
	docker-compose --context remote -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod logs -f

push:
	docker push ldej/acapy
	docker push ldej/issuer
	docker push ldej/tails-server