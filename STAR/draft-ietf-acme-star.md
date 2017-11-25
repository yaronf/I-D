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
    organization: Nokia
    email: thomas.fossati@nokia.com

normative:
  RFC2119:
  RFC3339:
  I-D.ietf-acme-acme:
  RFC7231:
  RFC7807:

informative:
  RFC7942:
  I-D.sheffer-acme-star-request:
  I-D.nir-saag-star:
  Topalovic:
    -: ta
    target: http://www.w2spconf.com/2012/papers/w2sp12-final9.pdf
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

--- abstract

This memo proposes an ACME extension to enable the issuance of short-term and automatically renewed certificates.

[RFC Editor: please remove before publication]

While the draft is being developed, the editor's version can be found at
https://github.com/yaronf/I-D/tree/master/STAR.

--- middle

# Introduction

The ACME protocol {{I-D.ietf-acme-acme}} automates the process of issuing a certificate to a Domain Name Owner (DNO).
Further ACME extensions [TODO-GENERIC-TOKEN] generalize the type of identifiers that an ACME server can issue certificates for.
Since the protocol defined in this document is agnostic as to the specific identifier type, we will use Identifier Owner (IO) instead of the more specific term DNO.

If the IO wishes to obtain a string of short-term certificates originating from the same private key (see {{Topalovic}} for the rationale), she must go through the whole ACME protocol each time a new short-term certificate is needed - e.g., every 2-3 days.
If done this way, the process would involve frequent interactions between the registration function of the ACME Certification Authority (CA) and the identity provider infrastructure (e.g.: DNS, web servers), therefore making the issuance of short-term certificates exceedingly dependent on the reliability of both.

This document presents an extension of the ACME protocol that optimizes this process by making short-term certificates first class objects in the ACME ecosystem.
Once the order for a string of short-term certificates is accepted, the CA is responsible for publishing the next certificate at an agreed upon URL before the previous one expires.  The IO can terminate the automatic renewal before the natural deadline, if needed - e.g., on key compromise.

For a more generic treatment of STAR certificates, readers are referred to {{I-D.nir-saag-star}}.

### Name Delegation Use Case

The proposed mechanism can be used as a building block of an efficient name-delegation protocol, for example one that exists between a CDN or a cloud provider and its users {{I-D.sheffer-acme-star-request}}, in a way that makes the delegator (i.e., the IO) in full control of the delegation by simply instructing the CA to stop the automatic renewal and letting the currently active certificate expire shortly thereafter.

## Terminology

IO
: Identifier Owner, the owner of an identifier, e.g. a domain name.

STAR
: Short-Term, Automatically Renewed X.509 certificates.

NDC
: Name Delegation Client, an entity to which the identifier owned by the IO is delegated for a limited time. This might be a CDN edge cache, a cloud provider's load balancer or Web Application Firewall (WAF).

## Conventions used in this document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in {{RFC2119}}.

# Protocol Flow

The following subsections describe the three main phases of the protocol:

- Bootstrap: the IO asks an ACME CA to create a short-term and automatically-renewed (STAR) certificate ({{proto-bootstrap}});
- Auto-renewal: the ACME CA periodically re-issues the short-term certificate and posts it
to a public URL ({{proto-auto-renewal}});
- Termination: the IO requests the ACME CA
to discontinue the automatic renewal of the certificate ({{proto-termination}}).

This diagram presents the entities that are (or may be) involved in the protocol and their interactions during the different phases.

~~~
                         Refresh
          . . . . . . . . . . . . . . . . . . . .  
      . '                                         ` v
   .----.        Bootstrap / Terminate         .---------.
   | IO |------------------------------------->| ACME CA |
   `----'                                      `---------'
      ^                  .- - -.                    ^
       ` . . . . . . . . : NDC : . . . . . . . . . '
            Request      `- - -'    Refresh
            Delegation
~~~

Note that there might be a distinct NDC entity (e.g., a CDN edge cache) that uses a separate channel to request the IO to set up a name delegation.  The protocol described in {{I-D.sheffer-acme-star-request}} might be used for this purpose.

## Bootstrap
{: #proto-bootstrap}

The IO, in its role as an
ACME client, requests the CA to issue a STAR certificate, i.e., one that:

- Has a short validity (e.g., 24 to 72 hours);
- Is automatically renewed by the CA for a certain period of time;
- Is downloadable from a (highly available) public link without requiring any special authorization.

Other than that, the ACME protocol flows as normal between IO and CA.
In particular, IO is responsible for satisfying the requested ACME
challenges until the CA is willing to issue the requested certificate.
Per normal ACME processing, the IO is given back an Order ID for the issued STAR
certificate to be used in subsequent interaction with the CA (e.g., if
the certificate needs to be terminated.)

The bootstrap phase ends when the IO obtains the OK from the ACME CA.

## Refresh
{: #proto-auto-renewal}

The CA automatically re-issues the certificate using the same CSR (and therefore the same name and public key) before it expires and publishes it to the URL that was returned to the IO at the end of the bootstrap phase.  The certificate user, which could be either the IO itself or a delegated third party, as described in {{I-D.sheffer-acme-star-request}}, obtains the certificate and uses it.

The refresh process ({{figprotorefresh}}) goes on until either:

- IO terminates the delegation, or
- Automatic renewal expires.

~~~~~~~~~~
     Certificate              ACME/STAR
        User                    Server
          |     Retrieve cert     |                     [...]
          |---------------------->|                      |
          |                       +------.              /
          |                       |      |             /
          |                       | Automatic renewal :
          |                       |      |             \
          |                       |<-----'              \
          |     Retrieve cert     |                      |
          |---------------------->|                   72 hours
          |                       |                      |
          |                       +------.              /
          |                       |      |             /
          |                       | Automatic renewal :
          |                       |      |             \
          |                       |<-----'              \
          |     Retrieve cert     |                      |
          |---------------------->|                   72 hours
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

The IO may request early termination of the STAR certificate by including
the Order ID in a certificate termination request to the ACME
interface, defined below.
After the CA receives and verifies the request, it shall:

- Cancel the automatic renewal process for the STAR certificate;
- Change the certificate publication resource to return an error
indicating the termination of the delegation to any external client.

Note that it is not necessary to explicitly revoke the short-term certificate.

~~~~~~~~~~
   STAR                    STAR                   ACME/STAR
   Client                  Proxy                  Server
     |                       |                       |
     |                       |  Terminate Order ID   |
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
     |                 Retrieve cert                 |
     +---------------------------------------------->|
     |                 Error: terminated             |
     |<----------------------------------------------+
     |                                               |
~~~~~~~~~~
{: #figprototerm title="Termination"}

# Protocol Details
{: #protocol-details}

This section describes the protocol details, namely the extensions
to the ACME protocol required to issue STAR certificates.

## ACME Extensions

This protocol extends the ACME protocol, to allow for recurrent orders.

### Extending the Order Resource

The Order resource is extended with the following attributes:

    {
        "recurrent": true,
        "recurrent-start-date": "2016-01-01T00:00:00Z",
        "recurrent-end-date": "2017-01-01T00:00:00Z",
        "recurrent-certificate-validity": 604800
    }

- recurrent: MUST be "true" for STAR certificates.
- recurrent-start-date: the earliest date of validity of the first certificate issued, in {{RFC3339}} format.
This attribute is optional. When omitted, the start date is as soon as authorization is complete.
- recurrent-end-date: the latest date of validity of the last certificate issued, in {{RFC3339}} format.
- recurrent-certificate-validity: the maximum validity period of each STAR certificate,
an integer that denotes a number of seconds.

These attributes are included in a POST message when creating the order, as part of the "payload" encoded object.
They are returned when the order has been created, and the ACME server
MAY adjust them at will, according to its local policy.

ACME defines the following values for the Order resource's status: "invalid", "pending", "processing", "valid".
In the case of recurrent orders, the status MUST be "valid" as long as STAR certificates are being issued.
We add a new status value, "canceled", see below.

### Canceling a Recurrent Order
{: #protocol-details-canceling}

An important property of the recurrent Order is that it can be canceled by the IO,
with no need for certificate
revocation. To cancel the Order, the ACME client sends a POST:

    POST /acme/order/1 HTTP/1.1
    Host: acme-server.example.org
    Content-Type: application/jose+json

    {
     "protected": base64url({
       "alg": "ES256",
       "kid": "https://example.com/acme/acct/1",
       "nonce": "5XJ1L3lEkMG7tR6pA00clA",
       "url": "https://example.com/acme/order/1"
     }),
     "payload": base64url({
       "status": "canceled"
     }),
     "signature": "H6ZXtGjTZyUnPeKn...wEA4TklBdh3e454g"
    }

The server MUST NOT issue any additional certificates for this Order, beyond the certificate that is available for collection at the time of deletion.  Immediately after the Order is canceled, the server SHOULD respond with 403 (Forbidden) to any requests to the certificate endpoint.  The response SHOULD provide additional information using a problem document {{RFC7807}} with type "urn:ietf:params:acme:error:recurrentOrderCanceled".

## Capability Discovery

In order to support the discovery of STAR capabilities, The directory object of an ACME STAR server MUST contain the following attributes inside the "meta" field:

- star-capable: boolean indicating STAR support.  An ACME STAR server MUST set this key to true.
- star-min-cert-validity: minimum acceptable value for recurrent-certificate-validity, in seconds.
- star-max-renewal: maximum delta between recurrent-end-date and recurrent-start-date, in seconds.

## Fetching the Certificates
{: #fetching-certificates}

The certificate is fetched from the certificate endpoint, as per {{I-D.ietf-acme-acme}}, Section 7.4.2.

~~~
   GET /acme/cert/asdf HTTP/1.1
   Host: acme-server.example.org
   Accept: application/pkix-cert

   HTTP/1.1 200 OK
   Content-Type: application/pem-certificate-chain
   Link: <https://example.com/acme/some-directory>;rel="index"
   Not-Before: Mon, 1 Feb 2016 00:00:00 GMT
   Not-After: Mon, 8 Feb 2016 00:00:00 GMT

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

The Server SHOULD include the "Not-Before" and "Not-After" headers. When they exist, they MUST be equal
to the respective fields inside the certificate. Their format is "HTTP-date" as defined in Section 7.1.1.2 of {{RFC7231}}.
Their purpose is to enable client implementations that do not parse the certificate.

To improve robustness, the next certificate MUST be made available by the ACME CA at the latest halfway through the lifetime of the currently active certificate.
It is worth noting that this has an implication in case of cancellation: in fact, from the time the next certificate is made available, the cancellation is not completely effective until the latter also expires.

The server MUST NOT issue any additional certificates for this Order beyond its recurrent-end-date.

Immediately after the Order expires, the server SHOULD respond with 403 (Forbidden) to any requests to the certificate endpoint.  The response SHOULD provide additional information using a problem document {{RFC7807}} with type "urn:ietf:params:acme:error:recurrentOrderExpired".

# Operational Considerations

## Certificate Transparency (CT) Logs

TBD: larger logs and how to deal with them.

# Implementation Status

Note to RFC Editor: please remove this section before publication,
including the reference to {{RFC7942}}.

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

The implementation is constructed around 3 elements: Client STAR for NDC,
Proxy STAR for IO and Server ACME for CA. The communication between
them is over an IP network and the HTTPS protocol.

The software of the implementation is available at: https://github.com/mami-project/lurk

The following subsections offer a basic description, detailed information
is available in https://github.com/mami-project/lurk/blob/master/proxySTAR_v1/README.md

### ACME Server with STAR extension

This is a fork of the Let's Encrypt Boulder project that implements an ACME compliant CA.
It includes modifications to extend the ACME protocol as it is specified in this draft,
to support recurrent orders and cancelling orders. 

The implementation understands the new "recurrent" attributes as part of the Certificate
issuance in the POST request for a new resource.
An additional process "renewalManager.go" has been included in parallel that reads
the details of each recurrent request, automatically produces a "cron" Linux based task
that issues the recurrent certificates, until the lifetime ends or the order is canceled.
This process is also in charge of maintaining a fixed URI to enable the NDC to download certificates,
unlike Boulder's regular process of producing a unique URI per certificate.

### Proxy STAR

The Proxy STAR, has a double role as ACME client and STAR Server. The former is a fork of the EFF
Certbot project that implements an ACME compliant client with the STAR extension.
The latter is a basic HTTP REST API server.

The proxy STAR understands the basic API request with a server. The current implementation
of the API is defined in draft-sheffer-acme-star-request-00. Registration or order cancellation
triggers the modified Certbot client that requests, or cancels, the recurrent generation
of certificates using the STAR extension over ACME protocol.
The URI with the location of the recurrent certificate is delivered to the STAR client as a response.

## Level of Maturity

This is a prototype. 

## Coverage

Client STAR is not included in this implementation, but done by direct HTTP request with any open HTTP REST API tool.
This is expected to be covered as part of {{I-D.sheffer-acme-star-request}} implementation.

This implementation completely covers Proxy STAR and Server ACME with STAR extension 

## Version Compatibility

The implementation is compatible with version draft-ietf-acme-star-00. 
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

No security is enabled (iptables default policies are "accept" and all rules removed)
in this implementation to simplify and test the protocol. 

## Contact Information

See author details below.

# IANA Considerations

[[RFC Editor: please replace XXXX below by the RFC number.]]

## New ACME Error Types

This document adds the following entry to the ACME Error Type registry:

| Type | Description | Reference |
|------|-------------|-----------|
| recurrentOrderCanceled | The short-term certificate is no longer available because the recurrent order has been explicitly canceled by the IO | RFC XXXX |
| recurrentOrderExpired | The short-term certificate is no longer available because the recurrent order has expired | RFC XXXX |

## New ACME Order Object Fields

This document adds the following entries to the ACME Order Object Fields registry:

| Field Name | Field Type | Configurable | Reference |
|------------|------------|--------------|-----------|
| recurrent | string | true | RFC XXXX |
| recurrent-start-date | string | true | RFC XXXX |
| recurrent-end-date | string | true | RFC XXXX |
| recurrent-certificate-validity | string | true | RFC XXXX |

# Security Considerations

TBD

# Acknowledgments

This work is partially supported by the European Commission under
Horizon 2020 grant agreement no. 688421 Measurement and Architecture
for a Middleboxed Internet (MAMI). This support does not imply endorsement.

Thanks to
Jon Peterson and
Martin Thomson
for helpful comments and discussions that have shaped this document.

--- back

# Document History

[[Note to RFC Editor: please remove before publication.]]

## draft-ietf-acme-star-02

- Discovery of STAR capabilities via the directory object
- Move from Domain Name Owner (DNO) to Identifier Owner (IO)

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
