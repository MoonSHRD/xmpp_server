package util

import (
    "encoding/base64"
    "encoding/hex"
    "golang.org/x/crypto/ripemd160"
    "os/exec"
    "strings"
    
    //"github.com/ethereum/go-ethereum/crypto/sha3"
)

type SignatureType uint8

const (
    SignatureType_EIP712 SignatureType = 0
    SignatureType_GETH   SignatureType = 1
    SignatureType_TREZOR SignatureType = 2
)

func CheckSign(msg,sign,pubKey string) (bool,error) {
    suc,err:=exec.Command("node","check_sign.js",msg,sign,pubKey).Output()
    suc_s:=string(suc)=="true"
    return suc_s,err
}

func AddrFromPrub(pub string) (string,error) {
    pubKey, err := base64.StdEncoding.DecodeString(pub)
    if err != nil {
        return "",err
    }
    hasher := ripemd160.New()
    hasher.Write(pubKey)
    hashBytes := hasher.Sum(nil)
    address:="0x"+hex.EncodeToString(hashBytes)
    return strings.ToLower(address),nil
}