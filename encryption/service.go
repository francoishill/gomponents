package encryption

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"

	passlib "gopkg.in/hlandau/passlib.v1"
)

type Service interface {
	NewRandomPassword() (string, error)
	HashPassword(password string) (string, error)
	VerifyPassword(password, hashedPassword string) error
}

func DefaultService() *defaultService { return &defaultService{} }

type defaultService struct{}

func (d *defaultService) NewRandomPassword() (string, error) {
	return d.generateRandStr(29, alphanumDictionary) + d.generateRandStr(1, upperDictionary) + d.generateRandStr(1, lowerDictionary) + d.generateRandStr(1, numberDictionary), nil
}

func (*defaultService) HashPassword(password string) (string, error) {
	//this uses default security policy and auto-rotation of hashes, read more at: https://github.com/hlandau/passlib/tree/v1.0.9
	hash, err := passlib.Hash(password)
	if err != nil {
		return "", errors.Errorf("Unable hash password, error: %s", err.Error())
	}

	return hash, nil
}

func (*defaultService) VerifyPassword(password, hashedPassword string) error {
	err := passlib.VerifyNoUpgrade(password, hashedPassword)
	if err != nil {
		return errors.Wrapf(err, "Password verification failed")
	}
	return nil
}

func (*defaultService) generateSecureRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (d *defaultService) generateMagicToken(s int) (string, error) {
	b, err := d.generateSecureRandomBytes(s)
	if err != nil {
		return "", err
	}

	hexStr := hex.EncodeToString(b)
	return hexStr, nil
}

const (
	alphanumDictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	alphaDictionary    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	upperDictionary    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerDictionary    = "abcdefghijklmnopqrstuvwxyz"
	numberDictionary   = "0123456789"
)

func (*defaultService) generateRandStr(size int, dictionary string) string {
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
