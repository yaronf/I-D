root := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

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

ifeq ($(strip $(kdrfc)),)
  kdrfc_check := $(shell gem list --installed kramdown-rfc2629 -v '>= 1.3.3')
  ifeq ($(kdrfc_check), true)
    kdrfc := kdrfc -3hp
  else
    $(error kdrfc not at the right version. To install it: 'gem install kramdown-rfc2629 -v '>= 1.3.3'')
  endif
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
csr_src += $(wildcard $(csr_dir)/*.cddl)
csr_src += $(wildcard $(csr_dir)/test/*.json)

# UTF-8/ASCII diagrams
art_dir := $(root)art
art_src := $(wildcard $(art_dir)/*.ascii-art)

COMMIT=origin/master

.PHONY: all latest submit clean checks

all latest: checks $(draft).txt $(draft).html

checks: $(csr_src)
	$(MAKE) -C $(csr_dir) check-example check-example-object test

submit: $(next).xml $(next).txt

$(draft).txt $(draft).html: $(art_src) $(csr_src)

idnits: $(next).txt
	$(idnits) $<

clean:
	-rm -f $(draft).txt $(draft).html $(draft).xml $(draft).prepped.xml
	-rm -f $(next).txt $(next).html
	-rm -f $(draft)-[0-9][0-9].xml
	-rm -f $(draft)-[0-9][0-9].txt
	-rm -f $(draft)-[0-9][0-9].md
	-rm -f metadata.min.js
	$(MAKE) -C $(csr_dir) clean

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

%.txt %.hmtl %.xml: %.md
	$(kdrfc) $<
	-mv $*.prepped.xml $*.xml

upload: $(draft).html $(draft).txt
	python upload-draft.py $(draft).html

spell: $(draft).md
	which aspell && aspell list < $< | sort | uniq | less
	# spell $< | sort | uniq | less
