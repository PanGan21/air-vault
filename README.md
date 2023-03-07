# air-vault

A poc for airdrops

## Airdrop system

WIN tokens will be distributed by an airdrop every X blocks to each user who has had FUD tokens deposited during this interval. A separate Golang service will airdrop WIN tokens by minting them directly to the depositor's address. The amount of airdropped tokens is equal to 5% of the average FUD token deposit over the last X blocks (X should be configurable)

### Contracts directory

The `contracts` directory is a hardhat project for local development of the the smart contracts required for the system. It includes the smart contracts and their unit tests.
<br/>

#### Run

In order to run the smart contracts tests run the following commands:
<br/>
`yarn install && yarn compile`
<br/>
`yarn test`

### Backend directory

The `backend` directory is the golang service which monitors the blockchain for events. It is a cli application that is based on environment variables specified in yaml format. The cli interface is included in the `cmd` directory. The smart contracts and their go bindings are expanded in the `contracts` directory using the [go-ethereum](https://github.com/ethereum/go-ethereum) library and the `abigen` tool. The core logic of the application is in `pkg` directory, which interacts with the `repository` directory as it's storage. The storage is only in memory for now, which is also the tradeoff of this poc. If the backend crashes it can miss events which could result to missing airdrop mints. In case of adding a persistent storage there could be a Filter of historical events on startup so the backend can synchronize with the missing events before accepting new ones.
<br/>

#### Solution

The solution is based on chunks. An airdrop storage has a current block number in the interval (1-100). There is a chunk of FUD Tokens balance added in the list with the number of blocks that each balance was in deposit for each user. The chunks are being updated for every message received. When there is an interval of blocks passed an airdrop is triggered for all users and the chunks are being cleaned. The last airdrop block number is also being stored, and is taken into account to decide when to trigger the airdrop.

The following illustrates the described solution:

```
airdrop storage => current block number in interval (1-100)
user starts with 0
at block number 1 deposits 10  => TotalFUD = 10 => airvault=(10*nil)
at block number 40 user deposits 10 => airvault=(10*40)+((10+10)*nil)
etc
```

#### Run

In order to run the backend:
<br/>
Build the backend: `make build`
<br/>
Run a ganache service: `make ganache`
<br/>
Deploy the smart contracts: `make deploy`
<br/>
Transfer some initial FUD Tokens to the demo user: `make transfer-fud`
<br/>
Start depositing and withdrawing FUD Tokens with this user until there is an airdrop using: `make deposit` or `make withdraw`

The private key, the account user and the amounts to be deposited and withdrawn are specified in `backend/config/local.yaml` where they can be changed for more combinations.

#### Private key handling

The private key is already passed in the application as an environment variable from a yaml file. In order to make it more secure, the private key can become a secret that is being unwrapped and passed through pipelines which integrate with some kind of vault.

### Improvements

- Backend testing: unit tests as well as integration tests should be written
- Database: Adding a persistent storage to the application could help not miss event.
- Filter logs: On startup a filter of previous logs could be fired. The filter should be fired for blocks newer than the last stored block in the database.
