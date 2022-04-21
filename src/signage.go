package main

import (
  "github.com/kataras/iris/v12"
  "encoding/json"
  "path"
  "fmt"
  "bytes"
  "github.com/bitly/go-simplejson"
  "html/template"
  "strings"
  "regexp"
  "strconv"
)

type Stream struct {
  Name string `json:"name"`
  Items []Item `json:"items"`
}

type Item struct {
  Contents string `json:"contents"`
  Duration int `json:"duration"`
  End string `json:"end"`
  Start string `json:"start"`
  Location string `json:"location"`
  Transition string `json:"transition"`
  Type string `json:"type"`
  Width int `json:"width"`
  Height int `json:"height"`
}

func stream(app *iris.Application, cacheDir string) {
  app.Get("/stream/{path:path}", func (c iris.Context){
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    err := abucket.Download(relativepath, path.Join(cacheDir, fullpath))
    if err != nil {
      fmt.Printf("Error stream download ---> %s\n", err)
    }

    reader, err := abucket.ReadStream(relativepath)
    if err != nil {
      fmt.Printf("Simplejson error reading stream: %s\n", err)
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(reader)
    s1 := buf.String() // Does a complete copy of the bytes in the buffer.
    m := map[string]interface{}{}
    err2 := json.Unmarshal([]byte(s1), &m)
    if err2 != nil {
      fmt.Printf("POST Error unmarshaling json %s\n", err2)
    }
    jj, _ := simplejson.NewJson([]byte(s1))
    s2, _ := jj.Get("data").MarshalJSON()
    myStream := Stream{}
    err = json.Unmarshal([]byte(s2), &myStream)
    if err != nil {
      fmt.Printf("json error parsing stream: %s\n", err)
    }
    pItems := []Item{}
    for _, v := range myStream.Items {
      if v.Type == "folder" {
        fmt.Printf("%+v\n", v)
        bi:=abucket.List(v.Location)
        for _, ii := range *bi {
          aitem := Item{}
          aitem.Location = path.Join("/", ii.Prefix(), ii.Path(), ii.Name())
          aitem.Duration = v.Duration
          aitem.Contents = v.Contents
          if strings.HasSuffix(ii.Name(), ".jpeg") || strings.HasSuffix(ii.Name(), ".jpg") ||  strings.HasSuffix(ii.Name(), ".png") {
            aitem.Type = "image"
          }
          if strings.HasSuffix(ii.Name(), ".mp4") {
            aitem.Type = "video"
          }
          if aitem.Type != "" {
            re, _ := regexp.Compile(`([\d]+)sec`)
            duration := re.FindStringSubmatch(ii.Name())
            if len(duration) >= 2 {
              d,_ := strconv.Atoi(duration[1])
              aitem.Duration = d
            }
            re2, _ := regexp.Compile(`Tit-([^\-]+)-`)
            title := re2.FindStringSubmatch(ii.Name())
            if len(title) >= 2 {
              aitem.Contents = fmt.Sprintf(`<div class="content-left"><h2><strong>%s</strong></h2></div>`, title[1])
            }
            pItems = append(pItems, aitem)
          }
        }
      } else {
        pItems = append(pItems, v)
      }
    }
    c.ViewData("stream", m)
    c.ViewData("items", pItems)
    c.View("stream.html")
    c.Next()
  })

  streamTemplate :=`
{{ range $index, $value := .Items }}
{{ if eq $value.Type "image" }}
<img class="mySlides" src="{{$value.Location}}" style="width:100%" data-duration="{{$value.Duration}}">
{{end }}
{{ if eq $value.Type "youtube" }}
<video
  id="vid{{$index}}"
    class="video-js vjs-default-skin mySlides"
    loop
    muted
    autoplay
    preload="auto"
    width="1024"  height="800"
  data-setup='{ "fluid": true, "techOrder": ["youtube"], "sources": [{ "type": "video/youtube", "src": "{{$value.Location}}"}] }'  data-duration="{{$value.Duration}}"></video>
{{end }}
{{ if eq $value.Type "video" }}
<video
  id="vid{{$index}}"
    class="video-js vjs-default-skin mySlides"
    loop
    muted
    autoplay
    preload="auto"
  data-setup='{ "fluid": true }'  data-duration="{{$value.Duration}}">
	<source src='{{$value.Location}}' type='video/mp4'>
</video>
{{end }}
{{ if eq $value.Type "frame" }}
<div class="iframe-container mySlides" data-duration="{{$value.Duration}}">
<iframe id="iframe{{$index}}" src="{{$value.Location}}"></iframe>
</div>
{{end }}
{{ if eq $value.Type "html" }}
 <div class="iframe-container mySlides" data-duration="{{$value.Duration}}">
{{$value.Contents}}
</div>
{{end }}
{{ end }}
`
  app.Get("/updatestream/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    err := abucket.Download(relativepath, path.Join(cacheDir, fullpath))
    if err != nil {
      fmt.Printf("Error stream download ---> %s\n", err)
    }

    reader, err := abucket.ReadStream(relativepath)
    if err != nil {
      fmt.Printf("Simplejson error reading stream: %s\n", err)
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(reader)
    s1 := buf.String() // Does a complete copy of the bytes in the buffer.
    m := map[string]interface{}{}
    err2 := json.Unmarshal([]byte(s1), &m)
    if err2 != nil {
      fmt.Printf("POST Error unmarshaling json %s\n", err2)
    }
    jj, _ := simplejson.NewJson([]byte(s1))
    s2, _ := jj.Get("data").MarshalJSON()
    myStream := Stream{}
    err = json.Unmarshal([]byte(s2), &myStream)
    if err != nil {
      fmt.Printf("json error parsing stream: %s\n", err)
    }
    pItems := []Item{}
    for _, v := range myStream.Items {
      if v.Type == "folder" {
        fmt.Printf("%+v\n", v)
        bi:=abucket.List(v.Location)
        for _, ii := range *bi {
          aitem := Item{}
          aitem.Location = path.Join("/", ii.Prefix(), ii.Path(), ii.Name())
          aitem.Duration = v.Duration
          aitem.Type = "image"
          pItems = append(pItems, aitem)
        }
      } else {
        pItems = append(pItems, v)
      }
    }

    t := template.New("stream")

    t, err = t.Parse(streamTemplate)
    if err != nil {
      return
    }


    data := struct{
      Items []Item
    }{
      Items: pItems,
    }

    var tpl bytes.Buffer
    if err := t.Execute(&tpl, data); err != nil {
      return
    }

    sss := tpl.String()
    fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> %s\n", myStream.Name)
    fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> %s\n", sss)
    mutex.Lock()
    for k := range Conn {
      fmt.Printf("************* %s\n", ConnProperties[k])
      //k.To().Emit("chat", sss)
      if ConnProperties[k]==myStream.Name {
        k.Emit("chat", sss)
      }
      //k.To(websocket.Broadcast).Emit("chat", "blablabla")
    }
    mutex.Unlock()
    c.Next()
  })

  app.Get("/showstreams", func (c iris.Context) {
    sss := ""
    mutex.Lock()
    for k := range Conn {
      sss = sss + fmt.Sprintf("************* %s\n", ConnProperties[k])
    }
    mutex.Unlock()
    c.Write([]byte(sss))
  })

  app.Get("/reloadstreams2/{stream:string}", func (c iris.Context) {
    mystream := c.Params().Get("stream")
    mutex.Lock()
    for k := range Conn {
      if ConnProperties[k]==mystream {
        k.Emit("reload", "")
      }
    }
    mutex.Unlock()
    c.Write([]byte("Done"))
    c.Next()
  })

  app.Get("/reloadstream/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    err := abucket.Download(relativepath, path.Join(cacheDir, fullpath))
    if err != nil {
      fmt.Printf("Error stream download ---> %s\n", err)
    }

    reader, err := abucket.ReadStream(relativepath)
    if err != nil {
      fmt.Printf("Simplejson error reading stream: %s\n", err)
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(reader)
    s1 := buf.String() // Does a complete copy of the bytes in the buffer.
    m := map[string]interface{}{}
    err2 := json.Unmarshal([]byte(s1), &m)
    if err2 != nil {
      fmt.Printf("POST Error unmarshaling json %s\n", err2)
    }
    jj, _ := simplejson.NewJson([]byte(s1))
    s2, _ := jj.Get("data").MarshalJSON()
    myStream := Stream{}
    err = json.Unmarshal([]byte(s2), &myStream)
    if err != nil {
      fmt.Printf("json error parsing stream: %s\n", err)
    }
    pItems := []Item{}
    for _, v := range myStream.Items {
      if v.Type == "folder" {
        fmt.Printf("%+v\n", v)
        bi:=abucket.List(v.Location)
        for _, ii := range *bi {
          aitem := Item{}
          aitem.Location = path.Join("/", ii.Prefix(), ii.Path(), ii.Name())
          aitem.Duration = v.Duration
          aitem.Type = "image"
          pItems = append(pItems, aitem)
        }
      } else {
        pItems = append(pItems, v)
      }
    }

    t := template.New("stream")

    t, err = t.Parse(streamTemplate)
    if err != nil {
      return
    }


    data := struct{
      Items []Item
    }{
      Items: pItems,
    }

    var tpl bytes.Buffer
    if err := t.Execute(&tpl, data); err != nil {
      return
    }

    sss := tpl.String()
    fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> %s\n", myStream.Name)
    fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>> %s\n", sss)
    mutex.Lock()
    for k := range Conn {
      fmt.Printf("************* %s\n", ConnProperties[k])
      //k.To().Emit("chat", sss)
      if ConnProperties[k]==myStream.Name {
        k.Emit("reload", sss)
      }
      //k.To(websocket.Broadcast).Emit("chat", "blablabla")
    }
    mutex.Unlock()
    c.Next()
  })
}

