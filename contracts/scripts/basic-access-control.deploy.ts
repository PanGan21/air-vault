import { ethers } from "hardhat";
import { BasicControlContract__factory, BasicControlContract } from "../types";

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Deploying contracts with the account:", deployer.address);
  console.log("Account balance:", (await deployer.getBalance()).toString());

  const BasicControlContract = (await ethers.getContractFactory(
    "BasicControlContract"
  )) as BasicControlContract__factory;

  const basicControlContract =
    (await BasicControlContract.deploy()) as BasicControlContract;

  console.log("Token address:", basicControlContract.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
