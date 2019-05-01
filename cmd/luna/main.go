package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("luna", "A tool for installing, finding, and publishing KSQL functions")

	search    = app.Command("search", "Search for a KSQL function")
	searchFor = search.Arg("function_name", "The name of the KSQL function to search for").Required().String()

	install          = app.Command("install", "Install a KSQL function")
	requirementsFile = install.Flag("r", "requirements file").File()
	installFunction  = install.Arg("function_name", "The name of the KSQL function to install").Required().String()
	installLocation  = install.Arg("install_location", "A directory to install the KSQL function in").Required().String()
)

type Artifact struct {
	Group   string
	Id      string
	Version string
	Url     string
}

func searchMaven(query string) ([]Artifact, error) {
	url := fmt.Sprintf("http://search.maven.org/solrsearch/select?q=ksql-udf-%s&rows=100&wt=json", query)
	resp, err := http.Get(url)
	// make sure the request to Maven was successful
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// make sure we can parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// parse the response
	json := string(body)
	numResults := gjson.Get(json, "response.numFound")
	artifacts := make([]Artifact, int(numResults.Int()))

	result := gjson.Get(json, "response.docs")
	i := 0
	result.ForEach(func(key, value gjson.Result) bool {
		coordinates := strings.Split(value.Get("id").String(), ":")
		group := coordinates[0]
		id := coordinates[1]
		version := value.Get("latestVersion").String()

		// construct the JAR url
		groupPath := strings.Replace(group, ".", "/", -1)
		url := fmt.Sprintf("https://search.maven.org/remotecontent?filepath=%s/%s/%s/%s-%s.jar", groupPath, id, version, id, version)

		artifacts[i] = Artifact{group, id, version, url}
		i++
		return true // keep iterating
	})
	return artifacts, nil
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
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
	if err != nil {
		return err
	}

	return nil
}

func printArtifacts(artifacts []Artifact) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Group", "Artifact", "Version"})

	for _, artifact := range artifacts {
		table.Append([]string{artifact.Group, strings.Replace(artifact.Id, "ksql-udf-", "", 1), artifact.Version})
	}
	table.Render() // Send output
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Search
	case search.FullCommand():
		artifacts, err := searchMaven(*searchFor)
		if err != nil {
			// error
			fmt.Println("Search failed")
		} else if len(artifacts) == 0 {
			// no results
			fmt.Println("No KSQL functions found")
		} else {
			printArtifacts(artifacts)
		}

	// Install KSQL function
	case install.FullCommand():
		if *requirementsFile != nil {
			fmt.Println("Install from file will be supported in a future version")
			return
		}
		artifacts, err := searchMaven(*installFunction)
		if err != nil {
			// error
			fmt.Println("Install failed")
		} else if len(artifacts) == 0 {
			// no results
			fmt.Println("No KSQL functions found")
		} else if len(artifacts) > 1 {
			fmt.Println("More than one artifact was found!")
			printArtifacts(artifacts)
		} else {
			artifact := artifacts[0]
			path := *installLocation
			fullPath := fmt.Sprintf("%s/%s-%s.jar", path, artifact.Id, artifact.Version)
			fmt.Println("Installing: ", artifact.Id+" to "+path)

			// Download the JAR
			err = downloadFile(fullPath, artifact.Url)
			if err != nil {
				fmt.Println("Could not save KSQL function")
				fmt.Println(err)
			}
		}

	}
}
