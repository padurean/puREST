package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

var pasetoV2 = paseto.NewV2()
var publicKey ed25519.PublicKey
var privateKey ed25519.PrivateKey

// GenerateOrLoadKeys ...
func GenerateOrLoadKeys() error {
	publicKeyFileName := "public_key"
	_, errPublic := os.Stat(publicKeyFileName)
	privateKeyFileName := "private_key"
	_, errPrivate := os.Stat(privateKeyFileName)

	bothExist := !os.IsNotExist(errPublic) && !os.IsNotExist(errPrivate)
	if !bothExist {
		var err error
		publicKey, privateKey, err = ed25519.GenerateKey(nil)
		if err != nil {
			return fmt.Errorf("error generating public and private keys: %v", err)
		}
		if err := writeKeyToFile([]byte(publicKey), publicKeyFileName); err != nil {
			return fmt.Errorf("error writing public key to file %s: %v", publicKeyFileName, err)
		}
		if err := writeKeyToFile([]byte(privateKey), privateKeyFileName); err != nil {
			return fmt.Errorf("error writing private key to file %s: %v", privateKeyFileName, err)
		}
		return nil
	}

	publicKeyBytes, err := readKeyFromFile(publicKeyFileName)
	if err != nil {
		return fmt.Errorf("error loading public key from file %s: %v", privateKeyFileName, err)
	}
	privateKeyBytes, err := readKeyFromFile(privateKeyFileName)
	if err != nil {
		return fmt.Errorf("error loading private key from file %s: %v", privateKeyFileName, err)
	}
	publicKey = ed25519.PublicKey(publicKeyBytes)
	privateKey = ed25519.PrivateKey(privateKeyBytes)

	return nil
}

func writeKeyToFile(key []byte, fileName string) error {
	keyHex := make([]byte, hex.EncodedLen(len(key)))
	hex.Encode(keyHex, key)
	return ioutil.WriteFile(fileName, keyHex, 0644)
}

func readKeyFromFile(fileName string) ([]byte, error) {
	keyBytesRead, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading from file %s: %v", fileName, err)
	}
	keyBytes := make([]byte, hex.DecodedLen(len(keyBytesRead)))
	_, err = hex.Decode(keyBytes, keyBytesRead)
	if err != nil {
		return nil, fmt.Errorf("error hex decoding key: %v", err)
	}
	return keyBytes, nil
}

// GenerateToken ...
func GenerateToken(userID int64, role Role) (string, time.Time, error) {
	expiration := time.Now().Add(24 * time.Hour)
	jsonToken := paseto.JSONToken{
		Expiration: expiration,
		Subject:    strconv.FormatInt(userID, 10),
	}
	jsonToken.Set("role", fmt.Sprintf("%d", role))
	footer := "puREST"
	token, err := pasetoV2.Sign(privateKey, jsonToken, footer)
	return token, expiration, err
}

// JSONToken ...
type JSONToken struct {
	UserID     int64
	Role       Role
	Expiration time.Time
}

// VerifyToken ...
func VerifyToken(token string) (*JSONToken, error) {
	var jsonToken paseto.JSONToken
	var footer string
	if err := pasetoV2.Verify(token, publicKey, &jsonToken, &footer); err != nil {
		return nil, err
	}
	if err := jsonToken.Validate(); err != nil {
		return nil, err
	}
	userID, err := strconv.ParseInt(jsonToken.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing token subject as user ID (i.e. int64): %v", err)
	}
	roleStr := jsonToken.Get("role")
	if roleStr == "" {
		return nil, errors.New("user role is missing from token")
	}
	role, err := ParseRole(roleStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing user role from token: %v", err)
	}
	return &JSONToken{
		UserID:     userID,
		Role:       role,
		Expiration: jsonToken.Expiration,
	}, nil
}
