---
title: JSON Web Token Best Current Practices
abbrev: JWT BCP
docname: draft-ietf-oauth-jwt-bcp-latest
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
    email: "dick.hardt@gmail.com"
  -
    ins: M.B. Jones
    name: Michael B. Jones
    organization: Microsoft
    email: "mbj@microsoft.com"
    uri: "http://self-issued.info/"

normative:
  RFC2119:
  RFC6979:
  RFC8259:
  RFC7515:
  RFC7516:
  RFC7518:
  RFC7519:
  RFC8037:
  RFC8174:

informative:
  RFC6749:
  RFC7517:

  Alawatugoda: DOI.10.1007/978-3-662-47854-7_6

  Kelsey: DOI.10.1007/3-540-45661-9_21
  
  Langkemper:
    author:
    -
      name: Sjoerd Langkemper
    title: "Attacking JWT Authentication"
    date: September 28, 2016
    target: https://www.sjoerdlangkemper.nl/2016/09/28/attacking-jwt-authentication/

  nist-sp-800-56a-r3:
    author:
    -
      name: Elaine Barker
    -
      name: Lily Chen
    -
      name: Sharon Keller
    -
      name: Allen Roginsky
    -
      name: Apostol Vassilev
    -
      name: Richard Davis
    title: "Recommendation for Pair-Wise Key Establishment Schemes Using Discrete Logarithm Cryptography, Draft NIST Special Publication 800-56A Revision 3"
    date: April 2018
    target: https://doi.org/10.6028/NIST.SP.800-56Ar3

  Valenta:
    author:
    -
      name: Luke Valenta
    -
      name: Nick Sullivan
    -
      name: Antonio Sanso
    -
      name: Nadia Heninger
    title: "In search of CurveSwap: Measuring elliptic curve implementations in the wild"
    date: March 29, 2018
    target: https://ia.cr/2018/298

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

  I-D.ietf-oauth-discovery:
  I-D.ietf-secevent-token:

--- abstract

JSON Web Tokens, also known as JWTs, are URL-safe JSON-based security tokens
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
Such attacks are the result of under-specified security mechanisms, as well as incomplete
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

{::boilerplate bcp14}

# Threats and Vulnerabilities

This section lists some known and possible problems with JWT implementations and deployments.
Each problem description is followed by references to one or more mitigations to those problems.

## Weak Signatures and Insufficient Signature Validation

Signed JSON Web Tokens carry an explicit indication of the signing algorithm,
in the form of the "alg" header parameter, to facilitate cryptographic agility.
This, in conjunction with design flaws in some libraries and applications, have led to several attacks:

* The algorithm can be changed to "none" by an attacker, and some libraries would trust
this value and "validate" the JWT without checking any signature.
* An "RS256" (RSA, 2048 bit) parameter value can be changed into
"HS256" (HMAC, SHA-256), and some libraries
would try to validate the signature using HMAC-SHA256 and using the RSA public key as the
HMAC shared secret.

For mitigations, see <xref target="algorithm-verification"/> and <xref target="appropriate-algorithms"/>.

## Weak symmetric keys

In addition, some applications sign tokens using a weak symmetric key and a keyed
MAC algorithm such as "HS256". In most cases, these keys are human memorable passwords
that are vulnerable to dictionary attacks [Langkemper].

For mitigations, see <xref target="key-entropy"/>.

## Multiplicity of JSON encodings

Previous versions of the JSON format {{RFC8259}} allowed several different character
encodings: UTF-8, UTF-16 and UTF-32. This is not the case anymore, with the latest
standard only allowing UTF-8. However older implementations may result in the JWT being
misinterpreted by its recipient, and this could be used by a malicious sender to bypass
the recipient's validation checks.

For mitigations, see <xref target="use-utf8"/>.

## Incorrect Composition of Encryption and Signature

Some libraries that decrypt a JWE-encrypted JWT to obtain a JWS-signed object
do not always validate the internal signature.

For mitigations, see <xref target="validate-crypto"/>.

## Insecure Use of Elliptic Curve Encryption

Per [Sanso], several JOSE libraries fail to validate their inputs correctly
when performing elliptic curve key agreement (the "ECDH-ES" algorithm).
An attacker that is able to send JWEs of its choosing that use invalid curve points and
observe the cleartext outputs resulting from decryption with the invalid curve points
can use this vulnerability to recover the recipient's private key.

For mitigations, see <xref target="validate-inputs"/>.

## Substitution Attacks ## {#substitution}

There are attacks in which one recipient will have a JWT intended for it
and attempt to use it at a different recipient that it was not intended for.
If not caught, these attacks can result in the attacker gaining access to resources
that it is not entitled to access.
For instance, if an OAuth 2.0 {{RFC6749}} access token is presented to an OAuth 2.0 protected resource
that it is intended for, that protected resource might then attempt to gain access to a different
protected resource by presenting that same access token to the different protected resource,
which the access token is not intended for.

For mitigations, see <xref target="validate-iss-sub"/> and <xref target="use-aud"/>.

## Cross-JWT Confusion

As JWTs are being used by more different protocols in diverse application areas, it becomes increasingly
important to prevent cases of JWT tokens that have been issued for one purpose
being subverted and used for another.
Note that this is a specific type of substitution attack.
If the JWT could be used in an application context in which it could be confused with other kinds of JWTs,
then mitigations MUST be employed to prevent these substitution attacks.

For mitigations, see <xref target="validate-iss-sub"/>, <xref target="use-aud"/>,
<xref target="use-typ"/>, and <xref target="preventing-confusion"/>.

# Best Practices # {#BP}

The best practices listed below should be applied by practitioners
to mitigate the threats listed in the preceding section.

## Perform Algorithm Verification ## {#algorithm-verification}

Libraries MUST enable the caller to specify a supported set of algorithms and
MUST NOT use any other algorithms when performing cryptographic operations.
The library MUST ensure that the "alg" or "enc" header specifies the same algorithm
that is used for the cryptographic operation.
Moreover, each key MUST be used with exactly one algorithm,
and this MUST be checked when the cryptographic operation is performed.

## Use Appropriate Algorithms ## {#appropriate-algorithms}

As Section 5.2 of {{RFC7515}} says, "it is an application decision which algorithms may
be used in a given context.  Even if a JWS can be successfully
validated, unless the algorithm(s) used in the JWS are acceptable to
the application, it SHOULD consider the JWS to be invalid."

Therefore, applications MUST only allow the use of cryptographically current algorithms
that meet the security requirements of the application.
This set will vary over time as new algorithms are introduced
and existing algorithms are deprecated due to discovered cryptographic weaknesses.
Applications MUST therefore be designed to enable cryptographic agility.

That said, if a JWT is cryptographically protected by a transport layer, such as TLS
using cryptographically current algorithms, there may be no need to apply another layer of
cryptographic protections to the JWT.
In such cases, the use of the "none" algorithm can be perfectly acceptable.
The "none" algorithm should only be used when the JWT is cryptographically protected by other means.
JWTs using "none" are often used in application contexts in which the content is optionally signed;
then the URL-safe claims representation and processing can be the same in both the signed and unsigned cases.
JWT libraries SHOULD NOT generate JWTs using "none" unless explicitly requested to do by the caller.

Applications SHOULD follow these algorithm-specific recommendations:

- Avoid all RSA-PKCS1 v1.5 encryption algorithms, preferring RSA-OAEP.
- ECDSA signatures require a unique random value for every message that is signed.
If even just a few bits of the random value are predictable across multiple messages then
the security of the signature scheme may be compromised. In the worst case,
the private key may be recoverable by an attacker. To counter these attacks,
JWT libraries SHOULD implement ECDSA using the deterministic approach defined in {{RFC6979}}.
This approach is completely compatible with existing ECDSA verifiers and so can be implemented
without new algorithm identifiers being required.

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
or other invalid points (see e.g. {{Valenta}}, Sec. 7.1).
Either the JWS/JWE library itself must validate these inputs before using them
or it must use underlying cryptographic libraries that do so (or both!).

ECDH-ES ephemeral public key (epk) inputs should be validated according to the recipient's
chosen elliptic curve. For the NIST prime-order curves P-256, P-384 and P-521, validation MUST
be performed according to Section 5.6.2.3.4 "ECC Partial Public-Key Validation Routine" of
NIST Special Publication 800-56A revision 3 [nist-sp-800-56a-r3].
Likewise, if the "X25519" or "X448" {{RFC8037}} algorithms are used,
then the security considerations in {{RFC8037}} apply.

## Ensure Cryptographic Keys have Sufficient Entropy {#key-entropy}

The Key Entropy and Random Values advice in Section 10.1 of {{RFC7515}} and
the Password Considerations in Section 8.8 of {{RFC7518}}
MUST be followed.
In particular, human-memorizable passwords MUST NOT be directly used
as the key to a keyed-MAC algorithm such as "HS256".
In particular, passwords should only be used to perform key encryption, rather than content encryption,
as described in Section 4.8 of {{RFC7518}}.
Note that even when used for key encryption, password-based encryption is still subject to brute-force attacks.

## Avoid Length-Dependent Encryption Inputs

Many encryption algorithms leak information about the length of the plaintext, with a varying amount of
leakage depending on the algorithm and mode of operation. Sensitive information, such as passwords,
SHOULD be padded before being encrypted. It is RECOMMENDED to avoid any compression of data before encryption
since such compression often reveals information about the plaintext. See {{Kelsey}} for general background
on compression and encryption, and {{Alawatugoda}} for a specific example of attacks on HTTP cookies.

## Use UTF-8 ## {#use-utf8}

{{RFC7515}}, {{RFC7516}}, and {{RFC7519}} all specify that UTF-8 be used for encoding and decoding JSON
used in Header Parameters and JWT Claims Sets. This is also in line with the latest JSON specification {{RFC8259}}.
Implementations and applications MUST do this, and not use or admit the use of
other Unicode encodings for these purposes.

## Validate Issuer and Subject ## {#validate-iss-sub}

When a JWT contains an "iss" (issuer) claim, the application MUST validate that the cryptographic keys
used for the cryptographic operations in the JWT belong to the issuer.
If they do not, the application MUST reject the JWT.

The means of determining the keys owned by an issuer is application-specific.
As one example, OpenID Connect {{OpenID.Core}} issuer values are "https" URLs
that reference a JSON metadata document that contains a "jwks_uri" value that is
an "https" URL from which the issuer's keys are retrieved as a JWK Set {{RFC7517}}.
This same mechanism is used by {{I-D.ietf-oauth-discovery}}.
Other applications may use different means of binding keys to issuers.

Similarly, when the JWT contains a "sub" (subject) claim, the application MUST validate that
the subject value corresponds to a valid subject and/or issuer/subject pair at the application.
This may include confirming that the issuer is trusted by the application.
If the issuer, subject, or the pair are invalid, the application MUST reject the JWT.

## Use and Validate Audience ## {#use-aud}

If the same issuer can issue JWTs that are intended for use by more than one relying party or application,
the JWT MUST contain an "aud" (audience) claim that can be used to determine whether the JWT
is being used by an intended party or was substituted by an attacker at an unintended party.
Furthermore, the relying party or application MUST validate the audience value
and if the audience value is not present or not associated with the recipient,
it MUST reject the JWT.

## Do Not Trust Received Claims

The "kid" (key ID) header is used by the relying application to perform key lookup. Applications
should ensure that this does not create SQL or LDAP injection vulnerabilities, by validating
and/or sanitizing the received value.

Similarly, blindly following a "jku" (JWK set URL) header, which may contain an arbitrary URL,
could result in server-side request forgery (SSRF) attacks. Applications should protect against such
attacks, e.g., by matching the URL to a whitelist of allowed locations,
and ensuring no cookies are sent in the GET request.

## Use Explicit Typing ## {#use-typ}

Confusion of one kind of JWT for another
can be prevented by having all the kinds of JWTs that could otherwise potentially be confused
include an explicit JWT type value and include checking the type value in their validation rules.
Explicit JWT typing is accomplished by using the "typ" header parameter.
For instance, the {{I-D.ietf-secevent-token}} specification uses the "application/secevent+jwt" media type
to perform explicit typing of Security Event Tokens (SETs).

Per the definition of "typ" in Section 4.1.9 of {{RFC7515}},
it is RECOMMENDED that the "application/" prefix be omitted from the "typ" value.
Therefore, for example, the "typ" value used to explicitly include a type for a SET
SHOULD be "secevent+jwt".
When explicit typing is employed for a JWT, it is RECOMMENDED that a media type name of the format
"application/example+jwt" be used, where "example" is replaced by the identifier for the specific kind of JWT.

When applying explicit typing to a Nested JWT, the "typ" header parameter containing the explicit type value
MUST be present in the inner JWT of the Nested JWT (the JWT whose payload is the JWT Claims Set).
The same "typ" header parameter value MAY be present in the outer JWT as well, to explicitly type the entire Nested JWT.

Note that the use of explicit typing may not achieve disambiguation from existing kinds of JWTs,
as the validation rules for existing kinds JWTs often do not use the "typ" header parameter value.
Explicit typing is RECOMMENDED for new uses of JWTs.

## Use Mutually Exclusive Validation Rules for Different Kinds of JWTs ## {#preventing-confusion}

Each application of JWTs defines a profile specifying the required and optional JWT claims
and the validation rules associated with them.
If more than one kind of JWT can be issued by the same issuer,
the validation rules for those JWTs MUST be written such that they are mutually exclusive,
rejecting JWTs of the wrong kind.
To prevent substitution of JWTs from one context into another, a number of strategies may be employed:

* Use explicit typing for different kinds of JWTs.
Then the distinct "typ" values can be used to differentiate between the different kinds of JWTs.
* Use different sets of required claims or different required claim values.
Then the validation rules for one kind of JWT will reject those with different claims or values.
* Use different sets of required header parameters or different required header parameter values.
Then the validation rules for one kind of JWT will reject those with different header parameters or values.
* Use different keys for different kinds of JWTs.
Then the keys used to validate one kind of JWT will fail to validate other kinds of JWTs.
* Use different "aud" values for different uses of JWTs from the same issuer.
Then audience validation will reject JWTs substituted into inappropriate contexts.
* Use different issuers for different kinds of JWTs.
Then the distinct "iss" values can be used to segregate the different kinds of JWTs.

Given the broad diversity of JWT usage and applications,
the best combination of types, required claims, values, header parameters, key usages, and issuers
to differentiate among different kinds of JWTs
will, in general, be application specific.
For new JWT applications, the use of explicit typing is RECOMMENDED.

# Security Considerations

This entire document is about security considerations when implementing and deploying JSON Web Tokens.

# IANA Considerations

This document requires no IANA actions.

# Acknowledgements

Thanks to Antonio Sanso for bringing the "ECDH-ES" invalid point attack to the attention
of JWE and JWT implementers. Tim McLean published the RSA/HMAC confusion attack.
Thanks to Nat Sakimura for advocating the use of explicit typing. Thanks to Neil Madden for his
numerous comments, and to Carsten Bormann, Brian Campbell, Brian Carpenter and Eric Rescorla for their reviews.

--- back

# Document History

[[ to be removed by the RFC editor before publication as an RFC ]]

## draft-ietf-oauth-jwt-bcp-05

- Genart review comments.

## draft-ietf-oauth-jwt-bcp-04

- AD review comments.

## draft-ietf-oauth-jwt-bcp-03

- Acknowledgements.

## draft-ietf-oauth-jwt-bcp-02

- Implemented WGLC feedback.

## draft-ietf-oauth-jwt-bcp-01

- Feedback from Brian Campbell.

## draft-ietf-oauth-jwt-bcp-00

- Initial WG draft. No change from the latest individual version.

## draft-sheffer-oauth-jwt-bcp-01

- Added explicit typing.

## draft-sheffer-oauth-jwt-bcp-00

- Initial version.
