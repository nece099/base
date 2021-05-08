/*
 *  Copyright (c) 2017, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package encrypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
)

/*
  telegram客户端使用的是pkcs1格式公钥证书，但GO语言没找到支持PKCS8格式的操作，
  参考
     http://blog.qiujinwu.com/2017/07/14/rsa/
     https://medium.com/@oyrxx/rsa%E7%A7%98%E9%92%A5%E4%BB%8B%E7%BB%8D%E5%8F%8Aopenssl%E7%94%9F%E6%88%90%E5%91%BD%E4%BB%A4-d3fcc689513f

  使用如下方法生成:
    openssl genrsa -out server.key 2048

    转成pcks8:
    openssl pkcs8 -topk8 -inform PEM -in server.key -outform pem -nocrypt -out server8.key
    openssl rsa -in server.key -pubout > public_pkcs8.pub

    转成pcks1:
    openssl rsa -in server.key -outform PEM -RSAPublicKey_out -out public_pkcs1.key
*/

var pkcs1PemPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvKLEOWTzt9Hn3/9Kdp/RdHcEhzmd8xXeLSpHIIzaXTLJDw8B
hJy1jR/iqeG8Je5yrtVabqMSkA6ltIpgylH///FojMsX1BHu4EPYOXQgB0qOi6kr
08iXZIH9/iOPQOWDsL+Lt8gDG0xBy+sPe/2ZHdzKMjX6O9B4sOsxjFrk5qDoWDri
oJorAJ7eFAfPpOBf2w73ohXudSrJE0lbQ8pCWNpMY8cB9i8r+WBitcvouLDAvmtn
TX7akhoDzmKgpJBYliAY4qA73v7u5UIepE8QgV0jCOhxJCPubP8dg+/PlLLVKyxU
5CdiQtZj2EMy4s9xlNKzX8XezE0MHEa6bQpnFwIDAQABAoIBACd+SGjfyursZoiO
MW/ejAK/PFJ3bKtNI8P++v9Enh8vF8swUBgMmzIdv93jZfnnD1mtT46kU6mXd3fy
FMunGVrjlwkLKET9MC8B5U46Es6T/H4fAA8KCzA+ywefOEnVA5pIsB7dIFFhyNDB
uO8zrBAFfsu+Y1KMlggsZaZGDXB/WVyUJDbEOMZstVx4uNhpcEgKYp28YQMP/yvv
dp4UgnTxXXXpDghzO5iqi5tUWY0p1lH2ii2OZBxEdqdDl7TirorhUDYIivyoe3B5
H30RNBRok/6w7W0WPyY2lSIcjd3cLPte6vx0QfBXVo2A6N9LTKAtAw3iWBp0x9NZ
N5p8OeECgYEA8QywXlM8nH5M7Sg2sMUYBOHA22O26ZPio7rJzcb8dlkV5gVHm+Kl
aDP61Uy8KoYABQ5kFdem/IQAUPepLxmJmiqfbwOIjfajOD3uVAQunFnDCHBWm4Uk
onbpdA5NlT/OUoSjIBemiBR/4CpDK1cEby/sg+EvQaGxqtedEe4xFmcCgYEAyFXe
MyAAOLpzmnCs9NYTTvMPofW8y+kLDodfbskl7M8q6l20VMo/E+g1gQ+65Aah901Z
/LKGi6HpzmHi5q9O2OJtqyI6FVwjXa07M5ueDbHcVKJw4hC9W0oHpMg8hqumPAWF
+MoN/Toy77p5LzoR30WUdhPvOAJPEL1p2a6r29ECgYEAiXfCEVkI5PqGZm2bmv4b
75TLhpJ8WwMSqms48Vi828V8Xpy+NOFxkVargv9rBBk9Y6TMYUSGH9Yr1AEZhBnd
RoVuPUJXmxaACPAQvetQpavvNR3T1od82AZWpvANQMONp7Oqz/+M4mhGcRHJEqti
hQJgsOk4KQbMqvChy/r6FZsCgYEAwyaqgkD9FkXC0UJLqWFUg8bQhqPcGwLUC34h
n8kAUbPpiU5omWQ+mATPAf8xvmkbo81NCJVb7W93U90U7ET/2NSRonCABkiwBtP2
ZKqGB68oA6YNspo960ytL38DPui80aFLxXQGtpPYBKEw5al6uXWNTozSrkvJe3QY
Rb4amdECgYBpGk7zPcK1TbJ++W5fkiory4qOdf0L1Zf0NbML4fY6dIww+dwMVUpq
FbsgCLqimqOFaaECU+LQEFUHHM7zrk7NBf7GzBvQ+qJx8zhJ66sFVox+IirBUyR9
Vh0+z5tIbFbKmYkO06NbeMlq87JexSlocPZtA3HMhEga5/0fHNHsNw==
-----END RSA PRIVATE KEY-----
`)

var pkcs8PemPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC8osQ5ZPO30eff
/0p2n9F0dwSHOZ3zFd4tKkcgjNpdMskPDwGEnLWNH+Kp4bwl7nKu1VpuoxKQDqW0
imDKUf//8WiMyxfUEe7gQ9g5dCAHSo6LqSvTyJdkgf3+I49A5YOwv4u3yAMbTEHL
6w97/Zkd3MoyNfo70Hiw6zGMWuTmoOhYOuKgmisAnt4UB8+k4F/bDveiFe51KskT
SVtDykJY2kxjxwH2Lyv5YGK1y+i4sMC+a2dNftqSGgPOYqCkkFiWIBjioDve/u7l
Qh6kTxCBXSMI6HEkI+5s/x2D78+UstUrLFTkJ2JC1mPYQzLiz3GU0rNfxd7MTQwc
RrptCmcXAgMBAAECggEAJ35IaN/K6uxmiI4xb96MAr88Undsq00jw/76/0SeHy8X
yzBQGAybMh2/3eNl+ecPWa1PjqRTqZd3d/IUy6cZWuOXCQsoRP0wLwHlTjoSzpP8
fh8ADwoLMD7LB584SdUDmkiwHt0gUWHI0MG47zOsEAV+y75jUoyWCCxlpkYNcH9Z
XJQkNsQ4xmy1XHi42GlwSApinbxhAw//K+92nhSCdPFddekOCHM7mKqLm1RZjSnW
UfaKLY5kHER2p0OXtOKuiuFQNgiK/Kh7cHkffRE0FGiT/rDtbRY/JjaVIhyN3dws
+17q/HRB8FdWjYDo30tMoC0DDeJYGnTH01k3mnw54QKBgQDxDLBeUzycfkztKDaw
xRgE4cDbY7bpk+KjusnNxvx2WRXmBUeb4qVoM/rVTLwqhgAFDmQV16b8hABQ96kv
GYmaKp9vA4iN9qM4Pe5UBC6cWcMIcFabhSSidul0Dk2VP85ShKMgF6aIFH/gKkMr
VwRvL+yD4S9BobGq150R7jEWZwKBgQDIVd4zIAA4unOacKz01hNO8w+h9bzL6QsO
h19uySXszyrqXbRUyj8T6DWBD7rkBqH3TVn8soaLoenOYeLmr07Y4m2rIjoVXCNd
rTszm54NsdxUonDiEL1bSgekyDyGq6Y8BYX4yg39OjLvunkvOhHfRZR2E+84Ak8Q
vWnZrqvb0QKBgQCJd8IRWQjk+oZmbZua/hvvlMuGknxbAxKqazjxWLzbxXxenL40
4XGRVquC/2sEGT1jpMxhRIYf1ivUARmEGd1GhW49QlebFoAI8BC961Clq+81HdPW
h3zYBlam8A1Aw42ns6rP/4ziaEZxEckSq2KFAmCw6TgpBsyq8KHL+voVmwKBgQDD
JqqCQP0WRcLRQkupYVSDxtCGo9wbAtQLfiGfyQBRs+mJTmiZZD6YBM8B/zG+aRuj
zU0IlVvtb3dT3RTsRP/Y1JGicIAGSLAG0/ZkqoYHrygDpg2ymj3rTK0vfwM+6LzR
oUvFdAa2k9gEoTDlqXq5dY1OjNKuS8l7dBhFvhqZ0QKBgGkaTvM9wrVNsn75bl+S
KivLio51/QvVl/Q1swvh9jp0jDD53AxVSmoVuyAIuqKao4VpoQJT4tAQVQcczvOu
Ts0F/sbMG9D6onHzOEnrqwVWjH4iKsFTJH1WHT7Pm0hsVsqZiQ7To1t4yWrzsl7F
KWhw9m0DccyESBrn/R8c0ew3
-----END PRIVATE KEY-----
`)

var pkcs1PemPublicKey = []byte(`
-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAvKLEOWTzt9Hn3/9Kdp/RdHcEhzmd8xXeLSpHIIzaXTLJDw8BhJy1
jR/iqeG8Je5yrtVabqMSkA6ltIpgylH///FojMsX1BHu4EPYOXQgB0qOi6kr08iX
ZIH9/iOPQOWDsL+Lt8gDG0xBy+sPe/2ZHdzKMjX6O9B4sOsxjFrk5qDoWDrioJor
AJ7eFAfPpOBf2w73ohXudSrJE0lbQ8pCWNpMY8cB9i8r+WBitcvouLDAvmtnTX7a
khoDzmKgpJBYliAY4qA73v7u5UIepE8QgV0jCOhxJCPubP8dg+/PlLLVKyxU5Cdi
QtZj2EMy4s9xlNKzX8XezE0MHEa6bQpnFwIDAQAB
-----END RSA PUBLIC KEY-----
`)

var pkcs8PemPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvKLEOWTzt9Hn3/9Kdp/R
dHcEhzmd8xXeLSpHIIzaXTLJDw8BhJy1jR/iqeG8Je5yrtVabqMSkA6ltIpgylH/
//FojMsX1BHu4EPYOXQgB0qOi6kr08iXZIH9/iOPQOWDsL+Lt8gDG0xBy+sPe/2Z
HdzKMjX6O9B4sOsxjFrk5qDoWDrioJorAJ7eFAfPpOBf2w73ohXudSrJE0lbQ8pC
WNpMY8cB9i8r+WBitcvouLDAvmtnTX7akhoDzmKgpJBYliAY4qA73v7u5UIepE8Q
gV0jCOhxJCPubP8dg+/PlLLVKyxU5CdiQtZj2EMy4s9xlNKzX8XezE0MHEa6bQpn
FwIDAQAB
-----END PUBLIC KEY-----
`)

// new cert
// var pkcs1PemPrivateKey = []byte(`
// -----BEGIN RSA PRIVATE KEY-----
// MIIEpAIBAAKCAQEA0NUJlRK/QEsukdNLXbFoip8SajmkzAdoMWF73zWeTKVNek7m
// cIRrNEEM9ZckaqxhgdWeUSL3SObkRxU9L3Pt/lEDZ/pi0sbYvjj9JWLLhZDYQHn/
// jeShyR/SrmwfPgUE6dyE2LqFWbRaMSJusGPidovS72zDeyaUQRIHfhJy23gLgH0A
// xq/kSWvH/m5SBbzWI8CNnrAB8iTOBc6viGXKEWB4ZF+4tqfwF2iV64J1H8NmCrn/
// 4jh2dTSsIr884x39XbizM/jRcOnJ8fW1U1gYpwBCrI4RlCQQXLD3dUyIijIOLqMc
// CbjZ5+6NXosz0wkwUnr+aAAupKOtXiYP8bSp9QIDAQABAoIBAGrtunNXXxA3rsfC
// TiPSVDoui0pS67wAyuwGA1xeYwjR12MaBUp1s0LVUCJsWpw4WdEWJXNcGQx+FUME
// cAjdLm564uiZv4I3iQGVwqEi/h0M9n3FOgJYoDKQldrzx+eEwGhSnr8ueltdSpVA
// ETdGXc3feIlZwppLPbw31BhMr/0Jg+kxGOzLaaPyAXj/8yzYXzMaL9mkDwcWsYAd
// hfLObTRujsD0/pGaYtEK98JKoIT2aOckSdFu7sqg4hN2VNn78O/p3/0Lv+EIFaF6
// xAh8RrleCQ/seULbWyDo44GqIF4l+81YccPX4W94//9VqEs66/faLeLG2/HUT+wJ
// /FIenN0CgYEA/XsoVNMF9lO6amhqnvL9Jn34UZRr3qxNOgRAjTx9OHekvZyr+CD2
// mPELcWGdhNcx7MP4OFBUZIa/wMZgZp0iMnvuvg0YvwDrX/3cRJy7YySTfRvpbNTx
// jCoJCnWbKmV32ST8ZR7fTalA56rW6D3+USVUBbbk4OXz2INaYONXdOsCgYEA0uhL
// oXKgpFbHn41iDxX8oeGwMk6bEhOMJ71M0MugjX1bH+NvT2oWOHqQjgciLyPvbG7S
// bv5vfhK96A4Ov1LRGi9A1o+ZFDGbDRtJGec2gwo8HKXwTf9EWkh9Tt1V4f9reHlP
// 75rZYrAhcvPbimn382jie2ORwiPXyRK3LMFWJJ8CgYEAjq15CS3yyDFe17BIe4m4
// lqcHVBwgD6mampJ0J0uqDFPEBfqfDb64L2RWlY5llLVwY533JPOKXT8/xemjr365
// FgOOYamLiU+iLVj+WByEmYyn/B7u6BSAle2/QwTpvxZ4PGDGNMEI3nTrlLsj1nu2
// n8RMJB9Le4/UDsX45FpzCtsCgYB3pOPSsK5EzB3ue0wXdsecJeXIhCMgPAqUOKUt
// BXcNDQH2sxTgHjSA0bbTe2R/DYmzH6Ms6BXjlUo6LE9dZePNUOUdUtTqScHFy6bK
// lQmtiM7VCaWq+ZaTCPBdHt6rmDQlYdxg9p0/iN9Q0NnISZkpcmSYzsFPOvoczQsw
// znTJzQKBgQCGFoorY7LNk8r8Qj+upooIcoVit8UzgWF+GsTijFaquks5dlRWojxQ
// 83UOQFyIgU3T01irt9/kdAMhJr4vRcVS4Q7W2qLF2Rqs8LMrFltn4bLgFmCy91iK
// vLroX+XeEhzU3H6E6DNKXSysH+z9HLDdmM53WHu7/aGqjhhKEQqbPg==
// -----END RSA PRIVATE KEY-----
// `)

// var pkcs8PemPrivateKey = []byte(`
// -----BEGIN PRIVATE KEY-----
// MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQ1QmVEr9ASy6R
// 00tdsWiKnxJqOaTMB2gxYXvfNZ5MpU16TuZwhGs0QQz1lyRqrGGB1Z5RIvdI5uRH
// FT0vc+3+UQNn+mLSxti+OP0lYsuFkNhAef+N5KHJH9KubB8+BQTp3ITYuoVZtFox
// Im6wY+J2i9LvbMN7JpRBEgd+EnLbeAuAfQDGr+RJa8f+blIFvNYjwI2esAHyJM4F
// zq+IZcoRYHhkX7i2p/AXaJXrgnUfw2YKuf/iOHZ1NKwivzzjHf1duLMz+NFw6cnx
// 9bVTWBinAEKsjhGUJBBcsPd1TIiKMg4uoxwJuNnn7o1eizPTCTBSev5oAC6ko61e
// Jg/xtKn1AgMBAAECggEAau26c1dfEDeux8JOI9JUOi6LSlLrvADK7AYDXF5jCNHX
// YxoFSnWzQtVQImxanDhZ0RYlc1wZDH4VQwRwCN0ubnri6Jm/gjeJAZXCoSL+HQz2
// fcU6AligMpCV2vPH54TAaFKevy56W11KlUARN0Zdzd94iVnCmks9vDfUGEyv/QmD
// 6TEY7Mtpo/IBeP/zLNhfMxov2aQPBxaxgB2F8s5tNG6OwPT+kZpi0Qr3wkqghPZo
// 5yRJ0W7uyqDiE3ZU2fvw7+nf/Qu/4QgVoXrECHxGuV4JD+x5QttbIOjjgaogXiX7
// zVhxw9fhb3j//1WoSzrr99ot4sbb8dRP7An8Uh6c3QKBgQD9eyhU0wX2U7pqaGqe
// 8v0mffhRlGverE06BECNPH04d6S9nKv4IPaY8QtxYZ2E1zHsw/g4UFRkhr/AxmBm
// nSIye+6+DRi/AOtf/dxEnLtjJJN9G+ls1PGMKgkKdZsqZXfZJPxlHt9NqUDnqtbo
// Pf5RJVQFtuTg5fPYg1pg41d06wKBgQDS6EuhcqCkVsefjWIPFfyh4bAyTpsSE4wn
// vUzQy6CNfVsf429PahY4epCOByIvI+9sbtJu/m9+Er3oDg6/UtEaL0DWj5kUMZsN
// G0kZ5zaDCjwcpfBN/0RaSH1O3VXh/2t4eU/vmtlisCFy89uKaffzaOJ7Y5HCI9fJ
// ErcswVYknwKBgQCOrXkJLfLIMV7XsEh7ibiWpwdUHCAPqZqaknQnS6oMU8QF+p8N
// vrgvZFaVjmWUtXBjnfck84pdPz/F6aOvfrkWA45hqYuJT6ItWP5YHISZjKf8Hu7o
// FICV7b9DBOm/Fng8YMY0wQjedOuUuyPWe7afxEwkH0t7j9QOxfjkWnMK2wKBgHek
// 49KwrkTMHe57TBd2x5wl5ciEIyA8CpQ4pS0Fdw0NAfazFOAeNIDRttN7ZH8NibMf
// oyzoFeOVSjosT11l481Q5R1S1OpJwcXLpsqVCa2IztUJpar5lpMI8F0e3quYNCVh
// 3GD2nT+I31DQ2chJmSlyZJjOwU86+hzNCzDOdMnNAoGBAIYWiitjss2TyvxCP66m
// ighyhWK3xTOBYX4axOKMVqq6Szl2VFaiPFDzdQ5AXIiBTdPTWKu33+R0AyEmvi9F
// xVLhDtbaosXZGqzwsysWW2fhsuAWYLL3WIq8uuhf5d4SHNTcfoToM0pdLKwf7P0c
// sN2YzndYe7v9oaqOGEoRCps+
// -----END PRIVATE KEY-----
// `)

// var pkcs1PemPublicKey = []byte(`
// -----BEGIN RSA PUBLIC KEY-----
// MIIBCgKCAQEA0NUJlRK/QEsukdNLXbFoip8SajmkzAdoMWF73zWeTKVNek7mcIRr
// NEEM9ZckaqxhgdWeUSL3SObkRxU9L3Pt/lEDZ/pi0sbYvjj9JWLLhZDYQHn/jeSh
// yR/SrmwfPgUE6dyE2LqFWbRaMSJusGPidovS72zDeyaUQRIHfhJy23gLgH0Axq/k
// SWvH/m5SBbzWI8CNnrAB8iTOBc6viGXKEWB4ZF+4tqfwF2iV64J1H8NmCrn/4jh2
// dTSsIr884x39XbizM/jRcOnJ8fW1U1gYpwBCrI4RlCQQXLD3dUyIijIOLqMcCbjZ
// 5+6NXosz0wkwUnr+aAAupKOtXiYP8bSp9QIDAQAB
// -----END RSA PUBLIC KEY-----
// `)

// var pkcs8PemPublicKey = []byte(`
// -----BEGIN PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0NUJlRK/QEsukdNLXbFo
// ip8SajmkzAdoMWF73zWeTKVNek7mcIRrNEEM9ZckaqxhgdWeUSL3SObkRxU9L3Pt
// /lEDZ/pi0sbYvjj9JWLLhZDYQHn/jeShyR/SrmwfPgUE6dyE2LqFWbRaMSJusGPi
// dovS72zDeyaUQRIHfhJy23gLgH0Axq/kSWvH/m5SBbzWI8CNnrAB8iTOBc6viGXK
// EWB4ZF+4tqfwF2iV64J1H8NmCrn/4jh2dTSsIr884x39XbizM/jRcOnJ8fW1U1gY
// pwBCrI4RlCQQXLD3dUyIijIOLqMcCbjZ5+6NXosz0wkwUnr+aAAupKOtXiYP8bSp
// 9QIDAQAB
// -----END PUBLIC KEY-----
// `)

// bca2c43964f3b7d1e7dfff4a769fd174770487399df315de2d2a47208cda5d32c90f0f01849cb58d1fe2a9e1bc25ee72aed55a6ea312900ea5b48a60ca51fffff1688ccb17d411eee043d8397420074a8e8ba92bd3c8976481fdfe238f40e583b0bf8bb7c8031b4c41cbeb0f7bfd991ddcca3235fa3bd078b0eb318c5ae4e6a0e8583ae2a09a2b009ede1407cfa4e05fdb0ef7a215ee752ac913495b43ca4258da4c63c701f62f2bf96062b5cbe8b8b0c0be6b674d7eda921a03ce62a0a49058962018e2a03bdefeeee5421ea44f10815d2308e8712423ee6cff1d83efcf94b2d52b2c54e4276242d663d84332e2cf7194d2b35fc5decc4d0c1c46ba6d0a6717
type RSACryptor struct {
	key *rsa.PrivateKey
}

// TODO(@benqi): 这里写死了pkcs1PemPrivateKey
func NewRSACryptor() *RSACryptor {
	// testPrivateKey
	block, _ := pem.Decode([]byte(pkcs1PemPrivateKey))
	if block == nil {
		panic("Invalid pemsKeyData: " + string(pkcs1PemPrivateKey))
		return nil
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("Failed to parse private key: " + err.Error())
	}

	return &RSACryptor{
		key: key,
	}
}

func (m *RSACryptor) Encrypt(b []byte) []byte {
	c := new(big.Int)
	c.Exp(new(big.Int).SetBytes(b), big.NewInt(int64(m.key.E)), m.key.N)

	e := c.Bytes()
	var size = len(e)
	if size < 256 {
		size = 256
	}

	res := make([]byte, size)
	copy(res, c.Bytes())

	return res
}

func (m *RSACryptor) Decrypt(b []byte) []byte {
	c := new(big.Int)
	c.Exp(new(big.Int).SetBytes(b), m.key.D, m.key.N)
	return c.Bytes()
}
