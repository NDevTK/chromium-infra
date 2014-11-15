// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cipd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"infra/cipd/internal/keys"
)

// Source of randomness for signing. Can be mocked in tests.
var signingEntropy = rand.Reader

// Sign generates a new signature block given a package data. Multiple such
// blocks may be serialized together with MarshalSignatureList. Resulting byte
// buffer should be appended to a package file. It is then can be read and
// deserialized by ReadSignatureList.
func Sign(data io.Reader, key *rsa.PrivateKey) (block SignatureBlock, err error) {
	// Hash.
	hash := sigBlockHash.New()
	_, err = io.Copy(hash, data)
	if err != nil {
		return
	}
	digest := hash.Sum(nil)

	// Sign the hash with the private key to get the signature.
	sig, err := rsa.SignPKCS1v15(signingEntropy, key, sigBlockHash, digest)
	if err != nil {
		return
	}
	keyFingerprint, err := keys.PublicKeyFingerprint(&key.PublicKey)
	if err == nil {
		block = SignatureBlock{
			HashAlgo:      sigBlockHashName,
			Hash:          digest,
			SignatureAlgo: sigBlockSigName,
			SignatureKey:  keyFingerprint,
			Signature:     sig,
		}
	}
	return
}

// MarshalSignature converts one SignatureBlock into byte buffer win PEM encoded
// signature data.
func MarshalSignature(s *SignatureBlock) ([]byte, error) {
	jsonBlob, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	out := pem.EncodeToMemory(&pem.Block{
		Type:  sigBlockPEMType,
		Bytes: jsonBlob,
	})
	return out, nil
}

// UnmarshalSignature takes byte buffer with PEM encoded signature and converts
// it to SignatureBlock. Passed byte buffer should contain one and only one
// PEM block.
func UnmarshalSignature(b []byte) (sig SignatureBlock, err error) {
	// Ensure it starts with PEM header. pem.Decode will skip garbage up to PEM
	// header. We want to be stricter.
	sigBlockPEMHeader := fmt.Sprintf("-----BEGIN %s-----\n", sigBlockPEMType)
	idx := bytes.Index(b, []byte(sigBlockPEMHeader))
	if idx != 0 {
		err = fmt.Errorf("Not a valid signature PEM block")
		return
	}

	pemBlock, rest := pem.Decode(b)
	if pemBlock == nil {
		err = fmt.Errorf("Not a valid signature PEM block: can't read PEM block")
		return
	}
	if pemBlock.Type != sigBlockPEMType {
		// PEM type already has been verified by validating header.
		panic("Impossible PEM block type")
		return
	}
	if len(rest) != 0 {
		err = fmt.Errorf("Not a valid signature PEM block: undecoded data left")
		return
	}

	err = json.Unmarshal(pemBlock.Bytes, &sig)
	return
}

// MarshalSignatureList joins a list of signatures into a single byte buffer.
// Note that this operation is not additive, serialized signature list ends with
// the magic footer. A signed package is a concatenation of a package data,
// generated by BuildPackage, and a packed list of signatures, generated by
// MarshalSignatureList.
func MarshalSignatureList(signatures []SignatureBlock) ([]byte, error) {
	// Empty list is special, no need to attach anything to the file.
	if len(signatures) == 0 {
		return []byte{}, nil
	}

	// Don't bother returning error. If writing to byte buffer fails, then the
	// process is doomed anyway.
	out := bytes.Buffer{}
	write := func(b []byte) {
		_, err := out.Write(b)
		if err != nil {
			panic("Failed to write to byte buffer")
		}
	}

	// Write PEM encoded JSON blobs.
	for _, s := range signatures {
		blob, err := MarshalSignature(&s)
		if err != nil {
			return nil, err
		}
		write(blob)
	}

	// Write the total size of PEMs.
	write([]byte(fmt.Sprintf("%d", len(out.Bytes()))))
	return out.Bytes(), nil
}

// ReadSignatureList finds serialized list of signatures at the tail of the
// file, reads and deserializes it. It also returns absolute offset of where
// this block was found.
func ReadSignatureList(r io.ReadSeeker) (blocks []SignatureBlock, offset int64, err error) {
	// Do not blindly seek into negative offset, bytes.Reader doesn't like it.
	total, err := r.Seek(0, os.SEEK_END)
	if err != nil {
		return
	}

	sigBlockPEMHeader := fmt.Sprintf("-----BEGIN %s-----\n", sigBlockPEMType)
	sigBlockPEMFooter := fmt.Sprintf("-----END %s-----\n", sigBlockPEMType)

	// Offset footer is no longer than 16 bytes (it's ascii string with int32).
	// We also need to read last PEM block footer to be sure that we are indeed
	// reading PEM blocks and not just some garbage.
	seekTo := total - int64(len(sigBlockPEMFooter)) - 16
	if seekTo < 0 {
		seekTo = 0
	}
	_, err = r.Seek(seekTo, os.SEEK_SET)
	if err != nil {
		return
	}
	tail, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	// Find last PEM block footer.
	idx := bytes.LastIndex(tail, []byte(sigBlockPEMFooter))
	if idx == -1 {
		// No PEM blocks at all. It's fine, just no signatures at all.
		offset = total
		return
	}

	// Read the combined PEMs length from the footer.
	pemsLen, err := strconv.ParseInt(string(tail[idx+len(sigBlockPEMFooter):]), 10, 32)
	if err != nil {
		err = fmt.Errorf("Not a valid signature block: bad offset")
		return
	}

	// We already found at least one PEM block...
	if pemsLen < int64(len(sigBlockPEMHeader)+len(sigBlockPEMFooter)) {
		err = fmt.Errorf("Not a valid signature block: offset if too small")
		return
	}

	// Where PEMs start.
	offset = seekTo + int64(idx+len(sigBlockPEMFooter)) - pemsLen
	if offset < 0 {
		err = fmt.Errorf("Not a valid signature block: negative offset")
		return
	}

	// Read them all.
	_, err = r.Seek(offset, os.SEEK_SET)
	if err != nil {
		return
	}
	pems := make([]byte, pemsLen)
	_, err = io.ReadFull(r, pems)
	if err != nil {
		return
	}

	// Read all PEM encoded blocks.
	for len(pems) != 0 {
		// Find first PEM block header and footer.
		start := bytes.Index(pems, []byte(sigBlockPEMHeader))
		if start != 0 {
			err = fmt.Errorf("Not a valid signature block: not a PEM block")
			return
		}
		end := bytes.Index(pems, []byte(sigBlockPEMFooter))
		if end == -1 {
			// Actually this is unreachable, since code above has verified that
			// there's a valid footer at the end of |pems|. So if some corrupted
			// footer creeps in inside, UnmarshalSignature call below will fail
			// first (by trying to decode non-valid PEM).
			err = fmt.Errorf("Not a valid signature block: no PEM footer")
			return
		}
		end += len(sigBlockPEMFooter)

		// Read it.
		var sig SignatureBlock
		sig, err = UnmarshalSignature(pems[:end])
		if err != nil {
			return
		}
		blocks = append(blocks, sig)

		// Move to the next block.
		pems = pems[end:]
	}

	return
}
