root := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

xml2rfc ?= xml2rfc --v3
xml2rfc-prep ?= $(xml2rfc) --preptool
kramdown-rfc2629 ?= kramdown-rfc2629 --v3
idnits ?= idnits
rfcdiff ?= rfcdiff

# check if goat and svgcheck are installed
goat ?= $(shell command -v goat)
ifeq ($(strip $(goat)),)
$(error goat (ascii art to svg converter) not found. To install goat: 'go get github.com/blampe/goat')
endif

svgcheck ?= $(shell command -v svgcheck)
ifeq ($(strip $(svgcheck)),)
$(error svgcheck not found. To install svgcheck: 'pip install svgcheck')
endif

draft := draft-ietf-acme-star-delegation
current_ver := $(shell git tag | grep "$(draft)" | tail -1 | sed -e"s/.*-//")
ifeq "${current_ver}" ""
next_ver ?= 00
else
next_ver ?= $(shell printf "%.2d" $$((1$(current_ver)-99)))
endif
next := $(draft)-$(next_ver)

# CSR templates
csr_dir := $(root)CSR-template
csr_src := $(wildcard $(csr_dir)/*.json)

# ASCII -> SVG art
art_dir := $(root)art
art_src := $(wildcard $(art_dir)/*.ascii-art)
art_svg := $(art_src:.ascii-art=.svg)

COMMIT=origin/master

.PHONY: all latest submit clean

all latest: $(draft).txt $(draft).html

submit: $(next).xml $(next).txt

$(draft).txt $(draft).html: $(art_svg) $(csr_src)

idnits: $(next).txt
	$(idnits) $<

clean:
	-rm -f $(draft).txt $(draft).html $(draft).xml
	-rm -f $(next).txt $(next).html
	-rm -f $(draft)-[0-9][0-9].xml
	-rm -f $(draft)-[0-9][0-9].txt
	-rm -f $(draft)-[0-9][0-9].md
	-rm -f metadata.min.js
	-rm -f $(art_svg)

%.svg: %.ascii-art
	@$(goat) $< | $(svgcheck) -r -o $@ 2>/dev/null || true

$(draft)-orig.md:
	-rm -rf $@
	git show origin/master:$(draft).md > $@

$(draft)-$(COMMIT).md:
	-rm -rf $@
	git show $(COMMIT):$(draft).md > $@

diff: $(draft).txt $(draft)-orig.txt
	$(rfcdiff) $(draft)-orig.txt $(draft).txt
	-rm -rf $(draft)-orig.*

diff-commit: $(draft).txt $(draft)-$(COMMIT).txt
	$(rfcdiff) $(draft)-$(COMMIT).txt $(draft).txt
	-rm -rf $(draft)-orig.*

$(next).md: $(draft).md
	sed -e"s/$(basename $<)-latest/$(basename $@)/" $< > $@

%.xml: %.md
	$(kramdown-rfc2629) $< > $@
	$(xml2rfc-prep) $@ -o $@

%.txt: %.xml
	$(xml2rfc) $< $@

%.html: %.xml
	$(xml2rfc) --html $< $@

upload: $(draft).html $(draft).txt
	python upload-draft.py $(draft).html

spell: $(draft).md
	which aspell && aspell list < $< | sort | uniq | less
	# spell $< | sort | uniq | less
