---
title: Support for Short-Term, Automatically-Renewed (STAR) Certificates in Automated Certificate Management Environment (ACME)
abbrev: ACME STAR
docname: draft-ietf-acme-star-latest
category: std

ipr: trust200902
area: Security
workgroup: ACME Working Group
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
  text-list-symbols: o-*+
  compact: yes
  subcompact: yes

author:
 -
    ins: Y. Sheffer
    name: Yaron Sheffer
    organization: Intuit
    email: yaronf.ietf@gmail.com
 -
    ins: D. Lopez
    name: Diego Lopez
    organization: Telefonica I+D
    email: diego.r.lopez@telefonica.com 
 -
    ins: O. Gonzalez de Dios
    name: Oscar Gonzalez de Dios
    organization: Telefonica I+D
    email: oscar.gonzalezdedios@telefonica.com
 -
    ins: A. Pastor Perales
    name: Antonio Agustin Pastor Perales
    organization: Telefonica I+D
    email: antonio.pastorperales@telefonica.com
 -
    ins: T. Fossati
    name: Thomas Fossati
    organization: ARM
    email: thomas.fossati@arm.com

normative:
  RFC2119:
  RFC3339:
  RFC7231:
  RFC7807:
  RFC8555:

informative:
  RFC7942:
  RFC7633:
  I-D.sheffer-acme-star-request:
  I-D.ietf-acme-star-delegation:
  I-D.nir-saag-star:
  Stark:
    -: ta
    target: http://crypto.stanford.edu/~dabo/pubs/abstracts/ssl-prefetch.html
    title: "The case for prefetching and prevalidating TLS server certificates"
    author:
      -
        ins: E. Stark
        name: Emily Stark
        org: Google
      -
        ins: L. Huang
        name: Lin-Shung Huang
        org: Carnegie Mellon University
      -
        ins: D. Israni
        name: Dinesh Israni
        org: Carnegie Mellon University
      -
        ins: C. Jackson
        name: Collin Jackson
        org: Carnegie Mellon University
      -
        ins: D. Boneh
        name: Dan Boneh
        org: Stanford University
    date: 2012
  Acer:
    -: ta
    target: https://acmccs.github.io/papers/p1407-acerA.pdf
    title: "Where the Wild Warnings Are: Root Causes of Chrome HTTPS Certificate Errors"
    seriesinfo:
      DOI: 10.1145/3133956.3134007
    author:
      -
        ins: M. E. Acer
        name: Mustafa Emre Acer
        org: Google
      -
        ins: E. Stark
        name: Emily Stark
        org: Google
      -
        ins: A. P. Felt
        name: Adrienne Porter Felt
        org: Google
      -
        ins: S. Fahl
        name: Sascha Fahl
        org: Leibniz University Hannover
      -
        ins: R. Bhargava
        name: Radhika Bhargava
        org: Purdue University
      -
        ins: B. Dev
        name: Bhanu Dev
        org: International Institute of Information Technology Hyderabad
      -
        ins: M. Braithwaite
        name: Matt Braithwaite
        org: Google
      -
        ins: R. Sleevi
        name: Ryan Sleevi
        org: Google
      -
        ins: P. Tabriz
        name: Parisa Tabriz
        org: Google
    date: 2017
  Topalovic:
    -: ta
    target: http://www.ieee-security.org/TC/W2SP/2012/papers/w2sp12-final9.pdf
    title: Towards Short-Lived Certificates
    author:
      -
        ins: E. Topalovic
        name: Emin Topalovic
        org: Stanford University
      -
        ins: B. Saeta
        name: Brennan Saeta
        org: Stanford University
      -
        ins: L. Huang
        name: Lin-Shung Huang
        org: Carnegie Mellon University
      -
        ins: C. Jackson
        name: Colling Jackson
        org: Carnegie Mellon University
      -
        ins: D. Boneh
        name: Dan Boneh
        org: Stanford University
    date: 2012
  OBrien:
    -: ta
    target: https://github.com/chromium/ct-policy
    title: Chromium Certificate Transparency Log Policy
    author:
      -
        ins: D. O'Brien
        name: Devon O'Brien
        org: Google
      -
        ins: R. Sleevi
        name: Ryan Sleevi
        org: Google
    date: 2017

--- abstract

Public-key certificates need to be revoked when they are compromised, that is, when the associated private key is exposed
to an unauthorized entity.
However the revocation process is often unreliable. An alternative to revocation is issuing a sequence
of certificates, each with a short validity period, and terminating this sequence upon compromise.
This memo proposes an ACME extension to enable the issuance of short-term and automatically renewed (STAR)
X.509 certificates.

[RFC Editor: please remove before publication]

While the draft is being developed, the editor's version can be found at
https://github.com/yaronf/I-D/tree/master/STAR.

--- middle

# Introduction

The ACME protocol {{RFC8555}} automates the process of issuing a certificate to a named entity
(an Identifier Owner or IdO). Typically, but not always, the identifier is a domain name.

If the IdO wishes to obtain a string of short-term certificates originating from the same private key (see {{Topalovic}} about why using short-lived certificates might be preferable to explicit revocation), she must go through the whole ACME protocol each time a new short-term certificate is needed - e.g., every 2-3 days.
If done this way, the process would involve frequent interactions between the registration function of the ACME Certification Authority (CA) and the identity provider infrastructure (e.g.: DNS, web servers), therefore making the issuance of short-term certificates exceedingly dependent on the reliability of both.

This document presents an extension of the ACME protocol that optimizes this process by making short-term certificates first class objects in the ACME ecosystem.
Once the Order for a string of short-term certificates is accepted, the CA is responsible for publishing the next certificate at an agreed upon URL before the previous one expires.  The IdO can terminate the automatic renewal before the negotiated deadline, if needed - e.g., on key compromise.

For a more generic treatment of STAR certificates, readers are referred to {{I-D.nir-saag-star}}.

### Name Delegation Use Case

The proposed mechanism can be used as a building block of an efficient
name-delegation protocol, for example one that exists between a CDN or a cloud
provider and its customers {{I-D.ietf-acme-star-delegation}}.  At any time,
the service customer (i.e., the IdO) can terminate the delegation by simply
instructing the CA to stop the automatic renewal and letting the currently
active certificate expire shortly thereafter.

Note that in the name delegation use case the delegated entity needs to access
the auto-renewed certificate without being in possession of the ACME account
key that was used for initiating the STAR issuance.

## Terminology

IdO
: Identifier Owner, the owner of an identifier, e.g.: a domain name, a telephone number.

STAR
: Short-Term and Automatically Renewed X.509 certificates.

## Conventions used in this document

{::boilerplate bcp14}

# Protocol Flow

The following subsections describe the three main phases of the protocol:

- Bootstrap: the IdO asks an ACME CA to create a short-term and automatically-renewed (STAR) certificate ({{proto-bootstrap}});
- Auto-renewal: the ACME CA periodically re-issues the short-term certificate and posts it to the star-certificate URL ({{proto-auto-renewal}});
- Termination: the IdO requests the ACME CA to discontinue the automatic renewal of the certificate ({{proto-termination}}).

## Bootstrap
{: #proto-bootstrap}

The IdO, in its role as an
ACME client, requests the CA to issue a STAR certificate, i.e., one that:

- Has a short validity, e.g., 24 to 72 hours. Note that the exact definition of "short" depends on the use case;
- Is automatically renewed by the CA for a certain period of time;
- Is downloadable from a (highly available) location.

Other than that, the ACME protocol flows as usual between IdO and CA.
In particular, IdO is responsible for satisfying the requested ACME challenges until the CA is willing to issue the requested certificate.
Per normal ACME processing, the IdO is given back an Order resource associated with the STAR certificate to be used in subsequent interaction with the CA (e.g., if
the certificate needs to be terminated.)

The bootstrap phase ends when the ACME CA updates the Order resource to include the URL for the issued STAR certificate.

## Refresh
{: #proto-auto-renewal}

The CA issues the initial certificate after the authorization completes successfully.
It then automatically re-issues the certificate using the same CSR (and
therefore the same identifier and public key) before the previous one expires, and publishes
it to the URL that was returned to the IdO at the end of the bootstrap phase.
The certificate user, which could be either the IdO itself or a delegated third
party, as described in {{I-D.ietf-acme-star-delegation}}, obtains the
certificate ({{fetching-certificates}}) and uses it.

The refresh process ({{figprotorefresh}}) goes on until either:

- IdO explicitly terminates the automatic renewal ({{proto-termination}}); or
- Automatic renewal expires.

~~~~~~~~~~
   Certificate             ACME/STAR
   User                    Server
   |     Retrieve cert     |                     [...]
   |---------------------->|                      |
   |                       +------.              /
   |                       |      |             /
   |                       | Automatic renewal :
   |                       |      |             \
   |                       |<-----'              \
   |     Retrieve cert     |                      |
   |---------------------->|            short validity period
   |                       |                      |
   |                       +------.              /
   |                       |      |             /
   |                       | Automatic renewal :
   |                       |      |             \
   |                       |<-----'              \
   |     Retrieve cert     |                      |
   |---------------------->|            short validity period
   |                       |                      |
   |                       +------.              /
   |                       |      |             /
   |                       | Automatic renewal :
   |                       |      |             \
   |                       |<-----'              \
   |                       |                      |
   |         [...]         |                    [...]
~~~~~~~~~~
{: #figprotorefresh title="Auto renewal"}


## Termination
{: #proto-termination}

The IdO may request early termination of the STAR certificate by sending a cancellation request to the Order resource, as described in {{protocol-details-canceling}}.
After the CA receives and verifies the request, it shall:

- Cancel the automatic renewal process for the STAR certificate;
- Change the certificate publication resource to return an error indicating the termination of the issuance;
- Change the status of the Order to "canceled".

Note that it is not necessary to explicitly revoke the short-term certificate.

~~~~~~~~~~
   Certificate                                     ACME/STAR
   User                    IdO                     Server
   |                       |                       |
   |                       |      Cancel Order     |
   |                       +---------------------->|
   |                       |                       +-------.
   |                       |                       |       |
   |                       |                       |  End auto renewal  
   |                       |                       |  Remove cert link
   |                       |                       |  etc.
   |                       |                       |       |
   |                       |         Done          |<------'
   |                       |<----------------------+
   |                       |                       |
   |                                               |
   |              Retrieve cert                    |
   +---------------------------------------------->|
   |              Error: autoRenewalCanceled       |
   |<----------------------------------------------+
   |                                               |
~~~~~~~~~~
{: #figprototerm title="Termination"}

# Protocol Details
{: #protocol-details}

This section describes the protocol details, namely the extensions
to the ACME protocol required to issue STAR certificates.

## ACME Extensions

This protocol extends the ACME protocol, to allow for automatically renewed Orders.

### Extending the Order Resource
{: #star-order-ext}

The Order resource is extended with a new "auto-renewal" object that MUST be present for STAR certificates.  The "auto-renewal" object has the following structure:

- start-date (optional, string): the earliest date of validity of the first certificate issued,
in {{RFC3339}} format.  When omitted, the start date is as soon as authorization is complete.
- end-date (required, string): the latest date of validity of the last certificate issued,
in {{RFC3339}} format.
- lifetime (required, integer): the maximum validity period of each STAR certificate, an integer that denotes a number of seconds.  This is a nominal value which does not include any extra validity time due to server or client adjustment (see below).
- lifetime-adjust (optional, integer): amount of "left pad" added to each STAR certificate, an integer that denotes a number of seconds.  The default is 0.  If present, the value of the notBefore field that would otherwise appear in the STAR certificates is pre-dated by the specified number of seconds.  See also {{operational-cons-clocks}} for why a client might want to use this control and {{computing-effective-cert-lifetime}} for how the effective certificate lifetime is computed.  The value reflected by the server, together with the value of the lifetime attribute, can be used by the client as a hint to configure its polling timer.
- allow-certificate-get (optional, boolean): see {{certificate-get-nego}}.

These attributes are included in a POST message when creating the Order, as part of the "payload" encoded object.
They are returned when the Order has been created, and the ACME server MAY adjust them at will, according to its local policy (see also {{capability-discovery}}).

The optional notBefore and notAfter fields defined in Section 7.1.3 of {{RFC8555}} MUST NOT be present in a STAR Order.
If they are included, the server MUST return an error with status code 400 "Bad
Request" and type "malformedRequest".

Section 7.1.6 of {{RFC8555}} defines the following values for the Order resource's status: "pending", "ready", "processing", "valid", and "invalid".
In the case of auto-renewal Orders, the status MUST be "valid" as long as STAR certificates are being issued.  We add a new status value: "canceled", see {{protocol-details-canceling}}.

A STAR certificate is by definition a mutable resource.  Instead of overloading the semantics of the "certificate" attribute, this document defines a new attribute "star-certificate" to be used instead of "certificate".

- star-certificate (optional, string):  A URL for the (rolling) STAR certificate that has been issued in response to this Order.

### Canceling an Auto-renewal Order
{: #protocol-details-canceling}

An important property of the auto-renewal Order is that it can be canceled by the IdO, with no need for certificate revocation. To cancel the Order, the ACME client sends a POST to the Order URL as shown in {{figcancelingstarorder}}.

~~~
  POST /acme/order/TOlocE8rfgo HTTP/1.1
  Host: example.org
  Content-Type: application/jose+json

  {
    "protected": base64url({
      "alg": "ES256",
      "kid": "https://example.com/acme/acct/evOfKhNU60wg",
      "nonce": "5XJ1L3lEkMG7tR6pA00clA",
      "url": "https://example.com/acme/order/TOlocE8rfgo"
    }),
    "payload": base64url({
      "status": "canceled"
    }),
    "signature": "H6ZXtGjTZyUnPeKn...wEA4TklBdh3e454g"
  }
~~~
{: #figcancelingstarorder title="Canceling an Auto-renewal Order"}

After a successful cancellation, the server MUST NOT issue any additional certificates for this Order.

Immediately after the Order is canceled, the server:

- MUST update the status of the Order resource to "canceled" and MUST set an appropriate "expires" date;
- MUST respond with 403 (Forbidden) to any requests to the star-certificate endpoint.  The response SHOULD provide
additional information using a problem document {{RFC7807}} with type "urn:ietf:params:acme:error:autoRenewalCanceled".

Issuing a cancellation for an Order that is not in "valid" state is not allowed.  A client MUST NOT send such a request, and a server MUST return an error response with status code 400 (Bad Request) and type "urn:ietf:params:acme:error:autoRenewalCancellationInvalid".

Explicit certificate revocation using the revokeCert interface (Section 7.6 of {{RFC8555}}) is not supported for STAR certificates.  A server receiving a revocation request for a STAR certificate MUST return an error response with status code 403 (Forbidden) and type "urn:ietf:params:acme:error:autoRenewalRevocationNotSupported".

## Capability Discovery
{: #capability-discovery}

In order to support the discovery of STAR capabilities, the "meta" field inside
the directory object defined in Section 9.7.6 of {{RFC8555}} is extended with a
new "auto-renewal" object.  The "auto-renewal" object MUST be present if the
server supports STAR.  Its structure is as follows:

- min-lifetime (required, integer): minimum acceptable value for auto-renewal lifetime, in seconds.
- max-duration (required, integer): maximum delta between the auto-renewal end-date and start-date, in seconds.
- allow-certificate-get (optional, boolean): see {{certificate-get-nego}}.

An example directory object advertising STAR support with one day min-lifetime and one year max-duration, and supporting certificate fetching with an HTTP GET is shown in {{figstardir}}.

~~~
  {
     "new-nonce": "https://example.com/acme/new-nonce",
     "new-account": "https://example.com/acme/new-account",
     "new-order": "https://example.com/acme/new-order",
     "new-authz": "https://example.com/acme/new-authz",
     "revoke-cert": "https://example.com/acme/revoke-cert",
     "key-change": "https://example.com/acme/key-change",
     "meta": {
       "terms-of-service": "https://example.com/acme/terms/2017-5-30",
       "website": "https://www.example.com/",
       "caa-identities": ["example.com"],
       "auto-renewal": {
         "min-lifetime": 86400,
         "max-duration":  31536000,
         "allow-certificate-get": true
       }
     }
  }
~~~
{: #figstardir title="Directory object with STAR support"}

## Fetching the Certificates
{: #fetching-certificates}

The certificate is fetched from the star-certificate endpoint with POST-as-GET
as per {{RFC8555}} Section 7.4.2, unless client and server have
successfully negotiated the "unauthenticated GET" option described in
{{certificate-get-nego}}.  In such case, the client can simply issue a GET to
the star-certificate resource without authenticating itself to the server as
illustrated in {{figunauthgetstarcert}}.

~~~
  GET /acme/cert/mAt3xBGaobw HTTP/1.1
  Host: example.org
  Accept: application/pem-certificate-chain

  HTTP/1.1 200 OK
  Content-Type: application/pem-certificate-chain
  Link: <https://example.com/acme/some-directory>;rel="index"
  Cert-Not-Before: Mon, 1 Feb 2016 00:00:00 GMT
  Cert-Not-After: Mon, 8 Feb 2016 00:00:00 GMT

  -----BEGIN CERTIFICATE-----
  [End-entity certificate contents]
  -----END CERTIFICATE-----
  -----BEGIN CERTIFICATE-----
  [Issuer certificate contents]
  -----END CERTIFICATE-----
  -----BEGIN CERTIFICATE-----
  [Other certificate contents]
  -----END CERTIFICATE-----
~~~
{: #figunauthgetstarcert title="Fetching a STAR certificate with unauthenticated GET"}

The Server SHOULD include the "Cert-Not-Before" and "Cert-Not-After" HTTP header fields in the response.
When they exist, they MUST be equal to the respective fields inside the end-entity certificate. Their format is "HTTP-date" as defined in Section 7.1.1.2 of {{RFC7231}}, protected with double-quotes to allow detecting and potentially recovering from situations where misbehaving intermediaries coalesce single value header fields (see Section 8.3.1 of {{RFC7231}}).
Their purpose is to enable client implementations that do not parse the certificate.

Following are further clarifications regarding usage of these header fields, as per {{RFC7231}} Sec. 8.3.1.
All apply to both headers.

* This header field is a single value, not a list.
* The header field is used only in responses to GET, HEAD and POST-as-GET requests, and only for MIME types that
denote public key certificates.
* Header field semantics are independent of context.
* The header field is not hop-by-hop.
* Intermediaries MAY insert or delete the value;
* If an intermediary inserts the value, it MUST ensure that the newly added value matches the corresponding value in the certificate.
* The header field is not appropriate for a Vary field.
* The header field is allowed within message trailers.
* The header field is not appropriate within redirects.
* The header field does not introduce additional security considerations. It discloses in a simpler form information
that is already available inside the credential.

To improve robustness, the next certificate MUST be made available by the ACME CA at the URL
pointed by "star-certificate" at the latest halfway through the lifetime of the currently active certificate.
It is worth noting that this has an implication in case of cancellation: in fact, from the time
the next certificate is made available, the cancellation is not completely effective until the latter
also expires.
To avoid the client accidentally entering a broken state, the notBefore of the "next" certificate MUST be set
so that the certificate is already valid when it is published at the "star-certificate" URL.  Note that the server
might need to increase the auto-renewal lifetime-adjust value to satisfy the latter requirement.
For further rationale on the need for adjusting the certificate validity, see {{operational-cons-clocks}}.

The server MUST NOT issue any additional certificates for this Order beyond its auto-renewal end-date.

Immediately after the Order expires, the server MUST respond with 403 (Forbidden) to any requests to the star-certificate endpoint.  The response SHOULD provide additional information using a problem document {{RFC7807}} with type "urn:ietf:params:acme:error:autoRenewalExpired". Note that the Order resource's state remains "valid", as per the base protocol.

## Negotiating an unauthenticated GET
{: #certificate-get-nego }

In order to enable the name delegation workflow defined in
{{I-D.ietf-acme-star-delegation}} as well as to increase the reliability of the
STAR ecosystem (see {{dependability}} for details), this document defines a
mechanism that allows a server to advertise support for accessing
star-certificate resources via unauthenticated GET (instead of, or in addition
to, POST-as-GET), and a client to enable this service with per-Order
granularity.

Specifically, a server states its availability to grant unauthenticated access
to a client's Order star-certificate by setting the allow-certificate-get
attribute to true in the auto-renewal object of the meta field inside the
Directory object:

- allow-certificate-get (optional, boolean): If this field is present and
  set to true, the server allows GET requests to star-certificate URLs.

A client states its will to access the issued star-certificate via
unauthenticated GET by adding an allow-certificate-get attribute to the
auto-renewal object of its Order and setting it to true.

- allow-certificate-get (optional, boolean): If this field is present and
  set to true, the client requests the server to allow unauthenticated GET to
  the star-certificate associated with this Order.

If the server accepts the request, it MUST reflect the attribute setting in the resulting Order object.

## Computing notBefore and notAfter of STAR Certificates
{: #computing-effective-cert-lifetime}

We define "nominal renewal date" the point in time when a new short-term
certificate for a given STAR Order is due.  It is a multiple of the Order's
auto-renewal lifetime that starts with the issuance of the first short-term
certificate and is upper-bounded by the Order's auto-renewal end-date
({{fignrd}}).

~~~
    T      - STAR Order's auto-renewal lifetime
    end    - STAR Order's auto-renewal end-date
    nrd[i] - nominal renewal date of the i-th STAR certificate


                 .- T -.   .- T -.   .- T -.   .__.
                /       \ /       \ /       \ /  end
    -----------o---------o---------o---------o----X-------> t
              nrd[0]    nrd[1]    nrd[2]    nrd[3]
~~~
{: #fignrd title="Nominal Renewal Date"}

The rules to determine the notBefore and notAfter values of the i-th STAR
certificate are as follows:

~~~
    notAfter  = min(nrd[i] + T, end)
    notBefore = nrd[i] - max(adjust_client, adjust_server)
~~~

Where "adjust_client" is the min between the auto-renewal lifetime-adjust value
("la"), optionally supplied by the client, and the auto-renewal lifetime of
each short-term certificate ("T"); "adjust_server" is the amount of padding
added by the ACME server to make sure that all certificates being published are
valid at the time of publication.  The server padding is a fraction f of T
(i.e., f * T with .5 <= f < 1, see {{fetching-certificates}}):

~~~
    adjust_client = min(T, la)
    adjust_server = f * T
~~~

Note that the ACME server MUST NOT set the notBefore of the first STAR
certificate to a date prior to the auto-renewal start-date.

### Example

Given a server that intends to publish the next STAR certificate halfway
through the lifetime of the previous one, and a STAR Order with the following
attributes:

~~~
     "auto-renewal": {
       "start-date": "2016-01-10T00:00:00Z",
       "end-date": "2016-01-20T00:00:00Z",
       "lifetime": 345600,          // 4 days
       "lifetime-adjust": 259200    // 3 days
     }
~~~

The amount of time that needs to be subtracted from each nominal renewal
date is 3 days -- i.e., max(min(345600, 259200), 345600 * .5).

The notBefore and notAfter of each short-term certificate are:

| notBefore | notAfter |
|-----------|----------|
| 2016-01-10T00:00:00Z | 2016-01-14T00:00:00Z |
| 2016-01-11T00:00:00Z | 2016-01-18T00:00:00Z |
| 2016-01-15T00:00:00Z | 2016-01-20T00:00:00Z |

The value of the notBefore is also the time at which the client should expect
the new certificate to be available from the star-certificate endpoint.

# Operational Considerations

## The Meaning of "Short Term" and the Impact of Skewed Clocks
{: #operational-cons-clocks }

"Short Term" is a relative concept, therefore trying to define a cut-off point that works in all cases would be a useless exercise.  In practice, the expected lifetime of a STAR certificate will be counted in minutes, hours or days, depending on different factors: the underlying requirements for revocation, how much clock synchronization is expected among relying parties and the issuing CA, etc.

Nevertheless, this section attempts to provide reasonable suggestions for the Web use case, informed by current operational and research experience.

Acer et al. {{Acer}} find that one of the main causes of "HTTPS error" warnings in browsers is misconfigured client clocks.  In particular, they observe that roughly 95% of the "severe" clock skews - the 6.7% of clock-related breakage reports which account for clients that are more than 24 hours behind - happen to be within 6-7 days.

In order to avoid these spurious warnings about a not (yet) valid server certificate, site owners could use the auto-renewal lifetime-adjust attribute to control the effective lifetime of their Web facing certificates.  The exact number depends on the percentage of the "clock-skewed" population that the site owner expects to protect - 5 days cover 97.3%, 7 days cover 99.6% as well as the nominal auto-renewal lifetime of the STAR Order.  Note that exact choice is also likely to depend on the kind of clients that is prevalent for a given site or app - for example, Android and Mac OS clients are known to behave better than Windows clients.  These considerations are clearly out of scope of the present document.

In terms of security, STAR certificates and certificates with OCSP must-staple {{RFC7633}} can be considered roughly equivalent if the STAR certificate's and the OCSP response's lifetimes are the same.  Given OCSP responses can be cached on average for 4 days {{Stark}}, it is RECOMMENDED that a STAR certificate that is used on the Web has an "effective" lifetime (excluding any adjustment to account for clock skews) no longer than 4 days.

## Impact on Certificate Transparency (CT) Logs

Even in the highly unlikely case STAR becomes the only certificate issuance model,
discussion with the IETF TRANS Working Group and Certificate Transparency (CT)
logs implementers suggests that existing CT Log Server implementations
are capable of sustaining the resulting 100-fold increase in ingestion
rate.  Additionally, such a future, higher load could be managed with a variety
of techniques (e.g., sharding by modulo of certificate hash, using "smart"
load-balancing CT proxies, etc.).  With regards to the increase in the log
size, current CT log growth is already being managed with schemes like Chrome's
Log Policy {{OBrien}} which allow Operators to define their log life-cycle; and
allowing the CAs, User Agents, Monitors, and any other interested entities to
build-in support for that life-cycle ahead of time.

## HTTP Caching and Dependability
{: #dependability}

When using authenticated POST-as-GET, the HTTPS endpoint from where the STAR
certificate is fetched can't be easily replicated by an on-path HTTP cache.
Reducing the caching properties of the protocol makes STAR clients increasingly
dependent on the ACME server availability.  This might be problematic given the
relatively high rate of client-server interactions in a STAR ecosystem.
Clients and servers should consider using the mechanism described in
{{certificate-get-nego}} to mitigate the risk.

When using unauthenticated GET to fetch the STAR certificate, the server SHALL
use the appropriate cache directives to set the freshness lifetime of the
response (Section 5.2 of {{!RFC7234}}) such that on-path caches will consider
it stale before or at the time its effective lifetime is due to expire.

# Implementation Status

Note to RFC Editor: please remove this section before publication,
including the reference to {{RFC7942}} and {{I-D.sheffer-acme-star-request}}.

This section records the status of known implementations of the
protocol defined by this specification at the time of posting of
this Internet-Draft, and is based on a proposal described in
{{RFC7942}}.  The description of implementations in this section is
intended to assist the IETF in its decision processes in
progressing drafts to RFCs.  Please note that the listing of any
individual implementation here does not imply endorsement by the
IETF.  Furthermore, no effort has been spent to verify the
information presented here that was supplied by IETF contributors.
This is not intended as, and must not be construed to be, a
catalog of available implementations or their features.  Readers
are advised to note that other implementations may exist.

According to {{RFC7942}}, "this will allow reviewers and working
groups to assign due consideration to documents that have the
benefit of running code, which may serve as evidence of valuable
experimentation and feedback that have made the implemented
protocols more mature.  It is up to the individual working groups
to use this information as they see fit".

## Overview

The implementation is constructed around 3 elements: STAR Client for the Name Delegation Client (NDC),
STAR Proxy for IdO and ACME Server for CA. The communication between
them is over an IP network and the HTTPS protocol.

The software of the implementation is available at: https://github.com/mami-project/lurk

The following subsections offer a basic description, detailed information
is available in https://github.com/mami-project/lurk/blob/master/proxySTAR_v2/README.md

### ACME Server with STAR extension

This is a fork of the Let's Encrypt Boulder project that implements an ACME compliant CA.
It includes modifications to extend the ACME protocol as it is specified in this draft,
to support recurrent Orders and cancelling Orders. 

The implementation understands the new "recurrent" attributes as part of the Certificate
issuance in the POST request for a new resource.
An additional process "renewalManager.go" has been included in parallel that reads
the details of each recurrent request, automatically produces a "cron" Linux based task
that issues the recurrent certificates, until the lifetime ends or the Order is canceled.
This process is also in charge of maintaining a fixed URI to enable the NDC to download certificates,
unlike Boulder's regular process of producing a unique URI per certificate.

### STAR Proxy

The STAR Proxy has a double role as ACME client and STAR Server. The former is a fork of the EFF
Certbot project that implements an ACME compliant client with the STAR extension.
The latter is a basic HTTP REST API server.

The STAR Proxy understands the basic API request with a server. The current implementation
of the API is defined in draft-ietf-acme-star-01. Registration or Order cancellation
triggers the modified Certbot client that requests, or cancels, the recurrent generation
of certificates using the STAR extension over ACME protocol.
The URI with the location of the recurrent certificate is delivered to the STAR client as a response.

## Level of Maturity

This is a prototype. 

## Coverage

A STAR Client is not included in this implementation, but done by direct HTTP request with any open HTTP REST API tool.
This is expected to be covered as part of the {{I-D.sheffer-acme-star-request}} implementation.

This implementation completely covers STAR Proxy and ACME Server with STAR extension. 

## Version Compatibility

The implementation is compatible with version draft-ietf-acme-star-01. 
The implementation is based on the Boulder and Certbot code release from 7-Aug-2017.

## Licensing
This implementation inherits the Boulder license (Mozilla Public License 2.0)
and Certbot license (Apache License Version 2.0 ).

## Implementation experience

To prove the concept all the implementation has been done with a self-signed CA,
to avoid impact on real domains. To be able to do it we use the FAKE_DNS property
of Boulder and static /etc/hosts entries with domains names.
Nonetheless this implementation should run with real domains. 

Most of the implementation has been made to avoid deep changes inside of Boulder
or Certbot, for example, the recurrent certificates issuance by the CA is based
on an external process that auto-configures the standard Linux "cron" daemon in the ACME CA server. 

The reference setup recommended is one physical host with 3 virtual machines,
one for each of the 3 components (client, proxy and server) and the connectivity based on host bridge.

Network security is not enabled (iptables default policies are "accept" and all rules removed)
in this implementation to simplify and test the protocol. 

## Contact Information

See author details below.

# IANA Considerations

[[RFC Editor: please replace XXXX below by the RFC number.]]

## New Registries

This document requests that IANA create the following new registries:

* ACME Order Auto Renewal Fields ({{iana-order-auto-renewal-registry}})
* ACME Directory Metadata Auto Renewal Fields ({{iana-metadata-auto-renewal-registry}})

All of these registries are administered under a Specification Required policy
{{!RFC8126}}.

## New Error Types

This document adds the following entries to the ACME Error Type registry:

| Type | Description | Reference |
|------|-------------|-----------|
| autoRenewalCanceled | The short-term certificate is no longer available because the auto-renewal Order has been explicitly canceled by the IdO | RFC XXXX |
| autoRenewalExpired | The short-term certificate is no longer available because the auto-renewal Order has expired | RFC XXXX |
| autoRenewalCancellationInvalid | A request to cancel a auto-renewal Order that is not in state "valid" has been received | RFC XXXX |
| autoRenewalRevocationNotSupported | A request to revoke a auto-renewal Order has been received | RFC XXXX |

## New fields in Order Objects

This document adds the following entries to the ACME Order Object Fields registry:

| Field Name | Field Type | Configurable | Reference |
|------------|------------|--------------|-----------|
| auto-renewal | object | true | RFC XXXX |
| star-certificate | string | false | RFC XXXX |

## Fields in the "auto-renewal" Object within an Order Object
{: #iana-order-auto-renewal-registry}

The "ACME Order Auto Renewal Fields" registry lists field names that are
defined for use in the JSON object included in the "auto-renewal" field of an
ACME order object.

Template:

* Field name: The string to be used as a field name in the JSON object
* Field type: The type of value to be provided, e.g., string, boolean, array of
  string
* Configurable: Boolean indicating whether the server should accept values
  provided by the client
* Reference: Where this field is defined

Initial contents: The fields and descriptions defined in {{star-order-ext}}.

| Field Name | Field Type | Configurable | Reference |
|------------|------------|--------------|-----------|
| start-date | string | true | RFC XXXX |
| end-date | string | true | RFC XXXX |
| lifetime | integer | true | RFC XXXX |
| lifetime-adjust | integer | true | RFC XXXX |
| allow-certificate-get | boolean | true | RFC XXXX |

## New fields in the "meta" Object within a Directory Object

This document adds the following entry to the ACME Directory Metadata Fields:

| Field Name | Field Type | Reference |
|------------|------------|-----------|
| auto-renewal | object | RFC XXXX |

## Fields in the "auto-renewal" Object within a Directory Metadata Object
{: #iana-metadata-auto-renewal-registry}

The "ACME Directory Metadata Auto Renewal Fields" registry lists field names
that are defined for use in the JSON object included in the "auto-renewal"
field of an ACME directory "meta" object.

Template:

* Field name: The string to be used as a field name in the JSON object
* Field type: The type of value to be provided, e.g., string, boolean, array of
  string
* Reference: Where this field is defined

Initial contents: The fields and descriptions defined in {{capability-discovery}}.

| Field Name | Field Type | Reference |
|------------|------------|--------------|
| min-lifetime | integer | RFC XXXX |
| max-duration | integer | RFC XXXX |
| allow-certificate-get | boolean | RFC XXXX |

## Cert-Not-Before and Cert-Not-After HTTP Headers
{: #iana-http-headers}

The "Message Headers" registry should be updated with the following additional values:

|   Header Field Name   | Protocol | Status   | Reference |
|-----------------------|----------|----------|-----------|
| Cert-Not-Before       | http     | standard | RFC XXXX, {{fetching-certificates}} |
| Cert-Not-After        | http     | standard | RFC XXXX, {{fetching-certificates}} |

# Security Considerations

## No revocation

STAR certificates eliminate an important security feature of PKI which
is the ability to revoke certificates.  Revocation allows the
administrator to limit the damage done by a rogue node or an adversary
who has control of the private key.  With STAR certificates, expiration
replaces revocation so there is a timeliness issue.  To that end, see
also the discussion on clock skew in {{operational-cons-clocks}}.

It should be noted that revocation also has timeliness issues, because
both CRLs and OCSP responses have nextUpdate fields that tell relying parties (RPs) how
long they should trust this revocation data.  These fields are typically
set to hours, days, or even weeks in the future.  Any revocation that
happens before the time in nextUpdate goes unnoticed by the RP.

More discussion of the security of STAR certificates is available in
{{Topalovic}}.

## Denial of Service Considerations

STAR adds a new attack vector that increases the threat of denial of
    service attacks, caused by the change to the CA's behavior. Each STAR
    request amplifies the resource demands upon the CA, where one Order
    produces not one, but potentially dozens or hundreds of certificates,
    depending on the auto-renewal "lifetime" parameter. An attacker
    can use this property to aggressively reduce the
    auto-renewal "lifetime" (e.g. 1 sec.) jointly with other ACME
    attack vectors identified in Sec. 10 of {{RFC8555}}. Other collateral impact is
    related to the certificate endpoint resource where the client can
    retrieve the certificates periodically. If this resource is external to
    the CA (e.g. a hosted web server), the previous attack will be reflected to
    that resource.

Mitigation recommendations from ACME still apply, but some of them need
    to be adjusted. For example, applying rate limiting to the initial
    request, by the nature of the auto-renewal behavior cannot solve the
    above problem. The CA server needs complementary mitigation and
    specifically, it SHOULD enforce a minimum value on
    auto-renewal "lifetime". Alternatively, the CA can set an
    internal certificate generation processes rate limit.

## Privacy Considerations

In order to avoid correlation of certificates by account, if unauthenticated
GET is negotiated ({{certificate-get-nego}}) the recommendation in Section 10.5
of {{RFC8555}} regarding the choice of URL structure applies, i.e. servers SHOULD
choose URLs of certificate resources in a non-guessable way, for example using
capability URLs {{?W3C.WD-capability-urls-20140218}}.

# Acknowledgments

This work is partially supported by the European Commission under
Horizon 2020 grant agreement no. 688421 Measurement and Architecture
for a Middleboxed Internet (MAMI). This support does not imply endorsement.

Thanks to
Ben Kaduk,
Richard Barnes,
Roman Danyliw,
Jon Peterson,
Eric Rescorla,
Ryan Sleevi,
Sean Turner,
Alexey Melnikov,
Adam Roach,
Martin Thomson and
Mehmet Ersue
for helpful comments and discussions that have shaped this document.

--- back

# Document History

[[Note to RFC Editor: please remove before publication.]]

## draft-ietf-acme-star-10

IESG processing:

- More clarity on IANA registration (Alexey);
- HTTP header requirements adjustments (Adam);

IANA expert review:

- wrap HTTP-date in DQUOTEs in header values


## draft-ietf-acme-star-09

Richard and Ryan's review resulted in the following updates:

- STAR Order and Directory Meta attributes renamed slightly and grouped under two brand new "auto-renewal" objects;
- IANA registration updated accordingly (note that two new registries have been added as a consequence);
- Unbounded pre-dating of certificates removed so that STAR certs are never issued with their notBefore in the past;
- Changed "recurrent" to "autoRenewal" in error codes;
- Changed "recurrent" to "auto-renewal" in reference to Orders;
- Added operational considerations for HTTP caches.

## draft-ietf-acme-star-08

- Improved text on interaction with CT Logs, responding to Mehmet Ersue's review.

## draft-ietf-acme-star-07

- Changed the HTTP headers names and clarified the IANA registration, following feedback
from the IANA expert reviewer

## draft-ietf-acme-star-06

- Roman's AD review

## draft-ietf-acme-star-05

- EKR's AD review
- A detailed example of the timing of certificate issuance and predating
- Added an explicit client-side parameter for predating
- Security considerations around unauthenticated GET

## draft-ietf-acme-star-04

- WG last call comments by Sean Turner
- revokeCert interface handling
- Allow negotiating plain-GET for certs
- In STAR Orders, use star-certificate instead of certificate

## draft-ietf-acme-star-03

- Clock skew considerations
- Recommendations for "short" in the Web use case
- CT log considerations

## draft-ietf-acme-star-02

- Discovery of STAR capabilities via the directory object
- Use the more generic term Identifier Owner (IdO) instead of Domain Name Owner (DNO)
- More precision about what goes in the order
- Detail server side behavior on cancellation

## draft-ietf-acme-star-01

- Generalized the introduction, separating out the specifics of CDNs.
- Clean out LURK-specific text.
- Using a POST to ensure cancellation is authenticated.
- First and last date of recurrent cert, as absolute dates. Validity of certs in seconds. 
- Use RFC7807 "Problem Details" in error responses.
- Add IANA considerations.
- Changed the document's title.

## draft-ietf-acme-star-00

- Initial working group version.
- Removed the STAR interface, the protocol between NDC and DNO. What remains is only
the extended ACME protocol.

## draft-sheffer-acme-star-02

- Using a more generic term for the delegation client, NDC.
- Added an additional use case: public cloud services.
- More detail on ACME authorization.

## draft-sheffer-acme-star-01

- A terminology section.
- Some cleanup.

## draft-sheffer-acme-star-00

- Renamed draft to prevent confusion with other work in this space.
- Added an initial STAR protocol: a REST API.
- Discussion of CDNI use cases.

## draft-sheffer-acme-star-lurk-00

- Initial version.
