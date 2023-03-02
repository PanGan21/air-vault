build:
	go build -o bin/airvault

run: build
	 ./bin/airvault

test:
	go test -v ./...

ganache:
	ganache-cli -d run west attitude bronze weapon goat spell coyote text image ignore lamp

compile-contracts:
	solc --optimize --abi ./contracts/solidity/AirVault.sol -o contracts/bin

generate-airvault-interface:
	abigen --abi ./contracts/bin/AirVault.abi --pkg contracts --out ./contracts/interfaces/air-vault/AirVault.go

generate-fud-interface:
	abigen --abi ./contracts/bin/FudToken.abi --pkg contracts --out ./contracts/interfaces/fud-token/FudToken.go

generate-win-interface:
	abigen --abi ./contracts/bin/WinToken.abi --pkg contracts --out ./contracts/interfaces/win-token/WinToken.go

.PHONY: build, run, test, ganache, compile-airvault, generate-contract-interface