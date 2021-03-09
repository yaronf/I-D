.DEFAULT_GOAL := test

SHELL := /bin/bash

cddl ?= cddl
ajv ?= npx ajv
jq ?= jq

CDDL_SCHEMA := template-schema.cddl
JSON_SCHEMA := template-schema.json

TESTFILES := $(wildcard test/*.json)

EXAMPLE := example-template.json

.PHONY: check-example
check-example: $(CDDL_SCHEMA) $(EXAMPLE)
	$(cddl) $< validate $(EXAMPLE)

EX_CFG_OBJ := example-configuration-object.json
EX_CFG_OBJ_CSR_TMPL := TEMP__config_csr_tmpl__.json
CLEANFILES += $(EX_CFG_OBJ_CSR_TMPL)

.PHONY: check-example-object
check-example-object: $(CDDL_SCHEMA) $(EX_CFG_OBJ_CSR_TMPL)
	$(cddl) $< validate $(EX_CFG_OBJ_CSR_TMPL)

$(EX_CFG_OBJ_CSR_TMPL): $(EX_CFG_OBJ)
	cat $< | $(jq) '."csr-template"' > $@

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

.PHONY: clean
clean: ; $(RM) $(CLEANFILES)
