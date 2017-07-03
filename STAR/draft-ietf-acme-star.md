---
title: Use of Short-Term, Automatically-Renewed (STAR) Certificates to Delegate Authority over Web Sites
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
  I-D.ietf-acme-acme:

informative:
  RFC6844:
  I-D.iab-web-pki-problems:
  I-D.cairns-tls-session-key-interface:
  I-D.erb-lurk-rsalg:
  I-D.ietf-acme-caa:
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
  CAB-CAA:
    -:
    target: "https://cabforum.org/2017/03/08/ballot-187-make-caa-checking-mandatory/"
    title: "Ballot 187 – Make CAA Checking Mandatory"
    date: March 2017
    author:
      -
        org: "CA/Browser Forum"

--- abstract

This memo proposes an ACME extension to enable the issuance of
short-term  and automatically renewed certificates.
This allows a domain name owner to delegate the use of
certificates to another party, while retaining the capability to
cancel this delegation at any time with no need to rely on certificate
revocation mechanisms.

[RFC Editor: please remove before publication]

While the draft is being developed, the editor's version can be found at
https://github.com/yaronf/I-D/tree/master/STAR.

--- middle

# Introduction: A Solution for the HTTPS CDN Use Case

A content provider (referred to in this document as Domain Name Owner,
DNO) has agreements in
place with one or more Content Delivery Networks (CDNs) that are
contracted to serve its content over HTTPS. The CDN terminates the
HTTPS connection at one of its edge cache servers and needs to present
its clients (browsers, set-top-boxes) a certificate whose name matches
the authority of the URL that is requested, i.e. that of the DNO.
However, many DNOs balk at sharing their long-term private keys with
another organization and, equally, CDN providers would rather not have
to handle other parties' long-term secrets.  This problem has been
discussed at the IETF under the LURK (limited use of remote keys) title.

This document proposes a solution to the above problem that involves
the use of short-term certificates with a DNO's name on them, and a
scheme for handling the naming delegation from the DNO to the CDN.
The generated short-term credentials are automatically renewed by an
ACME Certification Authority (CA) {{I-D.ietf-acme-acme}} and routinely
rotated by the CDN on its edge cache servers.  The DNO can end the
delegation at any time by simply instructing the CA to stop the
automatic renewal and let the certificate expire shortly thereafter.

Using short-term certificates makes revocation cheap and effective
{{Topalovic}} {{I-D.iab-web-pki-problems}} in case of key compromise
or of termination of the delegation; seamless certificate issuance and
renewal enable the level of workflow automation that is expected in
today's cloud environments.  Also, compared to other keyless-TLS
solutions {{I-D.cairns-tls-session-key-interface}}
{{I-D.erb-lurk-rsalg}}, the proposed approach doesn't suffer from
scalability issues or increase in connection setup latency, while
requiring virtually no changes to existing COTS caching software used
by the CDN.

This document describes the ACME extension. A companion document
[I-D.sheffer-acme-star-request] describes how the CDN can
request the DNO to initiate the protocol with the ACME server.

## Cloud Use Case

A similar use case is that of cloud infrastructure components, such as load balancers
and Web Application Firewalls (WAF). These components are typically provisioned with
the DNO's certificate, and similarly to the CDN use case, many organizations
would prefer to manage the private key only on their own cloud-based or
on-premise hosts, often on Hardware Security Modules (HSMs).

Here again, the STAR solution allows the DNO to delegate authority
over the domain to the cloud provider, with the ability to revoke
this authority at any time.

## Terminology

DNO
: Domain Name Owner, the owner of a domain that needs to be delegated.

NDC
: Name Delegation Consumer, the entity to which the domain name is delegated for a limited
time. This is often a CDN (in fact, readers may note the similarity of the two acronyms).

CDN
: Content Delivery Network, a widely distributed network
that serves the domain's web content to a wide audience at
high performance.

STAR
: Short-Term, Automatically Renewed X.509 certificates.

ACME
: The IETF Automated Certificate Management Environment, a certificate
management protocol.

CA
: A Certificate Authority that implements the ACME protocol.

## Conventions used in this document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in {{RFC2119}}.

# Protocol Flow

For clarity, we describe how the proposed ACME extension can be used
in a system that consists of
an NDC, an ACME Client (the DNO) and an ACME Server. Only the latter
part (ACME Client to ACME Server)
is in scope of this document.

The protocol flow can be split into two: a STAR interface, used by NDC
and DNO to agree on the name delegation, and the extended ACME
interface, used by DNO to obtain the short-term and automatically
renewed certificate from the CA, which is eventually consumed by the
NDC.  The latter is also used to terminate the delegation, if so needed.

Communication between
the NDC and the DNO (the STAR interface) is out of scope of this document.
It may take the form described in [I-D.sheffer-acme-star-request], some
other online protocol, or may even be through manual generation of the CSR.

The following subsections describe the three main phases of the protocol:

- Bootstrap: the DNO asks an ACME CA to create a corresponding short-term and auto-renewed (STAR) certificate,
possibly on a request from an NDC which is out of scope for this document;
- Auto-renewal: the ACME CA periodically re-issues the short-term certificate and posts it to a public URL ({{proto-auto-renewal}});
- Termination: the DNO (indirectly) stops name delegation by explicitly requesting the ACME CA to discontinue the automatic renewal of the certificate ({{proto-termination}}).

This diagram presents the entities involved in the protocol and their interactions during the different phases.

~~~~~~~~~~~
                                 +-----------------+
                                 |    STAR Proxy   |
                                 |      (DNO)      |
                       Bootstrap +-----------------+ Bootstrap
                     +---------->+  STAR  |  ACME  +-----------+
                     |           | Server | Client | Terminate |
                     |           +--------+--------+           |
                     |                                         v
                 +--------+                                +--------+
                 |  STAR  |            Refresh             |  ACME  |
                 | Client +------------------------------->| Server |
                 | (NDC)  |                                |  (CA)  |
                 +--------+                                +--------+
~~~~~~~~~~~~

## Bootstrap
{: #proto-bootstrap}

The DNO, in its role as an
ACME client, requests the CA to issue a STAR certificate, i.e., one that:

- Has a short validity (e.g., 24 to 72 hours);
- Is automatically renewed by the CA for a certain period of time;
- Is downloadable from a (highly available) public link without requiring any special authorization.

Other than that, the ACME protocol flows as normal between DNO and CA,
in particular DNO is responsible for satisfying the requested ACME
challenges until the CA is willing to issue the requested certificate.
Per normal ACME processing, the DNO is given back an Order ID for the issued STAR
certificate to be used in subsequent interaction with the CA (e.g., if
the certificate needs to be terminated.)

The bootstrap phase ends when the DNO obtains the OK from the ACME CA.

## Refresh
{: #proto-auto-renewal}

The CA automatically re-issues the certificate (using the same CSR)
before it expires and publishes it to the URL that the NDC has come to
know at the end of the bootstrap phase.  The NDC downloads and
installs it. This process goes on until either:

- DNO terminates the delegation, or
- Automatic renewal expires.

~~~~~~~~~~
        STAR                    ACME/STAR
        Client                  Server
          |     Retrieve cert     |                     [...]
          |<--------------------->|                      |
          |                       +------.              /
          |                       |      |             /
          |                       | Automatic renewal :
          |                       |      |             \
          |                       |<-----'              \
          |     Retrieve cert     |                      |
          |<--------------------->|                   72 hours
          |                       |                      |
          |                       +------.              /
          |                       |      |             /
          |                       | Automatic renewal :
          |                       |      |             \
          |                       |<-----'              \
          |     Retrieve cert     |                      |
          |<--------------------->|                   72 hours
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

The DNO may request early termination of the STAR certificate by including
the Order ID in a certificate termination request to the ACME
interface, defined below.
After the CA receives and verifies the request, it shall:

- Cancel the automatic renewal process for the STAR certificate;
- Change the certificate publication resource to return an error
indicating the termination of the delegation to external clients, including the NDC.

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

This section describes the protocol's details, namely the extensions
to the ACME protocol required to issue STAR certificates.



## ACME Extensions between Proxy and Server

This protocol extends the ACME protocol, to allow for recurrent orders.

### Extending the Order Resource

The Order resource is extended with the following attributes:

    {
        "recurrent": true,
        "recurrent-total-lifetime": 365, // requested lifetime of the
                                         // recurrent registration, in days
        "recurrent-certificate-validity": 7
           // requested validity of each certificate, in days
    }

These attributes are included in a POST message when creating the order, as part of the "payload" encoded object.
They are returned when the order has been created, and the ACME server
MAY adjust them at will, according to its local policy.

### Canceling a Recurrent Order

An important property of the recurrent Order is that it can be cancelled by the domain name owner,
with no need for certificate
revocation. We use the DELETE message to cancel the Order:

    DELETE /acme/order/1 HTTP/1.1
    Host: acme-server.example.org

Which returns:


    HTTP/1.1 202 Deleted


The server MUST NOT issue any additional certificates for this Order, beyond the certificate that is available for collection
at the time of deletion.

## Indicating Support of Recurrent Orders

ACME supports sending arbitrary extensions when creating an Order, and as a result, there is no need to explicitly
indicate support of this extension. The Proxy MUST verify that the "recurrent" attribute
was understood, as indicated by the "recurrent" attribute included in the created Order.
Since the standard ACME protocol does not allow to explicitly cancel a pending Order (the DELETE operation
above is an extension),
a Proxy that encounters an non-supporting server will probably let the Order
expire instead of following through with the authorization process.

## ACME Authorization

The DNO MUST restrict the authorizations it requests
from the ACME server to only those that cannot be
spoofed by a malicious NDC. In most cases the NDC
will have strong control of HTTP content under the
delegated domain, and therefore HTTPS-based authorization
MUST NOT be used. See also {{restricting-cdns}}.

## Fetching the Certificates
{: #fetching-certificates}

The certificate is fetched from the certificate endpoint, as per {{I-D.ietf-acme-acme}}, Sec. 7.4.2
"Downloading the Certificate".
The server MUST include an Expires header that indicates expiry of the specific certificate. When the certificate
expires, the client MAY assume that a newer certificate is already in place.

A certificate MUST be replaced by its successor at the latest halfway through its lifetime (the period between
its notBefore and notAfter times).

# Operational Considerations

## Certificate Transparency (CT) Logs

TBD: larger logs and how to deal with them.

# Security Considerations

## Restricting CDNs to the Delegation Mechanism
{: #restricting-cdns}

Currently there are no standard methods for the DNO to ensure that
the CDN cannot issue a certificate through mechanisms other than the one described here,
for the URLs under the CDN's control. For example, regardless of the STAR solution, a rogue CDN employee
can use the ACME protocol (or proprietary mechanisms used by various CAs) to create a fake certificate
for the DNO's content because ACME authorizes its requests using
information that may be under the adversary's control.

The best solution currently being worked on would consist of several related
configuration steps:

 * Make sure that the CDN cannot modify the DNS records for the domain.
 Typically this would mean that the content owner establishes a CNAME resource record
 from a subdomain into a CDN-managed domain.
 * Restrict certificate issuance for the domain to specific CAs that comply
 with ACME. This assumes
 universal deployment of CAA {{RFC6844}} by CAs, which is not the case
 yet. We note that the CA/Browser Forum has recently decided to
 require CAA checking {{CAB-CAA}}.
 * Deploy ACME-specific methods to restrict issuance to a specific authorization
 key which is controlled by the content owner {{I-D.ietf-acme-caa}}, and/or to specific
 ACME authorization methods.

This solution is recommended in general, even if an alternative to the
mechanism described here is used.

# Acknowledgments

This work is partially supported by the European Commission under
Horizon 2020 grant agreement no. 688421 Measurement and Architecture
for a Middleboxed Internet (MAMI). This support does not imply endorsement.

--- back

# Document History

[[Note to RFC Editor: please remove before publication.]]

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
