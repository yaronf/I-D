# Schema Validation

Use [ajv-cli](https://www.npmjs.com/package/ajv-cli) on MacOS.

## Validate the Schema

``ajv compile -s template-schema.json``

## Validate an Instance against the Schema

``ajv validate -s template-schema.json -d example-template.json``

