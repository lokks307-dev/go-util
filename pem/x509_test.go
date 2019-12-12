package pem

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/lokks307/pkcs8"
	"github.com/stretchr/testify/assert"
)

const (
	testCertPem = `-----BEGIN CERTIFICATE-----
MIIByDCCAW6gAwIBAgIRANe5mco0f7UWwpWksW4f/UYwCgYIKoZIzj0EAwIwKjEL
MAkGA1UEAxMCQ04xCzAJBgNVBAYTAmtyMQ4wDAYDVQQKEwVsb2trczAeFw0xOTEw
MTUwMTAwNTlaFw0yMDEwMTQwMTAwNTlaMDExEjAQBgNVBAMTCXRlc3QtbmFtZTEL
MAkGA1UEBhMCa3IxDjAMBgNVBAoTBWxva2tzMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEgT5mmLPtodHt1/IrVDQV9Cv4JMV5ET/wtKj2IWdZ2WOP6EzYNbW4iWHP
NQ9SE+yE3XlkRvXJ+1jGP+cDReaQGaNuMGwwIQYDVR0OBBoEGAW1qy6rCwvMH44x
ZPtzTRiaRRofM8vJEzAOBgNVHQ8BAf8EBAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIB
ATAjBgNVHSMEHDAagBhotF3YEXygRwDJeHGwycYHbLIuRZDs6DUwCgYIKoZIzj0E
AwIDSAAwRQIgQKw4XEGmX/nUcivfQAShcSi5fIYXy1/U1dDW4TX71OgCIQDb2m3+
4usQnMcTn4tQXSvYjoJ4J5aLZtvI1OWDC5dQEQ==
-----END CERTIFICATE-----`

	testTrueCAcertPem = `-----BEGIN CERTIFICATE-----
MIIBwDCCAWegAwIBAgIRALv1dtWfAcPkce7sXeDIBzQwCgYIKoZIzj0EAwIwKjEL
MAkGA1UEAxMCQ04xCzAJBgNVBAYTAmtyMQ4wDAYDVQQKEwVsb2trczAeFw0xOTEw
MTQwNzQ1MDBaFw0yMDEwMTMwNzQ1MDBaMCoxCzAJBgNVBAMTAkNOMQswCQYDVQQG
EwJrcjEOMAwGA1UEChMFbG9ra3MwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQm
zBMecRVlF/g+XyD+MUlaHBMw0mw/jIlvHGInC4AnQm4KiQkQj8K31w05EPZ4/vQ0
Zdr1KuiQaSAGLQGqrhFNo24wbDAhBgNVHQ4EGgQYaLRd2BF8oEcAyXhxsMnGB2yy
LkWQ7Og1MA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMCMGA1Ud
IwQcMBqAGGi0XdgRfKBHAMl4cbDJxgdssi5FkOzoNTAKBggqhkjOPQQDAgNHADBE
AiAPjyq+x1cpS/quxZTyMmb5HBz2GW6FXyqm3dwRl60dpQIgcxFTyoY7P/Gc8Ptz
1PB6KTQP6yJKGsLyd5ieY59Bn9o=
-----END CERTIFICATE-----`

	testFalseCAcertPem = `-----BEGIN CERTIFICATE-----
MIIB1DCCAXmgAwIBAgIRAN/nM+ZL7GV9gEt6ivJkrh4wCgYIKoZIzj0EAwIwMzEO
MAwGA1UEAxMFZmFsc2UxCzAJBgNVBAYTAnVzMRQwEgYDVQQKDAtsb2trc19mYWxz
ZTAeFw0xOTEwMTUwMTIwNDZaFw0yMDEwMTQwMTIwNDZaMDMxDjAMBgNVBAMTBWZh
bHNlMQswCQYDVQQGEwJ1czEUMBIGA1UECgwLbG9ra3NfZmFsc2UwWTATBgcqhkjO
PQIBBggqhkjOPQMBBwNCAAS2EcilTw8qwKebfD4AJDrGKIlZbubmajme5Et3dpll
yhfqpFGdP5i2z3HEXbRnzT9J5TFPHhFYLgVSM9f2KkOTo24wbDAhBgNVHQ4EGgQY
cx3szs2PEMQXjzKq5NY3ypZih4ozDJhqMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMB
Af8ECDAGAQH/AgEBMCMGA1UdIwQcMBqAGHMd7M7NjxDEF48yquTWN8qWYoeKMwyY
ajAKBggqhkjOPQQDAgNJADBGAiEAwtX7m9pskUj/Y+xPT8thR/LlPVrKWxADHR3k
GSn98xMCIQCvDHAHFPn6yJ+9u9/GMMr5vUXRAPKEgGMglDkAxzGhlg==
-----END CERTIFICATE-----`

	testCertDER = `MIIDwTCCAqkCFAFLClhfd7ogpF2ghCmasR6Zp5nrMA0GCSqGSIb3DQEBCwUAMIGc
MQswCQYDVQQGEwJLUjETMBEGA1UECAwKU29tZS1TdGF0ZTEQMA4GA1UEBwwHSW5j
aGVvbjEWMBQGA1UECgwNTG9ra3MzMDcgSW5jLjEcMBoGA1UECwwTUmVzZWFyY2gg
ZGVwYXJ0bWVudDERMA8GA1UEAwwIbG9ra3MuaW8xHTAbBgkqhkiG9w0BCQEWDmNh
dGh5QGxva2tzLmlvMB4XDTE5MTIxMjAxNDEzMloXDTI5MTIwOTAxNDEzMlowgZwx
CzAJBgNVBAYTAktSMRMwEQYDVQQIDApTb21lLVN0YXRlMRAwDgYDVQQHDAdJbmNo
ZW9uMRYwFAYDVQQKDA1Mb2trczMwNyBJbmMuMRwwGgYDVQQLDBNSZXNlYXJjaCBk
ZXBhcnRtZW50MREwDwYDVQQDDAhsb2trcy5pbzEdMBsGCSqGSIb3DQEJARYOY2F0
aHlAbG9ra3MuaW8wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtz+uv
XFJ9H99V+JV3XVGIxo4YCS8sQpHgQhfWbFbbcrKWRdMQS9Ftu2R45UBRpF9JTuoR
FqL7XcT7ZlgahdcxLX2W7usOIl4gvRgqMjQytaam9GJ1inkSYaRAgkDYNrsbxmBO
XGcIak+D6rWa9KpFxueDespHinkKjocXEeEhYTZFvMB4lxAMdVqMo8X+Y9QOlKSo
rjhlAOKByyJY67gjXBht7OJQvs1VoA84dgncZE4HIB6xaCk/rUnd8SyOoZJbNx+7
YOsimljnaH86HYdg5EYxoanjacN/yiYdgkbjAgR/rORAOXW4UZXuE7tzgGzID5Un
4j9pD98ymZdV8y47AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAE2UNkvHbYQo90xz
1t/eZcs1TfE3XlehU2tXFHFPmZTtlHtFn8D3dE/sWIjFSCLX0bo1TpRLAat80yd3
c8K4M/J9RAYqIbJa4dlrsoBO9kT1RQRvz9FYQjaGZORB+5IpofwEJ3AVikaTWGxO
vzL0+0oDz/zAalPEHouxEYqdpV++g3yGh+RWBkDKnge8XryIyfWU4rp1FkcZg7nc
6TKlAdPRPH90hOISZm8Uqi3mnihz0hxm35aBBVa3ikIoO68P5xLSelwJMwulGKpo
AJl1K2H5gcWCOhit6eK7TaGA5diBtiW6584yV21KAtBmLZvsUmYz3RKJdRtsVPd0
rwHPeG0=`

	testRSASecret = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEArc/rr1xSfR/fVfiVd11RiMaOGAkvLEKR4EIX1mxW23KylkXT
EEvRbbtkeOVAUaRfSU7qERai+13E+2ZYGoXXMS19lu7rDiJeIL0YKjI0MrWmpvRi
dYp5EmGkQIJA2Da7G8ZgTlxnCGpPg+q1mvSqRcbng3rKR4p5Co6HFxHhIWE2RbzA
eJcQDHVajKPF/mPUDpSkqK44ZQDigcsiWOu4I1wYbeziUL7NVaAPOHYJ3GROByAe
sWgpP61J3fEsjqGSWzcfu2DrIppY52h/Oh2HYORGMaGp42nDf8omHYJG4wIEf6zk
QDl1uFGV7hO7c4BsyA+VJ+I/aQ/fMpmXVfMuOwIDAQABAoIBAEEk463nCVeyQW+k
asjBJhUAbjNeBdst3CTUQMDx/B/lvj+KThAWipj5GjAhpFi1Ja2BMUNUW6OgwnqD
IlRWg4///8k5JMo4YVUd6leuV1gGMz7x65EoZDZaeEhhSVKAOOKxGFlmUouZ5NGJ
f6VjPApJAW6K8BOeN72YK9OetOVgPupN68nob+iFnqaEOCvqGgIAtUhn/9+bH/Ci
7QqaTs2trP/yS6h81y7XMFA98sg5ai2SpNG//x6eBBRw4e0N3AqWbcILrREaVT3R
kppXn7sR9Yi5zYUmzcWys+pcQxL1N1DyzISzfkB3PkbPcqBLEY5fjw0qOuW/VItf
C9qAt3kCgYEA1Cr0Ht3BtJL2BXd5FDiw8aaqHI5CIH0vZZalTO7oS7GHB9jzSKTg
jE88N0cYwLDF/dA5QBuRG5gH1vzck1uuyogLkZRFTyVxREaslhgL71aH48rTYwHQ
qIBJTbaX3Gn1eXLkdvLp3WSPSUKYDM9VKTsA/9dqiY/AZPzARMipiX0CgYEA0bhu
ankoMdlNkIzbP9iUtocHgJkZ1AUCx3iMAXE7eEjOOy5q7L0HP3RcaTE35rXTK6H7
LZ9+L007gqyq7z7E9Af2yMBx4gnnEsVpX3+oVEH/aYicn4BjE+KI/TO7AuUaZtzG
89uIKBysL7IVT51TJq/4VoCy0wVSKcvBYWB/ZBcCgYEA0WgDwmNFSKC0Sfj9fEPo
ANpqk/ykr8Re/3mMdT5n8C2sBMbQeCajqliaKkT13VmcUUMu/mM2+XE4a6zvWFHp
VuSn5mvdbJycCrNmrE3XmcZiISaTNOkZtPXJY/aQNHAwZEpNzEk9IdKaycf8osgQ
Wb1u4xUOhe9oCUSd0EXtb1ECgYEAlFjXke0953UFDtj0RgdXun1tayPhRz58Jsk4
j9Se8ojdiLNe4zMbK2GN9MLh/gpj45ti53TId4E0NU0aZL2L5+qyQHMQm4nRsE+A
KBNO6Lr+hpIh6BmS+//kUucCxBt3P4ewG9MQTv9pNRvNQ1HP/a/ABMBovignZHVC
xzTRJ+UCgYBFHuPYZbpARSRXbNvjUJBTuQpzlIABr0pt5Y4GogLv2jISKIvkPeW6
RGwWaZgw5eDnBt+rRFCBCRJ/iEU8mwmim9qa8lEpcdgQQbUn+JvTgTE7lcVdkFwu
cmvDeTkSmhZJbovElEzF9jDsXnp3FFbFuTLqzTKseb0WVpNux038Zg==
-----END RSA PRIVATE KEY-----`

	testRSASecretDer = `MIIEpAIBAAKCAQEArc/rr1xSfR/fVfiVd11RiMaOGAkvLEKR4EIX1mxW23KylkXT
EEvRbbtkeOVAUaRfSU7qERai+13E+2ZYGoXXMS19lu7rDiJeIL0YKjI0MrWmpvRi
dYp5EmGkQIJA2Da7G8ZgTlxnCGpPg+q1mvSqRcbng3rKR4p5Co6HFxHhIWE2RbzA
eJcQDHVajKPF/mPUDpSkqK44ZQDigcsiWOu4I1wYbeziUL7NVaAPOHYJ3GROByAe
sWgpP61J3fEsjqGSWzcfu2DrIppY52h/Oh2HYORGMaGp42nDf8omHYJG4wIEf6zk
QDl1uFGV7hO7c4BsyA+VJ+I/aQ/fMpmXVfMuOwIDAQABAoIBAEEk463nCVeyQW+k
asjBJhUAbjNeBdst3CTUQMDx/B/lvj+KThAWipj5GjAhpFi1Ja2BMUNUW6OgwnqD
IlRWg4///8k5JMo4YVUd6leuV1gGMz7x65EoZDZaeEhhSVKAOOKxGFlmUouZ5NGJ
f6VjPApJAW6K8BOeN72YK9OetOVgPupN68nob+iFnqaEOCvqGgIAtUhn/9+bH/Ci
7QqaTs2trP/yS6h81y7XMFA98sg5ai2SpNG//x6eBBRw4e0N3AqWbcILrREaVT3R
kppXn7sR9Yi5zYUmzcWys+pcQxL1N1DyzISzfkB3PkbPcqBLEY5fjw0qOuW/VItf
C9qAt3kCgYEA1Cr0Ht3BtJL2BXd5FDiw8aaqHI5CIH0vZZalTO7oS7GHB9jzSKTg
jE88N0cYwLDF/dA5QBuRG5gH1vzck1uuyogLkZRFTyVxREaslhgL71aH48rTYwHQ
qIBJTbaX3Gn1eXLkdvLp3WSPSUKYDM9VKTsA/9dqiY/AZPzARMipiX0CgYEA0bhu
ankoMdlNkIzbP9iUtocHgJkZ1AUCx3iMAXE7eEjOOy5q7L0HP3RcaTE35rXTK6H7
LZ9+L007gqyq7z7E9Af2yMBx4gnnEsVpX3+oVEH/aYicn4BjE+KI/TO7AuUaZtzG
89uIKBysL7IVT51TJq/4VoCy0wVSKcvBYWB/ZBcCgYEA0WgDwmNFSKC0Sfj9fEPo
ANpqk/ykr8Re/3mMdT5n8C2sBMbQeCajqliaKkT13VmcUUMu/mM2+XE4a6zvWFHp
VuSn5mvdbJycCrNmrE3XmcZiISaTNOkZtPXJY/aQNHAwZEpNzEk9IdKaycf8osgQ
Wb1u4xUOhe9oCUSd0EXtb1ECgYEAlFjXke0953UFDtj0RgdXun1tayPhRz58Jsk4
j9Se8ojdiLNe4zMbK2GN9MLh/gpj45ti53TId4E0NU0aZL2L5+qyQHMQm4nRsE+A
KBNO6Lr+hpIh6BmS+//kUucCxBt3P4ewG9MQTv9pNRvNQ1HP/a/ABMBovignZHVC
xzTRJ+UCgYBFHuPYZbpARSRXbNvjUJBTuQpzlIABr0pt5Y4GogLv2jISKIvkPeW6
RGwWaZgw5eDnBt+rRFCBCRJ/iEU8mwmim9qa8lEpcdgQQbUn+JvTgTE7lcVdkFwu
cmvDeTkSmhZJbovElEzF9jDsXnp3FFbFuTLqzTKseb0WVpNux038Zg==`

	// ASSUMPTION: this type of encrypted key file must be sent with pem phrases and header.
	// otherwise this cannot be stand-alone like DER format. Becasue the content is encrypted and
	// does not contain any information about encryption method.
	testRSASecretEnc = `-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-128-CBC,BEA59B5FC69611C81EB37DE7CBF052A8

S0JBZxoxBbWaDdrwl3itu1bhHOipbNigeXZsiUjU9ZVWblZyfoxBk14nY8IfpG4k
W9AjAmBOu8hy5xhcFyg+KnnPx08ZNnG6XRgk8H6Csf/T/+OfDkjB+TdQyCuAkIbN
ISR562O4bg3Nl7ZTI1frfkPu1tw7mXvJx9DBIPu0TvXAm3NIvfwMK/KS/sUSSq90
zwelB/ZwUGH7kHjVDh6Z9Y3i+M1KXAEHX7BJyOQvSJIOB+rKQ3po1iA2N689QJuO
dqYqaxtSuQfA6aUC2aU9N3U1sLPvLRzEcAIABBehguXeuq7syb7nG4W6cpuAIdq7
xSEKoSSgm1qGWcrX1aENT5/YdLu/wvYjhypWb1Ioa1jt5EjBeeTql89eQd26J9gu
4iH+8rpT9mrxmQc+eXV8u04POIcXsQuyLrAoHFeGe3Db+ivY3nrDuTou7tmXfSeu
ZtgtTkO+Z8iIgRqVcU3lbJZUA68N+/FJJqA82PYu3O891ogE+uskKYVz9K6BfaAc
1tXZ50OnOp2GeDwQJA7fs4ZtRHgOol7pYg9brgH3lFXLi8qSgneE+j3dmX0FEp78
wZJhgcqZXY5ek/D7nhNHEWP8SPi3MX8wOuSy8KeNf+Qop9lFjmbOFq5Zv2ZqKyVM
w1uJihsjKF0vL7GSxbnbsUQ2/ipyuqUsOga7becdGRgepLk0Xt2eZLwVQwhP/UZY
HuTrdcRJgsD1x3CfuxfkJ6Y3XepcSVEMOy46br7weOvJf4NQ+cHaB1DhOnXLtMex
s5PHfxz7KsjN/8U5bw0lriuicdHuCKBgjC2tz2DMoujZDGiR0vGGm0ipok72faQ2
jQ/g5OKhgdY4PhWAjiniy+iluEnHtQyNFeaX6Ky70n6cuQSKUsPi+LadAH6kJMtM
tAabGksnjEu+xeMQazpnGP60IivGqYtlVkUl0zPdlW2tEWJHr5i1NTVj/pnAWN5c
kZkcIVhD/bZgTNSqVuhCYQtpjvpeoKyzVd7Dds1nCm4XxfgOc7UGL+lCUolKo38f
GfnS/iaaEACd6WpF8snvarrDQUOJkyHdN0PPPFS1HC4sfD2Na1i7ghOwJu/7kyEh
Qd84Ww5alM9Hdb4b6feyT5g//MoP4omoA4A3mxBLqB+oCrRWca8/mPmbkjSepIP9
1KcKdFrwVe9yCEZNh6pb7KiMBdtGJ455ZanUKtwXtLFyz/gWpTuSLn2+cdzlnmBJ
Yl/0qHMTzPBPn0Y6OXL2XdHylbsgg6KgP7heAyxcrcK+P9bcfC62+fGaLMssyfSS
NEy9oqm3x74M3Yyo7qUmQn0biIaL8wlhCTa2oPo9HUeiZFhWHv+WFgpHMpZrmlkh
uu6X6pAe1H/smU5r/bFV/XZKDA80oDS3lnBceSiQiPMtY8GgbZleKgtDKYE4St9O
YoCO6XP3ID9p6kHfNP5drrieI/oqQaf+KJASiM6oWg51UNioA6JgMEHUbR8zZeM7
eqB0axk3riB8WkpORo2Cnaw8aLmPuyD0YspnVVhqJ3GqUREnjWkR6Qeqk2XlfUDc
LA75ztrfrDvq+hzIHtyaWnpNjzZXu6Da6IlgLbWdtaqFPMdEHeb7LMpWJafUgFk+
-----END RSA PRIVATE KEY-----`

	testRSASecretEncPkcs8 = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIFLTBXBgkqhkiG9w0BBQ0wSjApBgkqhkiG9w0BBQwwHAQI3hhFTiCsN5QCAggA
MAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAECBBDJ3UV5TjKjfp5wxEN0gZdnBIIE
0L1+YE23Fb/Mdi03yGtcWuouMtDpLVVbhNyKMM/OQTYezayXOoo8rjhzwU9VIGdu
LmNj7/fzsb6hxuJ6SZ5pb22ATWLGaxGZpT3UZHW7teenBNVhpJFQPebz4oqnXt4u
mFB/nEjdd7S5rUQF7WL7y1lEpnf1DtvcFvXIMsztAWfkHqXOCYRi/pjgyZOLQlmb
JN2wsfRlcJJ3raxdLoN209OZm467r0lKzTb1UKpcJ+vPAEktF2HdGHFaotjsH2xC
IgjkHJhxZIXdS9EVBYHCV5CQOKNCDyx8GtbIHb57fwUc1BR6AgPnFF+//RpN6pp8
C6JZo784TZYzSnQVY33vzxwJ1OQIL0RgLv7oSAZgl1UZSLW56qQgXrJW8IoD4UCb
bgRLk+ljJPuyO9fyDHR2VlTl3eQMNg4/lnEls3vUaphjamVdrynq9d7xa5q6FI3k
QcQUXrVhx3ZqI16xF/zjmiIRk0tzo9yaPV2UIhJn3pCkcWWQfEyMeZ1rLNQtLwtp
wQqLUsvhmW1BX8ZN0z46rPHs0lAijAttjMArVCdm/SnZw2JJ5jkrc7ba3zqOnzJD
hBFam2Wsg/C87/4YzjKqtjULFb1XqYKCTnSsGb7N3vHkTBa1HEafYMN3kMk9UArx
UBwFCryYj1/WJsGqX7Stiuu4HEXcv2m0EFmgmRUxE+iKCQM+ASS//MjciOGBgAdE
6qyh1c7lCkxixSPE4omtpArzueQ7cB0e2Ylt7q01Wuy/g8p93kQ1E8lqiYLQM2VL
L9V4D0Pxf2VaVGi1dUYW9OtVMRgZaL6zqNfrCIeWkO/W428UhxTgIpYrxZBe4ucu
9YAYm62BwGio9Nd/TacD3aPVzxA8Sk/TSdfoSfuAYS7gsxhdcNnS32R09E7sQ0SF
YG9gNbSvZBkFxmcF3/QwA5UAWY2hu5+cVkNiTSPjOuLDP6Trik7wnoii09hATMPO
aRtodHwS6yUl/o4KRKxV05mN4szQRGpfcvrJiTjDur0uIc1Hm3HoEE/xR/Gv3eQx
tNP3DSbFVqZgula3YqAzbwcUuyCTtwT3sfBQrOoXjmqopP9PYHq9ksQq3ua5dUlU
WUvDYmYD6efftWgM78h6xh9yGK2fVqYHYgJnYb6mjrqRjkvClPimEqeh9KWc2RCc
/cZU4ZcwFurcDXGtbNyglJPb3BsrBmiNdJnqg+z9IKsk1S5YeljLzitGBRo/7PRb
J801zRm3IlVFuJEY3FqRr7xTKNAwX+C7el6YNYpGR+LzU4MYNtU4iq1avy8LqeT7
l/xeGCJTPWSiS3dwhKKzm07Xkv8rYrot55gEy6U0DxAU/h2SitO9H221AdjKW9GE
zJamFNYgidUup0FLhYoMJvdxTae7ZRAVjLqsHVYFWmRhVOuA2SDiq1+3vgxfzFpV
vJa7X219nPML4I9IVfgv6yeMfCJDVfwBDk1C4v2WEB9DpEf62pJh7iM9Z54RjJ5q
RJCJzav6rcesvyXM02tUVDGv3lZk0/zqy4n5Qqv9RW2DwegxZRh2+U3n0giUNExX
dn/FcnjTwc718L6dktJKlPBV8XPcwpMo/eM8tiOZlbe0KyNL4fjaiUI00+b4zLUL
id2CcBLTfOickzFKIqkHXoWvWCDWvQo0QaCEW30p+c1f
-----END ENCRYPTED PRIVATE KEY-----`

	subCountry      = "kr"
	subOrganization = "lokks"
	subCommonName   = "test-name"
)

func TestPEM_PasrseX509Cert_Sucess(t *testing.T) {
	cert, parseErr := ParseX509Cert(testCertPem)
	assert.Nil(t, parseErr, "Certificate parsing failed")

	assert.Equal(t, subCountry, cert.Subject.Country[0], "They should be equal")
	assert.Equal(t, subOrganization, cert.Subject.Organization[0], "They should be equal")
	assert.Equal(t, subCommonName, cert.Subject.CommonName, "They should be equal")
}

func TestPEM_PasrseX509Cert_Fail(t *testing.T) {
	_, parseErr := ParseX509Cert("this is not a pem format")

	assert.NotNil(t, parseErr, "This case must make error but no error")
}

func TestPEM_VerifyCert_True(t *testing.T) {
	check := VerifyCert(testCertPem, testTrueCAcertPem)

	assert.True(t, check, "Verification must succeed")
}

func TestPEM_VerifyCert_False(t *testing.T) {
	check := VerifyCert(testCertPem, testFalseCAcertPem)

	assert.False(t, check, "Verification must fail")
}

func TestPEM_GetCertificateB64_Success(t *testing.T) {
	cert, parseErr := GetCertificateB64(testCertPem)

	assert.Nil(t, parseErr, "[PEM]: Certificate parsing failed")

	assert.Equal(t, subCountry, cert.Subject.Country[0], "[PEM]: They should be equal")
	assert.Equal(t, subOrganization, cert.Subject.Organization[0], "[PEM]: They should be equal")
	assert.Equal(t, subCommonName, cert.Subject.CommonName, "[PEM]: They should be equal")

	cert, parseErr = GetCertificateB64(testCertDER)
	assert.Nil(t, parseErr, "[DER]: Certificate parsing failed")
	assert.Equal(t, "lokks.io", cert.Subject.CommonName, "[PEM]: They should be equal")
}

func TestPEM_GetPrivateKey(t *testing.T) {
	sKey, err := GetPrivateKey(testRSASecret, "")
	assert.Nil(t, err, "[PrivKey]: Private key pem parsing failed")

	sKey, err = GetPrivateKey(testRSASecretDer, "")
	assert.Nil(t, err, "[PrivKey]: Private key der parsing failed")

	sKey, err = GetPrivateKey(testRSASecretEnc, "password")
	assert.Nil(t, err, "[PrivKey]: Private key pem with password parsing failed")

	testBlock, _ := DecodePEM(testRSASecretEncPkcs8)
	fmt.Println(testBlock.Type, "/", testBlock.Headers)
	fmt.Println(base64.StdEncoding.EncodeToString(testBlock.Bytes))
	testKey, _ := pkcs8.ParsePKCS8PrivateKey(testBlock.Bytes, []byte("password"))
	fmt.Println(x509.IsEncryptedPEMBlock(testBlock))
	fmt.Println(reflect.TypeOf(testKey))

	sKey, err = GetPrivateKey(testRSASecretEncPkcs8, "password")
	fmt.Println(sKey.(*rsa.PrivateKey).D)
	assert.Nil(t, err, "[PrivKey]: Private key der with password parsing failed")

}
