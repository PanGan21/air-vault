//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract WinToken is ERC20 {
    address public immutable minter;

    constructor() ERC20("WIN Token", "WIN") {
        minter = msg.sender;
    }

    function mint(address to, uint256 amount) public returns (bool) {
        require(msg.sender == minter, "only minter can mint new tokens");
        _mint(to, amount);
        return true;
    }
}
