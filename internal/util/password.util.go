package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Definisikan parameter untuk Argon2id.
// Anda bisa menyesuaikan ini sesuai kebutuhan keamanan dan performa server.
type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var params = &argonParams{
	memory:      64 * 1024, // 64 MB
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

// HashPassword mengenkripsi password menggunakan Argon2id.
func HashPassword(password string) (string, error) {
	// 1. Generate salt acak
	salt := make([]byte, params.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Hasilkan hash dari password
	hash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	// 3. Encode salt dan hash ke dalam format standar
	// Format: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	fullHash := fmt.Sprintf(format, argon2.Version, params.memory, params.iterations, params.parallelism, b64Salt, b64Hash)
	return fullHash, nil
}

// CheckPasswordHash membandingkan password plaintext dengan hash Argon2id.
func CheckPasswordHash(password, fullHash string) (bool, error) {
	// 1. Parse hash yang tersimpan
	parts := strings.Split(fullHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	// 2. Decode salt dan hash dari Base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// 3. Buat hash pembanding dari password yang diinput
	comparisonHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	// 4. Bandingkan hash (constant-time comparison untuk keamanan)
	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1 {
		return true, nil
	}
	return false, nil
}
