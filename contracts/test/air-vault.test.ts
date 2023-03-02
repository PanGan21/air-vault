import { expect } from "chai";
import { Signer } from "ethers";
import { ethers } from "hardhat";
import { AirVault, FudToken, WinToken } from "../types";

describe("AirVault", () => {
  let fudToken: FudToken;
  let winToken: WinToken;
  let airVault: AirVault;
  let owner: Signer;
  let user1: Signer;
  let user2: Signer;
  let users: Signer[];

  beforeEach(async () => {
    [owner, user1, user2, ...users] = await ethers.getSigners();

    const fudTokenFactory = await ethers.getContractFactory("FudToken");
    fudToken = await fudTokenFactory.connect(owner).deploy(1500000);

    const winTokenFactory = await ethers.getContractFactory("WinToken");
    winToken = await winTokenFactory.connect(owner).deploy();

    const airVaultFactory = await ethers.getContractFactory("AirVault");
    airVault = await airVaultFactory
      .connect(owner)
      .deploy(fudToken.address, winToken.address);
    await fudToken.connect(owner).transfer(await user1.getAddress(), 1000);
  });

  describe("deposit", () => {
    it("Should deposit FUD tokens successfully", async () => {
      await fudToken.connect(user1).approve(airVault.address, 100);
      await airVault.connect(user1).deposit(100);
      const balance = await fudToken.balanceOf(await user1.getAddress());
      const lockedBalance = await airVault.lockedBalanceOf(
        await user1.getAddress()
      );
      expect(balance).to.equal(900);
      expect(lockedBalance).to.equal(100);
    });

    it("Should fail to deposit without approval", async () => {
      await expect(airVault.connect(user1).deposit(100)).to.be.revertedWith(
        "ERC20: insufficient allowance"
      );
    });
  });

  describe("withdraw", () => {
    it("Should withdraw FUD tokens successfully", async () => {
      await fudToken.connect(user1).approve(airVault.address, 100);
      await airVault.connect(user1).deposit(100);
      await airVault.connect(user1).withdraw(50);
      const balance = await fudToken.balanceOf(await user1.getAddress());
      const lockedBlanace = await airVault.lockedBalanceOf(
        await user1.getAddress()
      );
      expect(balance).to.equal(950);
      expect(lockedBlanace).to.equal(50);
    });

    it("Should fail to withdraw mote than the deposited amount", async () => {
      await fudToken.connect(user1).approve(airVault.address, 100);
      await airVault.connect(user1).deposit(100);
      await expect(airVault.connect(user1).withdraw(150)).to.be.revertedWith(
        "insufficient balance"
      );
    });
  });
});
