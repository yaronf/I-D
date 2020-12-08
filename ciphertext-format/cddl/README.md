## Convert header to CBOR diagnostic notation

Install from https://crates.io/crates/cbor-diag-cli

Skip first 4 bytes, then parse:

```tail -c +5 sample.cbor | cbor-diag --seq```

## Validate that header matches the CDDL definition

Install from https://github.com/anweiss/cddl

```tail -c +5 sample.cbor | cddl validate --cddl var-header.cddl --stdin```

