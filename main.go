package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	flag "github.com/spf13/pflag"
	"github.com/xeipuuv/gojsonschema"
)

// LoadHTTPTimeout the default timeout for load requests
const LoadHTTPTimeout = 30 * time.Second

var (
	source string
	schema string
)

func init() {
	flag.StringVar(&source, "source", "", "reference to source YAML/JSON file (on local filesystem or via http(s))")
	flag.StringVar(&schema, "schema", "", "reference to JSON schema file (on local filesystem or via http(s))")
}

func main() {
	flag.Parse()

	if len(source) == 0 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("--source is required"))
		os.Exit(1)
	}

	if len(schema) == 0 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("--schema is required"))
		os.Exit(1)
	}

	rawSource, err := JSONDoc(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	rawSchema, err := JSONDoc(schema)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	documentLoader := gojsonschema.NewStringLoader(string(rawSource))
	schemaLoader := gojsonschema.NewStringLoader(string(rawSchema))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		os.Exit(1)
	}
}

// LoadFromFileOrHTTP loads the bytes from a file or a remote http server based on the path passed in
func LoadFromFileOrHTTP(path string) ([]byte, error) {
	return LoadStrategy(path, ioutil.ReadFile, loadHTTPBytes(LoadHTTPTimeout))(path)
}

// LoadStrategy returns a loader function for a given path or uri
func LoadStrategy(path string, local, remote func(string) ([]byte, error)) func(string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		return remote
	}
	return local
}

func loadHTTPBytes(timeout time.Duration) func(path string) ([]byte, error) {
	return func(path string) ([]byte, error) {
		client := &http.Client{Timeout: timeout}
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		defer func() {
			if resp != nil {
				if e := resp.Body.Close(); e != nil {
					log.Println(e)
				}
			}
		}()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("could not access document at %q [%s] ", path, resp.Status)
		}

		return ioutil.ReadAll(resp.Body)
	}
}

// JSONDoc loads a yaml document from either http or a file and converts it to json
func JSONDoc(path string) (json.RawMessage, error) {
	data, err := LoadFromFileOrHTTP(path)
	if err != nil {
		return nil, err
	}

	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
