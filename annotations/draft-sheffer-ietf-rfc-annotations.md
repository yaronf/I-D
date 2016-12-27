---
title: "Requesting Comments: Enabling Readers to Annotate RFCs"
abbrev: RFC Annotations
docname: draft-sheffer-ietf-rfc-annotations-latest
category: info
updates: 
obsoletes:

ipr: trust200902
area: General
workgroup:
keyword: Internet-Draft

stand_alone: yes
pi:
  rfcedstyle: yes
  toc: no
  tocindent: yes
  sortrefs: yes
  symrefs: yes
  strict: yes
  comments: yes
  inline: yes
  text-list-symbols: -o*+

author:
  -
    ins: Y. Sheffer
    name: Yaron Sheffer
    organization: Intuit
    email: yaronf.ietf@gmail.com

normative:
  RFC2119:
  RFC3933:

informative:

--- abstract

RFCs were initially intended as, literally, requests for comments. Since then,
they have turned into standards documents,
with a peculiar process to report errors and a highly
onerous process to actually
have the RFC modified/republished. Non-IETF participants are typically unaware of any way to provide
feedback to published RFCs, other than direct email to the listed authors.
This is very different from the way many web specifications are
developed today and arguably leads to the value of published RFCs diminishing over time.
This document proposes an experiment to remedy this situation through the deployment
of web annotations.

--- middle

# Introduction

IETF participants use the term "RFC" on a daily basis. We all know that "RFC" stands for
"Request for Comments". However the RFCs we publish are anything but requests for
comments. RFCs today are static documents that do not invite comments.
Acute readers who insist on providing
feedback will find the following text: "Information about the current status
of this document, any errata,
and how to provide feedback on it may be obtained at
http://www.rfc-editor.org/info/rfcXXXX." Once on this page, they will only find the
email address
of a working group, which may long be defunct.

We can do better than that. This document proposes, as a process experiment [rfc3933],
to enable web annotations on published RFCs. The target audience is non-IETF participants,
essentially the IETF's customers.
We discuss the advantages of such a system and the risks
associated with it.

## Conventions used in this document
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in <xref target="RFC2119"/>.

# Overview

We propose to enable, for an initial period of 1 year, annotations on published RFCs.
Document readers will be able to attach textual comments to published RFCs,
and these comments will be public, visible to all other readers who will also
be able to respond to them.

Specifically,
we recommend using the Hypothesis (https://hypothes.is/) system on our "tools" RFCs,
https://tools.ietf.org/html/rfcXXXX. We propose not to build any custom infrastructure around
this system but rather to use it as-is. When the experiment is done,
we will publish an experiment
report which will enable the IETF to decide whether this is of benefit for the
long term.

# Advantages

We foresee RFC annotations being used for a variety of purposes by RFC consumers, including:

* Providing feedback on correctness and pointing out errors. This is a much easier
   process than submitting errata, and as such would likely yield a larger number of corrections.
* Pointing out and even discussing implementation issues (annotation systems
   allow a user to "reply" to another user's comments).
* Linking to other standards and to implementations.
* Proposing ideas for and initiating discussion on "next generation" standards.

Other advantages are indirect:

* Improving the appearance of RFCs, bringing them more in line with people's expectations of web documents.
* Bringing in more people into the standards discussion, and eventually into the IETF.

# Potential Risks

The following section lists some of the issues and risks associated with this proposal,
along with a few concrete ways to mitigate some of them.

## Annotations can be improper and abusive

From a legal perspective, IETF deals with user-generated content continuously (Internet
drafts, mailing lists, wikis), so we know how to solve the problem.

However there can be a reputation cost, and in extreme cases people
may be driven away from a document because of defacement.
We might need to apply some after-the-fact
moderation to annotations, similarly to what we have now on the IETF discussion list.

## IPR issues around annotations

All public annotations made on Hypothesis are licensed under the Creative Commons CC0 license, which puts them explicitly in the public domain.

See also the Hypothesis Terms
of Service, https://hypothes.is/terms-of-service/. Note that Hypothesis itself
is a non-profit organization.

## Security and Privacy

Before they can annotate any pages, users need to register into Hypothesis. Pseudonyms
are explicitly allowed, but an email address must be provided.
Hypothesis does not currently support any federated login such as OpenID.

The Hypothesis TOS declares that they do not track users of the service. As far as the
we have seen, they only deploy a Google Analytics cookie.

Issue: can the GA cookie be disabled for particular URLs?

All traffic between the user's browser and Hypothesis is SSL-protected.

## Spam

So far spam has not been a problem with Hypothesis annotations, because users need to demonstrate a valid email address. If it ever becomes a problem, a process can be worked out where IETF volunteers monitor new annotations for spam, and the Hypothesis team removes it within a reasonable time.

## Long-term retention of annotations

If at the end of the experiment we choose to migrate to a different platform or to deploy
a private copy of Hypothesis, we will be able to use their documented API to retrieve
any extant annotations and store them into the new system.

## What if we build it and nobody comes

This would constitute a failure of the experiment, but would not have any other ill effects.

# Proposed Technical Solution

Technically, to enable annotations we simply need to add one line to each RFC published
on the "tools" site:

``
<script async defer src="https://hypothes.is/embed.js"></script>
``

A bit of additional code would be needed to display the IETF Note Well text if we choose to inform users that they implicitly agree to it.

RFC authors and WG participants can be alerted whenever their documents are annotated
using RSS and Atom feeds such as:
https://hypothes.is/stream.rss?uri=https://tools.ietf.org/html/rfc1149.

The Hypothesis system is open source, which means that it can be adopted to our needs
during the experiment or later.

# Trying it for Yourself

* Go to https://hypothes.is/, paste a link, e.g. https://tools.ietf.org/html/rfc1149 and press Annotate.
* Now open the sidebar to view existing public annotations.
* Highlight some text and right-click it. You will need to sign up for an account
  to create your own annotations. 

--- back

# Document History

## -01

* Minor changes after meeting with the Hypothesis team.

## draft-sheffer-ietf-rfc-annotations-00

Initial version.
