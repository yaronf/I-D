---
title: Use of Short-Term, Automatically-Renewed Certificates to address the LURK problem
abbrev: ACME STAR LURK
docname: draft-sheffer-acme-star-lurk-latest
date: 2016-09-17
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
    email: diego@telefonica.es
 -
    ins: O. Gonzalez de Dios
    name: Oscar Gonzalez de Dios
    organization: Telefonica I+D
    email: oscar.gonzalezdedios@telefonica.com
 -
    ins: T. Fossati
    name: Thomas Fossati
    organization: Nokia
    email: thomas.fossati@nokia.com

normative:
  RFC2119:
  I-D.ietf-acme-acme:

informative:
  I-D.iab-web-pki-problems:
  I-D.cairns-tls-session-key-interface:
  I-D.erb-lurk-rsalg:

--- abstract

This memo proposes two mechanisms that work in concert to address the LURK problem statement, allowing a third party (e.g., a content delivery network) to terminate TLS sessions on behalf of a domain name owner (e.g., a content provider).

The proposed mechanisms are:

1. An extension to the ACME protocol to enable the issuance of short-term and automatically renewed certificates, and
2.  A protocol that allows a domain name owner to delegate to a third party control over a certificate that bears its own name.

It should be noted that these are in fact independent building blocks that could be used separately to solve completely different problems.
--- middle

# A Solution for the HTTPS CDN Use Case

A content provider, and Domain Name Owner (DNO), has agreements in place with one or more Content Delivery Networks (CDN) that are contracted to serve its contents over HTTPS.  The CDN terminates the HTTPS connection at one of its edge cache servers and needs to present its clients (browsers, set-top-boxes) a certificate whose name matches the authority of the URL that is requested, i.e. that of the DNO.  However, many DNOs balk at sharing their long-term private keys with another organization and, equally, CDN providers would rather not have to handle other parties’ long-term secrets.

This document proposes a solution to the above problem that involves the use of short-term certificates with a DNO’s name on them, and a scheme for handling the naming delegation from the DNO to the CDN.  The generated short-term credentials are automatically renewed by an ACME Certification Authority (CA) and routinely rotated by the CDN on its edge cache servers.  The DNO can end the delegation at any time by simply instructing the CA to stop the automatic renewal and let the certificate expire shortly after.

Using short-term certificates makes revocation cheap and effective [TODO REF] in case of key compromise or of termination of the delegation; seamless certificate issuance and renewal enable the level of workflow automation that is expected in todays’ cloud environments.  Also, compared to other keyless-TLS solutions ({{I-D.cairns-tls-session-key-interface}}, {{I-D.erb-lurk-rsalg}}), the proposed approach doesn’t suffer from scalability issues or increase in connection setup latency, while requiring virtually no changes to existing COTS caching software in order to work.

(see section 3.2.1 of {{I-D.iab-web-pki-problems}}).

# Conventions used in this document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in {{RFC2119}}.

# Protocol Flow

The following subsections describe the preconditions ({{proto-preconditions}}), and the three main phases of the protocol flow:

- Bootstrap: the CDN requests DNO the delegation of a specific name and in turn DNO asks an ACME CA to create the corresponding short-term and auto-renewed (STAR) certificate ({{proto-bootstrap}});
- Auto-renewal: the ACME CA periodically re-issues the short-term certificate and posts it to a public URL ({{proto-auto-renewal}});
- Termination: the DNO (indirectly) ceases the name delegation by explicitly requesting the ACME CA to stop the automatic renewal of the certificate ({{proto-termination}}).



## Preconditions
{: #proto-preconditions}

The protocol assumes the following preconditions are met:

- A mutually authenticated channel between CDN and DNO pre-exists.  All LURK exchanges between CDN and DNO are run over this channel, which provides the guarantee that the LURK requests are authentic [^1]{:tf: source="tf"}.
- CDN and DNO have agreed on a “CSR template” to use, including at a minimum:
  - Subject name (e.g., “somesite.DNO.com”),
  - Validity (e.g., 24 to 72 hours),
  - Requested algorithms,
  - Key length,
  - Key usage. 

  The CDN is bound to use this template for every CSR created under the same delegation.
- DNO has registered through the ACME interface exposed by the Certificate Authority (CA) using the usual ACME registration procedure. The DNO shall, at the registration stage, query the ACME server for the supported STAR capabilities – for example: the minimum validity period of the issued certificate, the maximum duration of the automatic renewal process (either as a maximum number of renewal events, or as its maximum absolute life-span).

[^1]: Note that, under this assumption, the key used to authenticate the CDN to the DNO becomes a critical asset for the security of the proposed protocol, and that certain interactions (e.g., CSR submission) might require a stronger authentication mechanism.  For example, stacking a further authentication factor on top of CDN's LURK key would allow to distinguish an attacker that has only managed to successfully attack the CDN’s LURK key from the legit CDN.


## Bootstrap
{: #proto-bootstrap}

CDN generates a key-pair, wraps it into a Certificate Signing Request (CSR) according to the agreed CSR template, and sends it to the DNO over the LURK interface.  The DNO uses the CDN identity provided on the LURK channel to lookup the CSR template that applies to the requesting CDN and decides whether or not to accept the request.  (TBD: This is probably a case that would require a further authentication stage over the one provided by the mutual-authenticated LURK channel?)  Assuming everything is in order, it then “proxies” the CDN request to the ACME CA by means of the usual ACME application procedure. Specifically, DNO requests the CA a STAR certificate, i.e., one that:

- Has a short validity (e.g., 24 to 72 hours);
- Is automatically renewed by the CA for a certain period of time;
- Is downloadable from a (highly available) public link without requiring any special authorisation.

Other than that, the ACME protocol flows as normal between DNO and CA, in particular DNO is responsible for satisfying the requested ACME challenges until the CA is happy to issue the requested certificate.
The DNO is given back a unique identifier for the issued STAR certificate to be used in subsequent interaction with the CA (e.g., if the certificate needs to be terminated.)

Concurrently, an OK result has been sent back to the CDN with an endpoint to poll for completion of the certificate generation process. 

The bootstrap phase ends when the DNO obtains the OK from the ACME CA and posts the certificate’s URL to the “completion endpoint” where the CDN can retrieve it.  The information that is passed on to the CDN at this stage also includes details about how much time before the certificate expires can the CDN expect the replacement to be ready.

~~~~~~~~~~
TODO ASCII ART
~~~~~~~~~~
{: #figprotoboot title="Bootstrap"}


## Refresh
{: #proto-auto-renewal}

CA automatically re-issues the certificate (reusing the same CSR) before it expires and publishes it to the URL that the CDN has come to know at the end of the bootstrap phase.  The CDN downloads and installs it.  This process goes on until either:

- DNO terminates the delegation, or
- Automatic renewal expires.

~~~~~~~~~~
TODO ASCII ART
~~~~~~~~~~
{: #figprotorefresh title="Auto renewal"}


## Termination
{: #proto-termination}

DNO requests termination of the STAR certificate by including the previously obtained identifier in a STAR certificate termination request to the ACME interface.
After CA receives & verifies the request, it shall:

- Cancel the automatic renewal machinery for the LURK certificate;
- Change the certificate resource to return an error indicating the termination of the delegation to external clients, including the CDN;

Note that it is not necessary to explicitly revoke the short-term certificate.

~~~~~~~~~~
TODO ASCII ART
~~~~~~~~~~
{: #figprototerm title="Termination"}

# Security Considerations

- CDN's client certificate key is first order security asset and MUST be protected.  Absent 2FA/MFA, an attacker that can compromise the key might be able to obtain certificates bearing DNO’s identity.
- Consider collusion of two or more CDNs with contracts with the same DNO (?)

--- back
