//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./openzeppelin/contracts/SafeMath.sol";

import "./FudToken.sol";
import "./WinToken.sol";

contract AirVault {
    using SafeMath for uint256;

    IERC20 public immutable fudToken;
    WinToken public immutable winToken;

    mapping(address => uint256) public lockedBalance;

    event Deposited(
        address indexed account,
        uint256 amount,
        uint256 newBalance
    );
    event Withdrawn(
        address indexed account,
        uint256 amount,
        uint256 newBalance
    );

    constructor(address _fudToken, address _winToken) {
        fudToken = IERC20(_fudToken);
        winToken = WinToken(_winToken);
    }

    function deposit(uint256 amount) public returns (bool) {
        require(
            fudToken.transferFrom(msg.sender, address(this), amount),
            "transfer failed"
        );
        lockedBalance[msg.sender] = lockedBalance[msg.sender].add(amount);
        emit Deposited(msg.sender, amount, lockedBalance[msg.sender]);
        return true;
    }

    function withdraw(uint256 amount) public returns (bool) {
        require(lockedBalance[msg.sender] >= amount, "insufficient balance");
        require(fudToken.transfer(msg.sender, amount), "trasfer failed");
        lockedBalance[msg.sender] = lockedBalance[msg.sender].sub(amount);
        emit Withdrawn(msg.sender, amount, lockedBalance[msg.sender]);
        return true;
    }

    function lockedBalanceOf(address account) public view returns (uint256) {
        return lockedBalance[account];
    }
}
