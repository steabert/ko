package ko

import (
	"crypto/tls"
	"net/http"
)

// Certificates generated with
//   openssl req  -nodes -new -x509 -keyout srv.key -out srv.cert -subj '/CN=ko' -days 3650
// then copy src.cert to certPEMBlock, and srv.key to keyPEMPlock
const (
	certPEMBlock = `-----BEGIN CERTIFICATE-----
MIIC+zCCAeOgAwIBAgIUFsZi34N+TpVBwadUD6h2GmRAYtUwDQYJKoZIhvcNAQEL
BQAwDTELMAkGA1UEAwwCa28wHhcNMjEwMzA2MDgzMDU1WhcNMzEwMzA0MDgzMDU1
WjANMQswCQYDVQQDDAJrbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AM7h2koiNohXBFNCm3x4ZumAx0obhajuY5lmvqU7LzQRxEkF+T5szN6CqBKoRTfF
RCOHByeJnbxh7gDH/WbkzN+DsibvZElr8oVLdDixucIe+IkmSpOWgSdLfNF2R4y0
EOy69Lw9VdTKK1HzKO3edls/6kBLHFtf8BOGEeha+YaJ4RMS7uqY3n+fx2/gO2ld
4Z46i8U6J473DqKuBRcPvXOG0EJwjpkZpQ+qVUbYrungd0bcnO+TRmf1zE1FGOP6
CKYrDmKIDQW/UvLO8P1CAtghGn3tM2L334lkYT2e4lQFlhlUHPmMDiVszfAi+BDI
HYLoYHUkDlII3arpvxWCS/sCAwEAAaNTMFEwHQYDVR0OBBYEFDfgu44E60n6KyVn
LkqvN1jb0+vwMB8GA1UdIwQYMBaAFDfgu44E60n6KyVnLkqvN1jb0+vwMA8GA1Ud
EwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAMguMOH/Qfj6hhQ+4L1dwMTK
k4kcGxwXWfD71htPIWkgLmdrJVYP94erLtvKffYYn5lyEZC9YyX3XbdcT2C2XaPr
Z+tvkKN14Sxu9LSxFOPWhkxoYj+ld6tGoKbjbl4XW+03VdCeWvqdIpnj5mxCBUmc
WAhZY8+6MSjer45HrTuNv63N9kZyaPCmdgwFLYOcpBA4fjbRoaFRpO0SaQls1eV6
z4vrKVnv9HNXcs0IvVXgvp6FWJ1Nc6BuxJuhWBkGgsOjopMvATN3glGQus2IjPbO
KA8w3ksjySeTKU+/wW/CwRTrm8dyv8vYsb9G/G8F+0laABSPQJCx8VJXSdicmkQ=
-----END CERTIFICATE-----
`
	keyPEMBlock = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDO4dpKIjaIVwRT
Qpt8eGbpgMdKG4Wo7mOZZr6lOy80EcRJBfk+bMzegqgSqEU3xUQjhwcniZ28Ye4A
x/1m5Mzfg7Im72RJa/KFS3Q4sbnCHviJJkqTloEnS3zRdkeMtBDsuvS8PVXUyitR
8yjt3nZbP+pASxxbX/AThhHoWvmGieETEu7qmN5/n8dv4DtpXeGeOovFOieO9w6i
rgUXD71zhtBCcI6ZGaUPqlVG2K7p4HdG3Jzvk0Zn9cxNRRjj+gimKw5iiA0Fv1Ly
zvD9QgLYIRp97TNi99+JZGE9nuJUBZYZVBz5jA4lbM3wIvgQyB2C6GB1JA5SCN2q
6b8Vgkv7AgMBAAECggEBAK1oiIV9OgJ8FccIXLYvYeu1otZOTXG1KE0L3x82hbF6
dvHSjQGzRuH32JOS8jn2ItA4vVl5s3qVB18mQxQ9EjED/Y8/N+uHDQiHn4pqBk9d
kGu9aeNd0zIxxxT3tK+Ou2UCrGMgclJjh34weI0x3DlOULbFfqZkuyJSTa5amy8V
3wW3uYzuF5M5FjCFOFPx1+/hphXqI6yn5ZHzwrPnm2jzY9voTLcQR/ZXoeqxgivl
Pwnvqccrv3B7o1vVzBpdWjILB48q1SukVVV69FvvUUwo7Qxaho7scc4vh9syprTH
WpdtEnUFkcwIkBPh1FwLRhk+yfGl6yg3Zq3Qjw3uHHkCgYEA8gEkqptyUWqB8eaL
oWbD6s15y3TBhG477Cdu0GzbbLkigKosVwYhDXgalSQck5ntwf12alfDHgXeSriN
eGcpb1cZtDn6I6/mSwjJcFgQDQ545txDyBHL8ypTpP1mXWNtMSVJlZoKfpv1M2xD
6iL0W7ilirpcaAc3Phjexj+X1Q8CgYEA2ti6SRQEyRJyo5EUGaTAyMmnFmjJiAGd
30IoWoV6+OAp86um+nCKlmVnyuY3FM4DtExDEMd65pph3YxXlt48EzDC5P/fmVOR
C2kBrV619uiawrjp455ePuRB86chTqMMhU0IyCGOwS8mgYBrI7rpSo0cM2S0ymKS
9TJHYhIKklUCgYEAp/DzKRJG+wkWtHBxZciTHVcKto6H3QdCvld/J1Tj0UeJEhEG
RD4Uoew/RlCRJD0mKgFjM9lDpoocAW6hfnTY5FNlmxTA6hMfleK7KCN0wBrS/CLP
RwBSsKUm9tCDQTvGgtyFfDQyJDrGprDzUICBY0V4XBWIGwkm5QkNUDbBfzECgYEA
sQWmif5bcJoviQeNjsCqAMC9G29ftVg5b6KAKdjXBAGvbZ9nziTCta0JLCLUY0vR
y0H07dmuHGK8zwz6vNq/FXbX74zaPZhPNz+VT7vQzQySQvh4mNo9ufnBL2n9kzJo
qlsJw3kBlFqjdxV9lMVYeCl0qk3Hv/3EifpFq5qUWpECgYB1DSRo+NT+Pq/VG+vs
ncrNm2Aoz7XSnXl/OCCk6IexBdeu4ijIEYYV07AfiQ0lpATHflhSiAmxZbfKLFdC
M6qXSXy/Trt6w5yrXzzF+WSD7WxHh0B0EVgNuPth4FGwBds+qQg+Jji7t42XWLJI
aLT29qBp8Ca356pZ465NoNBLQw==
-----END PRIVATE KEY-----
`
)

// ListenAndServeTLS starts a server with self-signed certs
func ListenAndServeTLS(addr string, handler http.Handler) error {
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)

	var err error
	config.Certificates[0], err = tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
	if err != nil {
		return err
	}

	server := &http.Server{Addr: addr, Handler: handler, TLSConfig: config}
	return server.ListenAndServeTLS("", "")
}
