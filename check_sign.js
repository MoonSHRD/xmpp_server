// const EthCrypto = require('eth-crypto');
let nacl = require('tweetnacl');

const message   = B64ToUint8Array(process.argv[2]);
const signature      = B64ToUint8Array(process.argv[3]);
const pubKey      = B64ToUint8Array(process.argv[4]);

// const message   = B64ToUint8Array('/Jlx9qZFw7zthdxOi7aVO2uCnisLW+K9iMOJuH3pTgw=');
// const signature      = B64ToUint8Array('jt0Vacot6Wnpca/US9Le1YO6ZKXlGgA2dvBJnCSAbz8Yo4YiCxDWzfkr6aYbfvYDWJ9htUXVgzH+KJgJZhYSCQ==');
// const pubKey      = B64ToUint8Array('yupD8L1qCDXXlSsRWEFv+rheBCqZnEGR5dmFaSKoX5s=');

// const messageHash = EthCrypto.hash.keccak256(message);
// const signer = EthCrypto.recover(
//     signature,
//     messageHash
// );

// console.log(message,signature,pubKey);

const signer = nacl.sign.detached.verify(message,signature,pubKey).toString();

// console.log(signer);

process.stdout.write(signer);

function Uint8ArrayToB64(bytes) {
    return Buffer.from(bytes.buffer, bytes.byteOffset, bytes.byteLength).toString('base64');
}

function B64ToUint8Array(s) {
    return Buffer.from(s, 'base64')
}