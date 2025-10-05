package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Time        uint32
	Memory      uint32
	Parallelism uint8
	SaltLen     uint32
	KeyLen      uint32
}

var DefaultParams = Params{
	Time:        3,
	Memory:      64 * 1024,
	Parallelism: 1,
	SaltLen:     16,
	KeyLen:      32,
}

// If I understand correctly, the pepper is just appended? to the end of the SHA256 that we recieve from the
// client functions sort of like a signature so unless someone has the pepper, this encryption shouldn't be breakable? idk ~brtcrt
var pepper []byte

func init() {
	if v := os.Getenv("PASSWORD_PEPPER_BASE64"); v != "" { // omergotunusikim
		b, err := base64.StdEncoding.DecodeString(v)
		if err == nil && len(b) > 0 {
			pepper = b
			return
		}
	}
	pepper = nil
}
// We don't really need this anymore probably this was just to migrate the existing plain text passwords in the db. ~brtcrt
func Sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
// We recieve a SHA256 hashed password from the client.
// Then we hash again using argon2
// Then we encode to phc and we store that phc encoded string in the db
// So this is like 3-way hashing but not really (I don't count phc)
// The important thing is we never deal with plain text passwords ~brtcrt
func HashFromClientDigestHex(digestHex string, p Params) (string, error) {
	raw, err := hex.DecodeString(digestHex)
	if err != nil {
		return "", errors.New("invalid password digest")
	}
	material := append(raw, pepper...)

	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey(material, salt, p.Time, p.Memory, p.Parallelism, p.KeyLen)
	return phcEncode("argon2id", p, salt, hash), nil
}

// Seperate the phc encoded string and get the argon2 hash
// Argon2 the fucking SHA256 hash
// Compare the two hashes to see if they match ~brtcrt
func VerifyFromClientDigestHex(digestHex, encoded string) (bool, error) {
	alg, params, salt, want, err := phcDecode(encoded)
	if err != nil {
		return false, err
	}
	if alg != "argon2id" {
		return false, errors.New("unsupported alg")
	}
	raw, err := hex.DecodeString(digestHex)
	if err != nil {
		return false, errors.New("invalid password digest")
	}
	material := append(raw, pepper...)
	got := argon2.IDKey(material, salt, params.Time, params.Memory, params.Parallelism, uint32(len(want)))
	return subtle.ConstantTimeCompare(want, got) == 1, nil
}

func phcEncode(alg string, p Params, salt, hash []byte) string {
	return fmt.Sprintf("$%s$v=19$m=%d,t=%d,p=%d$%s$%s",
		alg, p.Memory, p.Time, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

// Decode the phc encoded password into various parts so we can cross-check them with the incoming password.
func phcDecode(s string) (alg string, p Params, salt, hash []byte, err error) {
	parts := strings.Split(s, "$")
	if len(parts) != 6 {
		return "", Params{}, nil, nil, errors.New("bad phc")
	}
	alg = parts[1]
	var v int
	if _, err = fmt.Sscanf(parts[2], "v=%d", &v); err != nil || v != 19 {
		return "", Params{}, nil, nil, errors.New("bad version")
	}
	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Time, &p.Parallelism); err != nil {
		return "", Params{}, nil, nil, errors.New("bad params")
	}
	if salt, err = base64.RawStdEncoding.DecodeString(parts[4]); err != nil {
		return "", Params{}, nil, nil, errors.New("bad salt")
	}
	if hash, err = base64.RawStdEncoding.DecodeString(parts[5]); err != nil {
		return "", Params{}, nil, nil, errors.New("bad hash")
	}
	return
}
