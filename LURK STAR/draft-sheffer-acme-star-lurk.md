---
title: Use of Short-Term, Automatically-Renewed (STAR) Certificates to address the LURK problem
abbrev: ACME STAR LURK
docname: draft-sheffer-acme-star-lurk-latest
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
  RFC6844:
  I-D.iab-web-pki-problems:
  I-D.cairns-tls-session-key-interface:
  I-D.erb-lurk-rsalg:
  I-D.landau-acme-caa:
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

This memo proposes two mechanisms that work in concert to address the LURK problem statement, allowing a third party (e.g., a content delivery network) to terminate TLS sessions on behalf of a domain name owner (e.g., a content provider).

The proposed mechanisms are:

1. An extension to the ACME protocol to enable the issuance of short-term and automatically renewed certificates, and
2.  A protocol that allows a domain name owner to delegate to a third party control over a certificate that bears its own name.

It should be noted that these are in fact independent building blocks that could be used separately to solve completely different problems.
--- middle

# A Solution for the HTTPS CDN Use Case

A content provider, and Domain Name Owner (DNO), has agreements in place with one or more Content Delivery Networks (CDN) that are contracted to serve its content over HTTPS.  The CDN terminates the HTTPS connection at one of its edge cache servers and needs to present its clients (browsers, set-top-boxes) a certificate whose name matches the authority of the URL that is requested, i.e. that of the DNO.  However, many DNOs balk at sharing their long-term private keys with another organization and, equally, CDN providers would rather not have to handle other parties' long-term secrets. This problem has been discussed at the IETF under the LURK (limited use of remote keys) title.

This document proposes a solution to the above problem that involves the use of short-term certificates with a DNO's name on them, and a scheme for handling the naming delegation from the DNO to the CDN.  The generated short-term credentials are automatically renewed by an ACME Certification Authority (CA) {{I-D.ietf-acme-acme}} and routinely rotated by the CDN on its edge cache servers.  The DNO can end the delegation at any time by simply instructing the CA to stop the automatic renewal and let the certificate expire shortly after.

Using short-term certificates makes revocation cheap and effective {{Topalovic}} {{I-D.iab-web-pki-problems}} in case of key compromise or of termination of the delegation; seamless certificate issuance and renewal enable the level of workflow automation that is expected in today's cloud environments.  Also, compared to other keyless-TLS solutions {{I-D.cairns-tls-session-key-interface}} {{I-D.erb-lurk-rsalg}}, the proposed approach doesn't suffer from scalability issues or increase in connection setup latency, while requiring virtually no changes to existing COTS caching software used by the CDN.

# Conventions used in this document

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in {{RFC2119}}.

TODO: glossary.

# Protocol Flow

The protocol flow can be split into two: a LURK interface, used by CDN and DNO to agree on the name delegation, and the ACME STAR interface, used by DNO to obtain the short-term and automatically renewed certificate from the CA, which is eventually consumed by the CDN.  The latter is also used to terminate the delegation, if so needed.

The following subsections describe the preconditions ({{proto-preconditions}}), and the three main phases of the protocol:

- Bootstrap: the CDN requests from the DNO the delegation of a specific name and in turn DNO asks an ACME CA to create the corresponding short-term and auto-renewed (STAR) certificate ({{proto-bootstrap}});
- Auto-renewal: the ACME CA periodically re-issues the short-term certificate and posts it to a public URL ({{proto-auto-renewal}});
- Termination: the DNO (indirectly) stops name delegation by explicitly requesting the ACME CA to discontinue the automatic renewal of the certificate ({{proto-termination}}).



## Preconditions
{: #proto-preconditions}

The protocol assumes the following preconditions are met:

- A mutually authenticated channel between CDN and DNO pre-exists.  This is called "LURK channel" and all LURK protocol exchanges between CDN and DNO are run over it.  It provides the guarantee that the LURK requests and responses are authentic [^1]{:tf: source="tf"}.
- CDN and DNO have agreed on a "CSR template" to use, including at a minimum:
  - Subject name (e.g., "somesite.DNO.com"),
  - Validity (e.g., 24 to 72 hours),
  - Requested algorithms,
  - Key length,
  - Key usage.

  The CDN is required to use this template for every CSR created under the same delegation.
- DNO has registered through the ACME interface exposed by the Certificate Authority (CA) using the usual ACME registration procedure. The DNO shall, at the registration stage, query the ACME server for the supported STAR capabilities – for example: the minimum validity period of the issued certificate, the maximum duration of the automatic renewal process (either as a maximum number of renewal events, or as its maximum absolute life-span).

[^1]: Note that, under this assumption, the key used to authenticate the CDN to the DNO becomes a critical asset for the security of the proposed protocol, and that certain interactions (e.g., CSR submission) might require a stronger authentication mechanism.  For example, stacking a further authentication factor on top of CDN's LURK key would allow to distinguish an attacker that has only managed to successfully attack the CDN's LURK key from the legitimate CDN.


## Bootstrap
{: #proto-bootstrap}

CDN (LURK client) generates a key-pair, wraps it into a Certificate Signing Request (CSR) according to the agreed CSR template, and sends it to the DNO (LURK server) over the pre-established LURK channel.  The DNO uses the CDN identity provided on the LURK channel to look up the CSR template that applies to the requesting CDN and decides whether or not to accept the request.  (TBD: This is probably a case that would require a further authentication stage over the one provided by the mutual-authenticated LURK channel?)  Assuming everything is in order, it then "forwards" the CDN request to the ACME CA by means of the usual ACME application procedure. Specifically, DNO, in its role as a ACME/STAR client, requests the CA a STAR certificate, i.e., one that:

- Has a short validity (e.g., 24 to 72 hours);
- Is automatically renewed by the CA for a certain period of time;
- Is downloadable from a (highly available) public link without requiring any special authorisation.

Other than that, the ACME protocol flows as normal between DNO and CA, in particular DNO is responsible for satisfying the requested ACME challenges until the CA is willing to issue the requested certificate.
The DNO is given back a unique identifier for the issued STAR certificate to be used in subsequent interaction with the CA (e.g., if the certificate needs to be terminated.)

Concurrently, a 202 response has been sent back to the CDN with an endpoint to poll for completion of the certificate generation process. 

The bootstrap phase ends when the DNO obtains the OK from the ACME CA and posts the certificate's URL to the "completion endpoint" where the CDN can retrieve it.  The information that is passed on to the CDN at this stage also includes details about how much time before the certificate expires can the CDN expect the replacement to be ready.

~~~~~~~~~~
                     ...........................
LURK                 :  LURK          ACME/STAR:              ACME/STAR
Client               :  Server         Client  :               Server
  |                  :    |               |    :                  |
  |                  :    |               |   ACME registration   |
  +-------.          :    |               |<--------------------->|
  |       |          :    |               |   STAR capabilities   |
  |   generate CSR   :    |               |    :                  |
  |       |          :    |               |    :                  |
  |<------'          :    |               |    :                  |
  |                  :    |               |    :                  |
  |     Request new  :    |               |    :                  |
  +---------------------->|               |    :                  |
  |     cert for CSR :    |               |    :                  |
  |                  :    +-------.       |    :                  |
  |                  :    |       |       |    :                  |
  |                  :    |   Verify CSR  |    :                  |
  |                  :    |       |       |    :                  |
  |                  :    +<------'       |    :                  |
  |   Accepted, poll at   |               |    :                  |
  |<----------------------+               |    :                  |
  |    "completion URL"   |- - - - - - - >|    Application for    |
  |                  :    |               +---------------------->|
  |                  :    |               |    STAR certificate   |
  |                  :    |               |    :                  |
  |  GET "completion URL" |               |    :  Challenge       |
  |<--------------------->|               |<--------------------->|
  |   202, in progress    |               |    :  Response        |
  |                  :    |               |    :                  |
  |                  :    |               |  Finalize/Certificate |
  |                  :    |               |<----------------------+
  |  GET "completion URL" |< - - - - - - -|    : + STAR Id        |
  +---------------------->|               |    :                  |
  |                  :    |               |    :                  |
  |  200, certificate URL |               |    :                  |
  |<----------------------+               |    :                  |
  |   and other metadata  |               |    :                  |
  |                  :    |               |    :                  |
                     `.........................'
~~~~~~~~~~
{: #figprotoboot title="Bootstrap"}


## Refresh
{: #proto-auto-renewal}

The CA automatically re-issues the certificate (using the same CSR) before it expires and publishes it to the URL that the CDN has come to know at the end of the bootstrap phase.  The CDN downloads and installs it.  This process goes on until either:

- DNO terminates the delegation, or
- Automatic renewal expires.

~~~~~~~~~~
        LURK                    ACME/STAR
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

DNO requests termination of the STAR certificate by including the previously obtained identifier in a STAR certificate termination request to the ACME interface.
After CA receives and verifies the request, it shall:

- Cancel the automatic renewal process for the LURK certificate;
- Change the certificate publication resource to return an error indicating the termination of the delegation to external clients, including the CDN;

Note that it is not necessary to explicitly revoke the short-term certificate.

~~~~~~~~~~
   LURK                   ACME/STAR               ACME/STAR
   Client                  Client                  Server
     |                       |                       |
     |                       |  Terminate STAR Id    |
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

This section describes the protocol's details. We start with the LURK API between the LURK Client and the ACME Proxy.
Then we describe a few extensions to the ACME protocol running between the ACME Proxy and the ACME Server.

## LURK API

This API allows the LURK Client to request a STAR certificate via the Proxy, using a previously agreed-upon CSR template.

The API consists of a single resource, "registration". A new Registration is created with a POST and then the
Registration instance is polled to obtain its details.

### Creating a Registration

To create a registration, use:

    POST /lurk/registration
    Host: example.com
    Content-Type: application/json
    {
        "csr": "...", // CSR in PEM format
        "lifetime": 365 // requested registration lifetime in days,
                        // between 1 and 1095
    }

Upon success, the call returns the new Registration resource. 

    HTTP/1.1 201 Created
    Replay-Nonce: D8s4D2mLs8Vn-goWuPQeKA
    Location: https://example.com/lurk/registration/567

### Polling the Registration

The returned Registration can be polled until the information is available from the ACME server.

    GET /lurk/registration/567
    Host: example.com

While still polling the server, this returns:

    HTTP/1.1 200 OK
    Retry-After: 10 // in seconds; this is an optional header
    {
        "status": "pending"
    }

When the operation is completed, the Proxy returns:

    HTTP/1.1 200 OK
    {
        "status": "valid", // or "failed"
        "expires": "2018-09-09T14:09:00Z", // expiration of this response
                                           // and the Registration resource
        "lifetime": 365, // lifetime of the registration in days,
                         //  possibly less than requested
        "certificates": "https://example-server.com/certificates/A251A3"
    }

The "expires" field applies to the registration resource itself, and may be as small as a few minutes.
It is unrelated to the order's lifetime which is measured in days or longer. The "certificates" attribute
contains a URL of the certificate pull endpoint, see {{fetching-certificates}}.

## ACME Extensions between Proxy and Server

We propose to extend the ACME protocol slightly, by allowing recurrent orders.

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
They are returned when the order has been created, possibly with adjusted values.

### Canceling a Recurrent Order

An important property of the recurrent order is that it can be cancelled by the domain name owner,
with no need for certificate
revocation. We use the DELETE message for that:

    DELETE /acme/order/1 HTTP/1.1
    Host: example-server.com

Which returns:


    HTTP/1.1 202 Deleted


The server MUST NOT issue any additional certificates for this Order, beyond the certificate that is available for collection
at the time of deletion.

### Indicating Support of Recurrent Orders

ACME supports sending arbitrary extensions when creating an Order, and as a result, there is no need to explicitly
indicate support of this extension. The Proxy MUST verify that the "recurrent" attribute
was understood, as indicated by the "recurrent" attribute included in the created Order.
Since the standard ACME protocol does not allow to explicitly cancel a pending Order (the DELETE operation
above is an extension),
an unhappy Proxy will probably let the Order
expire instead of following through with the authorization process.

## Fetching the Certificates
{: #fetching-certificates}

The certificate is fetched from the certificate endpoint, as per {{I-D.ietf-acme-acme}}, Sec. 7.4.2
"Downloading the Certificate".
The server MUST include an Expires header that indicates expiry of the specific certificate. When the certificate
expires, the client MAY assume that a newer certificate is already in place.

A certificate MUST be replaced by its successor at the latest 24 hours before its "Not After" time.

# Security Considerations

- CDN's client certificate key is first order security asset and MUST be protected.  Absent 2FA/MFA, an attacker that can compromise the key might be able to obtain certificates bearing DNO's identity.
- Consider collusion of two or more CDNs with contracts with the same DNO (?)

## Restricting CDNs to the Delegation Mechanism

Currently there are no standard methods for the DNO to ensure that
the CDN cannot issue a certificate through mechanisms other than the one described here,
for the URLs under the CDN's control. For example, regardless of the STAR solution, a rogue CDN employee
can use the ACME protocol (or proprietary mechanisms used by various CAs) to create a fake certificate
for the DNO's content.

The best solution currently being worked on would consist of several related
configuration steps:

 * Make sure that the CDN cannot modify the DNS records for the domain.
 Typically this would mean that the content owner establishes a CNAME resource record
 from a subdomain into a CDN-managed domain.
 * Restrict certificate issuance for the domain to specific CAs that comply
 with ACME. This assumes
 universal deployment of CAA {{RFC6844}} by CAs, which is not the case yet.
 * Deploy ACME-specific methods to restrict issuance to a specific authorization
 key which is controlled by the content owner {{I-D.landau-acme-caa}}, and/or to specific
 ACME authorization methods.

This solution is recommended in general, even if an alternative to the
mechanism described here is used.

# Acknowledgments

This work is partially supported by the European Commission under Horizon 2020 grant agreement no. 688421 Measurement and Architecture for a Middleboxed Internet (MAMI). This support does not imply endorsement.

--- back
