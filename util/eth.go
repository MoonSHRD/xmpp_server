package util

import "os/exec"

func CheckSign(msg,sign string) (string,error) {
    //path,err:=exec.LookPath("./check_sign.js")
    //if err!=nil {
    //    return "",err
    //}
    suc,err:=exec.Command("node","check_sign.js",msg,sign).Output()
    //suc,err:=exec.Command("ls").Output()
    suc_s:=string(suc)
    return suc_s,err
}
