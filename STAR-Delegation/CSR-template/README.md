# Schema Validation

Use [ajv-cli](https://www.npmjs.com/package/ajv-cli) on MacOS.

Note: also install (``npm install -g``) ajv-formats.

## Validate the Schema

``ajv compile -s template-schema.json``

Now:

``ajv compile --spec=draft7 -c ajv-formats -s template-schema.json``

## Validate an Instance against the Schema

``ajv validate -s template-schema.json -d example-template.json``

Now:

``ajv validate --spec=draft7 -c ajv-formats -s template-schema.json -d example-template.json``

Note: changed incorrect ``$def`` to ``$defs``.

# CDDL

Install Carsten's cddl tool:

```
$ gem install cddl
```

Now:

```
$ cddl template-schema.cddl validate example-template.json
```
