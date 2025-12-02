# go-multiotpwork

This package is for work with MultiOTP(https://github.com/multiOTP/multiotp/releases) cli application.

<h2>Functions</h2>

<h3>GetMultiOTPTokenURL</h3>

```
func GetMultiOTPTokenURL(user string, multiOTPBinPath string, descrString string) ([]byte, error)
```

Get MultiOTP user's totpURL.

* descrString - is for Token description(prefix), may be empty string
* user - user's login
* multiOTPBinPath - full path to multiotp binary

Due to multiotp console tools throw Exit codes every time need to check err.ExitCode, because err will be always set
* exit status 17: is success for '-urllink' cmd

Result example -  
```
otpauth://totp/<descrString, default=multiOTP>:<NAME>%20<SURNANME>?secret=<BASE32 SEED>&digits=6&period=30
```

<h3>delMultiOTPUser</h3>

```
func delMultiOTPUser(multiOTPBinPath string, user string) error
```

Delete MultiOTP User

If user doesn't exist - returns noting(not error)

* user - user's login
* multiOTPBinPath - full path to multiotp binary

Due to multiotp console tools throw Exit codes every time need to check err.ExitCode, because err will be always set.
* 12 INFO: User successfully deleted: is success for '-delete user' cmd
* 19 INFO: Requested operation successfully done: is success for '-delete user' cmd
* 21 ERROR: User doesn't exist: not error

<h3>resyncMultiOTPUsers</h3>

```
func resyncMultiOTPUsers(multiOTPBinPath string) error
```

Resync MultiOTP LDAP Users

Due to multiotp console tools throw Exit codes every time need to check err.ExitCode, because err will be always set.
* 19 INFO: Requested operation successfully done: is success for '-delete user' cmd

<h3>ReissueMultiOTPQR</h3>

```
func ReissueMultiOTPQR(multiOTPBinPath string, user string) error
```

Reissue MultiOTP QR

* user - user's login
* multiOTPBinPath - full path to multiotp binary

- First del user from MultiOTP db
- Second resync MultiOTP db to get same user back with new QR generated(May take some time to resync(depend of users number))

<h3>GenerateMultiOTPQRPng</h3>

Generate PNG QR

```
func GenerateMultiOTPQRPNG(multiOTPBinPath string, user string, qrCodesPath string) error
```

MutliOTP CLI command: 
```
multiotp -qrcode <USER> <FULL PATH TO OUTPUT PNG FILE>
```

* user - user's login
* multiOTPBinPath - full path to multiotp binary
* qrCodesPath - dir to save qr png-files

Due to multiotp console tools throw Exit codes every time need to check err.ExitCode, because err will be always set.
* 16 INFO: QRcode successfully created