//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

contract BasicControlContract {
    address public admin;

    event AdminRegistered(address indexed _from);

    constructor() {
        admin = msg.sender;
        emit AdminRegistered(admin);
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin");
        _;
    }

    function publicFunction() external {}

    function privateFunction1() external onlyAdmin {}

    function privateFunction2() external onlyAdmin {}
}
