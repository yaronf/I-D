xml2rfc ?= xml2rfc
kramdown-rfc2629 ?= XML_RESOURCE_ORG_PREFIX=https://unicorn-wg.github.io/idrefs kramdown-rfc2629
idnits ?= idnits
rfcdiff ?= rfcdiff

draft := draft-sheffer-ietf-rfc-annotations
current_ver := $(shell git tag | grep "$(draft)" | tail -1 | sed -e"s/.*-//")
ifeq "${current_ver}" ""
next_ver ?= 00
else
next_ver ?= $(shell printf "%.2d" $$((1$(current_ver)-99)))
endif
next := $(draft)-$(next_ver)

COMMIT=origin/master

.PHONY: all latest submit clean
.INTERMEDIATE: %.md2

all latest: $(draft).txt $(draft).html

submit: $(next).xml $(next).txt

idnits: $(next).txt
	$(idnits) $<

clean:
	-rm -f $(draft).txt $(draft).html
	-rm -f $(next).txt $(next).html
	-rm -f $(draft)-[0-9][0-9].xml
	-rm -f $(draft)-[0-9][0-9].md
	-rm -f $(draft)-[0-9][0-9].txt
	-rm -f *.md2

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

%.txt: %.xml
	$(xml2rfc) $< --text

%.html: %.xml
	$(xml2rfc) --html $<

upload: $(draft).html $(draft).txt
	python upload-draft.py $(draft).html

spell: $(draft).md
	spell $< | sort | uniq | less
