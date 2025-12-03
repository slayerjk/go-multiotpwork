package main

import (
	"flag"
	"fmt"
	"os"

	multiotp "github.com/slayerjk/go-multiotpwork"
)

func main() {
	opt := flag.String("o", "none", "t=getTokenURL; d=delUser; r=resyncLdapUsers; rq=reissueUsersQR; p=generatePNG")
	multiOTPBinPath := flag.String("m", "/usr/local/bin/multiotp/multiotp.php", "full path to multiotp binary")
	qrCodesPath := flag.String("q", "/etc/multiotp/qrcodes", "qr codes full path")
	user := flag.String("u", "user", "user to generate qr")
	descrString := flag.String("ds", "TEST", "token")
	flag.Parse()

	switch *opt {
	case "none":
		fmt.Fprint(os.Stdout, "no key for -o set\n")
		os.Exit(0)
	case "t":
		result, err := multiotp.GetMultiOTPTokenURL(*user, *multiOTPBinPath, *descrString)
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "DONE, tokenURL: %s", string(result))
	case "d":
		err := multiotp.DelMultiOTPUser(*multiOTPBinPath, *user)
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, "DONE deleting user")
	case "r":
		err := multiotp.ResyncMultiOTPUsers(*multiOTPBinPath)
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, "DONE resyncing LDAP users")
	case "rq":
		err := multiotp.ReissueMultiOTPQR(*multiOTPBinPath, *user)
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, "DONE reissuing QR for user")
	case "p":
		err := multiotp.GenerateMultiOTPQRPng(*multiOTPBinPath, *user, *qrCodesPath)
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, "DONE generating png QR for user")
	}
}
