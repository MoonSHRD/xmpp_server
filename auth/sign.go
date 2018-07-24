package auth

import (
    "github.com/ortuman/jackal/stream"
    "encoding/base64"
    "github.com/ortuman/jackal/xml"
    "github.com/ortuman/jackal/util"
    "fmt"
    "strings"
    "github.com/ortuman/jackal/model"
    "github.com/ethereum/go-ethereum/crypto"
)

type signState int

const (
    startSignState signState = iota
    challengedSignState
)

type Sign struct {
   stm           stream.C2S
   state         signState
   username      string
   authenticated bool
}

type signParameters struct {
    username  string
    realm     string
    nonce     []byte
    cnonce    string
    nc        string
    qop       string
    servType  string
    digestURI string
    response  string
    charset   string
    authID    string
    signature []byte
    pubKey    []byte
}

func (r *signParameters) setParameter(p string) {
    key, val := util.SplitKeyAndValue(p, '=')
    
    // strip value double quotes
    val = strings.TrimPrefix(val, `"`)
    val = strings.TrimSuffix(val, `"`)
    
    switch key {
    case "username":
        r.username = val
    case "realm":
        r.realm = val
    case "nonce":
        r.nonce = []byte(val)
    case "cnonce":
        r.cnonce = val
    case "nc":
        r.nc = val
    case "qop":
        r.qop = val
    case "serv-type":
        r.servType = val
    case "digest-uri":
        r.digestURI = val
    case "response":
        r.response = val
    case "charset":
        r.charset = val
    case "authzid":
        r.authID = val
    case "signature":
        r.signature = []byte(val)
    case "pubKey":
        r.pubKey = []byte(val)
    }
}

func (sig *Sign) parseParameters(str string) *signParameters {
    params := &signParameters{}
    s := strings.Split(str, ",")
    for i := 0; i < len(s); i++ {
        params.setParameter(s[i])
    }
    return params
}

func NewSign(stm stream.C2S) *Sign {
    return &Sign{
        stm: stm,
        state: startSignState,
    }
}

// Mechanism returns authenticator mechanism name.
func (sig *Sign) Mechanism() string {
    return "SIGN"
}

// Username returns authenticated username in case
// authentication process has been completed.
func (sig *Sign) Username() string {
    return sig.username
}

// Authenticated returns whether or not user has been authenticated.
func (sig *Sign) Authenticated() bool {
    return sig.authenticated
}

// UsesChannelBinding returns whether or not plain authenticator
// requires channel binding bytes.
func (sig *Sign) UsesChannelBinding() bool {
    return false
}

func (sig *Sign) ProcessElement(elem xml.XElement) error {
    if sig.Authenticated() {
        return nil
    }
    switch elem.Name() {
    case "auth":
        switch sig.state {
        case startSignState:
            return sig.handleStart(elem)
        }
    case "response":
        switch sig.state {
        case challengedSignState:
            return sig.handleChallenged(elem)
        }
    }
    return ErrSASLNotAuthorized
}

func (sig *Sign) handleStart(elem xml.XElement) error {
    domain := sig.stm.Domain()
    nonce := base64.StdEncoding.EncodeToString(util.RandomBytes(32))
    chnge := fmt.Sprintf(`realm="%s",nonce="%s",qop="auth",charset=utf-8,algorithm=eth-sign`, domain, nonce)
    
    respElem := xml.NewElementNamespace("challenge", nonSaslNamespace)
    respElem.SetText(chnge)
    sig.stm.SendElement(respElem)
    
    sig.state = challengedSignState
    return nil
}

func (sig *Sign) handleChallenged(elem xml.XElement) error {
    if len(elem.Text()) == 0 {
        return ErrSASLMalformedRequest
    }
    params := sig.parseParameters(elem.Text())
    
    // validate realm
    //if params.realm != sig.stm.Domain() {
    //    return ErrSASLNotAuthorized
    //}
    // validate nc
    //if params.nc != "00000001" {
    //    return ErrSASLNotAuthorized
    //}
    // validate qop
    if params.qop != "auth" {
        return ErrSASLNotAuthorized
    }
    // validate serv-type
    if len(params.servType) > 0 && params.servType != "xmpp" {
        return ErrSASLNotAuthorized
    }
    // validate digest-uri
    //if !strings.HasPrefix(params.digestURI, "xmpp/") || params.digestURI[5:] != sig.stm.Domain() {
    //    return ErrSASLNotAuthorized
    //}
    
    //validate pub_key
    //crypto.UnmarshalPubkey(params.pubKey)
    
    //key,_:=crypto.GenerateKey()
    //key1:=string(crypto.FromECDSA(key))
    ////pri:="b27a276db9c01d272116f337ddd02b4aa7b2d5869ff5687e5929005196e480fc"
    //pub:="0x06ef2f0b4be72a8ecce6b2adcda1aad4c91fccf1fe8e1574e07446e47caf106234581ce02e0e328f7d450b648ef40a7f9a203c848893ca66ca0119403ab481e1"
    //addr:="0xfb951431c04241d6c82b5e0edfcd82ca592e6bab"
    //
    ////fmt.Print(fefe)
    //
    //pub_ec,err:=crypto.UnmarshalPubkey([]byte(pub))
    //if err !=nil {
    //    fmt.Print(err)
    //}
    //if crypto.PubkeyToAddress(*pub_ec).String()!=addr {
    //    return ErrSASLNotAuthorized
    //}
    
    //crypto.Ch
    
    pub_ec,err:=crypto.DecompressPubkey(params.pubKey)
    if err !=nil {
        fmt.Println(err)
    }
    
    pub_ec,err=crypto.UnmarshalPubkey(params.pubKey)
    if err !=nil {
        fmt.Println(err)
    }
    if crypto.PubkeyToAddress(*pub_ec).String()!=params.username {
        return ErrSASLNotAuthorized
    }
    
    //validate sign
    if !crypto.VerifySignature(params.pubKey,params.nonce,params.signature) {
        return ErrSASLNotAuthorized
    }
    
    //// validate user
    //user, err := storage.Instance().FetchUser(params.username)
    //if err != nil {
    //	return err
    //}
    //if user == nil {
    //	return ErrSASLNotAuthorized
    //}
    //jid:=jid2.JID{domain:"localhost"}
    user:=new(model.User)//{"govno","123"}
    user.Username=params.username
    user.Password=""
    
    ////validate response
    //clientResp := d.computeResponse(params, user, true)
    //if clientResp != params.response {
    //	return ErrSASLNotAuthorized
    //}
    
    // authenticated... compute and send server response
    sig.stm.SendElement(xml.NewElementNamespace("success", nonSaslNamespace))
    
    sig.username = user.Username
    sig.authenticated=true
    return nil
}

// Reset resets plain authenticator internal state.
func (sig *Sign) Reset() {
    sig.state = startSignState
    sig.username = ""
    sig.authenticated = false
}