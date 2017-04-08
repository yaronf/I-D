---
title: JSON Web Token Best Current Practices
abbrev: JWT BCP
docname: draft-sheffer-oauth-jwt-bcp-latest
category: bcp
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
  RFC2119:
  RFC7159:
  RFC7515:
  RFC7519:

informative:
  Langkemper:
    author:
    -
      name: Sjoerd Langkemper
    title: "Attacking JWT Authentication"
    date: September 28, 2016
    target: https://www.sjoerdlangkemper.nl/2016/09/28/attacking-jwt-authentication/

  Sanso:
    author:
    -
      name: Antonio Sanso
    title: "Critical Vulnerability Uncovered in JSON Encryption"
    date: March 13, 2017
    target: https://blogs.adobe.com/security/2017/03/critical-vulnerability-uncovered-in-json-encryption.html

--- abstract

This is the abstract.

--- middle

# Introduction

JSON Web Tokens, also known as JWT {{RFC7519}}, are JSON structures that
typically contain identity information, and are signed and/or encrypted.
The JWT specification has seen rapid adoption because it encapsulates
security-relevant information in one, easy to protect location, and because
it is easy to implement using widely-available tools.

Since the JWT specification was published, there have been several widely published
attacks on implementations.
Such attacks are the result of underspecified security mechanisms, as well as incomplete
implementations.

The goal of this document is to enable implementations to deploy JWT securely.

These are minimum recommendations for the use of JWT in the vast majority of implementation
and deployment scenarios.  Other specifications that reference this document can have
stricter requirements related to one or more aspects of the format, based on their
particular circumstances; when that is the case, implementers are advised to adhere
to those stricter requirements.  Furthermore, this document provides a floor, not a ceiling,
so stronger options are always allowed (e.g., depending on differing evaluations of the
importance of cryptographic strength vs. computational load).

Community knowledge about the strength of various algorithms and feasible attacks can
change quickly, and experience shows that a Best Current Practice (BCP) document about
security is a point-in-time statement. Readers are advised to seek out any errata or
updates that apply to this document.

## Target Audience

The target of this document is:

* Implementors of JWT libraries,
* Implementors of code that uses such libraries (to the extent that some mechanisms may
not be provided by libraries, or until they are), and
* Developers of specifications that rely on JWT, both inside and outside the IETF.

## Conventions used in this document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in <xref target="RFC2119"/>.

# Threats and Vulnerabilities

We start by listing some of the known attacks on JWT implementations.

## Weak Signatures and Insufficient Signature Validation

Signed JSON Web Tokens carry an explicit indication of the signing algorithm,
in the form of the "alg" attribute. This, in conjunction with design flaws in some
libraries, have led to several attacks:

* The algorithm can be changed to "none" by an attacker, and some libraries would trust
this value and "validate" the JWT without checking any signature.
* A RS256 (RSA, 2048 bit) can be changed into HS256 (HMAC, SHA-256), and some libraries
would try to validate the signature using HMAC-SHA256 and using the RSA public key as the
HMAC shared secret.

In addition, some implementations sign tokens using a weak symmetric key and a keyed
MAC algorithm such as HS256. In most cases, these keys are human memorable password
that are vulnerable to dictionary attacks [Langkemper].

Many practitioners are not aware that JSON {{RFC7159}} allows several different character
encodings: UTF-8, UTF-16 and UTF-32. As a result the signed JWT may be misunderstood
by its recipient.

## Incorrect Composition of Encryption and Signature

Libraries that decrypt a JWE-encrypted JWT to obtain a JWS-signed object, do not
always validate the internal signature.

## Insecure Use of Elliptic Curve Encryption

According to [Sanso], several JOSE libraries fail to validate their inputs correctly
when performing elliptic curve decryption (the ECDH-ES algorithm). An attacker who
is able to send multiple version of the ciphertext and observe the resulting cleartext
can use this vulnerability to recover the recipient's private key.

## Recipient Confusion

"aud" claim - do we want to discuss it or is it too OAUTH-specific?

## Cross-JWT Confusion

As JWT is being used by more and more different protocols, it becomes increasingly
important to prevent cases of a JWT token that has been signed for one purpose,
being subverted and used for another.

Unfortunately the JWT specification does not include one standardized attribute
that can be used to distinguish between different applications. The "critical"
attribute of [RFC7515] only points to critical headers (attributes), and therefore cannot
be used as a distinguishing value.
Neither can the "typ" attribute be used, since it is defined as
a Media Type, and in fact the highly generic values "JOSE" [RFC7515] or "JWT" [RFC7519]
are valid. 

# Best Practices

## Signature Verification

Libraries MUST allow the caller to specify a supported algorithm (or a short list
of algorithms), and
MUST NOT use any other algorithms when validating the signature. The library MUST ensure
that the "alg" header is in fact the same algorithm that is used for the signature.
Moreover, each key
MUST be associated with exactly one algorithm, and this MUST be checked when
the signature is validated.

<!-- See draft-ietf-oauth-jwsreq-13, sec. 6.1 and 6.2. -->

## Distinguishing between Different Applications of JWTs

Libraries MUST allow callers to specify the "typ" header parameter of each JWT, both when
creating a JWT and when parsing and validating it.
Applications SHOULD use this value to distinguish between different uses, to
ensure that a signed statement cannot be misinterpreted. For example, a single organization
may sign JWTs with different "typ" values:

* "JWT+id" for identity tokens.
* "JWT+set+risc" for RISC events.
* "JWT+set+oidc-logout" for OID Connect Logout events.

Specifications that define JWT applications SHOULD define a unique "typ" value.

# IANA Considerations

This document requires no IANA actions.

# Acknowledgements

TBD.

--- back

# Document History

## draft-sheffer-oauth-jwt-bcp-00

- Initial version.
