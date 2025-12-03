package multiotp

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

/*
multiotp -qrcode user png_file_name.png
multiotp -update-pin user pin
multiotp -remove-token user

# reissue token
1) multiotp -delete user
2) multiotp -ldap-users-sync

# get totpURL
multiotp -urllink user
# otpauth://totp/multiOTP:<NAME>%20<SURNANME>?secret=<BASE32 SEED>&digits=6&period=30
*/

// Get MultiOTP user's totpURL
// descrString is for Token description(prefix), may be empty string
// Result example# otpauth://totp/<descrString, default=multiOTP>:<NAME>%20<SURNANME>?secret=<BASE32 SEED>&digits=6&period=30
func GetMultiOTPTokenURL(user string, multiOTPBinPath string, descrString string) ([]byte, error) {
	// define command to get TOTP URL for user
	cmd := exec.Command(multiOTPBinPath, "-urllink", user)
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// exit status 17: is success for '-urllink' cmd
	out, err := cmd.Output()

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			// 17 INFO: UrlLink successfully created
			// 21 ERROR: User doesn't exist
			switch {
			case err.ExitCode() == 21:
				return nil, fmt.Errorf("%s doesn't exist", user)
			case err.ExitCode() != 17:
				return nil, err
			}
		}
		return nil, fmt.Errorf("multiotp exec err: \n\t%v", err)
	}

	// check output is what expected
	patternOTPAuth := regexp.MustCompile(`^otpauth:`)
	if !patternOTPAuth.Match(out) {
		return nil, fmt.Errorf("mutliotp command doesn't match '^otpauth://', output:\n\t\t%s", out)
	}

	// if no descrString assigned, return original multiOTP url
	if len(descrString) == 0 {
		return out, nil
	}

	patternDescr := regexp.MustCompile(`^otpauth:\/\/\w+\/(\w+):`)
	subStringToReplace := patternDescr.FindStringSubmatch(string(out))

	// if no subString found(index 1 of FindStringSubmatch result), return original multiOTP url
	if len(subStringToReplace) < 2 {
		return out, nil
	}

	// making replace with descrString
	result := strings.Replace(string(out), subStringToReplace[1], descrString, 1)
	if len(result) == 0 {
		return nil, fmt.Errorf("mutliotp command doesn't match '^otpauth://', output:\n\t\t%s", out)
	}

	return []byte(result), nil
}

// Delete MultiOTP User
// If user doesn't exist - returns noting(not error)
func DelMultiOTPUser(multiOTPBinPath string, user string) error {
	// define command to delete user
	cmd := exec.Command(multiOTPBinPath, "-delete", user)
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// 12 INFO: User successfully deleted: is success for '-delete user' cmd
	// OR
	// 19 INFO: Requested operation successfully done: is success for '-delete user' cmd
	// 21 ERROR: User doesn't exist: not error
	_, err := cmd.Output()

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			switch {
			case err.ExitCode() == 21:
				return nil
			case err.ExitCode() != 12 && err.ExitCode() != 19:
				return err
			}
		}
		return fmt.Errorf("multiotp exec err: \n\t%v", err)
	}

	return nil
}

// Resync MultiOTP Users
func ResyncMultiOTPUsers(multiOTPBinPath string) error {
	// define command to delete user
	cmd := exec.Command(multiOTPBinPath, "-ldap-users-sync")
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// 19 INFO: Requested operation successfully done: is success for '-delete user' cmd
	_, err := cmd.Output()

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			if err.ExitCode() != 19 {
				return err
			}
		}
		return fmt.Errorf("multiotp exec err: \n\t%v", err)
	}

	return nil
}

// Reissue MultiOTP QR
func ReissueMultiOTPQR(multiOTPBinPath string, user string) error {
	// first del user from MultiOTP db
	err := DelMultiOTPUser(multiOTPBinPath, user)
	if err != nil {
		return fmt.Errorf("reissue qr: failed to del user:\n\t%v", err)
	}

	// second resync MultiOTP db to get same user back with new QR generated
	// may take some time to resync(depend of users number)
	err = ResyncMultiOTPUsers(multiOTPBinPath)
	if err != nil {
		return fmt.Errorf("reissue qr: failed to resync users:\n\t%v", err)
	}

	return nil
}

// Generate PNG QR
// multiotp -qrcode <USER> <FULL PATH TO OUTPUT PNG FILE>
func GenerateMultiOTPQRPng(multiOTPBinPath string, user string, qrCodesPath string) error {
	// form png file full path
	qrFullPath := fmt.Sprintf("%s/%s.png", qrCodesPath, user)

	cmd := exec.Command(multiOTPBinPath, "-qrcode", user, qrFullPath)
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// 16 INFO: QRcode successfully created
	_, err := cmd.Output()

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			if err.ExitCode() != 16 {
				if err.ExitCode() == 21 {
					return fmt.Errorf("21 ERROR: User doesn't exist")
				}
				return fmt.Errorf("unknown err code: %v", err)
			}
		}
		return fmt.Errorf("multiotp exec err: \n\t%v", err)
	}

	return nil
}
