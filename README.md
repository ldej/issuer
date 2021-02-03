# Issuer

## Checkout

```
$ git clone --recursive git@github.com:ldej/issuer.git
```

## Running locally

### Start a VON-network ledger

[github.com/bcgov/von-network](https://github.com/bcgov/von-network)

Start 4 Indy nodes and the von-webserver. The von-webserver has a web interface at [localhost:9000](http://localhost:9000) which allows you to browse the transactions in the blockchain.

```shell script
$ git clone https://github.com/bcgov/von-network
$ cd von-network
$ ./manage start --logs
```

### Start a Tails server

[github.com/bcgov/indy-tails-server](https://github.com/bcgov/indy-tails-server)

Start a Tails server for the revocation registry tails files.

```shell script
$ git clone https://github.com/bcgov/indy-tails-server
$ cd indy-tails-server
$ ./docker/manage start
```

### Create an environment file

```shell
$ cat .env.local
AGENT_WALLET_SEED=<some-32-char-wallet-seed>
LABEL=<name-of-your-application>
ACAPY_ENDPOINT_PORT=8000
ACAPY_ENDPOINT_URL=http://localhost:8000/
ACAPY_ADMIN_PORT=11000
LEDGER_URL=http://172.17.0.1:9000
ISSUER_PORT=8080
WALLET_NAME=<wallet-name>
WALLET_KEY=<secret>
```

### Start

```shell
$ make start-db
$ make start-local
$ make logs-local
```

## ACA-py docker image

The ACA-py docker image is made with the [acapy.dockerfile](./docker/acapy.dockerfile). It is a custom image where libindy is installed and the postgres plugin is installed as a wallet storage backend. I could only install the postgres plugin with the `indy-sdk` repository, that's why it is a git submodule. `aries-cloudagent-python` is included in this repo as a submodule, so I can run the latest ~master~, I mean _main_ branch.

## Issuer docker image

The issuer docker image is used for both building and running the Go application.

## nginx and certbot

I used [this blog post](https://medium.com/@pentacent/nginx-and-lets-encrypt-with-docker-in-less-than-5-minutes-b4b8a60d3a71) as a source of inspiration for getting the easiest set up to work. That's also where `init-letsencrypt.sh` comes from.

## docker-compose

I tried to understand the [aries-cloudagent-python/deploymentModel.md](https://github.com/hyperledger/aries-cloudagent-python/blob/main/docs/deploymentModel.md), but it was too much to read. The two examples at the bottom ([indy-email-verification](https://github.com/bcgov/indy-email-verification
) and [iiwbook](https://github.com/bcgov/iiwbook)) helped me get in the right direction with the `docker-compose.yml` file.

## Deployment

This issuer is deployed on Digital Ocean using the cheapest pre-installed docker droplet. Apparently the `ufw` firewill is enabled by default.

https://www.digitalocean.com/docs/networking/firewalls/resources/troubleshooting/

## TODO

- Automate deployment using Github Actions
- Push docker images to a registry (which one?, the cheapest one!)
- Add functionality for issuing credentials
- Add a frontend