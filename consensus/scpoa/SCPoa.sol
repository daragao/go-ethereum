pragma solidity ^0.4.23;

interface SCPoa {
  function getSigner(address _signerAddr) view external returns (bool);
}


/*
 * Signers was created just as an example of an implementation of the SCPoa interface
 */
contract Signers {
  mapping (address => bool) public signers;

  constructor(address[] _signers) public {
    for(uint i=0; i<_signers.length; i++) {
      signers[_signers[i]] = true;
    }
  }

  function getSigner(address _signerAddr) view external returns (bool) {
    return signers[_signerAddr];
  }
}
