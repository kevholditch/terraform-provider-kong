package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKongCertificate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testCreateCertificateConfig, testCert1, testKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongCertificateExists("kong_certificate.certificate"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "certificate", testCert1+"\n"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "private_key", testKey1+"\n"),
				),
			},
			{
				Config: fmt.Sprintf(testUpdateCertificateConfig, testCert2, testKey2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongCertificateExists("kong_certificate.certificate"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "certificate", testCert2+"\n"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "private_key", testKey2+"\n"),
				),
			},
		},
	})
}

func TestAccKongCertificateImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongCertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testCreateCertificateConfig, testCert1, testKey1),
			},

			resource.TestStep{
				ResourceName:      "kong_certificate.certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongCertificateDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient

	certificates := getResourcesByType("kong_certificate", state)

	if len(certificates) != 1 {
		return fmt.Errorf("expecting only 1 certificate resource found %v", len(certificates))
	}

	response, err := client.Certificates().GetById(certificates[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get certificate by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("certificate %s still exists, %+v", certificates[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongCertificateExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*config).adminClient.Certificates().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("certificate with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateCertificateConfig = `
resource "kong_certificate" "certificate" {
	certificate  = <<EOF
%s
EOF
	private_key =  <<EOF
%s
EOF
}
`
const testUpdateCertificateConfig = `
resource "kong_certificate" "certificate" {
	certificate  = <<EOF
%s
EOF
	private_key = <<EOF
%s
EOF
}
`

const (
	testKey1 = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDXi8zFDavAN7fl
RJO2G4oLj6NIT86BJnzM3XqtGl6pvfp0bo9so+h/0HhGtnIh7Je4BL7PGsv5BSdg
6EDDZDXZn/ZDe3jje+Ee1sfn98H+1mjTDlm0U2it/cWaZ+a8GEhPwidWyI1AeS4O
XrM5VcbcmVoIPRr5z0iJBJ2LRY0L/rPVsJOGbT1WPsFjFMZrc7GAixrjjk0jovKr
X7hK/Sj3vRGjZsP6CDEVIeEJQOMDvup9YCbgPTofn3gH9PoKFQdA0CbuGsdKiUOF
bfrp2MFP8933WsCGhUbpsdDd/Vt9JvpFlR0aCShv0wIWhvzv8WQ5MUWFQvKTmXDv
tdnRqXTvAgMBAAECggEARP648zKnEYZEVR0Ycyhpjb3StGjnXyvksucKR7KzLn5j
VzW0rz/gQlmGxovMCNPk1MCgG0cml3Vw33I4mNLQ8fJkL8GsNpUGwIpbvwLtlBcp
wrVLPY+daGRdBknP79GOBAnP8dWMcWDYvzzM/cNZPm/QA+cbZW9WdpWFoHkI5xdw
LIuI4q7jsjbRHEK5Mc055EDC5hwkmZrjVi9hsqk2De5AYaXhX1kqPopOuKvxJjff
H0QGvE285HFi0qdDYcEzpFEgzBUMjU144bqyl+s7VLhuJ302VVsUJZbNq0m1h9xN
/KSqK6Q3d3ZgsP9fWpq/H9N1NcxpmJbxyET+3P4lEQKBgQDvHzYRpn2bCgPrZaF7
oCbFcdBPdRc+GSTuNzZLYl5e/dGLD/+soxOobZzH9HmhpMkOwpZX4ajz7/5y3GTH
OJ9nyuh52b1weOAWOodmUrMqwRcxNNtkUE6LFMVdatgGpC8ihvoMfDM3GK91Hdbr
bs1Ws7tfIkftad3nwM4YO62oQwKBgQDmwpZ1sXgDUxswWq9vvTy4TDOkRdcEO2J/
8Aphut8RJ/itT2cm5WKWeNkQvyC1Gs4PMGnqU2LWG+31YPB1GAk3a+OP7GN7WewS
rMBnmgdoCjDqXOBg15KOn2m23Ac6a6dBGI9pvtQfhX0e59izcSdA3ZwboFfYdmMI
hYn49uO75QKBgQCXYQnouJ7R1NBQaKGHUwbYfkni03yoWmCv0hI0PQ0DU+ohADra
/s5GFUZoq5OIynpiNrvY3MoJzAgojO/b0zPPEHyGD1tHZa5vRBRNqdM1INJe21h8
s/5VPAwKLMafxbb1Q7/uwX3mxmDlYsOZfibOWbAn9NrWKOxLeBrA6p7wYwKBgF/x
J31ne+5l7zf7fFWI6GX3yMDUCMHJrvpiYu6fM39+jvX/vXN+i67kL9u2m3Kw4luO
VXsHkGBU3GrZEyCcDbjtMn/0WKhAitZ43MY2VD39frjyRJf/CQAjZ2CPurGfcLqv
63Cb1rYEWjEvU/nHYfqmKPGTiPKGxkYUv3izrZvBAoGAVyi8N8LtYhG+T52tVhXM
0wopUsU5NkZjuTof8lG0oAUv6B0I6osacsB9AFv1Ai8Y6hBpZlLRALKDLu/pLVrj
GGRb4/tQrQz7/cn73RiefOha1TzDSEEfq237/eBy7dqk75C3nfLNtu6vEFuAgQbS
dTLBTUGXDSVpcySM5XwoVuM=
-----END PRIVATE KEY-----`

	testCert1 = `-----BEGIN CERTIFICATE-----
MIIDKjCCAhICCQDV2QFIWXZy0DANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJH
QjENMAsGA1UECAwEQ0FNQjESMBAGA1UEBwwJQ2FtYnJpZGdlMRQwEgYDVQQKDAtr
ZXZob2xkaXRjaDEPMA0GA1UEAwwGZ29rb25nMB4XDTE5MDExNDIxMTgxMVoXDTI5
MDExMTIxMTgxMVowVzELMAkGA1UEBhMCR0IxDTALBgNVBAgMBENBTUIxEjAQBgNV
BAcMCUNhbWJyaWRnZTEUMBIGA1UECgwLa2V2aG9sZGl0Y2gxDzANBgNVBAMMBmdv
a29uZzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANeLzMUNq8A3t+VE
k7YbiguPo0hPzoEmfMzdeq0aXqm9+nRuj2yj6H/QeEa2ciHsl7gEvs8ay/kFJ2Do
QMNkNdmf9kN7eON74R7Wx+f3wf7WaNMOWbRTaK39xZpn5rwYSE/CJ1bIjUB5Lg5e
szlVxtyZWgg9GvnPSIkEnYtFjQv+s9Wwk4ZtPVY+wWMUxmtzsYCLGuOOTSOi8qtf
uEr9KPe9EaNmw/oIMRUh4QlA4wO+6n1gJuA9Oh+feAf0+goVB0DQJu4ax0qJQ4Vt
+unYwU/z3fdawIaFRumx0N39W30m+kWVHRoJKG/TAhaG/O/xZDkxRYVC8pOZcO+1
2dGpdO8CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAP6xjv2nqMb9NmyUPz6bGlLNq
8lqUE4zWK61YS6P3BinRIswwDfUg42eMcafebOBgyc34yLBSbKF9paDupuI/xcyk
ySQk48vSGYAuo0wlN8YAmf6SC7tkfk7PL8uVl8bblDREk+D28UzEMNMA4ScCoYtQ
21G2HUhMonRI+MGRtbaVmc14XXjpPww29W6s5nxuG5MaGWd6wkIL7pmHmVBSN2QK
RQzGLmfi0TxOiCNCb9fArIaxlXYfR/yBoV/NdEKrFdpQg3pxKNKu0+IYJHomJDpZ
+Hr3Nf7YNDiX/eCuG//beQaE2H4A9/K7i15szIbv/inpIkcx7z5eIGULR7Hykw==
-----END CERTIFICATE-----`

	testKey2 = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDIuH4e8oJC3D9I
E1PXHuE+CDyUwM0/UWMbuZXCu1LslE8kY10WN60qmvUc+JcLXKTyo1tfhsglcDEj
ac0tl7WYAF0MbLQRSuO2XTMcCByMXxt4CclNE6irl075y2lZxgvvu+dokpGSDcK7
cQv6rYKwPuHuyHwFBwmN1+jW/aLtRXY6+RyWTFtBjvO1Dis/JBMVGishGwg16vG0
iLoN1ijlVC44zeGjvgnD1RuhVP2MYOLWbRucDSQWT2ihAitm0oB+HXLI34hlSKQk
Km/7fYNg5JhhyRm27V+UPHxQZX0p32nvyzQSImkQjgNqXlGaHpBxPLhrycTYSccz
xskw5AAHAgMBAAECggEAC50hLwk5IEU+JB16LGhNABnZ54GAXrmG0oPadnoAQlAS
hDh7ml1+V4i/xf1cP/wMTz3Ee7KAwinLrhDhHlI436Klv6JKiPWcV2DtJCIUBMs7
+6YVT7BjggovN6TdY6Rh7G0i5poUci75pX7VApwng6sfx4EyK9hMZTio8Eectb3E
7HHpbXDmUf5KLW0d0pY0pynk38wvStdpsIchy38fGUVWs7ICXX9zrJMLFO2XqToP
Bc76zzrZwQ8AOtG5N58HnpRmJxyHR3TbTWzB5NGTUEtdab3aR6S+aU0ZNDUexc4t
ru9kvyWhpq2xUxCjTV5buYAD/sS/kz0ZGPXDE1RSgQKBgQDnvOiFEfL4OqLsnGnE
92FQ+hbgX7TeKDNkM2wyjFpS4lfJFoKXs+zuaTRaGr7nT6kHk1qRdRDxor7nmFZP
/UZNsTvSsKlG90fsZ3wyGZnUn27EDYjK3ZQkw7CYVsecjdSghr9R4Ai8qR11dcyh
KO503h3I1m52TIBI4v2P9etEVQKBgQDdvEDNkbG07K/L+eew1P2twomAHA/P7Nv4
/DH504PDFQyZc8x/+n7vVdqRdcxTWK+9BrfDhahq5XinQcnUTOgQZXQ1dgFvm6F6
A216Q/NWDj3Qvm1nUmSz7O8Y0K71ea68pKct5/K7xonwMxYewQIgcexWSBMMHdWf
ZE5aGhEu6wKBgALmRjKECvo4IZm8S0Z/oeQVfPvJtsWg0yPJ0OzA/NBUrKyDb5be
VXSWNGV8OC9Eu+SwX026nT+ovaLAMTRHAv4t3YXBWBzfMzMuCEvILjiO1h/122RO
aXAcUrVVQKIg1Cw+A17O4s0ZgJjbpHfPv0wPC2hb5n3sbx56WJnYhd0RAoGAGfGQ
03ycgkK/Pup6hWImXFJKrEacQwO/qR446rYo0IaB9uJppv+0ImS1MhfKVGYoCcHj
dmkJD5jRySAVcVWeQWzdb/PiryBSgGguQwP4ULVR3j6pplEpnzO1tf2UuvwFyeSp
+uEMsZPFR+lySR7kFM+/i0jbkatb905RLJGoOQkCgYBdTQp7f8/pZ/t7W3mngHLb
vYeBsTS5/MsiNHJFDO5iGbjx7a+BJNbFyMczq4UOUJDDkycTyfcRaIlfYTLrzYp0
mFxZ8nw7YTuIpXNYhDbYgcreqMH05z65dCrNLisBx3LYC/Vi/W34k0lbgr3So0Fr
Y3BwgknnyHoyRu+wW0xmcg==
-----END PRIVATE KEY-----`

	testCert2 = `-----BEGIN CERTIFICATE-----
MIIC/jCCAeYCCQCpd2Zgh3uOqDANBgkqhkiG9w0BAQsFADBBMQswCQYDVQQGEwJH
QjENMAsGA1UECAwEQ0FNQjESMBAGA1UEBwwJQ2FtYnJpZGdlMQ8wDQYDVQQDDAZn
b2tvbmcwHhcNMTkwMTE0MjEyODAzWhcNMjkwMTExMjEyODAzWjBBMQswCQYDVQQG
EwJHQjENMAsGA1UECAwEQ0FNQjESMBAGA1UEBwwJQ2FtYnJpZGdlMQ8wDQYDVQQD
DAZnb2tvbmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDIuH4e8oJC
3D9IE1PXHuE+CDyUwM0/UWMbuZXCu1LslE8kY10WN60qmvUc+JcLXKTyo1tfhsgl
cDEjac0tl7WYAF0MbLQRSuO2XTMcCByMXxt4CclNE6irl075y2lZxgvvu+dokpGS
DcK7cQv6rYKwPuHuyHwFBwmN1+jW/aLtRXY6+RyWTFtBjvO1Dis/JBMVGishGwg1
6vG0iLoN1ijlVC44zeGjvgnD1RuhVP2MYOLWbRucDSQWT2ihAitm0oB+HXLI34hl
SKQkKm/7fYNg5JhhyRm27V+UPHxQZX0p32nvyzQSImkQjgNqXlGaHpBxPLhrycTY
Scczxskw5AAHAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAHRwGXeFncVy0S9zeNR0
G3YHLJ8CLRVFBAIhVDh/KUcjeZfLTK9byqCizzSUfe87jpdoQ0h/f9ikgWefx7qo
FGvj8kVhwRVqP6omS5ko1OXTDcCZiIoQjcdiB80l5b4K0rLAxEpi9jues7eqgrJU
VLzE3KPJgC0vDSyCZKCPj0yi1lyVd5sW1XPWJBfyd7Sgdvpt8TDc4lcPV/KZ5C+d
MU9gkS6sckqx5jED6t3tAR29sdMs7j2+PjeQAG8wLYPj/4l0k6tJ9Rr6DHp+XRfq
Q3ZTKuoA+/0S3D+KmwXmFW2lwY2y0DQxXbwzl8M4F6tBUQAbD0jA/htPpI5uw6+R
6+8=
-----END CERTIFICATE-----`
)
