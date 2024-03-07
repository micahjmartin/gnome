package sarcophagus

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

const chunkSize int = 50000
const fileheader = "sarcophagus"

type Vault struct {
	f *os.File
	k [32]byte

	files map[string]int64
}

const (
	time        uint32 = 1
	memory             = 64 * 1024
	parallelism        = 1
)

func Open(path, password string) (*Vault, error) {
	var key [32]byte
	akey := argon2.IDKey([]byte(password), []byte("aaaaaaaa"), time, memory, parallelism, 32)
	copy(key[:], akey)
	v := &Vault{
		k:     key,
		files: make(map[string]int64),
	}
	var err error
	v.f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	h, err := v.unseal(nil)
	if err != nil {
		return nil, fmt.Errorf("invalid vault")
	}
	if len(h) != len(fileheader) || string(h) != fileheader {
		return nil, fmt.Errorf("invalid vault")
	}

	for {
		name, err := v.unseal(nil)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read name: %v", err)
		}
		off, _ := v.f.Seek(0, io.SeekCurrent)
		v.files[string(name)] = off
		if err = v.skip(); err != nil {
			return nil, fmt.Errorf("cannot skip contents: %v", err)
		}
	}
	return v, nil
}

// Read the given entry from the vault and write it to string
func (v *Vault) ReadFile(name string, w io.Writer) error {
	off, ok := v.files[name]
	if !ok {
		return fmt.Errorf("file '%s' does not exist", name)
	}
	if off2, err := v.f.Seek(off, io.SeekStart); off != off2 || err != nil {
		return fmt.Errorf("could not seek: %v", err)
	}
	_, err := v.unseal(w)
	return err
}
func (v *Vault) Exists(name string) bool {
	_, ok := v.files[name]
	return ok
}

// List all the files in the Sarcophagus
func (v *Vault) Files() []string {
	res := make([]string, 0, len(v.files))
	for n := range v.files {
		res = append(res, n)
	}
	return res
}

// Close the vault, saving all the contents
func (v *Vault) Close() error {
	return v.f.Close()
}

/* INTERNAL FUNCTIONS */

type header struct {
	size  uint64
	nonce [24]byte
	next  bool // If the next bit is set, it mean the next chunk in the vault belongs to the same entry
}

// Read the nonce and size for this entry
func (v *Vault) readHeader() (hdr header, err error) {
	b := make([]byte, 8)
	if _, err = v.f.Read(b); err != nil {
		return
	}
	if _, err = v.f.Read(hdr.nonce[:]); err != nil {
		return
	}
	sz := binary.BigEndian.Uint64(b)
	hdr.next = sz>>63 == 1
	hdr.size = sz & 0x7fffffffffffffff // 63 bits
	return
}

// Easy unseal
func (v *Vault) unseal(w io.Writer) ([]byte, error) {
	for {
		off, _ := v.f.Seek(0, io.SeekCurrent)
		hdr, err := v.readHeader()
		if err != nil {
			return nil, err
		}
		// Safety check
		if hdr.size > uint64(chunkSize)*2 {
			return nil, fmt.Errorf("invalid sarcophagus")
		}
		b := make([]byte, hdr.size+secretbox.Overhead)
		if _, err := v.f.Read(b); err != nil {
			return nil, err
		}
		res, ok := secretbox.Open(nil, b, &hdr.nonce, &v.k)
		if !ok {
			return nil, fmt.Errorf("could not unseal at %d", off)
		}
		if w != nil {
			w.Write(res)
		}
		if !hdr.next {
			return res, nil
		}
	}
}

// Skip an entry in the vault (used for reading filenames but skipping content)
func (v *Vault) skip() error {
	for {
		hdr, err := v.readHeader()
		if err != nil {
			return err
		}
		// Jump to the end of this entry
		_, err = v.f.Seek(int64(hdr.size)+secretbox.Overhead, io.SeekCurrent)
		if err != nil {
			return err
		}
		if !hdr.next {
			break
		}
	}
	return nil
}
