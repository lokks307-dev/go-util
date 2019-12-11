package pem

import (
	"testing"

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
