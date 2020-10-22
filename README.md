# Issuer

## Checkout

```
$ git clone --recursive git@github.com:ldej/issuer.git
```

## Postgres

Running aca-py with `--wallet-storage-type postgres_storage`

`OSError: libindystrgpostgres.so: cannot open shared object file: No such file or directory`

```
$ sudo apt install -y cargo libzmq3-dev
$ cd indy-sdk/experimental/plugins/postgres_storage
$ cargo build
$ export LD_LIBRARY_PATH=/home/laurencedejong/projects/aries-developer/indy-sdk/experimental/plugins/postgres_storage/target/debug
```

## Create a new public DID with indy-cli

```
$ sudo apt install -y indy-cli

$ indy-cli --config cliconfig.json
indy> pool create buildernet gen_txn_file=pool_transactions_builder_genesis
(only the first time)

indy> pool connect buildernet
Would you like to read it? (y/n)
(select y)
Would you like to accept it? (y/n)
(select y)

indy> wallet create issuer key=issuer
indy> wallet open issuer key=issuer
indy> did new (seed=<32 character secret seed> optional)
 -> Go to https://selfserve.sovrin.org/ and enter the did and verkey
 -> {"statusCode":200,"headers":{"Access-Control-Allow-Origin":"*"},"body":"{\"statusCode\": 200, \"DjPCcebRjN4XRA2F7gR8hw\": {\"status\": \"Success\", \"statusCode\": 200, \"reason\": \"Successfully wrote NYM identified by DjPCcebRjN4XRA2F7gR8hw to the ledger with role ENDORSER\"}}"}
indy> ledger get-nym did=DjPCcebRjN4XRA2F7gR8hw
indy> did use DjPCcebRjN4XRA2F7gR8hw
indy> ledger schema name=MyFirstSchema version=1.0 attr_names=FirstName,LastName,Address,Birthdate,SSN
```

## Using indy-cli with postgres

```
LD_LIBRARY_PATH=/home/laurencedejong/projects/aries-developer/indy-sdk/experimental/plugins/postgres_storage/target/debug indy-cli --config cliconfig.json
indy> load-plugin library=/home/laurencedejong/projects/aries-developer/indy-sdk/experimental/plugins/postgres_storage/target/debug/libindystrgpostgres.so initializer=postgresstorage_init
indy> wallet create wallet_psx key storage_type=postgres_storage storage_config={"url":"localhost:5432"} storage_credentials={"account":"postgres","password":"mysecretpassword","admin_account":"postgres","admin_password":"mysecretpassword"}
indy> wallet open wallet_psx key storage_credentials={"account":"postgres","password":"mysecretpassword","admin_account":"postgres","admin_password":"mysecretpassword"}
```