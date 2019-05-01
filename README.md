# luna
A tool for installing, finding, and sharing KSQL functions. The name is a WIP. I may go with something like `ksqlfn`,
`ksqlpkg`, etc.

# Development
This is pretty garbage right now.. But it's a start!
```bash
# build the binary
$ go build -o luna cmd/luna/main.go

# search for a UDF named 'sentiment'
$ ./luna search sentiment
+------------------+--------------------+---------+
|      GROUP       |      ARTIFACT      | VERSION |
+------------------+--------------------+---------+
| com.mitchseymour | sentiment-analysis | 0.2.0   |
+------------------+--------------------+---------+

# install a UDF named 'dialogflow'. 
# TODO: differentiate UDFs with the same name that are published using different group names
$ ./luna install dialogflow /some/path
```

# How it currently works
Upload your UDF to Maven Central using the following naming convention:

```bash
<your-package>:ksql-udf-<udf-name>:<version>
```

Example:

```
com.example:ksql-udf-sentiment-analysis:0.2.0
```

The CLI will then be able to find your artifact by searching Maven whenever a user wants tries to search for or install your artifact.


# Other ideas for how this could work
- Show we implement a corresponding short + versioned URL service, like [gopkg](http://labix.org/gopkg.in)? Then, maybe users wouldn't have to publish to Maven Central, but simply use github releases to build and version a JAR?
    ```bash
    $ ./luna install ksqlpkg.in/mitch-seymour/sentiment-analysis.v1
    ```
- Or should we implement a specfile, where people could publish their function as long as they provide some metadata? e.g.
  ```
  # sentiment.ksqlspec
  name = "sentiment-analysis"
  version = "0.1.0"
  author = "Mitch"
  license = "MIT"
  url = "https://search.maven.org/remotecontent?filepath=com/mitchseymour/ksql-udf-sentiment-analysis/0.2.0/ksql-udf-sentiment-analysis-0.2.0.jar"
  ```
  
  and then..
  ```bash
   ./luna publish sentiment.ksqlspec
   ```
   
   Which would update some database we maintain, which is also used for search UDFs via `luna search`
 - ???

