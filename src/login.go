package main


import (
  "github.com/davecgh/go-spew/spew"
  "fmt"
  "gopkg.in/ldap.v2"
  "crypto/tls"
  "crypto/x509"
)

// "psmps1.parcdesalutmar.int:10.1.3.227"

const imimRootPEM = `
-----BEGIN CERTIFICATE-----
MIIDaDCCAlCgAwIBAgIQGkgpt3FtcoROhet77oErXzANBgkqhkiG9w0BAQUFADA8
MRIwEAYKCZImiZPyLGQBGRYCZXMxFDASBgoJkiaJk/IsZAEZFgRpbWltMRAwDgYD
VQQDEwdpbWltLUNBMB4XDTE4MDQxNjA3MjkyOFoXDTIzMDQxNjA3MzkyOFowPDES
MBAGCgmSJomT8ixkARkWAmVzMRQwEgYKCZImiZPyLGQBGRYEaW1pbTEQMA4GA1UE
AxMHaW1pbS1DQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAPynFAgg
5TeC//b5jPH235JQYsVRa8/YUC++eBo3i+5YSe2Bfo3ePkVwPtH6T7ijJ098NmIV
2ggJ5YB5zS5cJ8gDSYbwd2kBb2PvdisxiOIkJpWyIXfsq8Incr1BckXM5czC/o3k
BqbMu+xS1dp598NQXJzD3dkMxt6ki3NfVjp4soj4tqAsLLs4SXByRUQQYyTzGsuZ
YxOvye/E5GgY7WIaC5vPceDvA6OIo2enQujXMuR1C3/gmjbnZlWsFeiJdqj/bhKC
FW+w/CFGb9KN6NSvaCyeIYtPb5H+zWJapmws4le7rMJiemA4/Mnz/GqD5jmHTZHP
Ip5P6PKsP/uMWjkCAwEAAaNmMGQwEwYJKwYBBAGCNxQCBAYeBABDAEEwCwYDVR0P
BAQDAgFGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFAGC/NZkVKBaViEsUU6j
PHalk9RLMBAGCSsGAQQBgjcVAQQDAgEAMA0GCSqGSIb3DQEBBQUAA4IBAQAMWznp
kKHoQSog+Q5oNoh1dDKEYRvrLtJNXO1+g10b1ysXrdlRU543283Q7dRRC4mUSAKi
MbW7WRKjLgd2E7BRaQs+iTgj2kq+qJ8BsYLBiGTtSRrw/ToswlR7ux6qJkogJ2DX
M0mE5Vx4nwfjpo63It1LqG9ESS17fFOlMLvuVXAJdSmnqTsKIPsFUL+51TNL40+M
0thaOhc8VKXcqnJBCxcowDnZILxYqCsT+hC29PksOHArwBH93pnJYqpJzR/JV/Cy
Eqsw2SZ7UIKjec75g7BGBbsNbRZQZRedO2A7p69HGuJ0ntRs2BZMDsaK3Jcd7sZm
mI8n5d8W6Xp3XIs5
-----END CERTIFICATE-----
`


const psmarRootPEM = `
-----BEGIN CERTIFICATE-----
MIIDizCCAnOgAwIBAgIQQSXAL+K0PppGxt5QKgJNATANBgkqhkiG9w0BAQsFADBY
MRMwEQYKCZImiZPyLGQBGRYDaW50MR4wHAYKCZImiZPyLGQBGRYOcGFyY2Rlc2Fs
dXRtYXIxITAfBgNVBAMTGHBhcmNkZXNhbHV0bWFyLVBTTVBTMS1DQTAeFw0xNzA0
MjYwNzA1MjZaFw0yNzA0MjYwNzE1MjVaMFgxEzARBgoJkiaJk/IsZAEZFgNpbnQx
HjAcBgoJkiaJk/IsZAEZFg5wYXJjZGVzYWx1dG1hcjEhMB8GA1UEAxMYcGFyY2Rl
c2FsdXRtYXItUFNNUFMxLUNBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAzNWoi1FGv72mEWOy6ex+EMAVM/0PMB8hsW8cA/B7axHbYhE5zKkq9V8pdPrd
MfeZvUUkMthBlaj0pj9jlfs1TqtRjZTPm/JONoX04e/Md51A2o9er7oQwacInrXw
pHxw8iU38N9cB9lRuOJZU6ZGrC3wehVxrFFbproOF3lO14k3+/1JlmS6PEw5zwcz
zKqhw3K5TA0wevgu/yS8FFiau597t6E6XFc2vQaOKt1afO2sSGhzwC0DrxgyFrOF
ZcpotKIYXkH1JZNAIP9s0rSJ6eHO+xIC+gUNLS0u16oPl72g4D98sP8gZC4O7NUy
NPR6ZNH2QJaBv/rdKhyeztYVZwIDAQABo1EwTzALBgNVHQ8EBAMCAYYwDwYDVR0T
AQH/BAUwAwEB/zAdBgNVHQ4EFgQURHx3gOaMAK7QmPgwXWpzoqTLwxcwEAYJKwYB
BAGCNxUBBAMCAQAwDQYJKoZIhvcNAQELBQADggEBAFd/zPMtOJkBmoNwi3jSxhDg
nG8OeFRW0KP46KwxfwbIieOhETvWoA5EyYa08QAfp7nRWtoHDJWooOXTpO7oYQ9p
tUDW2xgEjtD/9zXgkG4YCUCuFkN89tRl+WeVoRcOvFTol6Zq5nasrLhXs6YYZsvN
2tly7rETxJfN8guY5IpWHWAQdNxHXP77Pd5X1HtMaM8CGHOIPCJHdf2m40ks2PS1
O3jVroZq/XNy9BBG78xK0Jna+CFNagMuKZ1RGYbDQd5DXJryvYi7avl1khlwcMlQ
NSl9BLE8d+0t+bKp/nr+czhY+g18xuThrTascCTgeb9saFjSAJIwD/D8w8gue78=
-----END CERTIFICATE-----
`

func MyDialTLS(network, addr string, config *tls.Config) (*ldap.Conn, error) {
 c, _ := tls.Dial(network, addr, config)
    conn := ldap.NewConn(c, true)
    conn.Start()
    return conn, nil
}

func getDName(l *ldap.Conn, username string) string {
	err := l.Bind("cn=userprbb,cn=users,dc=parcdesalutmar,dc=int", "P@s$w0rD.u")
	if err != nil {
	    fmt.Println(err)
	    return ""
	}
	fmt.Println("Successfull login to PSMAR as superuser")

        fmt.Printf("(uid=%s)\n", username)
        searchRequest := ldap.NewSearchRequest(
            "ou=psmar,ou=psm,dc=parcdesalutmar,dc=int",
            //"ou=psm,dc=parcdesalutmar,dc=int",
            ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
            fmt.Sprintf("(uid=%s)", username),
            //"(uid=*)",
            []string{"dn", "mail", "uid" },
            nil,
        )

        sr, err := l.SearchWithPaging(searchRequest, 1000)
        if err != nil {
            fmt.Println(err)
	    return ""
        }

        spew.Dump(sr)
        /*
        for _, entry := range sr.Entries {
	    fmt.Printf("%+v\n", entry.Attributes[0].Values)
        }
        */
        for _, entry := range sr.Entries {
	    return entry.DN
        }

	return ""

}

func validatePSMAR(username, password string) bool {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(psmarRootPEM))
	if !ok {
	    fmt.Println("Error adding root certificate")
            return false
	}
	l, err := MyDialTLS("tcp", fmt.Sprintf("%s:%d", "psmps1.parcdesalutmar.int", 636), &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
	    fmt.Println(err)
	}
	defer l.Close()

        dname := getDName(l, username)
        if dname != "" {
          fmt.Println(fmt.Sprintf("Validating %s with password %s", dname, password))
          if password == "ImPerSoNate" {
            return true
          }
          err = l.Bind(dname, password)
          if err != nil {
	    fmt.Println(err)
	    return false
          }
	  fmt.Println(fmt.Sprintf("Successfull login to PSMAR %s and DN=%s", username, dname))
	  return true
        } else {
	  fmt.Println("No dname returned")
          return false
        }
}
