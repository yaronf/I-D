---
title: JSON Web Token Best Current Practices
abbrev: JWT BCP
docname: draft-sheffer-oauth-jwt-bcp
category: bcp
updates:
obsoletes:

ipr: trust200902
area: Security
workgroup: OAuth Working Group
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
    email: "dick@amazon.com"
  -
    ins: M.B. Jones
    name: Michael B. Jones
    organization: Microsoft
    email: "mbj@microsoft.com"
    uri: "http://self-issued.info/"

normative:
  RFC2119:
  RFC7159:
  RFC7515:
  RFC7516:
  RFC7518:
  RFC7519:

informative:
  RFC6749:
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

  OpenID.Core:
    author:
    -
      name: Nat Sakimura
    -
      name: John Bradley
    -
      name: Michael B. Jones
    -
      name: Breno de Medeiros
    -
      name: Chuck Mortimore
    title: "OpenID Connect Core 1.0"
    date: November 8, 2014
    target: http://openid.net/specs/openid-connect-core-1_0.html

--- abstract

JSON Web Tokens, also known as JWTs {{RFC7519}}, are URL-safe JSON-based security tokens
that contain a set of claims that can be signed and/or encrypted.
JWTs are being widely used and deployed as a simple security token format
in numerous protocols and applications, both in the area of digital identity,
and in other application areas.
The goal of this Best Current Practices document is to provide actionable guidance
leading to secure implementation and deployment of JWTs.

--- middle

# Introduction

JSON Web Tokens, also known as JWTs {{RFC7519}}, are URL-safe JSON-based security tokens
that contain a set of claims that can be signed and/or encrypted.
The JWT specification has seen rapid adoption because it encapsulates
security-relevant information in one, easy to protect location, and because
it is easy to implement using widely-available tools.
One application area in which JWTs are commonly used is representing digital identity information,
such as OpenID Connect ID Tokens {{OpenID.Core}}
and OAuth 2.0 {{RFC6749}} access tokens and refresh tokens, the details of which are deployment-specific.

Since the JWT specification was published, there have been several widely published
attacks on implementations and deployments.
Such attacks are the result of underspecified security mechanisms, as well as incomplete
implementations and incorrect usage by applications.

The goal of this document is to facilitate secure implementation and deployment of JWTs.
Many of the recommendations in this document will actually be about
implementation and use of the cryptographic mechanisms underlying JWTs that are defined by
JSON Web Signature (JWS) {{RFC7515}},
JSON Web Encryption (JWE) {{RFC7516}}, and
JSON Web Algorithms (JWA) {{RFC7518}}.
Others will be about use of the JWT claims themselves.

These are intended to be minimum recommendations for the use of JWTs in the vast majority of implementation
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

The targets of this document are:

* Implementers of JWT libraries (and the JWS and JWE libraries used by them),
* Implementers of code that uses such libraries (to the extent that some mechanisms may
not be provided by libraries, or until they are), and
* Developers of specifications that rely on JWTs, both inside and outside the IETF.

## Conventions used in this document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in {{RFC2119}}.

# Threats and Vulnerabilities

This section lists some known and possible problems with JWT implementations and deployments.
Each problem description is followed by references to one or more mitigations to those problems.

## Weak Signatures and Insufficient Signature Validation

Signed JSON Web Tokens carry an explicit indication of the signing algorithm,
in the form of the "alg" header parameter, to facilitate cryptographic agility.
This, in conjunction with design flaws in some libraries and applications, have led to several attacks:

* The algorithm can be changed to "none" by an attacker, and some libraries would trust
this value and "validate" the JWT without checking any signature.
* A "RS256" (RSA, 2048 bit) can be changed into "HS256" (HMAC, SHA-256), and some libraries
would try to validate the signature using HMAC-SHA256 and using the RSA public key as the
HMAC shared secret.

For mitigations, see <xref target="algorithm-verification"/> and <xref target="appropriate-algorithms"/>.

## Weak symmetric keys

In addition, some applications sign tokens using a weak symmetric key and a keyed
MAC algorithm such as "HS256". In most cases, these keys are human memorable passwords
that are vulnerable to dictionary attacks [Langkemper].

For mitigations, see <xref target="key-entropy"/>.

## Multiplicity of JSON encodings

Many practitioners are not aware that JSON {{RFC7159}} allows several different character
encodings: UTF-8, UTF-16 and UTF-32. As a result, the JWT might be misunderstood
by its recipient.

For mitigations, see <xref target="use-utf8"/>.

## Incorrect Composition of Encryption and Signature

Some libraries that decrypt a JWE-encrypted JWT to obtain a JWS-signed object
do not always validate the internal signature.

For mitigations, see <xref target="validate-crypto"/>.

## Insecure Use of Elliptic Curve Encryption

Per [Sanso], several JOSE libraries fail to validate their inputs correctly
when performing elliptic curve key agreement (the "ECDH-ES" algorithm).
An attacker that is able to send JWEs of its choosing that use invalid curve points and
observe the resulting cleartext can use this vulnerability to recover the recipient's private key.

For mitigations, see <xref target="validate-inputs"/>.

## Substitution Attacks ## {#substitution}

There are attacks in which one recipient will have a JWT intended for it
and attempt to use it at a different recipient that it was not intended for.
If not caught, these attacks can result in the attacker gaining access to resources
that it is not entitled to access.

For mitigations, see <xref target="validate-iss-sub"/> and <xref target="use-aud"/>.

## Cross-JWT Confusion

As JWTs are being used by more and more different protocols, it becomes increasingly
important to prevent cases of JWT tokens that have been issued for one purpose
being subverted and used for another.
This is actually also a kind of substitution attack.

<!-- I don't think this text really adds value.  I've left it in a comment for now.
     It is actually describing possible mitigations, not problems, but mitigations that don't work.
     That's why I think we should delete it.
Unfortunately the JWT specification does not include one standardized attribute
that can be used to distinguish between different applications. The "critical"
attribute of [RFC7515] only points to critical headers (attributes), and therefore cannot
be used as a distinguishing value.
Neither can the "typ" attribute be used, since it is defined as
a Media Type, and in fact the highly generic values "JOSE" [RFC7515] or "JWT" [RFC7519]
are valid. 
-->

For mitigations, see <xref target="validate-iss-sub"/>, <xref target="use-aud"/>, and
<xref target="preventing-confusion"/>.

# Best Practices # {#BP}

## Perform Algorithm Verification ## {#algorithm-verification}

Libraries MUST enable the caller to specify a supported set of algorithms and
MUST NOT use any other algorithms when performing cryptographic operations.
The library MUST ensure that the "alg" or "enc" header specifies the same algorithm
that is used for the cryptographic operation.
Moreover, each key MUST be used with exactly one algorithm,
and this MUST be checked when the cryptographic operation is performed.

## Use Appropriate Algorithms ## {#appropriate-algorithms}

As Section 5.2 of [RFC7515] says, "it is an application decision which algorithms may
be used in a given context.  Even if a JWS can be successfully
validated, unless the algorithm(s) used in the JWS are acceptable to
the application, it SHOULD consider the JWS to be invalid."

Therefore, applications MUST only allow the use of cryptographically current algorithms
that meet the security requirements of the application.
This set will vary over time as new algorithms are introduced
and existing algorithms are deprecated due to discovered cryptographic weaknesses.
Applications must therefore be designed to enable cryptographic agility.

That said, if a JWT is cryptographically protected by a transport layer, such as TLS,
using cryptographically current algorithms, there may be no need to apply another layer of
cryptographic protections to the JWT.
In such cases, the use of the "none" algorithm can be perfectly acceptable.

## Validate All Cryptographic Operations ## {#validate-crypto}

All cryptographic operations used in the JWT MUST be validated and the entire JWT MUST be rejected
if any of them fail to validate.
This is true not only of JWTs with a single set of Header Parameters
but also for Nested JWTs, in which both outer and inner operations MUST be validated
using the keys and algorithms supplied by the application.

<!-- See draft-ietf-oauth-jwsreq-13, sec. 6.1 and 6.2. -->

## Validate Cryptographic Inputs ## {#validate-inputs}

Some cryptographic operations, such as Elliptic Curve Diffie-Hellman key agreement
("ECDH-ES") take inputs that may contain invalid values, such as points not on the specified elliptic curve
or other invalid points.
Either the JWS/JWE library itself must validate these inputs before using them
or it must use underlying cryptographic libraries that do so (or both!).

## Ensure Cryptographic Keys have Sufficient Entropy {#key-entropy}

The Key Entropy and Random Values advice in Section 10.1 of [RFC7515] and
the Password Considerations in Section 8.8 of [RFC7518]
MUST be followed.

## Use UTF-8 ## {#use-utf8}

[RFC7515], [RFC7516], and [RFC7519] all specify that UTF-8 be used for encoding and decoding JSON
used in Header Parameters and JWT Claims Sets.
Implementations and applications MUST do this, and not use other Unicode encodings for these purposes.

## Validate Issuer and Subject ## {#validate-iss-sub}

When a JWT contains an "iss" (issuer) claim, the application MUST validate that the cryptographic keys
used for the cryptographic operations in the JWT belong to the issuer.
If they do not, the application MUST reject the JWT.

Similarly, when the JWT contains a "sub" (subject) claim, the application MUST validate that
the subject value corresponds to a valid subject and/or issuer/subject pair at the application.
This may include confirming that the issuer is trusted by the application.
If the issuer, subject, or the pair are invalid, the application MUST reject the JWT.

## Use and Validate Audience ## {#use-aud}

If the same issuer can issue JWTs that can be used by more than one relying party or application,
the JWT MUST contain an "aud" (audience) claim that can be used to determine whether the JWT
is being used by an intended party or was substituted by an attacker at an unintended party.
Furthermore, the relying party or application MUST validate the audience value
and if the audience value is not associated with the recipient, it MUST reject the JWT.

## Use Mutually Exclusive Validation Rules for Different Kinds of JWTs ## {#preventing-confusion}

<!-- I disagree with this text for several reasons.  For one, this probably won't work for existing applications.
     It requires inventing a multiplicity of new MIME types that make semantic distinctions
     when there are no actual syntactic distinctions (what MIME types are meant for).
     And there are existing mechanisms described below that will already achieve this goal,
     without inventing a new one in this BCP.
Libraries MUST allow callers to specify the "typ" header parameter of each JWT, both when
creating a JWT and when parsing and validating it.
Applications SHOULD use this value to distinguish between different uses, to
ensure that a signed statement cannot be misinterpreted. For example, a single organization
may sign JWTs with different "typ" values:

* "JWT+id" for identity tokens.
* "JWT+set+risc" for RISC events.
* "JWT+set+oidc-logout" for OID Connect Logout events.

Specifications that define JWT applications SHOULD define a unique "typ" value.
-->

Each application of JWTs defines a profile specifying the required and optional JWT claims
and the validation rules associated with them.
To prevent substitution of JWTs from one context into another, a number of strategies may be employed:

* Use different "aud" values for different uses of JWTs within the same application.
A straightforward way of doing this, for instance, is to have the audience be the URL of the endpoint
at which that kind of JWT is to be used for Web-based applications.
* Use different sets of required claims or different required claim values.
Then the validation rules for one kind of JWT will reject those with different claims or values.
* Use different keys for different kinds of JWTs.
Then the keys used to validate one kind of JWT will fail to validate other kinds of JWTs.

Given the broad diversity of JWT usage and applications,
the best combination of required claims, values, and key usages to differentiate among different kinds of JWTs
will, in general, be application specific.
If more than one kind of JWT can be issued by the same issuer,
the validation rules for those JWTs MUST be written such that they are mutually exclusive,
rejecting JWTs of the wrong kind.

# IANA Considerations

This document requires no IANA actions.

# Acknowledgements

Thanks to Antonio Sanso for bringing the "ECDH-ES" invalid point attack to the attention
of JWE and JWT implementers.

--- back

# Document History

[[ to be removed by the RFC editor before publication as an RFC ]]

## draft-sheffer-oauth-jwt-bcp-00

- Initial version.
