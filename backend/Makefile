build:
	go build -o bin/airvault

run: build
	 ./bin/airvault run

test:
	go test -v ./...

ganache:
	# ganache-cli --account="0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef,1000000000000000000000000000000000000000" --mnemonic="range pear quit paddle harvest glory insect tissue erupt spray sport child"
	ganache-cli --mnemonic="range pear quit paddle harvest glory insect tissue erupt spray sport child"

compile-contracts:
	solc --optimize --abi ./contracts/solidity/AirVault.sol -o contracts/bin --overwrite
	solc --optimize --bin ./contracts/solidity/AirVault.sol -o contracts/bin --overwrite

generate-airvault-interface:
	abigen --abi=./contracts/bin/AirVault.abi --bin=./contracts/bin/AirVault.bin --pkg=contracts --out ./contracts/interfaces/air-vault/AirVault.go
	
generate-fud-interface:
	abigen --abi=./contracts/bin/FudToken.abi --bin=./contracts/bin/FudToken.bin --pkg=contracts --out ./contracts/interfaces/fud-token/FudToken.go

generate-win-interface:
	abigen --abi ./contracts/bin/WinToken.abi --bin=./contracts/bin/WinToken.bin --pkg contracts --out ./contracts/interfaces/win-token/WinToken.go

setup-contracts: compile-contracts generate-airvault-interface generate-fud-interface generate-win-interface

deploy:
	go run main.go deploy

deposit:
	./bin/airvault deposit

transfer-fud:
	./bin/airvault transfer-fud

.PHONY: build, run, test, ganache, compile-airvault, generate-contract-interface, setup-contracts, deploy, deposit, transfer-fud