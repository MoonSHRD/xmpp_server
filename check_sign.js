const EthCrypto = require('eth-crypto');

const message   = process.argv[2];
const signature      = process.argv[3];

const messageHash = EthCrypto.hash.keccak256(message);
const signer = EthCrypto.recover(
    signature,
    messageHash
);

process.stdout.write(signer);