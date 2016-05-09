---
title: "Delegating ACME Certificates to a CDN" 
abbrev: ACME Delegation
docname: draft-sheffer-lurk-acme-delegation-latest
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
    email: yaronf.ietf@gmail.com

normative:
  RFC2119:

informative:
  I-D.draft-landau-acme-caa-00:
  I-D.draft-cairns-tls-session-key-interface-01:

--- abstract

An organization that owns web content often prefers to delegate hosting of this content
to a Content Delivery Network (CDN). To serve HTTP content securely, it needs to be
protected with TLS. This document proposes a way for the CDN
to request constrained certificates so that it can serve web content under
the content owner's name, without having the owner's
long term certificate.

--- middle

# Introduction

Content owners frquently prefer a Content Delivery Network (CDN)
to host their content. CDNs have large scale networks, and are often
able to serve large scale traffic to a global audience.

To protect this traffic, the CDN uses HTTPS and presents a certificate that
is typically required to bear the content owner's name. However, many content owners
balk at sharing their long-term private keys with another organization. This
document proposes a way for the CDN to obtain
short-term credentials, allowing the content owner to revoke this authority at short
notice.

We note that there are alternative solutions to this problem:

* The CDN
could contact the content owner on each TLS handshake and have the content
owner assist in completing the TLS handshake, such a solution is described in
e.g. {{I-D.cairns-tls-session-key-interface}}.
* We could extend ACME by allowing the content owner to share an authorization
"ticket" with the CDN, with the CDN using it to obtain short-term
certificates directly from the ACME server. This alternative is possibly easier
to deploy than the one described in this document, but
it would require a non-trivial change to the ACME protocol. 

## Conventions used in this document
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in <xref target="RFC2119"/>.

# Overview

We define the interaction between the CDN and the content owner, where the CDN
requests a short-term certificate periodically, and the content owner obtains it
on the CDN's behalf and returns it to the CDN.

# Advantages

* Compared with solutions that require the CDN to have the content owner sign each
handshake, this solution does not require the content owner to set up its own
scalable infrastructure.
* Moreover, the need to scale the content owner's web service could result in the content
owner ending up by sharing the private keys with the CDN and abdicating its responsibility
for its own security.

# LURK Operations

This section lists the operations that the content owner needs to provide the CDN.

## Request a Certificate

```
POST ./well-known/lurk/certificate/1234 HTTP/1.1
```

with a text/plain content type, and a message body that contains a public key in PEM format.
There is no negotiation of key type (RSA or ECDSA), key length or validity dates,
and the client and server
must pre-coordinate these details.

Similarly, the server MUST be able to determine the FQDN to be included in the certificate
based on the authenticated client's identity.

The number at the end of the URI is a request ID, which MAY be generated randomly.

The server responds with one of:

* A "200 OK" status code, and response body containing the certificate, with the content type:
application/x-509-server-cert.
* A "201 Accepted" status code if the certificate is not yet ready. The CDN should poll
the content owner periodically (see below), but not more often than once every 5 seconds.
* Other responses if the request is not acceptable or not allowed.

## Poll for a Certificate

```
GET ./well-known/lurk/certificate/1234 HTTP/1.1
```

The server responds with one of:

* A "200 OK" status code, and response body containing the certificate, with the content type:
application/x-509-server-cert.
* A "204 No Content" status code if the certificate is not yet ready.
* Other responses if the request is not acceptable or not allowed.

Access to these resources MUST be protected by TLS.

The client cannot assume that the sever will cache the certificate beyond a few seconds
after it is first fetched.

Both requests MUST be authenticated, using one of the following methods:

* Mutual TLS authentication with a client certificate. This is the RECOMMENDED option.
* TLS with preshared secret authentication or TLS-SRP.
* TLS with HTTP-Basic or Digest authentication.

# Security Considerations

This section presents additional considerations beyond those strictly required by the protocol.

## Certificate Details

* It is RECOMMENDED to restrict the certificate's scope as much as possible.
Specifically, the certificate request SHOULD specify restrictive Key Usage.
* The RECOMMENDED validity period for certificates provisioned using this
mechanism is 3 days.

## Restricting CDNs to the Delegation Mechanism

Currently there are no standard methods for the content owner to ensure that
the CDN cannot issue a certificate through mechanisms other than the one described here.
The best solution currently being worked on would consist of several related
configurations:

 * Make sure that the CDN cannot modify the DNS records for the domain.
 Typically this would mean that the content owner established a CNAME record
 from a subdomain into a CDN-managed domain.
 * Restrict certificate issuance for the domain to specific CAs that comply
 with ACME. This assumes
 universal deployment of CAA by CAs, which is not the case yet.
 * Deploy ACME-specific methods to restrict issuance to a specific authorization
 key which is controlled by the content owner {{I-D.landau-acme-caa}}.
 
--- back

# Document History

## draft-sheffer-lurk-acme-delegation-00

Initial version.
