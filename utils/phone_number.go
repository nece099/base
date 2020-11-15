package utils

import (
	"errors"
	"fmt"

	"github.com/ttacon/libphonenumber"
)

type phoneNumberUtil struct {
	phoneNumber *libphonenumber.PhoneNumber
}

func MakePhoneNumberUtil(number, region string) (*phoneNumberUtil, error) {
	var (
		pnumber *libphonenumber.PhoneNumber
		err     error
	)

	if number == "" {
		err = errors.New("phone number empty")
		return nil, err
	}

	// Android客户端手机号格式为: 8611111111111, Parse结果为invalid country code
	// 转换成+8611111111111，再进行Parse
	if region == "" && number[:1] != "+" {
		number = "+" + number
	}

	// fmt.Println(number)
	// check phone invalid
	pnumber, err = libphonenumber.Parse(number, region)
	// fmt.Println(pnumber)
	if err != nil {
		// fmt.Println(err)
		err = errors.New(fmt.Sprintf("invalid phone number: %v", err))
	} else {
		if !libphonenumber.IsValidNumber(pnumber) {
			err = errors.New("invalid phone number")
		}
	}

	if err != nil {
		return nil, err
	} else {
		return &phoneNumberUtil{pnumber}, nil
	}
}

func (p *phoneNumberUtil) GetNormalizeDigits() string {
	// DB里存储归一化的phone
	return libphonenumber.NormalizeDigitsOnly(libphonenumber.Format(p.phoneNumber, libphonenumber.E164))
}

func (p *phoneNumberUtil) GetRegionCode() string {
	return libphonenumber.GetRegionCodeForNumber(p.phoneNumber)
}

// Check number
// 客户端发送的手机号格式为: "+86 111 1111 1111"，归一化
func CheckAndGetPhoneNumber(number string) (phoneNumber string, err error) {
	var (
		pnumber *phoneNumberUtil
	)

	pnumber, err = MakePhoneNumberUtil(number, "")
	if err != nil {
		return
	}

	return pnumber.GetNormalizeDigits(), nil
}
