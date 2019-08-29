package main

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateCsrDesc ...
func ValidateCsrDesc(tmpl, desc gojsonschema.JSONLoader) (*gojsonschema.Result, error) {
	res, err := gojsonschema.Validate(tmpl, desc)
	if err != nil {
		return nil, err
	}

	if res.Valid() {
		return nil, nil
	}

	return res, errors.New("validation failed")
}

// DoCsrFromDesc ...
func DoCsrFromDesc(csrDescFileName string) error {
	return errors.New("TODO(tho) DoCsrFromDesc")
}

// build <-tmpl csr-template.json> <-in csr.json> [-out csr.pem]
func build(csrTmplFileName, csrDescFileName string) {
	csrTmpl := gojsonschema.NewReferenceLoader(csrTmplFileName)
	csrDesc := gojsonschema.NewReferenceLoader(csrDescFileName)

	res, err := ValidateCsrDesc(csrTmpl, csrDesc)
	if err != nil {
		for _, desc := range res.Errors() {
			fmt.Println("-> ", desc)
		}
		return
	}

	fmt.Println("CSR description is valid")

	err = DoCsrFromDesc(csrDescFileName)
	if err != nil {
		fmt.Println("CSR creation failed:", err.Error())
	}
}

// verify <-tmpl csr-template.json> <-in csr.pem>
func verify() {
	fmt.Println("TODO(tho) verify")
}

func main() {
	build("file://./csr-template.json", "file://./csr.json")
}
