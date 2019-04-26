# luna
A tool for installing, finding, and sharing KSQL functions. The name is a WIP. I may go with something like `ksqlfn`,
`ksqlpkg`, etc.

# Development
This is pretty garbage right now.. But it's a start!
```bash
# build the binary
$ go build -o luna cmd/luna/main.go

# search for a UDF named 'sentiment'. Note: we are not parsing the response from Maven
# Central yet, so this will print a JSON string
$ ./luna search sentiment

# install a UDF named 'dialogflow'. Note: the artifact URL is currently hardcoded!
# this just shows we can download a JAR from Maven Central
$ ./luna install dialogflow
```

