package modules

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"

	"go.starlark.net/starlark"

	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
)

var chunkSize = 4096

// Implement the crypto package
// https://docs.realm.pub/user-guide/eldritch#crypto

// Converts a JSON string to a Starlark Dict
// https://docs.realm.pub/user-guide/eldritch#cryptofrom_json
func cryptoFromJson(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// Check the number of arguments
	var str starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &str); err != nil {
		return nil, err
	}

	// Parse the JSON string into a Go map
	var data interface{}
	if err := json.Unmarshal([]byte(str.GoString()), &data); err != nil {
		return nil, err
	}

	return ToStarlarkValue(data)
}

func cryptoToJson(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// Check the number of arguments
	var val starlark.Value
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &val); err != nil {
		return nil, err
	}

	// Parse the JSON string into a Go map
	v, err := ToGolangValue(val)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return starlark.String(string(res)), err
}

func cryptoAesEncrypt(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	var dst starlark.String
	var key starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 3, &src, &dst, &key); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("crypto.aes_encrypt_file not impemented")
}

func cryptoAesDecrypt(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var src starlark.String
	var dst starlark.String
	var key starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 3, &src, &dst, &key); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("crypto.aes_decrypt_file not impemented")
}

func cryptoEncodeB64(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var str starlark.String
	var typ starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &str, &typ); err != nil {
		return nil, err
	}
	enc := base64.StdEncoding
	switch typ {
	case "STANDARD_NO_PAD":
		enc = base64.RawStdEncoding
	case "URL_SAFE":
		enc = base64.URLEncoding
	case "URL_SAFE_NO_PAD":
		enc = base64.RawURLEncoding
	case "STANDARD":
		fallthrough
	case "":
		break
	default:
		return nil, fmt.Errorf("invalid enconding selected '%s'", typ)
	}
	return starlark.String(enc.EncodeToString([]byte(str.GoString()))), nil
}

func cryptoDecodeB64(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var str starlark.String
	var typ starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &str, &typ); err != nil {
		return nil, err
	}
	enc := base64.StdEncoding
	switch typ {
	case "STANDARD_NO_PAD":
		enc = base64.RawStdEncoding
	case "URL_SAFE":
		enc = base64.URLEncoding
	case "URL_SAFE_NO_PAD":
		enc = base64.RawURLEncoding
	case "STANDARD":
		fallthrough
	case "":
		break
	default:
		return nil, fmt.Errorf("invalid enconding selected '%s'", typ)
	}
	res, err := enc.DecodeString(str.GoString())
	if err != nil {
		return nil, err
	}
	return starlark.String(string(res)), nil
}

func cryptoHashFile(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var file starlark.String
	var algo starlark.String
	if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &file, &algo); err != nil {
		return nil, err
	}
	var alg hash.Hash
	switch algo {
	case "MD5":
		alg = crypto.MD5.New()
	case "SHA1":
		alg = crypto.SHA1.New()
	case "SHA256":
		alg = crypto.SHA256.New()
	case "SHA512":
		alg = crypto.SHA512.New()
	default:
		return nil, fmt.Errorf("invalid algorithm selected '%s'", algo)
	}

	f, err := os.Open(file.GoString())
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, chunkSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		alg.Write(buf[:n])
	}
	res := hex.EncodeToString(alg.Sum(nil))
	return starlark.String(res), nil
}

var Crypto = Module{
	"from_json":        starlark.NewBuiltin("", cryptoFromJson),
	"to_json":          starlark.NewBuiltin("", cryptoToJson),
	"aes_encrypt_file": starlark.NewBuiltin("", cryptoAesEncrypt),
	"aes_decrypt_file": starlark.NewBuiltin("", cryptoAesDecrypt),
	"encode_b64":       starlark.NewBuiltin("", cryptoEncodeB64),
	"decode_b64":       starlark.NewBuiltin("", cryptoDecodeB64),
	"hash_file":        starlark.NewBuiltin("", cryptoHashFile),
}
