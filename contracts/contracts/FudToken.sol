//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract FudToken is ERC20 {
    uint public constant MAX_SUPPLY = 1500000;

    constructor(uint256 supply) ERC20("FUD Token", "FUD") {
        require(supply <= MAX_SUPPLY, "max supply exceeded");
        _mint(msg.sender, supply);
    }

    function mint(address to, uint256 amount) public {
        require(totalSupply() + amount <= MAX_SUPPLY, "max supply exceeded");
        _mint(to, amount);
    }
}
