
build:
	docker build . -t ldej.nl/issuer

start-db:
	docker-compose up --scale issuer=0

provision:
	# aca-py provision --arg-file ./arguments.provision.yaml
	docker run --net=host -it ldej.nl/issuer /bin/bash -c "aca-py provision --arg-file ./arguments.provision.yaml"

start:
	docker-compose up