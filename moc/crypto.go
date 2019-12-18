package moc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"math/big"

	"github.com/lokks307/pkcs8"
	"software.sslmate.com/src/go-pkcs12"
)

// MOdern Cryptography library for TethysCore

func GetCertificate(dataB64 string) (*x509.Certificate, error) {
	data := ParseDataToDer(dataB64)

	if data == nil {
		return nil, errors.New("Can't decode data")
	}

	var cert *x509.Certificate
	cert, parseErr := x509.ParseCertificate(data)

	if parseErr != nil {
		return nil, parseErr
	}

	return cert, nil
}

func DecodePFX(der []byte, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	privKey, cert, err := pkcs12.Decode(der, password)
	if err != nil {
		return nil, nil, err
	}

	return privKey, cert, nil
}

func DecodePFXB64(pfxDataB64 string, password string) (privateKey interface{}, certificate *x509.Certificate, err error) {
	res := ParseDataToDer(pfxDataB64)

	return DecodePFX(res, password)
}

func GetPrivateKey(dataB64, password string) (key interface{}, err error) {
	derBytes := ParseDataToDer(dataB64, password)

	// cannot check pkcs1 or pkcs8 type in der format. so try pkcs8 first and try pkcs1 again
	privKey, _, err := pkcs8.ParsePrivateKey(derBytes, []byte(password))
	if err != nil {
		pkcs1Key, parseErr := x509.ParsePKCS1PrivateKey(derBytes)

		if parseErr != nil {
			return nil, parseErr
		}

		return pkcs1Key, nil
	}

	return privKey, nil
}

func GetPublicKey(dataB64 string) (public interface{}, err error) {
	cert, parseErr := GetCertificate(dataB64)

	if parseErr != nil {
		return nil, parseErr
	}

	return cert.PublicKey, nil
}

func DoSign(msg []byte, key interface{}) ([]byte, error) {
	rng := rand.Reader
	var signature []byte
	var err error
	hashed := sha256.Sum256(msg)
	err = nil

	switch key.(type) {
	case *rsa.PrivateKey:
		signature, err = rsa.SignPKCS1v15(rng, key.(*rsa.PrivateKey), crypto.SHA256, hashed[:])
	case *ecdsa.PrivateKey:
		var r *big.Int
		var s *big.Int
		r, s, err = ecdsa.Sign(rng, key.(*ecdsa.PrivateKey), hashed[:])
		signature = r.Bytes()
		signature = append(signature, s.Bytes()...)
	case ed25519.PrivateKey:
		signature = ed25519.Sign(key.(ed25519.PrivateKey), msg)
	default:
		signature = nil
		err = errors.New("Unsupported type of crypto method")
	}

	return signature, err
}

func Sign(msg []byte, skeyPem, pass string) ([]byte, error) {
	privateKey, err := GetPrivateKey(skeyPem, pass)

	if err != nil {
		return nil, err
	}

	return DoSign(msg, privateKey)
}

func Verify(msg, sigBytes []byte, certPem string) bool {
	publicKey, err := GetPublicKey(certPem)

	if err != nil {
		return false
	}

	var result bool

	switch publicKey.(type) {
	case *rsa.PublicKey:
		hashed := sha256.Sum256(msg)
		err := rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], sigBytes)
		if err != nil {
			result = false
		} else {
			result = true
		}
	case *ecdsa.PublicKey:
		halfSigLen := len(sigBytes) / 2
		r := new(big.Int)
		r.SetBytes(sigBytes[:halfSigLen])

		s := new(big.Int)
		s.SetBytes(sigBytes[halfSigLen:])

		result = ecdsa.Verify(publicKey.(*ecdsa.PublicKey), msg, r, s)

	case ed25519.PublicKey:
		result = ed25519.Verify(publicKey.(ed25519.PublicKey), msg, sigBytes)
	default:
		result = false
	}

	return result
}