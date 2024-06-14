package mt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"math"
	"math/big"
	"strconv"

	"github.com/lokks307/go-util/bytesbuilder"
	"lukechampine.com/blake3"
)

func getZeroPadded(key string) []byte {
	bkey := []byte(key)
	pkey := make([]byte, 16)
	copy(pkey[0:], bkey[0:16])
	return pkey
}

var EGHIS_IV = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func EncAESCBC(ptext string, key []byte, iv []byte) (string, error) {
	block, keyerr := aes.NewCipher(key)
	if keyerr != nil {
		return "", keyerr
	}

	lenText := len(ptext)
	numBlock := int(math.Ceil(float64(lenText) / 16.0))
	pad := byte(0)

	if lenText%16 == 0 {
		numBlock++
	}

	lenCText := numBlock * 16
	pad = byte(lenCText - lenText)

	var b bytes.Buffer
	b.WriteString(ptext)
	b.Write(bytes.Repeat([]byte{byte(pad)}, int(pad)))

	src := b.Bytes()
	dst := make([]byte, len(src))

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(dst, src)

	return hex.EncodeToString(dst), nil
}

func DecAESCBC(hexctext string, key []byte, iv []byte) (string, error) {
	ctext, decerr := hex.DecodeString(hexctext) // 128bit
	if decerr != nil {
		return "", decerr
	}

	var block cipher.Block

	block, keyerr := aes.NewCipher(key)
	if keyerr != nil {
		return "", keyerr
	}

	lenPtext := len(ctext)
	if lenPtext == 0 {
		return "", errors.New("ciphertext length 0")
	}

	ptext := make([]byte, lenPtext)

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ptext, ctext)

	padLen := int(ptext[lenPtext-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen < 16 && padLen > 0 && bytes.HasSuffix(ptext, ref) {
		return string(ptext[:lenPtext-padLen]), nil
	}

	return "", errors.New("padding error")
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 possibilities
	letterIdxBits = 6                                                      // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
)

func SecureRandomAlphaString(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)

	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = SecureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result)
}

// SecureRandomBytes returns the requested number of bytes using crypto/rand
func SecureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

func SecureRandomHex(length int) string { // in byte
	return hex.EncodeToString(SecureRandomBytes(length))
}

func UnpaddingPKCS7Data(plain []byte) []byte {
	lenPlain := len(plain)
	if lenPlain == 0 {
		return []byte{}
	}

	padLen := int(plain[lenPlain-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen < 16 && padLen > 0 && bytes.HasSuffix(plain, ref) {
		return plain[:lenPlain-padLen]
	}

	return []byte{}
}

func IsValidPubilcKeyFormat(pkStr string) bool {
	pkRaw, err := base64.StdEncoding.DecodeString(pkStr)
	if err != nil {
		return false
	}

	if _, err := x509.ParsePKIXPublicKey(pkRaw); err != nil {
		return false
	}

	return true
}

func DecryptDasom(ciphertext, key string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	for i := 0; i < len(ciphertextBytes); i += aes.BlockSize {
		cipherBlock.Decrypt(ciphertextBytes[i:i+aes.BlockSize], ciphertextBytes[i:i+aes.BlockSize])
	}

	plaintext := string(ciphertextBytes)
	return plaintext, nil
}

func EncryptDasom(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)

	ecb := newECBEncrypter(block)

	ecb.CryptBlocks(plaintextBytes, plaintextBytes)

	ciphertext := base64.StdEncoding.EncodeToString(plaintextBytes)
	return ciphertext, nil
}

func newECBEncrypter(cipherBlock cipher.Block) cipher.BlockMode {
	return &ecbEncrypter{cipherBlock: cipherBlock}
}

type ecbEncrypter struct {
	cipherBlock cipher.Block
}

func (e *ecbEncrypter) BlockSize() int {
	return e.cipherBlock.BlockSize()
}

func (e *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%e.cipherBlock.BlockSize() != 0 {
		//logrus.Error("invalid cipher block size")
		return
	}
	for i := 0; i < len(src); i += e.cipherBlock.BlockSize() {
		e.cipherBlock.Encrypt(dst[i:i+e.cipherBlock.BlockSize()], src[i:i+e.cipherBlock.BlockSize()])
	}
}

func Sha256(b ...interface{}) string {
	bb := bytesbuilder.NewBuilder()
	bb.Append(b...)
	hash := sha256.Sum256(bb.GetBytes())
	return hex.EncodeToString(hash[:])
}

func Blake3(b ...interface{}) string {
	bb := bytesbuilder.NewBuilder()
	bb.Append(b...)
	hash := blake3.Sum256(bb.GetBytes())
	return hex.EncodeToString(hash[:])
}

func GetRandInt64(max int64) int64 {
	r, e := rand.Int(rand.Reader, big.NewInt(max))
	if e != nil || r == nil {
		return 0
	}

	return r.Int64()
}

const (
	LETTERS = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func PasswordGenerator(plen int) string {
	pwd := ""
	lenLetter := len(LETTERS)

	rs := make([]byte, plen)
	_, _ = rand.Read(rs)

	for i := 0; i < plen; i++ {
		pos := int(rs[i]) % lenLetter
		pwd += string(LETTERS[pos])
	}

	return pwd
}

func GenRandomHex(n int) string {
	b := make([]byte, n/2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func HashPassword(pw, salt string) string {
	comb := []byte(pw + salt)
	sum := sha256.Sum256(comb)

	return hex.EncodeToString(sum[:])
}

func VerifyPassword(src, target, salt string) bool {
	return target == HashPassword(src, salt)
}

func MakeUID(email string, time int64) string {
	comb := []byte(email + strconv.FormatInt(time, 10))
	uid := sha256.Sum256(comb)

	return hex.EncodeToString(uid[:])
}
