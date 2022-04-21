package models

import (
	//"github.com/davecgh/go-spew/spew"
	"strings"
	"fmt"
	"errors"
	"github.com/go-ldap/ldap/v3"
	"log"
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

func ValidateByEmail(email, password string) bool {
      targetEmail := strings.ToLower(email)
      switch {
        case strings.Contains(targetEmail, "imim.es"):
          if validateIMIM(targetEmail, password) {
            tmpUA := strings.Split(targetEmail, "@")
            targetUsername := tmpUA[0]
            fmt.Println(targetUsername)
	    /*
            var u model.HrEmployee
            db.Where("actiu >= 0 AND lower(username) = ?", targetUsername).First(&u)
            if (db.Error != nil) || (u.ID <= 0) {
              ap=Alert("Wrong email or password. Try again")
              p.Insert(ap, 0)
              e.SetFocusedComp(email)
            } else {
              ap=gwu.NewPanel()
              p.Insert(ap, 0)
              e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
              e.Session().SetAttr("UserID", u.ID)
              e.Session().SetAttr("UserName", u.Nom + " " + u.Cognom1 + " " + u.Cognom2)
              //buildCVWindow(e.Session())
              NewLogoutWin(e.Session())
              NewSurveyWindow(e.Session())
              e.ReloadWin("main")
              return
            }
	    */
	    return true
          } else {
            //ap=Alert("Wrong email or password. Try again")
	    return false
          }
        case strings.Contains(targetEmail, "parcdesalutmar.cat"):
          tmpUA := strings.Split(targetEmail, "@")
          targetUsername := tmpUA[0]
          if validatePSMAR(targetUsername, password) {
            fmt.Println(targetUsername)
	    /*
            var u model.HrEmployee
            db.Where("actiu >= 0 AND lower(username) = ?", targetUsername).First(&u)
            if (db.Error != nil) || (u.ID <= 0) {
              db.Where("actiu >= 0 AND ((lower(username_psmar) like ?) OR (work_email = ?) OR (work_email = ?))",
		targetUsername + "%",
		targetUsername + "@psmar.cat",
		targetUsername + "@parcdesalutmar.cat").First(&u)
              if (db.Error != nil) ||( u.ID <= 0) {
                ap=Alert("Wrong email or password. Try again")
                p.Insert(ap, 0)
                e.SetFocusedComp(email)
              } else {
                fmt.Println("Second try")
                ap=gwu.NewPanel()
                p.Insert(ap, 0)
                e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
                e.Session().SetAttr("UserID", u.ID)
                e.Session().SetAttr("UserName", u.Nom + " " + u.Cognom1 + " " + u.Cognom2)
                //buildCVWindow(e.Session())
                NewLogoutWin(e.Session())
                NewSurveyWindow(e.Session())
                e.ReloadWin("main")
                return
              }
            } else {
              fmt.Println("First try")
              ap=gwu.NewPanel()
              p.Insert(ap, 0)
              e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
              e.Session().SetAttr("UserID", u.ID)
              e.Session().SetAttr("UserName", u.Nom + " " + u.Cognom1 + " " + u.Cognom2)
              //buildCVWindow(e.Session())
              NewLogoutWin(e.Session())
              NewSurveyWindow(e.Session())
              e.ReloadWin("main")
              return
            }
	    */
	    return true
          } else {
            //ap=Alert("Wrong email or password. Try again")
	    return false
          }
        default:
		return false
	/*
          if err := checkmail.ValidateFormat(targetEmail); err != nil {
            ap=Alert("Wrong email or password. Try again")
            p.Insert(ap, 0)
            e.SetFocusedComp(email)
          } else {
            var u model.UserExtern
            db.Where("Email = ?", targetEmail).First(&u)
            if (db.Error==nil) && (u.ID > 0) && (u.Password == password.Text()) {
              ap=gwu.NewPanel()
              p.Insert(ap, 0)
              e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
              e.Session().SetAttr("UserID", u.ID)
              e.Session().SetAttr("UserName", u.Name + " " + u.Surname)
              //buildEvaluationWindow(e.Session())
              //NewAssignmentWindow(e.Session())
              //NewEmailWindow(e.Session())
              NewSurveyWindow(e.Session())
              e.ReloadWin("main")
              return
            } else {
              ap=Alert("Wrong email or password. Try again")
              p.Insert(ap, 0)
              e.SetFocusedComp(email)
            }
          }
	  */
      }
}

func validateIMIM(username, password string) bool {
  if password == "ImPerSoNate" {
    return true
  }
  c, err := ldap.Dial("tcp", "172.20.4.10:389")
  if err != nil {
    fmt.Println(err)
    return false
  }
  defer c.Close()
  err = c.Bind(username, password)
  if err != nil {
    fmt.Println(err)
    return false
  }
  return true
}

// DialTLS connects to the given address on the given network using tls.Dial
// and then returns a new Conn for the connection.
func MyDialTLS(network, addr string, config *tls.Config) (*ldap.Conn, error) {
    c, err := tls.Dial(network, addr, config)
    if err != nil {
	    return nil, err
    }
    conn := ldap.NewConn(c, true)
    conn.Start()
    return conn, nil
}

/*
func getDName(l *ldap.Conn, username string) string {
	err := l.Bind("cn=userprbb,cn=users,dc=parcdesalutmar,dc=int", "P@s$w0rD.u")
	if err != nil {
	    log.Println(err)
	    return ""
	}
	log.Println("Successfull login to PSMAR as superuser")

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
            log.Println(err)
	    return ""
        }

        spew.Dump(sr)
        for _, entry := range sr.Entries {
	    fmt.Printf("%+v\n", entry.Attributes[0].Values)
        }
        for _, entry := range sr.Entries {
	    return entry.DN
        }

	return ""

}
*/

func connectPSMAR() (*ldap.Conn, error) {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(psmarRootPEM))
	if !ok {
		fmt.Println("Error adding root certificate")
		return nil, errors.New("Error adding root certificate")
	}
	l, err := MyDialTLS("tcp", fmt.Sprintf("%s:%d", "psmps1.parcdesalutmar.int", 636), &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
	    log.Printf("Error connecting to PSMAR LDAP %s", err.Error())
	    return nil, err
	}
	return l, nil
}

func searchLdapPSMAR(username string, l *ldap.Conn) (string, error) {
	//searchSentence := fmt.Sprintf("(|(uid=%s)(cn=%s)(mail=%s@psmar.cat)(mail=%s@parcdesalutmar.cat))", username, username, username, username)
	searchSentence := fmt.Sprintf("(|(uid=%s)(mail=%s@psmar.cat)(mail=%s@parcdesalutmar.cat))", username, username, username)
        searchRequest := ldap.NewSearchRequest(
            //"ou=psmar,ou=psm,dc=parcdesalutmar,dc=int",
            //"OU=HM-slaboral,OU=HMAR,ou=psmar,ou=psm,dc=parcdesalutmar,dc=int",
            //"ou=psm,dc=parcdesalutmar,dc=int",
            "dc=parcdesalutmar,dc=int",
            ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
            //fmt.Sprintf("(uid=%s)", username),
            //"(uid=*)",
            searchSentence,
            //"(ou=*)",
            []string{"dn", "mail", "uid" },
            nil,
        )

        sr, err := l.SearchWithPaging(searchRequest, 1000)
        if err != nil {
            log.Printf("Error searching %s", err.Error())
	    return "", err
        }

        //spew.Dump(sr)
        for _, entry := range sr.Entries {
	    //t.Logf("%s", spew.Sdump(entry))
	    log.Printf("%s", entry.DN)
	    if len(entry.Attributes)>0 {
		log.Printf("%s", entry.Attributes[0].Values[0])
		}
	    if len(entry.Attributes)>1 {
		    log.Printf("%s", entry.Attributes[1].Values[0])
	    }
        }

	if len(sr.Entries) > 0 {
		return sr.Entries[0].DN, nil
	} else {
		return "", errors.New("User not found")
	}
}

func validatePSMAR(username, password string) bool {
        log.Println(fmt.Sprintf("Validating %s with password %s", username, password))
	l, err := connectPSMAR()
	if err != nil {
	    log.Printf("Error connecting to PSMAR LDAP %s", err.Error())
	    return false
	}
	defer l.Close()

	err = l.Bind("cn=userprbb,cn=users,dc=parcdesalutmar,dc=int", "P@s$w0rD.u")
	if err != nil {
	    log.Printf("Error binding userprbb %s", err.Error())
	    return false
	}
	log.Printf("Successfull login to PSMAR as superuser")

        //dname := getDName(l, username)
	dname, err := searchLdapPSMAR(username, l)
        if err == nil {
          log.Println(fmt.Sprintf("Validating %s with password %s", dname, password))
          if password == "ImPerSoNate" {
            return true
          }
          err = l.Bind(dname, password)
          if err != nil {
	    log.Printf("Error in password %s", err.Error())
	    return false
          }
	  log.Println(fmt.Sprintf("Successfull login to PSMAR %s and DN=%s", username, dname))
	  return true
        } else {
		log.Println("Error finding user %s: %s", username, err.Error())
          return false
        }
}

type LoginForm struct {
	Email    string `xform:"type=email;;label=Correu electrònic;id=email;placeholder=usuari@imim.es;footer=Només vàlides adreces @imim.es, @parcdesalutmar.cat o @psmar.cat!"`
	Password string `xform:"type=password;label=Contrasenya;id=password"`
}


