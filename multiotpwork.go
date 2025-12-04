package multiotp

import (
	"fmt"
	"os"
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

// Check file path exists
// func checkPath(path string) (bool, error) {
// 	os.Stat()
// }

// checking path exists
func checkPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// is dir writable
func isWritable(path string) (bool, error) {
	// check if path exists
	if !checkPath(path) {
		return false, fmt.Errorf("%s doesn't exists", path)
	}

	tmpFile := "tmpfile"

	file, err := os.CreateTemp(path, tmpFile)
	if err != nil {
		return false, err
	}

	defer os.Remove(file.Name())
	defer file.Close()

	return true, nil
}

// check multiotp bin is executable
func checkMultiOTPBin(multiOTPBinPath string) (bool, error) {
	// check if multiOTPBinPath exists
	if !checkPath(multiOTPBinPath) {
		return false, fmt.Errorf("%s: path doesn't exist", multiOTPBinPath)
	}

	// trying to run bin
	cmd := exec.Command(multiOTPBinPath, "-version")

	out, err := cmd.Output()
	if err, ok := err.(*exec.ExitError); ok {
		if err.ExitCode() == 19 {
			return true, nil
		}
	}

	return false, fmt.Errorf("check errors: %v ; %s", err, out)
}

// Get MultiOTP user's totpURL
// descrString is for Token description(prefix), may be empty string
// Result example# otpauth://totp/<descrString, default=multiOTP>:<NAME>%20<SURNANME>?secret=<BASE32 SEED>&digits=6&period=30
func GetMultiOTPTokenURL(user string, multiOTPBinPath string, descrString string) ([]byte, error) {
	// check path exists
	if ok, err := checkMultiOTPBin(multiOTPBinPath); !ok {
		return nil, fmt.Errorf("multiotpBinPath doesn't exist or isn't executable: %v", err)
	}

	// define command to get TOTP URL for user
	cmd := exec.Command(multiOTPBinPath, "-urllink", user)
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// exit status 17: is success for '-urllink' cmd
	out, err := cmd.Output()

	if err, ok := err.(*exec.ExitError); ok {
		if err.ExitCode() != 17 {
			if err.ExitCode() == 21 {
				return nil, fmt.Errorf("21 ERROR: User doesn't exist")
			}
			return nil, fmt.Errorf("unknown err code: %v ; %s", err, out)
		}
	}

	// check output is what expected
	patternOTPAuth := regexp.MustCompile(`^otpauth:`)
	if !patternOTPAuth.Match(out) {
		return nil, fmt.Errorf("mutliotp command doesn't match '^otpauth://', output: %s", out)
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
		return nil, fmt.Errorf("mutliotp command doesn't match '^otpauth://', output: %s", out)
	}

	return []byte(result), nil
}

// Delete MultiOTP User
// If user doesn't exist - returns noting(not error)
func DelMultiOTPUser(multiOTPBinPath string, user string) error {
	// check path exists
	if ok, err := checkMultiOTPBin(multiOTPBinPath); !ok {
		return fmt.Errorf("multiotpBinPath doesn't exist or isn't executable: %v", err)
	}

	// define command to delete user
	cmd := exec.Command(multiOTPBinPath, "-delete", user)
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// 12 INFO: User successfully deleted: is success for '-delete user' cmd
	// OR
	// 19 INFO: Requested operation successfully done: is success for '-delete user' cmd
	// 21 ERROR: User doesn't exist: not error
	out, err := cmd.Output()

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			switch {
			case err.ExitCode() == 21:
				return fmt.Errorf("21 ERROR: User doesn't exist")
			// 12 INFO: User successfully deleted
			case err.ExitCode() == 12:
				return nil
			// 19 INFO: Requested operation successfully done
			case err.ExitCode() == 19:
				return nil
			default:
				return fmt.Errorf("unknown err code: %v ; %s", err, out)
			}
		}
	}

	return nil
}

// Resync MultiOTP Users
func ResyncMultiOTPUsers(multiOTPBinPath string) error {
	// check path exists
	if ok, err := checkMultiOTPBin(multiOTPBinPath); !ok {
		return fmt.Errorf("multiotpBinPath doesn't exist or isn't executable: %v", err)
	}

	// define command to delete user
	cmd := exec.Command(multiOTPBinPath, "-ldap-users-sync")
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// 19 INFO: Requested operation successfully done: is success for '-delete user' cmd
	out, err := cmd.Output()

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			if err.ExitCode() != 19 {
				return fmt.Errorf("unknown err code: %v ; %s", err, out)
			}
		}
	}

	return nil
}

// Reissue MultiOTP QR
func ReissueMultiOTPQR(multiOTPBinPath string, user string) error {
	// first del user from MultiOTP db
	err := DelMultiOTPUser(multiOTPBinPath, user)
	if err != nil {
		return fmt.Errorf("reissue qr: failed to del user: %v", err)
	}

	// second resync MultiOTP db to get same user back with new QR generated
	// may take some time to resync(depend of users number)
	err = ResyncMultiOTPUsers(multiOTPBinPath)
	if err != nil {
		return fmt.Errorf("reissue qr: failed to resync users: %v", err)
	}

	return nil
}

// Generate PNG QR
// multiotp -qrcode <USER> <FULL PATH TO OUTPUT PNG FILE>
func GenerateMultiOTPQRPng(multiOTPBinPath string, user string, qrCodesPath string) error {
	// check path exists
	if ok, err := checkMultiOTPBin(multiOTPBinPath); !ok {
		return fmt.Errorf("multiotpBinPath doesn't exist or isn't executable: %v", err)
	}

	// check qrCodesPath is writable
	if _, err := isWritable(qrCodesPath); err != nil {
		return fmt.Errorf("%s is not writable or does't exist: %v", qrCodesPath, err)
	}

	// form png file full path
	qrFullPath := fmt.Sprintf("%s/%s.png", qrCodesPath, user)

	cmd := exec.Command(multiOTPBinPath, "-qrcode", user, qrFullPath)
	// due to multiotp console tools throw Exit codes every time
	// need to check err.ExitCode, because err will be always
	// 16 INFO: QRcode successfully created
	out, err := cmd.Output()

	if err, ok := err.(*exec.ExitError); ok {
		if err.ExitCode() != 16 {
			if err.ExitCode() == 21 {
				return fmt.Errorf("21 ERROR: User doesn't exist")
			}
			return fmt.Errorf("unknown err code: %v ; %s", err, out)
		}
	}

	return nil
}
