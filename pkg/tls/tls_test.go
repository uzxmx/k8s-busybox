package tls

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var certPem = `-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgIIE31FZVaPXTUwDQYJKoZIhvcNAQEFBQAwSTELMAkGA1UE
BhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2dsZSBJbnRl
cm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMTI5MTMyNzQzWhcNMTQwNTI5MDAwMDAw
WjBpMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
TW91bnRhaW4gVmlldzETMBEGA1UECgwKR29vZ2xlIEluYzEYMBYGA1UEAwwPbWFp
bC5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfRrObuSW5T7q
5CnSEqefEmtH4CCv6+5EckuriNr1CjfVvqzwfAhopXkLrq45EQm8vkmf7W96XJhC
7ZM0dYi1/qOCAU8wggFLMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAa
BgNVHREEEzARgg9tYWlsLmdvb2dsZS5jb20wCwYDVR0PBAQDAgeAMGgGCCsGAQUF
BwEBBFwwWjArBggrBgEFBQcwAoYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcy
LmNydDArBggrBgEFBQcwAYYfaHR0cDovL2NsaWVudHMxLmdvb2dsZS5jb20vb2Nz
cDAdBgNVHQ4EFgQUiJxtimAuTfwb+aUtBn5UYKreKvMwDAYDVR0TAQH/BAIwADAf
BgNVHSMEGDAWgBRK3QYWG7z2aLV29YG2u2IaulqBLzAXBgNVHSAEEDAOMAwGCisG
AQQB1nkCBQEwMAYDVR0fBCkwJzAloCOgIYYfaHR0cDovL3BraS5nb29nbGUuY29t
L0dJQUcyLmNybDANBgkqhkiG9w0BAQUFAAOCAQEAH6RYHxHdcGpMpFE3oxDoFnP+
gtuBCHan2yE2GRbJ2Cw8Lw0MmuKqHlf9RSeYfd3BXeKkj1qO6TVKwCh+0HdZk283
TZZyzmEOyclm3UGFYe82P/iDFt+CeQ3NpmBg+GoaVCuWAARJN/KfglbLyyYygcQq
0SgeDh8dRKUiaW3HQSoYvTvdTuqzwK4CXsr3b5/dAOY8uMuG/IAR3FgwTbZ1dtoW
RvOTa8hYiU6A475WuZKyEHcwnGYe57u2I2KbMgcKjPniocj4QzgYsVAVKW3IwaOh
yE+vPxsiUkvQHdO2fojCkY8jg70jxM+gu59tPDNbw3Uh/2Ij310FgTHsnGQMyA==
-----END CERTIFICATE-----`

var secret = `
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUR1akNDQXFLZ0F3SUJBZ0lJRTMxRlpWYVBYVFV3RFFZSktvWklodmNOQVFFRkJRQXdTVEVMTUFrR0ExVUUKQmhNQ1ZWTXhFekFSQmdOVkJBb1RDa2R2YjJkc1pTQkpibU14SlRBakJnTlZCQU1USEVkdmIyZHNaU0JKYm5SbApjbTVsZENCQmRYUm9iM0pwZEhrZ1J6SXdIaGNOTVRRd01USTVNVE15TnpReldoY05NVFF3TlRJNU1EQXdNREF3CldqQnBNUXN3Q1FZRFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0F3S1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ3d04KVFc5MWJuUmhhVzRnVm1sbGR6RVRNQkVHQTFVRUNnd0tSMjl2WjJ4bElFbHVZekVZTUJZR0ExVUVBd3dQYldGcApiQzVuYjI5bmJHVXVZMjl0TUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFZlJyT2J1U1c1VDdxCjVDblNFcWVmRW10SDRDQ3Y2KzVFY2t1cmlOcjFDamZWdnF6d2ZBaG9wWGtMcnE0NUVRbTh2a21mN1c5NlhKaEMKN1pNMGRZaTEvcU9DQVU4d2dnRkxNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUZCUWNEQWpBYQpCZ05WSFJFRUV6QVJnZzl0WVdsc0xtZHZiMmRzWlM1amIyMHdDd1lEVlIwUEJBUURBZ2VBTUdnR0NDc0dBUVVGCkJ3RUJCRnd3V2pBckJnZ3JCZ0VGQlFjd0FvWWZhSFIwY0RvdkwzQnJhUzVuYjI5bmJHVXVZMjl0TDBkSlFVY3kKTG1OeWREQXJCZ2dyQmdFRkJRY3dBWVlmYUhSMGNEb3ZMMk5zYVdWdWRITXhMbWR2YjJkc1pTNWpiMjB2YjJOegpjREFkQmdOVkhRNEVGZ1FVaUp4dGltQXVUZndiK2FVdEJuNVVZS3JlS3ZNd0RBWURWUjBUQVFIL0JBSXdBREFmCkJnTlZIU01FR0RBV2dCUkszUVlXRzd6MmFMVjI5WUcydTJJYXVscUJMekFYQmdOVkhTQUVFREFPTUF3R0Npc0cKQVFRQjFua0NCUUV3TUFZRFZSMGZCQ2t3SnpBbG9DT2dJWVlmYUhSMGNEb3ZMM0JyYVM1bmIyOW5iR1V1WTI5dApMMGRKUVVjeUxtTnliREFOQmdrcWhraUc5dzBCQVFVRkFBT0NBUUVBSDZSWUh4SGRjR3BNcEZFM294RG9GblArCmd0dUJDSGFuMnlFMkdSYkoyQ3c4THcwTW11S3FIbGY5UlNlWWZkM0JYZUtrajFxTzZUVkt3Q2grMEhkWmsyODMKVFpaeXptRU95Y2xtM1VHRlllODJQL2lERnQrQ2VRM05wbUJnK0dvYVZDdVdBQVJKTi9LZmdsYkx5eVl5Z2NRcQowU2dlRGg4ZFJLVWlhVzNIUVNvWXZUdmRUdXF6d0s0Q1hzcjNiNS9kQU9ZOHVNdUcvSUFSM0Znd1RiWjFkdG9XClJ2T1RhOGhZaVU2QTQ3NVd1Wkt5RUhjd25HWWU1N3UySTJLYk1nY0tqUG5pb2NqNFF6Z1lzVkFWS1czSXdhT2gKeUUrdlB4c2lVa3ZRSGRPMmZvakNrWThqZzcwanhNK2d1NTl0UEROYnczVWgvMklqMzEwRmdUSHNuR1FNeUE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  tls.key: ""
kind: Secret
metadata:
  name: test
type: Opaque
`

var expectedOutput = `Subject: CN=mail.google.com,O=Google Inc,L=Mountain View,ST=California,C=US
Issur CommonName: Google Internet Authority G2
Subject CommonName: mail.google.com
DNSNames: [mail.google.com]
EmailAddresses: []
IPAddresses: []
NotBefore: 2014-01-29 13:27:43 +0000 UTC
NotAfter: 2014-05-29 00:00:00 +0000 UTC
`

func TestRunFromPemFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "cert")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.WriteString(certPem)

	c := NewController()
	buf := &bytes.Buffer{}
	c.writer = buf
	c.fromPemFile = file.Name()
	if err = c.Run(); err != nil {
		t.Fatal(err)
	}

	if buf.String() != expectedOutput {
		t.Fatalf("\nExpected:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}

func TestRunFromSecretFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "secret")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.WriteString(secret)

	c := NewController()
	buf := &bytes.Buffer{}
	c.writer = buf
	c.fromSecretFile = file.Name()
	if err = c.Run(); err != nil {
		t.Fatal(err)
	}

	if buf.String() != expectedOutput {
		t.Fatalf("\nExpected:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}
