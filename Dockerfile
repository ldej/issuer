FROM ubuntu:18.04

RUN apt-get update && apt-get install -y gnupg2 software-properties-common python3-pip cargo libzmq3-dev libsodium-dev pkg-config libssl-dev
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 68DB5E88 && \
    add-apt-repository "deb https://repo.sovrin.org/sdk/deb bionic master" && \
    apt-get update && \
    apt-get install -y libindy

# Build libindystrgpostgres.so for connecting to postgres
# requires cargo libzmq3-dev libsodium-dev pkg-config libssl-dev
ADD indy-sdk /indy-sdk
RUN cd /indy-sdk/experimental/plugins/postgres_storage && cargo build

# Install ACA-py
# requires python3-pip
COPY requirements.txt /
RUN pip3 install -r /requirements.txt

# Add genesis file and ACA-py configuration
COPY pool_transactions_builder_genesis /
COPY arguments.*.yaml /

# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

# Announce location of libindystrgpostgres.so
ENV LD_LIBRARY_PATH /indy-sdk/experimental/plugins/postgres_storage/target/debug

CMD /wait && aca-py start --arg-file ./arguments.start.yaml