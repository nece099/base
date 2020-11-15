/*
	Base is https://github.com/vbatts/go-google-authenticator
*/
package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"image"
	"image/png"
	"math"
	"reflect"
	"time"

	"github.com/skip2/go-qrcode"
)

var (
	Dim               = 256 // dimensions for the QrCode()
	Debug             = false
	Issuer            = "go-google-authenticator"
	ErrCipherNotFound = errors.New("cipher not supported")
)

func getTs() int64 {
	un := float64(time.Now().UnixNano()) / float64(1000) / float64(30)
	return int64(math.Floor(un))
}

// sah256 can't use
// func GenSecretKey(cipher string) (string, error) {
func GenSecretKey() (string, error) {
	cipher := "sha1"
	var hmac_hash hash.Hash
	switch cipher {
	case "sha1":
		hmac_hash = sha1.New()
	case "sha256":
		hmac_hash = sha256.New()
	default:
		return "", ErrCipherNotFound
	}
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, getTs())
	if err != nil {
		return "", err
	}
	h := hmac.New(func() hash.Hash { return hmac_hash }, buf.Bytes())
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func QrCode(account, issuer, key string) string {
	if len(issuer) > 0 {
		Issuer = issuer
	}
	otp_str := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		Issuer, account, base32.StdEncoding.EncodeToString([]byte(key)), Issuer)
	// use our qrcode
	// return fmt.Sprintf("https://chart.googleapis.com/chart?chs=%dx%d&cht=qr&choe=UTF-8&chl=%s",
	// 	Dim, Dim, url.QueryEscape(otp_str))
	ret_str, err := MakeQRCode2(otp_str, Dim)
	if err != nil {
		Log.Errorf("MakeQRCode failed, err = %v", err)
		return ""
	}
	return ret_str
}

func MakeQRCode2(instr string, imgSize int) (string, error) {

	var pngbyte []byte
	pngbyte, err := qrcode.Encode(instr, qrcode.Low, imgSize)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(pngbyte))
	if err != nil {
		return "", err
	}

	var subImg image.Image = img

	buf := new(bytes.Buffer)
	err = png.Encode(buf, subImg)
	if err != nil {
		return "", err
	}

	pngbyte = buf.Bytes()

	imgb64 := Base64Encode2(pngbyte)
	pngstr := "data:image/png;base64," + string(imgb64)

	return pngstr, nil
}

const (
	base64Table2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table2)

func Base64Encode2(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

type Authenticator struct {
	interval  int
	SecretKey []byte
	hash      hash.Hash
}

/*
Get a new Authenticator.
from a string of the ToTP salt,
and a bool of whether to use sha256 (default is sha1)

Also defaults to 30 second interval
*/
// sah256 can't use
// func New(salt string, twofiftysix bool) Authenticator {
func New(userKey string) Authenticator {
	twofiftysix := false
	var h hash.Hash
	if twofiftysix {
		h = sha256.New()
	} else {
		h = sha1.New()
	}

	return Authenticator{
		interval:  30, // always use 30 sec
		hash:      h,
		SecretKey: bytes.NewBufferString(userKey).Bytes(),
	}
}

/*
Construct an new Hmac hash with the secret and hashlib (sha1 or sha256)
*/
func (a Authenticator) Hmac() hash.Hash {
	h := a.copyHash(a.hash)
	return hmac.New(func() hash.Hash { return h }, a.SecretKey)
}

/*
Convenience to get the token for the present time
*/
func (a Authenticator) GetCodeCurrent() ([]int, error) {
	nowUnix := time.Now().Unix()
	codes := []int{}
	for i := -1; i <= 1; i++ {
		code, _, err := a.GetCode(i, nowUnix)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

/*
Generate the Time-based One Time Passcode.

c = -1 :: the previous code
c = 0  :: the current code
c = 1  :: the next code

the returns are: code, seconds to expire, error
*/
func (a Authenticator) GetCode(c int, now int64) (int, int64, error) {
	t_chunk := (now / int64(a.interval)) + int64(c)

	buf_in := bytes.Buffer{}
	err := binary.Write(&buf_in, binary.BigEndian, int64(t_chunk))
	if err != nil {
		return 0, 0, err
	}

	h := a.Hmac()
	h.Write(buf_in.Bytes())
	sum := h.Sum(nil)
	offset := sum[len(sum)-1] & 0xF
	code_sect := sum[offset : offset+4]
	if Debug {
		fmt.Printf("sum:\t\t%t\n", sum)
		fmt.Printf("last:\t\t%t\n", sum[len(sum)-1])
		fmt.Printf("offset:\t\t%t\n", offset)
		fmt.Printf("code_sect:\t%t\n", code_sect)
		fmt.Printf("code_sect:\t%#v\n", code_sect)
	}
	var code int32
	buf_out := bytes.NewBuffer(code_sect)
	err = binary.Read(buf_out, binary.BigEndian, &code)
	if err != nil {
		return 0, 0, err
	}
	if Debug {
		fmt.Printf("unpacked code:\t%#v\n", code)
		fmt.Printf("unpacked code:\t%b\n", code)
	}
	code = code & 0x7FFFFFFF
	if Debug {
		fmt.Printf("sig bit:\t%#v\n", code)
		fmt.Printf("sig bit:\t%b\n", code)
	}
	code = code % 1000000
	if Debug {
		fmt.Printf("mod1000000:\t%#v\n", code)
		fmt.Printf("mod1000000:\t%b\n", code)
	}

	i := int64(a.interval)
	x := (((now + i) / i) * i) - now
	if Debug {
		fmt.Printf("expires:\t%d\n", x)
	}
	return int(code), x, nil
}

func (a Authenticator) Authenticate(token string) bool {
	codes, err := a.GetCodeCurrent()
	if err != nil {
		Log.Error("err = %v", err)
		return false
	}
	for _, code := range codes {
		sCode := fmt.Sprintf("%06d", code)
		if token == sCode {
			return true
		}
	}
	return false
}

func (a Authenticator) copyHash(src hash.Hash) hash.Hash {
	typ := reflect.TypeOf(src)
	val := reflect.ValueOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	elem := reflect.New(typ).Elem()
	elem.Set(val)
	return elem.Addr().Interface().(hash.Hash)
}
