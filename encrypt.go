// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package madmin

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"sync"

	"github.com/secure-io/sio-go"
	"github.com/secure-io/sio-go/sioutil"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

// IsEncrypted reports whether data is encrypted.
func IsEncrypted(data []byte) bool {
	if len(data) <= 32 {
		return false
	}
	b := data[32]
	return b == pbkdf2AESGCM || b == argon2idAESGCM || b == argon2idChaCHa20Poly1305
}

// EncryptData encrypts the data with an unique key
// derived from password using the Argon2id PBKDF.
//
// The returned ciphertext data consists of:
//
//	salt | AEAD ID | nonce | encrypted data
//	 32      1         8      ~ len(data)
func EncryptData(password string, data []byte) ([]byte, error) {
	salt := sioutil.MustRandom(32)

	var (
		id     byte
		err    error
		stream *sio.Stream
	)
	if FIPSEnabled() {
		key := pbkdf2.Key([]byte(password), salt, pbkdf2Cost, 32, sha256.New)
		stream, err = sio.AES_256_GCM.Stream(key)
		if err != nil {
			return nil, err
		}
		id = pbkdf2AESGCM
	} else {
		argon2Mu.Lock()
		key := argon2.IDKey([]byte(password), salt, argon2idTime, argon2idMemory, argon2idThreads, 32)
		argon2Mu.Unlock()
		if sioutil.NativeAES() {
			stream, err = sio.AES_256_GCM.Stream(key)
			if err != nil {
				return nil, err
			}
			id = argon2idAESGCM
		} else {
			stream, err = sio.ChaCha20Poly1305.Stream(key)
			if err != nil {
				return nil, err
			}
			id = argon2idChaCHa20Poly1305
		}
	}

	nonce := sioutil.MustRandom(stream.NonceSize())

	// ciphertext = salt || AEAD ID | nonce | encrypted data
	cLen := int64(len(salt)+1+len(nonce)+len(data)) + stream.Overhead(int64(len(data)))
	ciphertext := bytes.NewBuffer(make([]byte, 0, cLen)) // pre-alloc correct length

	// Prefix the ciphertext with salt, AEAD ID and nonce
	ciphertext.Write(salt)
	ciphertext.WriteByte(id)
	ciphertext.Write(nonce)

	w := stream.EncryptWriter(ciphertext, nonce, nil)
	if _, err = w.Write(data); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return ciphertext.Bytes(), nil
}

// argon2Mu is used to control concurrent use of argon2,
// which is very cpu/ram intensive.
// Running concurrent operations most often provides no benefit anyway,
// since it already uses 32 threads.
var argon2Mu sync.Mutex

// ErrMaliciousData indicates that the stream cannot be
// decrypted by provided credentials.
var ErrMaliciousData = sio.NotAuthentic

// ErrUnexpectedHeader indicates that the data stream returned unexpected header
var ErrUnexpectedHeader = errors.New("unexpected header")

// DecryptData decrypts the data with the key derived
// from the salt (part of data) and the password using
// the PBKDF used in EncryptData. DecryptData returns
// the decrypted plaintext on success.
//
// The data must be a valid ciphertext produced by
// EncryptData. Otherwise, the decryption will fail.
func DecryptData(password string, data io.Reader) ([]byte, error) {
	// Parse the stream header
	var hdr [32 + 1 + 8]byte
	if _, err := io.ReadFull(data, hdr[:]); err != nil {
		if errors.Is(err, io.EOF) {
			// Incomplete header, return unexpected header
			return nil, ErrUnexpectedHeader
		}
		return nil, err
	}
	salt, id, nonce := hdr[0:32], hdr[32:33], hdr[33:41]

	var (
		err    error
		stream *sio.Stream
	)
	switch id[0] {
	case argon2idAESGCM:
		argon2Mu.Lock()
		key := argon2.IDKey([]byte(password), salt, argon2idTime, argon2idMemory, argon2idThreads, 32)
		argon2Mu.Unlock()
		stream, err = sio.AES_256_GCM.Stream(key)
	case argon2idChaCHa20Poly1305:
		argon2Mu.Lock()
		key := argon2.IDKey([]byte(password), salt, argon2idTime, argon2idMemory, argon2idThreads, 32)
		argon2Mu.Unlock()
		stream, err = sio.ChaCha20Poly1305.Stream(key)
	case pbkdf2AESGCM:
		key := pbkdf2.Key([]byte(password), salt, pbkdf2Cost, 32, sha256.New)
		stream, err = sio.AES_256_GCM.Stream(key)
	default:
		err = errors.New("madmin: invalid encryption algorithm ID")
	}
	if err != nil {
		return nil, err
	}

	return io.ReadAll(stream.DecryptReader(data, nonce, nil))
}

const (
	argon2idAESGCM           = 0x00
	argon2idChaCHa20Poly1305 = 0x01
	pbkdf2AESGCM             = 0x02
)

const (
	argon2idTime    = 1
	argon2idMemory  = 64 * 1024
	argon2idThreads = 4
	pbkdf2Cost      = 8192
)
