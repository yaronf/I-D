---
title: An ACME Profile for Generating Delegated STAR Certificates
abbrev: ACME STAR Delegation
docname: draft-ietf-acme-star-delegation-latest
category: std

ipr: trust200902
area: Security
workgroup: ACME
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
  consensus: false

author:
 -
    ins: Y. Sheffer
    name: Yaron Sheffer
    organization: Intuit
    email: yaronf.ietf@gmail.com
 -
    ins: D. López
    name: Diego López
    organization: Telefonica I+D
    email: diego.r.lopez@telefonica.com
 -
    ins: A. Pastor Perales
    name: Antonio Agustín Pastor Perales
    organization: Telefonica I+D
    email: antonio.pastorperales@telefonica.com
 -
    ins: T. Fossati
    name: Thomas Fossati
    organization: ARM
    email: thomas.fossati@arm.com

informative:
  json-schema-07:
    -: ta
    target: https://datatracker.ietf.org/doc/html/draft-handrews-json-schema-validation-01
    title: "JSON Schema Validation: A Vocabulary for Structural Validation of JSON"
    author:
      -
        ins: A. Wright
        name: Austin Wright
      -
        ins: H. Andrews
        name: Henry Andrews
      -
        ins: G. Luff
        name: Geraint Luff
    date: 2018

--- abstract

This memo proposes a profile of the ACME protocol that allows the owner of an
identifier (e.g., a domain name) to delegate to a third party access to a
certificate associated with said identifier.  A primary use case is that of a
CDN (the third party) terminating TLS sessions on behalf of a content provider
(the owner of a domain name).  The presented mechanism allows the owner of the
identifier to retain control over the delegation and revoke it at any time by
cancelling the associated STAR certificate renewal with the ACME CA.  Another
key property of this mechanism is it does not require any modification to the
deployed TLS ecosystem.

--- middle

# Introduction

This document is a companion document to {{!RFC8739}}.  To avoid duplication,
we give here a bare-bones description of the motivation for this solution.  For
more details and further use cases, please refer to the introductory sections
of {{!RFC8739}}.

An Identifier Owner (IdO) has agreements
in place with one or more NDC (Name Delegation Consumer) to use and attest its
identity.

In the primary use case the IdO is a content provider, and we consider a Content Delivery Network (CDN) provider contracted to
serve the content over HTTPS.  The CDN terminates the HTTPS connection at
one of its edge cache servers and needs to present its clients (browsers,
mobile apps, set-top-boxes) a certificate whose name matches the authority of
the URL that is requested, i.e., that of the IdO.  Understandably, some IdOs may balk at sharing their long-term private keys with another organization and,
equally, delegates would rather not have to handle other parties' long-term
secrets. Other relevant use cases are discussed in {{further-use-cases}}.

This document describes a profile of the ACME protocol {{!RFC8555}} that allows
the NDC to request from the IdO, acting as a profiled ACME server, a certificate for
a delegated identity - i.e., one belonging to the IdO.  The IdO then uses the
ACME protocol (with the extensions described in {{!RFC8739}}) to request
issuance of a STAR certificate for the same delegated identity.  The generated
short-term certificate is automatically renewed by the ACME Certification
Authority (CA), periodically fetched by the NDC and used to terminate HTTPS
connections in lieu of the IdO.  The IdO can end the delegation at any time by
simply instructing the CA to stop the automatic renewal and letting the
certificate expire shortly thereafter.

In case the delegated identity is a domain name, this document also provides a
way for the NDC to inform the IdO about the CNAME mappings that need to be
installed in the IdO's DNS zone to enable the aliasing of the delegated name,
thus allowing the complete name delegation workflow to be handled using a
single interface.

While the primary use case we address is delegation of STAR certificates,
the mechanism proposed here accommodates any certificate managed with
the ACME protocol. See {{non-star-delegation}} for details.

We note that other ongoing efforts address the problem of certificate delegation for TLS connections, specifically {{?I-D.ietf-tls-subcerts}} and {{?I-D.mglt-lurk-tls13}}. Compared to these other solutions, the current draft does not introduce additional latency to the TLS connection, nor does it require changes to the TLS network stack of either the client or the server.

## Terminology

IdO
: Identifier Owner, the owner of an identifier (e.g., a domain
  name) that needs to be delegated.

NDC
: Name Delegation Consumer, the entity to which the domain name is
  delegated for a limited time.  This is a CDN in the primary use
  case (in fact, readers may note the symmetry of the two
  acronyms).

CDN
: Content Delivery Network, a widely distributed network that
  serves the domain's web content to a wide audience at high
  performance.

STAR
: Short-Term, Automatically Renewed X.509 certificates.

ACME
: The IETF Automated Certificate Management Environment, a
  certificate management protocol.

CA
: A Certificate Authority that implements the ACME protocol. Synonymous with "ACME server".

## Conventions used in this document

{::boilerplate bcp14}

# Protocol Flow

This section presents the protocol flow.  For completeness, we include the ACME
profile proposed in this draft as well as the extended ACME protocol described
in {{!RFC8739}}.

## Preconditions
{: #proto-preconditions}

The protocol assumes the following preconditions are met:

- The IdO exposes an ACME server interface to the NDC(s) comprising the account
  management interface;
- The NDC has registered an ACME account with the IdO;
- NDC and IdO have agreed on a "CSR template" to use, including at a minimum:
  subject name (e.g., `somesite.example.com`), requested algorithms and key
  length, key usage, extensions (e.g., TNAuthList). The NDC is required to use
  this template for every CSR created under the same delegation;
- IdO has registered an ACME account with the Certificate Authority (CA)

Note that even if the IdO implements the ACME server role, it is not acting as
a CA: in fact, from the point of view of the certificate issuance process, the
IdO only works as a "policing" forwarder of the NDC's key-pair and is
responsible for completing the identity verification process towards the ACME
server.

## Overview

The interaction between the NDC and the IdO is governed by the profiled ACME
workflow detailed in {{sec-profile}}.  The interaction between the IdO and the
CA is ruled by ACME STAR {{!RFC8739}} as well as any other ACME extension that
applies (e.g., {{?I-D.ietf-acme-authority-token-tnauthlist}} for STIR).

The outline of the combined protocol is as follow ({{fig-endtoend}}):

- NDC sends an order Order1 for the delegated identifier to IdO;
- IdO creates an Order1 resource in state `ready` with a `finalize` URL;
- NDC immediately sends a finalize request (which includes the CSR) to the IdO;
- IdO verifies the CSR according to the agreed upon CSR template;
- If the CSR verification fails, Order1 is moved to an `invalid` state and
  everything stops;
- If the CSR verification is successful, IdO moves  Order1 to state
  `processing`, and sends a new Order2 (using its own account) for the delegated
  identifier to the CA;
- If the ACME STAR protocol fails, Order2 moves to `invalid` and the same state
  is reflected in Order1 (i.e., the NDC Order);
- If the ACME STAR run is successful (i.e., Order2 is `valid`), IdO copies the
  `star-certificate` URL from Order2 to Order1 and updates the Order1 state to
  `valid`.

The NDC can now download, install and use the short-term certificate bearing
the name delegated by the IdO.  This can continue until the STAR certificate
expires or the IdO decides to cancel the automatic renewal process with the CA.

Note that the interactive identifier authorization phase described in Section
7.5 of {{RFC8555}} is suppressed on the NDC-IdO side because the delegated
identity contained in the CSR presented to the IdO is validated against the
configured CSR template ({{sec-profile-dele-config}}).  Therefore, the NDC
sends the finalize request, including the CSR, to the IdO immediately after
Order1 has been acknowledged.  The IdO SHALL buffer a (valid) CSR until the
Validation phase completes successfully.

~~~ goat
{::include art/e2e-flow.ascii-art}
~~~
{: #fig-endtoend title="End to end STAR delegation flow"}

## Delegated Identity Profile
{: #sec-profile}

This section defines a profile of the ACME protocol, to be used between the NDC
and IdO.

### Delegation Configuration
{: #sec-profile-dele-config}

#### Account Object Extensions

An NDC identifies itself to the IdO as an ACME account.  The IdO can delegate
multiple names to a NDC, and these configurations are described through
`delegation` objects associated with the NDC's Account object on the IdO.

As shown in {{fig-account-object}}, the ACME account resource on the IdO is
extended with a new `delegations` attribute:

- delegations (required, string): A URL from which a list of delegations
  configured for this account can be fetched via a POST-as-GET request.

~~~
{
  "status": "valid",
  "contact": [
    "mailto:delegation-admin@ido.example"
  ],
  "termsOfServiceAgreed": true,
  "orders": "https://example.com/acme/orders/rzGoeA",
  "delegations": "https://acme.ido.example/acme/delegations/adFqoz"
}
~~~
{: #fig-account-object title="Example Account object with delegations"}

#### Delegation Objects

This profile extends the ACME resource model with a new read-only delegation
object that represents a delegation configuration that applies to a given NDC.

A delegation object contains the CSR template (see {{sec-csr-template}}) that
applies to that delegation, and optionally any related CNAME mapping for the
delegated identifiers.  Its structure is as follows:

- csr-template (required, object): CSR template as defined in
  {{sec-csr-template}}.
- cname-map (optional, object): a map of FQDN pairs.  In each pair, the name is
  the delegated identifier, the value is the corresponding IdO name that is
  aliased in the IdO's zone file to redirect the resolvers to the delegated
  entity.  Both names and values MUST be FQDNs with a terminating '.'.
  This field is only meaningful for identifiers of type `dns`.

An example delegation object is shown in {{fig-configuration-object}}.

~~~
{::include CSR-template/example-configuration-object.json}
~~~
{: #fig-configuration-object title="Example Delegation Configuration object"}

In order to indicate which specific delegation applies to the requested
certificate a new `delegation` attribute is added to the identifier in the
Order object on the NDC-IdO side (see {{sec-profile-order-journey}}).  The
value of this attribute is the URL pointing to the delegation configuration
object that is to be used for this certificate request.  If the `delegation`
attribute in the Order object contains a URL that does not correspond to a
configuration available to the requesting NDC, the IdO MUST return an error
response with status code 403 (Forbidden) and type
`urn:ietf:params:acme:error:unknownDelegation`.

### Order Object Transmitted from NDC to IdO and to ACME Server
{: #sec-profile-order-journey}

The Order object created by the NDC:

- MUST have the delegated name as the identifier value with a `delegation`
  attribute indicating the configuration used for the identifier.

Besides, when delegation is for a STAR certificate, the Order:

- MUST NOT contain the `notBefore` and `notAfter` fields;
- MUST contain an `auto-renewal` object and inside it, the fields
  listed in Section 3.1.1 of {{!RFC8739}}.

~~~
POST /acme/new-order HTTP/1.1
Host: acme.ido.example
Content-Type: application/jose+json

{
  "protected": base64url({
    "alg": "ES256",
    "kid": "https://acme.ido.example/acme/acct/evOfKhNU60wg",
    "nonce": "5XJ1L3lEkMG7tR6pA00clA",
    "url": "https://acme.ido.example/acme/new-order"
  }),
  "payload": base64url({
    "identifiers": [
      {
        "type": "dns",
        "value": "abc.ndc.ido.example.",
        "delegation":
           "https://acme.ido.example/acme/delegations/adFqoz/2"
      }
    ],
    "auto-renewal": {
      "end-date": "2020-04-20T00:00:00Z",
      "lifetime": 345600,          // 4 days
      "allow-certificate-get": true
    }
  }),
  "signature": "H6ZXtGjTZyUnPeKn...wEA4TklBdh3e454g"
}
~~~

The Order object that is created on the IdO:

- MUST start in the `ready` state;
- MUST contain an `authorizations` array with zero elements;
- MUST contain the indicated `delegation` configurations.

Besides, when delegation is for a STAR certificate, the Order:

- MUST NOT contain the `notBefore` and `notAfter` fields.

~~~
{
  "status": "ready",
  "expires": "2019-05-01T00:00:00Z",

  "identifiers": [
   {
     "type": "dns",
     "value": "abc.ndc.ido.example.",
     "delegation":
        "https://acme.ido.example/acme/delegations/adFqoz/2"
   }
  ],

  "auto-renewal": {
    "end-date": "2020-04-20T00:00:00Z",
    "lifetime": 345600,
    "allow-certificate-get": true
  },

  "authorizations": [],

  "finalize": "https://acme.ido.example/acme/order/TO8rfgo/finalize"
}
~~~

The Order is then finalized by the NDC supplying the CSR containing the
delegated identifiers.  The IdO checks the provided CSR against the template
that applies to each delegated identifier, as described in
{{sec-csr-template-syntax}}.  If the CSR fails validation for any of the
identifiers, the IdO MUST return an error response with status code 403
(Forbidden) and an appropriate type, e.g., `rejectedIdentifier` or `badCSR`.
The error response SHOULD contain subproblems (Section 6.7.1 of {{RFC8555}})
for each failed identifier.  If the CSR is successfully validated, the Order
object status moves to `processing` and the twin ACME protocol instance is
initiated on the IdO-CA side.

The Order object created by the IdO:

- MUST copy the identifiers sent by the NDC and strip the `delegation`
  attribute;

Besides, when delegation is for a STAR certificate, the Order:

- MUST carry a copy of the `auto-renewal` object sent by the NDC and augment it
  with an `allow-certificate-get` attribute set to true.

Instead, when the delegation is for a non-STAR certificate, the Order:

- MUST include the `allow-certificate-get` attribute set to true.

When the validation of the identifiers has been successfully completed and the
certificate has been issued by the CA, the IdO:

- MUST move its Order resource status to `valid`.

Besides, when delegation is for a STAR certificate, the IdO:

- MUST copy the `star-certificate` field from the STAR Order.  The latter
  indirectly includes (via the NotBefore and NotAfter HTTP headers) the renewal
  timers needed by the NDC to inform its certificate reload logic.

Instead, when the delegation is for a non-STAR certificate, the IdO:

- MUST copy the `certificate` field from the Order, as well as `notBefore`
  and `notAfter` if these fields exist.

~~~
{
  "status": "valid",
  "expires": "2019-05-01T00:00:00Z",

  "identifiers": [
   {
     "type": "dns",
     "value": "abc.ndc.ido.example.",
     "delegation":
        "https://acme.ido.example/acme/delegations/adFqoz/2"
   }
  ],

  "auto-renewal": {
    "end-date": "2020-04-20T00:00:00Z",
    "lifetime": 345600,
    "allow-certificate-get": true
  },

  "authorizations": [],

  "finalize": "https://acme.ido.example/acme/order/TO8rfgo/finalize",

  "star-certificate": "https://acme.ca.example/acme/order/yTr23sSDg9"
}
~~~

If an identifier object of type `dns` was included, the IdO can add the
corresponding CNAME records to its zone, e.g.:

~~~
   abc.ndc.ido.example. CNAME abc.ndc.example.
~~~

### Capability Discovery

In order to help a client to discover support for this profile, the directory
object of an ACME server MUST contain the following attribute in the `meta`
field:

- delegation-enabled: boolean flag indicating support for the profile
  specified in this memo.  An ACME server that supports this delegation profile
  MUST include this key, and MUST set it to true.
  

The `delegation-enabled` flag may be specified regardless of the existence or
setting of the `auto-renewal` flag.

### On Cancellation

It is worth noting that cancellation of the ACME STAR certificate is a
prerogative of the IdO.  The NDC does not own the relevant account key on the
ACME server, therefore it can't issue a cancellation request for the STAR cert.
Potentially, since it holds the STAR certificate's private key, it could request the
revocation of a single STAR certificate.  However, STAR explicitly disables the
revokeCert interface.

## Delegation of Non-STAR Certificates
{: #non-star-delegation}

The mechanism defined here can be used to delegate regular ACME certificates
whose expiry is not "short term".

To allow delegation of non-STAR certificates, this document allows use of
`allow-certificate-get` directly in the Order object and independently of the
`auto-renewal` object, so that the NDC can fetch the certificate without having
to authenticate into the ACME server.

The following differences exist between STAR and non-STAR certificate delegation:

* With STAR certificates, the `star-certificate` field is copied by the IdO;
  with non-STAR certificates, the `certificate` field is copied.
* The `auto-renewal` object is not used (either in the request or response) for
  non-STAR certificates. The field `allow-certificate-get` MUST be included in
  the order object, and its value MUST be `true`.
* The `notBefore` and `notAfter` order fields are omitted only in STAR
  certificates.

When delegating a non-STAR certificate, standard certificate revocation still
applies. The ACME certificate revocation endpoint is explicitly unavailable for
STAR certificates but it is available for all other certificates. We note that
according to Sec. 7.6 of {{RFC8555}}, the revocation endpoint can be used with
either the account keypair, or the certificate keypair. In other words, the NDC
would be able to revoke the certificate. However, given the trust relationship
between NDC and IdO expected by the delegation trust model
({{sec-trust-model}}) as well as the lack of incentives for the NDC - which,
doing so, would create a self-inflicted DoS - this does not represent a
security risk.

## Proxy Behavior

There are cases where the ACME Delegation flow should be proxied, such as the
use case described in {{sec-cdni-dele}}. This section describes the behavior of
such proxies.

An entity implementing the IdO server role - an "ACME Delegation server" - can
decide, on a per-identity case, whether to act as a proxy into another ACME
Delegation server, or to behave as an IdO and obtain a certificate directly.
The determining factor is whether it can successfully be authorized by
the ACME server for the identity associated with the certificate request.

The identities supported by each server and the disposition for each of them
are preconfigured.

Following is the proxy's behavior for each of the messages exchanged in the
ACME Delegation process:

* New-order request:
  * The complete `identifiers` object MUST be copied as-is.
  * Similarly, the `auto-renewal` object MUST be copied as-is.
* New-order response:
  * The `status`, `expires`, `authorizations`, `identifiers` and `auto-renewal`
    attributes/objects MUST be copied as-is.
  * The `finalize` URL is rewritten, so that the `finalize` request will be
    made to the proxy.
  * Similarly, the `Location` header MUST be rewritten to point to an Order object on the proxy.
  * And similarly, any `Link` relations.
* Get Order response:
  * The `status`, `expires`, `authorizations`, `identifiers` and `auto-renewal`
    attributes/objects MUST be copied as-is.
  * Similarly, the `star-certificate` URL MUST be copied as-is.
  * The `finalize` URL is rewritten, so that the `finalize` request will be
    made to the proxy.
  * The `Location` header MUST be rewritten to point to an Order object on the proxy.
  * Any `Link` relations MUST be rewritten to point to the proxy.
* Finalize request:
  * The CSR MUST be copied as-is.
* Finalize response:
  * The `Location` header, `Link` relations and the `finalize` URLs are rewritten as for Get Order.

We note that all the above messages are authenticated, and therefore each proxy
must be able to authenticate any subordinate server.

# CSR Template
{: #sec-csr-template}

The CSR template is used to express and constrain the shape of the CSR that the
NDC uses to request the certificate.  The CSR is used for every certificate
created under the same delegation.  Its validation by the IdO is a critical
element in the security of the whole delegation mechanism.

Instead of defining every possible CSR attribute, this document takes a
minimalist approach by declaring only the minimum attribute set and deferring
the registration of further, more specific, attributes to future documents.

## Template Syntax
{: #sec-csr-template-syntax}

The template is a JSON document. Each field denotes one of:

* A mandatory field, where the template specifies the literal value of that
  field. This is denoted by a literal string, such as
  `client1.ndc.ido.example.com`.
* A mandatory field, where the content of the field is defined by the client.
  This is denoted by `**`.
* An optional field, where the client decides whether the field is included in
  the CSR and if so, what its value is. This is denoted by `*`.

The NDC MUST NOT include in the CSR any fields, including any extensions, unless they are specified in the
template.

The structure of the template object is described using CDDL {{!RFC8610}} notation:

~~~
{::include CSR-template/template-schema.cddl}
~~~

The `subject` field and its subfields are mapped into the `subject` field of the CSR, as per {{RFC5280}}, Sec. 4.1.2.6. Other extension fields of the CSR template are mapped into the CSR according to the table in {{csr-template-registry}}.

When the CSR is received by the IdO, it MUST verify that the CSR is consistent
with the template that the IdO sent earlier. The IdO MAY enforce additional
constraints, e.g. by restricting field lengths.

## Example

The CSR template in {{fig-csr-template}} represents one possible CSR template
governing the delegation exchanges provided in the rest of this document.

~~~
{::include CSR-template/example-template.json}
~~~
{: #fig-csr-template title="Example CSR template"}

The template syntax is defined in {{csr-template-schema}}.

# Further Use Cases
{: #further-use-cases}

This non-normative section describes additional use cases that use STAR certificate
delegation in non-trivial ways.

## CDNI

{{?I-D.ietf-cdni-interfaces-https-delegation}} discusses several solutions
addressing different delegation requirements for the CDNI (CDN Interconnection)
environment.  This section discusses two of the stated requirements in the
context of the STAR delegation workflow.

This section uses specifically CDNI terminology, e.g. "uCDN" and "dCDN", as defined in {{?RFC7336}}.

### Multiple Parallel Delegates

In some cases the content owner (IdO) would like to delegate authority over a
web site to multiple NDCs (CDNs).  This could happen if the IdO has agreements
in place with different regional CDNs for different geographical regions, or if
a "backup" CDN is used to handle overflow traffic by temporarily altering some
of the CNAME mappings in place.  The STAR delegation flow enables this use case
naturally, since each CDN can authenticate separately to the IdO (via its own
separate account) specifying its CSR, and the IdO is free to allow or deny each
certificate request according to its own policy.

### Chained Delegation
{: #sec-cdni-dele}
In other cases, a content owner (IdO) delegates some domains to a large CDN
(uCDN), which in turn delegates to a smaller regional CDN, dCDN.  The IdO has a
contractual relationship with uCDN, and uCDN has a similar relationship with
dCDN.  However IdO may not even know about dCDN.

If needed, the STAR protocol can be chained to support this use case: uCDN
could forward requests from dCDN to IdO, and forward responses back to dCDN.
Whether such proxying is allowed is governed by policy and contracts between
the parties.

A mechanism is necessary at the interface between uCDN and dCDN by which the
uCDN can advertise:

- The namespace that is made available to the dCDN to mint its delegated names;
- The policy for creating the key material (allowed algorithms, minimum key
  lengths, key usage, etc.) that the dCDN needs to satisfy.

Note that such mechanism is provided by the CSR template.

#### Two-Level Delegation in CDNI

A User Agent (UA), browser or set-top-box, wants to fetch the video resource at
the following URI: `https://video.cp.example/movie`.  Redirection between
Content Provider (CP), upstream, and downstream CDNs is arranged as a
CNAME-based aliasing chain as illustrated in {{fig-cdni-dns-redirection}}.

~~~ goat
{::include art/cdni-dns-redirection.ascii-art}
~~~
{: #fig-cdni-dns-redirection title="DNS Redirection"}

Unlike HTTP based redirection, where the original URL is supplanted by the one
found in the Location header of the 302 response, DNS redirection is completely
transparent to the User Agent.  As a result, the TLS connection to the dCDN
edge is done with an SNI equal to the `host` in the original URI - in the
example, `video.cp.example`.  So, in order to successfully complete the
handshake, the landing dCDN node has to be configured with a certificate whose
subjectAltName matches `video.cp.example`, i.e., a Content Provider's name.

{{fig-cdni-flow}} illustrates the cascaded delegation flow that allows dCDN to
obtain a STAR certificate that bears a name belonging to the Content Provider
with a private key that is only known to the dCDN.

~~~ goat
{::include art/cdni-delegation.ascii-art}
~~~
{: #fig-cdni-flow title="Two levels delegation in CDNI"}

uCDN is configured to delegate to dCDN, and CP is configured to delegate to uCDN, both as defined in {{sec-profile-dele-config}}.

1. dCDN requests CDNI path metadata to uCDN;
2. uCDN replies with, among other CDNI metadata, the STAR delegation
   configuration, which includes the delegated Content Provider's name;
3. dCDN creates a key-pair and the CSR with the delegated name.  It then places
   an order for the delegated name to uCDN;
4. uCDN forwards the received order to the Content Provider (CP);
5. CP creates an order for a STAR certificate and sends it to the CA.  The
   order also requests unauthenticated access to the certificate resource;
6. After all authorizations complete successfully, the STAR certificate is
   issued;
7. CP notifies uCDN that the STAR certificate is available at the order's
   star-certificate URL;
8. uCDN forwards the information to dCDN.  At this point the ACME signalling is
   complete;
9. dCDN requests the STAR certificate using unauthenticated GET from the ACME
   server;
10. the CA returns the certificate.  Now dCDN is fully configured to handle
    HTTPS traffic in-lieu of the Content Provider.

Note that 9. and 10. repeat until the delegation expires or is terminated.

## STIR

As a second use case, we consider the delegation of credentials in the STIR
ecosystem  {{?I-D.ietf-stir-cert-delegation}}.

In the STIR `delegated` mode, a service provider SP2 - the NDC - needs to sign
PASSPorT’s {{?RFC8225}} for telephone numbers (e.g., TN=+123) belonging to
another service provider, SP1 - the IdO.  In order to do that, SP2 needs a STIR
certificate, and private key, that includes TN=+123 in the TNAuthList
{{?RFC8226}} certificate extension.

In details ({{fig-stir-flow}}):

1. SP1 and SP2 agree on the configuration of the delegation - in particular,
   the CSR template that applies;
2. SP2 generates a private/public key-pair and sends a CSR to SP1 requesting
   creation of a certificate with: SP1 name, SP2 public key, and a TNAuthList
   extension with the list of TNs that SP1 delegates to SP2.  (Note that the
   CSR sent by SP2 to SP1 needs to be validated against the CSR template
   agreed upon in step 1.);
3. SP1 sends an Order for the CSR to the CA;
4. Subsequently, after the required TNAuthList authorizations are successfully
   completed, the CA moves the Order to a "valid" state; at the same
   time the star-certificate endpoint is populated.
5. The Order contents are forwarded from SP1 to SP2 by means of the paired
   "delegation" Order.
6. SP2 dereferences the star-certificate URL in the Order to fetch the rolling
   STAR certificate bearing the delegated identifiers.

~~~ goat
{::include art/stir-delegation.ascii-art}
~~~
{: #fig-stir-flow title="Delegation in STIR"}

As shown, the STAR delegation profile described in this document applies
straightforwardly, the only extra requirement being the ability to instruct the
NDC about the allowed TNAuthList values.  This can be achieved by a simple
extension to the CSR template.

# IANA Considerations

[[RFC Editor: please replace XXXX below by the RFC number.]]

## New ACME Identifier Object Fields

This document requests that IANA create the following new registry under the
Automated Certificate Management Environment (ACME) Protocol:

* ACME Identifier Object Fields

This registry is administered under a Specification Required policy
{{!RFC8126}}.

The "ACME Identifier Object Fields" registry lists field names that are
defined for use in the ACME identifier object.

Template:

* Field name: The string to be used as a field name in the JSON object
* Field type: The type of value to be provided, e.g., string, boolean, array of
  string
* Reference: Where this field is defined

| Field Name | Field Type | Reference |
|------------|------------|-----------|
| type | string | Section 7.1.3 of RFC 8555 |
| value | string | Section 7.1.3 of RFC 8555 |
| delegation | string | RFC XXXX |

Note: this registry was not created at the time {{RFC8555}} was standardized
likely because it was not anticipated that the identifier object would be
extended.  It is retrospectively introduced to record the status quo and allow
controlled extensibility of the identifier object.

## New Fields in the "meta" Object within a Directory Object

This document adds the following entries to the ACME Directory Metadata Fields registry:

| Field Name | Field Type | Reference |
|------------|------------|-----------|
| delegation-enabled | boolean | RFC XXXX |

## New Fields in the Order Object

This document adds the following entries to the ACME Order Object Fields registry:

| Field Name | Field Type | Configurable | Reference |
|------------|------------|--------------|-----------|
| allow-certificate-get | boolean | true | RFC XXXX |

## New Fields in the Account Object

This document adds the following entries to the ACME Account Object Fields registry:

| Field Name | Field Type | Requests | Reference |
|------------|------------|----------|-----------|
| delegations | string | none | RFC XXXX |

Note that the `delegations` field is only reported by ACME servers that have
`delegation-enabled` set to true in their meta Object.

## New Error Types

This document adds the following entries to the ACME Error Type registry:

| Type | Description | Reference |
|------|-------------|-----------|
| unknownDelegation | An unknown configuration is listed in the `delegations` attribute of the request Order | RFC XXXX |

## CSR Template Extensions
{: #csr-template-registry }

IANA is requested to establish a registry "STAR Delegation CSR Template
Extensions", with "Expert Review" as its registration procedure.

Each extension registered must specify:

* An extension name.
* An extension syntax, as a reference to a JSON Schema document that defines this extension.
* The extension's mapping into an X.509 certificate extension.

The initial contents of this registry are the extensions defined by the JSON
Schema document in {{csr-template-schema}}.

| Extension Name   | Extension Syntax            | Mapping to X.509 Certificate Extension                       |
| ---------------- | --------------------------- | ------------------------------------------------------------ |
| keyUsage         | See {{csr-template-schema}} | {{!RFC5280}}, Sec. 4.2.1.3                                   |
| extendedKeyUsage | See {{csr-template-schema}} | {{!RFC5280}}, Sec. 4.2.1.12                                  |
| subjectAltName   | See {{csr-template-schema}} | {{!RFC5280}}, Sec. 4.2.1.6 (note that only specific name formats are allowed: IPv4 address, IPv6 address, DNS name, email address) |

# Security Considerations

## Trust Model
{: #sec-trust-model}

The ACME trust model needs to be extended to include the trust relationship
between NDC and IdO.  Note that once this trust link is established, it
potentially becomes recursive.  Therefore, there has to be a trust relationship
between each of the nodes in the delegation chain; for example, in case of
cascading CDNs this is contractually defined.  Note that using standard
{{?RFC6125}} identity verification there are no mechanisms available to the IdO
to restrict the use of the delegated name once the name has been handed over to
the first NDC.

## Delegation Security Goal

Delegation introduces a new security goal: only an NDC that has been authorised
by the IdO, either directly or transitively, can obtain a certificate with an
IdO identity.

From a security point of view, the delegation process has five separate parts:

1. Enabling a specific third party (the intended NDC) to submit requests for
   delegated certificates;
2. Making sure that any request for a delegated certificate matches the
   intended "shape" in terms of delegated identities as well as any other
   certificate metadata, e.g., key length, x.509 extensions, etc.;
3. Serving the certificate back to the NDC;
4. A process for handling revocation of the delegation;
5. A process for handling revocation of the certificate itself.

The first part is covered by the NDC's ACME account that is administered by the
IdO, whose security relies on the correct handling of the associated key pair.
When a compromise of the private key is detected, the delegate MUST use the
account deactivation procedures defined in Section 7.3.6 of {{!RFC8555}}.

The second part is covered by the act of checking an NDC's certificate request
against the intended CSR template.  The steps of shaping the CSR template
correctly, selecting the right CSR template to check against the presented CSR,
and making sure that the presented CSR matches the selected CSR template are
all security relevant.

The third part builds on the trust relationship between NDC and IdO that is
responsible for correctly forwarding the certificate URL from the Order
returned by the ACME server.

The fourth part is associated with the ability of the IdO to unilaterally
remove the delegation object associated with the revoked identity, therefore
disabling any further NDC requests for such identity.  Note that, in more
extreme circumstances, the IdO might decide to disable the NDC account
thus entirely blocking any further interaction.

The fifth is covered by two different mechanisms, depending on the nature of
the certificate.  For STAR, the IdO shall use the cancellation interface
defined in Section 2.3 of {{RFC8739}}. For non-STAR, the certificate revocation
interface defined in Section 7.6 of {{RFC8555}}).

## New ACME Channels

Using the model established in Section 10.1 of {{!RFC8555}}, we can decompose
the interactions of the basic delegation workflow as shown in
{{fig-sec-channels}}.

~~~ goat
{::include art/channels.ascii-art}
~~~
{: #fig-sec-channels title="Delegation Channels Topology"}

The considerations regarding the security of the ACME Channel and Validation
Channel discussed in {{!RFC8555}} apply verbatim to the IdO/ACME server leg.
The same can be said for the ACME channel on the NDC/IdO leg.  A slightly
different set of considerations apply to the ACME Channel between NDC and ACME
server, which consists of a subset of the ACME interface comprising two API
endpoints: the unauthenticated certificate retrieval and, potentially, non-STAR
revocation via certificate private key.  No specific security considerations
apply to the former, but the privacy considerations in Section 6.3 of
{{!RFC8739}} do.  With regards to the latter, it should be noted that there is
currently no means for an IdO to disable authorising revocation based on
certificate private keys.  So, in theory, an NDC could use the revocation API
directly with the ACME server, therefore bypassing the IdO.  The NDC SHOULD NOT
directly use the revocation interface exposed by the ACME server unless failing
to do so would compromise the overall security, for example if the certificate
private key is compromised and the IdO is not currently reachable.

All other security considerations from {{!RFC8555}} and {{!RFC8739}} apply
as-is to the delegation topology.

## Restricting CDNs to the Delegation Mechanism

When a web site is delegated to a CDN, the CDN can in principle modify the web
site at will, create and remove pages. This means that a malicious or breached
CDN can pass the ACME (as well as common non-ACME) HTTPS-based validation
challenges and generate a certificate for the site. This is true regardless of
whether the CNAME mechanisms defined in the current document is used or not.

In some cases, this is the desired behavior: the domain owner trusts the CDN to
have full control of the cryptographic credentials for the site. The current
document however assumes that the domain owner only wants to delegate
restricted control, and wishes to retain the capability to cancel the CDN's
credentials at a short notice.

The following is a possible mitigation when the IdO wishes to ensure that a
rogue CDN cannot issue unauthorized certificates:

- The domain owner makes sure that the CDN cannot modify the DNS records for
  the domain.  The domain owner should ensure it is the only entity authorized
  to modify the DNS zone. Typically, it establishes a CNAME resource record
  from a subdomain into a CDN-managed domain.
- The domain owner uses a CAA record {{!RFC8659}} to restrict certificate
  issuance for the domain to specific CAs that comply with ACME and are known
  to implement {{!RFC8657}}.
- The domain owner uses the ACME-specific CAA mechanism {{!RFC8657}} to
  restrict issuance to a specific account key which is controlled by it, and
  MUST require "dns-01" as the sole validation method.

We note that the above solution may need to be tweaked depending on the exact
capabilities and authorisation flows supported by the selected CAs.

# Acknowledgments

We would like to thank the following people who contributed significantly to this document with their review comments and design proposals: Roman Danyliw, <contact fullname="Frédéric" asciiFullname="Frederic"/> Fieau, Sanjay Mishra, Jon Peterson, Ryan Sleevi, Emile Stephan.

This work is partially supported by the European Commission under Horizon 2020
grant agreement no. 688421 Measurement and Architecture for a Middleboxed
Internet (MAMI). This support does not imply endorsement.

--- back


# Document History

[[Note to RFC Editor: please remove before publication.]]

## draft-ietf-acme-star-delegation-05

* Detailed AD review by Roman Danyliw.
* Some comments that were left unaddressed in Ryan Sleevi's review.
* Numerous other edits for clarity and consistency.

## draft-ietf-acme-star-delegation-04

- Delegation of non-STAR certificates.
- More IANA clarity, specifically on certificate extensions.
- Add delegation configuration object and extend account and order objects
  accordingly.
- A lot more depth on Security Considerations.

## draft-ietf-acme-star-delegation-03

- Consistency with the latest changes in the base ACME STAR document,
  e.g. star-delegation-enabled capability renamed and moved.
- Proxy use cases (recursive delegation) and the definition of proxy behavior.
- More detailed analysis of the CDNI and STIR use cases, including
  sequence diagrams.

## draft-ietf-acme-star-delegation-02

- Security considerations: review by Ryan Sleevi.
- CSR template simplified: instead of being a JSON Schema document itself,
  it is now a simple JSON document which validates to a JSON Schema.

## draft-ietf-acme-star-delegation-01

- Refinement of the CDNI use case.
- Addition of the CSR template (partial, more work required).
- Further security considerations (work in progress).

## draft-ietf-acme-star-delegation-00

- Republished as a working group draft.

## draft-sheffer-acme-star-delegation-01

- Added security considerations about disallowing CDNs from issuing
  certificates for a delegated domain.

## draft-sheffer-acme-star-delegation-00

- Initial version, some text extracted from draft-sheffer-acme-star-requests-02

# CSR Template Schema
{: #csr-template-schema}

Following is a JSON Schema definition of the CSR template. The syntax used is that of draft 7 of JSON Schema, which is documented in {{json-schema-07}}. Note that later versions of this (now expired) draft describe later versions of the JSON Schema syntax. At the time of writing, a stable reference for this syntax is not yet available, and we have chosen to use the draft version which is currently best supported by tool implementations.

While the CSR template must follow the syntax defined here, neither the IdO nor
the NDC are expected to validate it at run-time.

~~~
{::include CSR-template/template-schema.json}
~~~
