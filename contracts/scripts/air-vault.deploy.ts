import { ethers } from "hardhat";
import {
  FudToken__factory,
  WinToken__factory,
  AirVault__factory,
  FudToken,
  WinToken,
  AirVault,
} from "../types";

async function main() {
  const [deployer, user1] = await ethers.getSigners();

  const initialSupply = 1000000;

  console.log(
    "Deploying FudToken contract with the account:",
    deployer.address
  );
  console.log("Account balance:", (await deployer.getBalance()).toString());

  const FudTokenContract = (await ethers.getContractFactory(
    "FudToken"
  )) as FudToken__factory;
  const fudTokenContract = (await FudTokenContract.deploy(
    initialSupply
  )) as FudToken;
  await fudTokenContract.deployed();
  console.log("FudToken address:", fudTokenContract.address);

  console.log(
    "Deploying WinToken contract with the account:",
    deployer.address
  );
  console.log("Account balance:", (await deployer.getBalance()).toString());

  const WinTokenContract = (await ethers.getContractFactory(
    "WinToken"
  )) as WinToken__factory;
  const winTokenContract = (await WinTokenContract.deploy()) as WinToken;
  await winTokenContract.deployed();

  console.log("WinToken address:", winTokenContract.address);

  console.log(
    "Deploying AirVault contract with the account:",
    deployer.address
  );
  console.log("Account balance:", (await deployer.getBalance()).toString());

  const AirVaultContract = (await ethers.getContractFactory(
    "AirVault"
  )) as AirVault__factory;
  const airVaultContract = (await AirVaultContract.deploy(
    fudTokenContract.address,
    winTokenContract.address
  )) as AirVault;
  await airVaultContract.deployed();

  console.log("AirVault address:", airVaultContract.address);

  await fudTokenContract.transfer(user1.address, 10000);
  const user1FudBalance = await fudTokenContract.balanceOf(user1.address);
  console.log(
    `user1 with address ${user1.address} has ${user1FudBalance} FUD Tokens`
  );
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
