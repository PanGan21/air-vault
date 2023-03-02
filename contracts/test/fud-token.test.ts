import { expect } from "chai";
import { Signer } from "ethers";
import { ethers } from "hardhat";
import { FudToken, FudToken__factory } from "../types";

describe("FudToken", () => {
  let fudTokenFactory: FudToken__factory;
  let fudToken: FudToken;
  let owner: Signer;
  let user1: Signer;
  let user2: Signer;

  const INITIAL_SUPPLY = 1000000;
  const MAX_SUPPLY = 1500000;

  before(async () => {
    [owner, user1, user2] = await ethers.getSigners();
    fudTokenFactory = await ethers.getContractFactory("FudToken", owner);
  });

  describe("deployment", () => {
    it("Should set the correct name and symbol", async () => {
      fudToken = await fudTokenFactory.deploy(INITIAL_SUPPLY);
      expect(await fudToken.name()).to.equal("FUD Token");
      expect(await fudToken.symbol()).to.equal("FUD");
    });

    it("Should set the total supply to the initial supply", async () => {
      expect(await fudToken.totalSupply()).to.equal(INITIAL_SUPPLY);
    });

    it("Should not allow an initial supply greater that the max supply", async () => {
      await expect(fudTokenFactory.deploy(MAX_SUPPLY + 1)).to.be.revertedWith(
        "max supply exceeded"
      );
    });
  });

  describe("minting", () => {
    it("Should allow the owner to mint tokens", async () => {
      await fudToken.mint(await user1.getAddress(), 1000);
      expect(await fudToken.totalSupply()).to.equal(INITIAL_SUPPLY + 1000);
      expect(await fudToken.balanceOf(await user1.getAddress())).to.equal(1000);
    });

    it("Should not allow minting more than the max supply", async () => {
      await expect(
        fudToken.mint(await user1.getAddress(), MAX_SUPPLY + 1)
      ).to.be.revertedWith("max supply exceeded");
    });
  });
});
