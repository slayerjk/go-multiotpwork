package main

import (
	"flag"
	"fmt"
	"os"

	multiotp "github.com/slayerjk/go-multiotpwork"
)

func main() {
	opt := flag.String("o", "none", "t=getTokenURL; d=delUser; r=resyncLdapUsers; rq=reissueUserQR; p=generatePNG; i=ldap-user-info")
	multiOTPBinPath := flag.String("m", "/usr/local/bin/multiotp/multiotp.php", "full path to multiotp binary")
	qrCodesPath := flag.String("q", "/etc/multiotp/qrcodes", "qr codes full path, needs for '-o p'")
	user := flag.String("u", "user", "user name")
	descrString := flag.String("ds", "TEST", "token description")
	flag.Parse()

	switch *opt {
	case "none":
		fmt.Println("no key for -o set")
		os.Exit(0)
	case "t":
		result, err := multiotp.GetMultiOTPTokenURL(*user, *multiOTPBinPath, *descrString)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("DONE, tokenURL: %s\n", string(result))
	case "d":
		err := multiotp.DelMultiOTPUser(*multiOTPBinPath, *user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("DONE deleting user")
	case "r":
		err := multiotp.ResyncMultiOTPUsers(*multiOTPBinPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("DONE resyncing LDAP users")
	case "rq":
		err := multiotp.ReissueMultiOTPQR(*multiOTPBinPath, *user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("DONE reissuing QR for user")
	case "p":
		err := multiotp.GenerateMultiOTPQRPng(*multiOTPBinPath, *user, *qrCodesPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("DONE generating png QR for user")
	case "i":
		out, err := multiotp.GetLdapUserInfo(*multiOTPBinPath, *user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("User LDAP Info:\n", out)
	default:
		fmt.Println("not valid value for -o")
	}
}
