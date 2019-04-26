package main

import (
  "os"
  "log"
  "fmt"
//  "encoding/json"
  "io/ioutil"
  "io"
  "net/http"
  "strings"
  "gopkg.in/alecthomas/kingpin.v2"
)

var (
  app      = kingpin.New("luna", "A tool for installing, finding, and publishing KSQL functions")

  search    = app.Command("search", "Search for a KSQL function")
  searchFor = search.Arg("function_name", "The name of the KSQL function to search for").Required().String()

  install          = app.Command("install", "Install a KSQL function")
  requirementsFile = install.Flag("r", "requirements file").File()
  installFunction  = install.Arg("function_name", "The name of the KSQL function to install").Strings()
)

func searchMaven(query string) (string, error) {
    url := fmt.Sprintf("http://search.maven.org/solrsearch/select?q=ksql-udf-%s&rows=100&wt=json", query)
    resp, err := http.Get(url)
    // make sure the request to Maven was successful
    if err != nil {
        log.Fatalln(err)
        return "", err
    }
    // make sure we can parse the response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
        return "", err
    }
    return string(body), nil
}

func downloadFile(filepath string, url string) (err error) {

  // Create the file
  out, err := os.Create(filepath)
  if err != nil  {
    return err
  }
  defer out.Close()

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Check server response
  if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("bad status: %s", resp.Status)
  }

  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    return err
  }

  return nil
}

func main() {
  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
  // Search
  case search.FullCommand():
    resp, err := searchMaven(*searchFor)
    if err != nil {
        fmt.Println("Search failed:", resp)
    } else {
        // TODO: parse the JSON response
        fmt.Println(resp)
    }

  // Install KSQL function
  case install.FullCommand():
    if *requirementsFile != nil {
        fmt.Println("TODO: Install from file")
        return
    }
    coordinates := strings.Join(*installFunction, "-")
    fmt.Println("Install:", coordinates)
    // TODO: update path to a configurable KSQL directory and remove hardcoded artifact URL
    downloadFile("ksql-udf-dialogflow-0.1.0.jar", "http://search.maven.org/remotecontent?filepath=com/mitchseymour/ksql-udf-dialogflow/0.1.0/ksql-udf-dialogflow-0.1.0.jar")
  }
}
