import { expect } from "chai";
import type {
  BasicControlContract,
  BasicControlContract__factory,
} from "../types";
import { ethers } from "hardhat";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";

describe("BasicControlContract", () => {
  let basicControlContract: BasicControlContract;
  let signers: SignerWithAddress[];
  let eventAdmin: string;

  before(async () => {
    signers = await ethers.getSigners();

    const BasicControlContract = (await ethers.getContractFactory(
      "BasicControlContract",
      signers[0]
    )) as BasicControlContract__factory;

    basicControlContract =
      (await BasicControlContract.deploy()) as BasicControlContract;
    await basicControlContract.deployed();

    basicControlContract.on("AdminRegistered", (c: string) => {
      eventAdmin = c;
    });

    //required for the event https://github.com/nomiclabs/hardhat/issues/1692#issuecomment-904904674
    await new Promise((res) => setTimeout(() => res(null), 5000));
  });

  it("should set the admin and have the event fired", async () => {
    const admin = await basicControlContract.admin();

    expect(admin).to.equals(signers[0].address);
    expect(eventAdmin).to.equals(signers[0].address);
  });

  it("should allow accessing functions with admin", async () => {
    await expect(basicControlContract.publicFunction()).not.to.be.reverted;
    await expect(basicControlContract.privateFunction1()).not.to.be.reverted;
    await expect(basicControlContract.privateFunction2()).not.to.be.reverted;
  });

  it("should revert when accessing not allowed function only", async () => {
    await expect(basicControlContract.connect(signers[1]).publicFunction()).not
      .to.be.reverted;
    await expect(
      basicControlContract.connect(signers[1]).privateFunction1()
    ).to.be.revertedWith("Only admin");
    await expect(
      basicControlContract.connect(signers[1]).privateFunction2()
    ).to.be.revertedWith("Only admin");
  });
});
