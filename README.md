# YAML-JSONSchema

Validate YAML or JSON against [JSON Schema](http://json-schema.org/) using CLI tool.

# Installation

`go get github.com/estambakio/yaml-jsonschema`

# Usage

`yaml-jsonschema --source test_data.yaml --schema test_schema.json`

Both `source` and `schema` **can be JSON or YAML files**.

Both `source` and `schema` can be **paths on local filesystem** or **http(s):// links**, therefore YAML/JSON files can be validated against remote schema.

# Motivation

YAML files are de-facto standard for defining configuration in cloud-native environment, e.g. [Docker](https://www.docker.com/), [Kubernetes](https://kubernetes.io/) resorces, [Helm](https://helm.sh/) charts etc. These configuration files tend to have tremendous amount of variables and are difficult to test and reason about. One possible solution is to at least validate configuration files against some sort of schema, for example validate generated `values.yaml` in Helm chart. The most throrough standard in this area is [JSON schema](http://json-schema.org/). And the easiest way to validate files is to use some kind of CLI tool.

While JSON schema is a popular standard for defining JSON structure, there are not many CLI tools which can easily validate YAML file against JSON schema.

There's JSON schema validation library called [go-jsonschema](https://github.com/xeipuuv/gojsonschema) but it works only with JSON. There's also Go library [yaml](https://github.com/ghodss/yaml) which can convert YAML to JSON.

**The idea is**: YAML is a superset of JSON, therefore if a tool works with YAML it means it should work also with JSON, because JSON is a subset of YAML. As a result, we can combine aforementioned libraries and validate both YAML and JSON files against JSON schema. If YAML file is supplied then it's converted to JSON and fed to `go-jsonschema` validator. If JSON file is supplied then convertion is no-op (by design of YAML) and `go-jsonschema` validates this JSON file.
