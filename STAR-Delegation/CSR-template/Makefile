.DEFAULT_GOAL := test

SHELL := /bin/bash

cddl ?= cddl
ajv ?= npx ajv

CDDL_SCHEMA := template-schema.cddl
JSON_SCHEMA := template-schema.json

TESTFILES := $(wildcard test/*.json)

EXAMPLE := example-template.json

.PHONY: check-example
check-example: $(CDDL_SCHEMA) $(EXAMPLE)
	$(cddl) $< validate $(EXAMPLE)

test-cddl-schema: $(CDDL_SCHEMA) $(TESTFILES)
	@echo "## testing against CDDL schema ($<)"
	@for f in $(TESTFILES); do \
		$(cddl) $< validate $$f &> /dev/null ; \
		case $$f in \
		*GOOD_*) [ $$? -eq 0 ] && echo "[OK] $$f" || echo ">> [KO] $$f" ;; \
		*FAIL_*) [ $$? -ne 0 ] && echo "[OK] $$f" || echo ">> [KO] $$f" ;; \
		esac ; \
	done

test-json-schema: $(JSON_SCHEMA) $(TESTFILES)
	@echo "## testing against JSON schema ($<)"
	@for f in $(TESTFILES); do \
		$(ajv) validate --spec=draft7 -c ajv-formats -s $< -d $$f &> /dev/null ; \
		case $$f in \
		*GOOD_*) [ $$? -eq 0 ] && echo "[OK] $$f" || echo ">> [KO] $$f" ;; \
		*FAIL_*) [ $$? -ne 0 ] && echo "[OK] $$f" || echo ">> [KO] $$f" ;; \
		esac ; \
	done

test: test-cddl-schema test-json-schema

.PHONY: test
