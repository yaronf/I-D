---
title: JSON Web Token Best Current Practices
abbrev: JWT BCP
docname: draft-sheffer-oauth-jwt-bcp-latest
category: std
updates:
obsoletes:

ipr: trust200902
area: General
workgroup:
keyword: Internet-Draft

stand_alone: yes
pi:
  rfcedstyle: yes
  toc: yes
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
    email: "yaronf.ietf@gmail.com"
  -
    ins: D. Hardt
    name: Dick Hardt
    organization: Amazon
  -
    ins: M. Jones
    name: Michael B. Jones
    organization: Microsoft

normative:
  RFC2104:
  RFC2119:

informative:
  RFC6454:

--- abstract

This is the anstract.

--- middle


# Introduction

JSON Web Tokens, also known as JWT [], are JSON structures that typically contain identity information, and are signed and/or encrypted. The JWT specification has seen rapid adoption because it encapsulates security-relevant information in one, easy to protect location, and it is easy to implement using widely-available tools.

Since JWT was published there have been several widely published attacks on implementations. Such attacks are the result of underspecified security mechanisms, as well as incomplete implementations.

The goal of this document is to enable implementations to deploy JWT securely.

These are minimum recommendations for the use of JWT in the vast majority of implementation and deployment scenarios.  Other specifications that reference this document can have stricter requirements related to one or more aspects of the format, based on their particular circumstances; when that is the case, implementers are advised to adhere to those stricter requirements.  Furthermore, this document provides a floor, not a ceiling, so stronger options are always allowed (e.g., depending on differing evaluations of the importance of cryptographic strength vs. computational load).

Community knowledge about the strength of various algorithms and feasible attacks can change quickly, and experience shows that a Best Current Practice (BCP) document about security is a point-in-time statement.  Readers are advised to seek out any errata or updates that apply to this document.

## Target Audience

The target of this document is:

* Implementors of JWT libraries
* Implementors of code that uses such libraries (to the extent that some mechanisms may not be provided by libraries, or until they are)
* Developers of specifications that rely on JWT.

# Threats and Vulnerabilities

## Insufficient Signature Validation

"none" and symmetric hashes

See draft-ietf-oauth-jwsreq-13, sec. 6.1 and 6.2.

UTF-8 vs. UTF-16?

## Recipient Confusion

"aud" claim

## Insecure ECC Encryption

No validation of input point

## Cross-JWT Confusion

"critical" and "type"

# Best Practices

