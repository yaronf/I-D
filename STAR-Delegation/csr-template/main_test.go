package main

import (
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestValidateCsrDesc_Fixed_Names(t *testing.T) {
	// A very restrictive template which fixes the exact shape
	// of allowed subject names.  No deviations are allowed.
	csrTmplString := `{
	  "$id": "./",
	  "allOf": [
	    {
	      "$ref": "file://./acme-star-csr-template-v1.0.0.json"
	    },
	    {
	      "properties": {
		"subject": {
		  "properties": {
		    "distinguished_name": {
		      "$comment": "constrain RDN to exactly CN=abc.ndc.dno.example",
		      "minProperties": 1,
		      "propertyNames": {
			"const": "common_name"
		      },
		      "properties": {
			"common_name": {
			  "const": "abc.ndc.dno.example"
			}
		      }
		    },
		    "alternative_names": {
		      "$comment": "constrain SAN choice to exactly DNS:abc.ndc.dno.example",
		      "minItems": 1,
		      "maxItems": 1,
		      "items": {
			"propertyNames": {
			  "const": "dns"
			},
			"properties": {
			  "dns": {
			    "const": "abc.ndc.dno.example"
			  }
			}
		      }
		    }
		  }
		}
	      }
	    }
	  ]
	}`

	type args struct {
		csrDescString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"The golden CSR description",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example"
				    },
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			false,
		},
		{
			"Extra RDN attribute fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example",
				      "country_name": "XY"
				    },
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Empty RDN fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {},
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Wrong RDN type fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "country_name": "abc.ndc.dno.example"
				    },
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Wrong RDN value fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example.extra"
				    },
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Extra SAN attribute fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example"
				    },
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example"
				      },
				      {
					"rfc822": "root@ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Empty SAN fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example"
				    },
				    "alternative_names": []
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Wrong SAN type fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example"
				    },
				    "alternative_names": [
				      {
					"rfc822": "root@ndc.dno.example"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
		{
			"Wrong SAN value fails validation",
			args{
				`{
				  "subject": {
				    "distinguished_name": {
				      "common_name": "abc.ndc.dno.example"
				    },
				    "alternative_names": [
				      {
					"dns": "abc.ndc.dno.example.extra"
				      }
				    ]
				  },
				  "key": {
				    "algorithm": "ecdsa",
				    "size": 384
				  }
				}`,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ValidateCsrDesc(
				gojsonschema.NewStringLoader(csrTmplString),
				gojsonschema.NewStringLoader(tt.args.csrDescString),
			)

			if (err != nil) != tt.wantErr {
				if err != nil && res != nil {
					for i, ve := range res.Errors() {
						t.Logf("unexpected validation error [%d]: %s", i, ve)
					}
				}
				t.Errorf("ValidateCsrDesc() <tc='%s'> err=%v", tt.name, err)
			}
		})
	}
}
