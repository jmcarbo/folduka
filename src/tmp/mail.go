package main

import (
  "os"
  "io/ioutil"
  "strings"
  "path/filepath"
  "fmt"
  "time"
  "io"
  "github.com/emersion/go-message/charset"
  "golang.org/x/text/encoding/charmap"
  "github.com/emersion/go-imap/client"
  "github.com/emersion/go-message/mail"
)


func init() {
  charset.RegisterEncoding("iso-8859-7", charmap.ISO8859_7)
  charset.RegisterEncoding("iso-8859-8", charmap.ISO8859_8)
  charset.RegisterEncoding("windows-1253", charmap.Windows1253)
  charset.RegisterEncoding("windows-1254", charmap.Windows1254)
  charset.RegisterEncoding("windows-1257", charmap.Windows1257)
  charset.RegisterEncoding("windows-1258", charmap.Windows1258)
  charset.RegisterEncoding("windows-1255", charmap.Windows1255)
  charset.RegisterEncoding("windows-1256", charmap.Windows1256)
}

func main() {
}

func LoginImap() *client.Client {
  c2, err2 := client.DialTLS("hermes4.imim.es:993", nil)
  if err2 != nil {
    fmt.Printf("%s\n", err2)
    return nil
  }
  //defer c2.Logout()
  fmt.Println("Connected")
  // Login
  if err := c2.Login("imim/jmcarbo", "t0kKKZgN5i"); err != nil {
    fmt.Println(err)
    return nil
  }
  fmt.Println("Logged in")
  return c2
}

func LogoutImap(c *client.Client) {
  c.Logout()
}



func ParseMessageDate(r io.Reader) time.Time {
 m, err := mail.CreateReader(r)
 if err != nil {
     fmt.Println("Reading Message Date")
     return time.Time{}
 }

  fmt.Printf("%s\n", m.Header.Header["Date"])
  if len(m.Header.Header["Date"]) == 0 {
    return time.Now() 
  }
  dateStringArray := strings.Split(m.Header.Header["Date"][0], " boundary")
  dateString := dateStringArray[0]
  const longForm = "Jan 2, 2006 at 3:04pm (MST)"
// "15 Nov 2016 13:12:06 +0000"
  const longForm2 = "02 Jan 2006 15:04:05 -0700"
  const longForm3 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
  const longForm4 = "Mon, 2 Jan 2006 15:04:05 -0700"
  const longForm5 = "Mon, 02 Jan 2006 15:04:05 MST"
  const longForm7 = "Mon, 02 Jan 2006 15:04 MST"
  const longForm6 = "Mon, 02 Jan 2006 15:04:05 UT"
  // "1 Sep 2017 06:43:44 -0400"
  const longForm8 = "2 Jan 2006 15:04:05 -0700"
  // "Sat, 14 Apr 2018 16:31:20 +0200"
  //t, err := time.Parse(time.RFC3339, m.Header.Header["Date"][0])
  var t time.Time
  t, err = time.Parse(time.RFC1123Z, dateString)
   if err != nil {
     t, err = time.Parse(time.RFC822Z, dateString)
     if err != nil {
       t, err = time.Parse(longForm2, dateString)
       if err != nil {
         t, err = time.Parse(longForm3, dateString)
         if err != nil {
           t, err = time.Parse(longForm4, dateString)
           if err != nil {
             t, err = time.Parse(longForm5, dateString)
             if err != nil {
               t, err = time.Parse(longForm6, dateString)
               if err != nil {
                 t, err = time.Parse(longForm7, dateString)
                 if err != nil {
                   t, err = time.Parse(longForm8, dateString)
                   if err != nil {
                     fmt.Println(err)
		     t = time.Now()
                   }
                 }
               }
             }
           }
         }
       }
     }
   }
   fmt.Println(t)
/* 
     for {
         p, err := m.NextPart()
         if err == io.EOF {
             break
         } else if err != nil {
             log.Fatal(err)
         }
 
         //t, _, _ := p.Header.ContentType()
         //log.Println("A part with type", t)
         log.Printf("%+v\n", p.Header)
     }
*/
  return t
}

func ImportMessages() {
  var i int64
	fmt.Println("Importing messages:")
	err := filepath.Walk("/Volumes/Travel/lbotet@imim.es", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
                /*
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
                */
                if info.Name() != ".DS_Store" && filepath.Ext(info.Name()) != ".labels" && !info.IsDir() {
		  fmt.Printf("visited file or dir: %q\n", path)
                  //labelsBytes, err := ioutil.ReadFile(path+".labels")
                  if err != nil {
                    //log.Fatal(err)
		    return nil
                  }
                  /*
		  labels:= strings.Split(strings.TrimPrefix(strings.TrimSuffix(string(labelsBytes), "]\n"),"[")," ")
                  label := "INBOX"
                  seen := true
                  for _, l := range labels {
			switch l {
                          case "SENT":
                            label = "SENT"
                          case "UNREAD":
                            seen = false
                          default:
                            if strings.Contains(l, "Label_") {
                              label = myLabels[l] //strings.Replace(l, "CATEGORY_", "", 1)
                            }
			}
                  }
                  */
                  messageBytes, err := ioutil.ReadFile(path)
                  if err != nil {
                    fmt.Printf("%s\n", err)
                    return nil
                  }
                  //SaveMessage(label, messageBytes, seen)
		  i++
                }
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %v\n", err)
		return
	}

	fmt.Printf("Seen %d files\n", i)
}
