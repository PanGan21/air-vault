import { expect } from "chai";
import { Signer } from "ethers";
import { ethers } from "hardhat";
import { WinToken } from "../types";

describe("WinToken", () => {
  let winToken: WinToken;
  let owner: Signer;
  let user1: Signer;
  let user2: Signer;

  beforeEach(async () => {
    [owner, user1, user2] = await ethers.getSigners();
    const winTokenFactory = await ethers.getContractFactory("WinToken", owner);
    winToken = await winTokenFactory.deploy();
    await winToken.deployed();
  });

  describe("deployment", () => {
    it("Should have the correct name and symbol", async () => {
      expect(await winToken.name()).to.equal("WIN Token");
      expect(await winToken.symbol()).to.equal("WIN");
    });
  });

  describe("mint", () => {
    it("Should be able to mint the new tokens", async () => {
      await winToken.connect(owner).mint(await user1.getAddress(), 100);
      expect(await winToken.balanceOf(await user1.getAddress())).to.equal(100);
    });

    it("Should not be able to mint new tokens by non owner", async () => {
      await expect(
        winToken.connect(user1).mint(await user2.getAddress(), 100)
      ).to.be.revertedWith("only minter can mint new tokens");
    });
  });
});
