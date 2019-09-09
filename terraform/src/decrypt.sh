#!/usr/bin/expect -f

# This script use keybase to decrypt secret.
# Usage: ./decrypt.sh $passphrase $secret_file
#   passphrase: keybase passphrase
#   secret_file: path to the file containing the encrypting secret
#
# You can store the output of the decrypt file in a variable that you can
# then store in a secret management, for example:
# SECRET_ACCESS_KEY=$(/home/builder/src/decrypt.sh ${KEYBASE_PASSPHRASE} ${MY_ENCRYPTED_SECRET_FILE})

set PASSPHRASE [lindex $argv 0];
set SECRET_FILE [lindex $argv 1];
set timeout -1
spawn keybase pgp decrypt -i ${SECRET_FILE}
expect {
	"Reason: PGP Decryption:" {
		send -- "${PASSPHRASE}\r"
		expect eof
	}
}
