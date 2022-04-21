package main

import (
  "github.com/satori/go.uuid"
  "github.com/looplab/fsm"
  "folduka/bucket"
  "path"
  "path/filepath"
  "errors"
  "fmt"
  "encoding/json"
  "strings"
  "io/ioutil"
  "regexp"
)


type WorkflowDefinition struct  {
  Name string
  Pattern string
  InitialState string
  InitialEvent string
  Events []fsm.EventDesc
  Callbacks map[string]string
}

type WorkflowInstance struct  {
  Definition WorkflowDefinition
  UUIDStr string
  CurrentState string
  OriginalPath string
  CurrentPath string
  WorkflowFilename string
  CurrentUsername string
}

func (wf *WorkflowInstance) processEvent(e *fsm.Event) {
}

func (wf *WorkflowInstance) Run(eventStr string, abucket bucket.Bucket) {
  fmt.Printf("Running workflow %s\n", wf.UUIDStr)
  fmt.Printf("%+v\n", wf)
  fmt.Printf("Event %s\n", eventStr)
  for _, e:= range wf.Definition.Events {
    if e.Name == eventStr {
      inState := false
      for _, s := range e.Src {
        if s == wf.CurrentState {
          inState = true
          break
        }
      }
      if inState {
        actions, err := getActions(abucket, wf.CurrentPath)
        if err != nil {
          fmt.Printf("Error err: %s\n", err)
        }

        for callback, v := range wf.Definition.Callbacks {
          if strings.HasPrefix(callback, fmt.Sprintf("enter_%s", e.Dst)) {
            for _, a := range actions {
              if a.Id == v {
                a.RunWorkflowInstance(wf, abucket)
              }
            }
          }
        }
        wf.CurrentState = e.Dst
        b, _ := json.Marshal(wf)
        err = abucket.WriteStream(wf.WorkflowFilename, strings.NewReader(string(b)), 0644)
        if err != nil {
          fmt.Printf("%s\n", err)
        }
      }
    }
  }
}


func getWorkflow(abucket bucket.Bucket, apath, eventStr, username string) (*[]WorkflowInstance, error) {
  wfis := []WorkflowInstance{}
  curfilestat, err := abucket.Stat(apath)
  if err != nil {
    return nil, err
  }
  // look for an instantiated workflows
  workflowpath := ""
  if curfilestat.IsDir() {
    workflowpath = path.Join(apath, "_config", "workflows")
  } else {
    workflowpath = path.Join(filepath.Dir(apath), "_config", "workflows")
  }


  targetmatch := "_folder"
  if !curfilestat.IsDir() {
    targetmatch = filepath.Base(apath)
  }
  els := abucket.List(workflowpath)
  for _, k := range *els {
    if strings.HasPrefix(k.Name(), targetmatch) {
      workflowinstancefile := path.Join(workflowpath, k.Name())
      bb, err := abucket.ReadStream(workflowinstancefile)
      if err != nil {
        fmt.Println(err)
        return nil, err
      } else {
        var wfi WorkflowInstance
        b, _ := ioutil.ReadAll(bb)
        if err := json.Unmarshal([]byte(b), &wfi); err != nil {
          fmt.Printf("error parsing workflowinstance %s\n", err)
          return nil, err
        }
        wfi.CurrentUsername = username
        wfis = append(wfis, wfi)
      }
    }
  }

  // return workflow instances
  if len(wfis) > 0 {
    return &wfis, nil
  }


  // look for a matching workflow definition
  workflowdefpath := ""
  if curfilestat.IsDir() {
    workflowdefpath = path.Join(apath, "_config", "workflows.json")
  } else {
    workflowdefpath = path.Join(filepath.Dir(apath), "_config", "workflows.json")
  }
  bb, err :=abucket.ReadStream(workflowdefpath)
  if err != nil {
    fmt.Println(err)
  } else {
    wfds := []WorkflowDefinition{}
    b, _ := ioutil.ReadAll(bb)
    if err := json.Unmarshal([]byte(b), &wfds); err != nil {
      fmt.Printf("error parsing workflowdefinition %s\n", err)
      return nil, err
    }
    for _, wfd := range wfds {
      match, _ := regexp.MatchString(wfd.Pattern, apath)
      if match && wfd.InitialEvent == eventStr {
        var wfi WorkflowInstance
        wfi.Definition = wfd
        wfi.CurrentState = wfd.InitialState
        wfi.UUIDStr = uuid.NewV4().String()
        wfi.OriginalPath = apath
        wfi.CurrentPath = apath
        b, _ := json.Marshal(&wfi)

        abucket.Make(workflowpath)
        wfiTargetFilename := path.Join(workflowpath, filepath.Base(apath) + "." + wfi.UUIDStr  + ".workflow")
        wfi.WorkflowFilename = wfiTargetFilename
        fmt.Printf("¡¡¡¡¡¡¡¡¡ Writing %s\n", wfiTargetFilename)
        err := abucket.WriteStream(wfiTargetFilename, strings.NewReader(string(b)), 0644)
        if err != nil {
          fmt.Printf("%s\n", err)
        }
        wfi.CurrentUsername = username
        wfis = append(wfis, wfi)
      }
    }
  }

  // return workflow instances
  if len(wfis) > 0 {
    return &wfis, nil
  }
  //
  return nil, errors.New("Workflow not found")
}



