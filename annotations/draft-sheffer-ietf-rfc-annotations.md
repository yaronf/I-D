
--- Abstract

RFCs were first construed as, literally, requests for comments. Since then,
they have become standards, with a peculiar process to report errors and a highly
onerous process to actually
have the RFC modified. Non-IETF participants are typically unaware of any way to provide
feedback to published RFCs, other than direct email to the listed authors.
This is very different from the way many web specifications are
developed today and arguably leads to the value of publised RFCs diminishing over time.
This document proposes an experiment to remedy this situation through the deployment
of web annotations.

--- middle

# Introduction

IETF participants use the term RFC on a daily basis. We all know that "RFC" stands for
"Request for Comments". However the RFCs we publish are anything but requests for
comments. RFCs today are static documents, and acute readers who insist on providing
feedback will find the following text: "Information about the current status
of this document, any errata,
and how to provide feedback on it may be obtained at
http://www.rfc-editor.org/info/rfcXXXX." Once on this page, they will only find the
email address
of a working group, which may long be defunct.

We can do better than that. This document proposes, as a process experiment [rfc3933],
to enable web annotations on published RFCs. The target audience is non-IETF participants,
essntially our customers. We discuss the advantages of such a system and the risks
associated with it.

# Overview

We propose to enable, for a period of 1 year, annotations on published RFCs. Specifically,
we recommend using the Hypothesis (https://hypothes.is/) system on our "tools" RFCs,
https://tools.ietf.org/html/rfcXXXX. We propose not to build any custom infrastructure around
this system but rather to use it as-is. When the experiment is done,
we will publish an experiment
report which will enable the IETF to decide whether this is a useful process for the
long term.

# Advantages

We foresee RFC annotations being used for a variety of purposes by RFC consumers, including:

 * Providing feedback on correctness and pointing out errors. This is a much easier
process than submiting errata, and as such would likely yield a larger number of corrections.
 * Pointing out and even discussing implementation issues (annotation systems
allow a user to "reply" to another user's comments).
 * Linking to other standards and to implementations.
 * Proposing and even initiating discussion on "next generation" standards.

Other advantages are indirect:

 * Improving the appearance of RFCs, bringing them more in line with what people expect
from web documents.
 * Bringing in more people into the standards discussion, and eventually into the IETF.

# Potential Risks

The following lists some of the issues and risks associated with this proposal,
along with a few concrete ways to mitigate some of them.

## Annotations can be abused

From a legal perspective, IETF deals with user-generated content continuously (Internet
drafts, mailing lists, wikis), so we know how to solve the problem.

However there can be a reputation cost, and we might need to apply some after-the-fact
moderation to annotations, similarly to what we have now on the IETF discussion list.

## IPR issues around annotations

All annotations on Hypothesis are in the public domain. See also the Hypothesis Terms
of Service, https://hypothes.is/terms-of-service/.

## Security and Privacy

Login and accounts (OpenID).

Tracking.

SSL support.

## Long-term retention of annotations

Use API to migrate annots.

## We build it and nobody comes

This would constitute a failure of the experiment, but would not have any other ill efects.

# Proposed Technical Solution

Technically, to enable annotations we simply need to add one line to each RFC published
on the "tools" site:

``
<script async defer src="https://hypothes.is/embed.js"></script>
``

RFC authors and WG participants can be alerted on annotations to their documents
using RSS and Atom feeds such as:
https://hypothes.is/stream.rss?uri=https://tools.ietf.org/html/rfc1149.

The Hypothesis system is open source, so it can be adopted to our needs later on.

# Trying it for Yourself

* Go to https://hypothes.is/ and paste a link, e.g. https://tools.ietf.org/html/rfc1149.
* Now open the sidebar to view existing public annotations.
* Highlight some text and right-click it. You will need to open an account to create your own
annotations. 
