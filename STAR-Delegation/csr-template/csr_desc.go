package main

// CsrDesc is the internalised JSON CSR decription
type CsrDesc struct {
	Subject    Subject      `json:"subject"`
	Key        Key          `json:"key"`
	Extensions *[]Extension `json:"extensions"`
}

// DistinguishedName ...
type DistinguishedName struct {
	CommonName  string `json:"common_name,omitempty"`
	CountryName string `json:"country_name,omitempty"`
}

// AlternativeName ...
type AlternativeName map[string]string

// ValidAlternativeNames ...
var ValidAlternativeNames = [...]string{"dns", "rfc822", "ipv4", "ipv6", "uri"}

// Subject ...
type Subject struct {
	DistinguishedName *DistinguishedName `json:"distinguished_name,omitempty"`
	AlternativeNames  []AlternativeName  `json:"alternative_names"`
}

// Key ...
type Key struct {
	Algorithm string `json:"algorithm"`
	Size      uint   `json:"size"`
}

// Extension ...
type Extension string
