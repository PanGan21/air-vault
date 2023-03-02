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
  const [deployer] = await ethers.getSigners();

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

  console.log("AirVault address:", airVaultContract.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
