---
title: TLS Server Identity Pinning with Tickets
abbrev: Pinning Tickets
docname: draft-sheffer-tls-pinning-ticket-latest
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
  RFC2104:
  RFC2119:
  I-D.ietf-tls-tls13:
  RFC5077:

informative:
  RFC6454:
  RFC6962:
  RFC6982:
  RFC7258:
  RFC7469:
  RFC7507:
  I-D.perrin-tls-tack:
  Oreo:
    title: "Firm Grip Handshakes: A Tool for Bidirectional Vouching"
    date: 2012
    author:
      - ins: O. Berkman
      - ins: B. Pinkas
      - ins: M. Yung
    seriesinfo: Cryptology and Network Security, pp. 142-157
  Netcraft:
    author:
    -
      name: Paul Mutton
    title: "HTTP Public Key Pinning: You're doing it wrong!"
    date: March 30, 2016
    target: http://news.netcraft.com/archives/2016/03/30/http-public-key-pinning-youre-doing-it-wrong.html

--- abstract

Fake public-key certificates are an ongoing problem for users of TLS. Several solutions
have been proposed, but none is currently in wide use.
This document proposes to extend TLS with opaque tickets, similar to those being
used for TLS session resumption, as a way to pin the server's identity.
That is, to ensure the client that it is connecting to the right server even
in the presence of corrupt certificate authorities and fake certificates. The main advantage
of this solution is that no manual management actions are required.

--- middle

# Introduction
The weaknesses of the global PKI system are by now widely known. Essentially,
any valid CA may issue a certificate for any organization without the
organization's approval (a misissued or "fake" certificate),
and use the certificate to impersonate the organization.
There are many attempts
to resolve these weaknesses, including Certificate Transparency (CT) {{RFC6962}}, HTTP Public Key
Pinning (HPKP) {{RFC7469}}, and TACK {{I-D.perrin-tls-tack}}.
CT requires cooperation of a large
portion of the hundreds of
extant certificate authorities (CAs) before it can be used "for real", in enforcing
mode. It is noted that the relevant industry forum (CA/Browser Forum) is indeed pushing for such
extensive adoption.
TACK has some similarities to the current proposal, but work on it seems to have stalled.
{{tack}} compares our proposal to TACK.
HPKP is a standard,
but so far has proven hard to deploy (see {{hpkp}}).
This proposal augments these mechanisms
with a much easier to implement and deploy solution for server identity pinning, by
reusing some of the mechanisms behind TLS session resumption.

When a client first connects to a server, the server responds with a ticket and a
committed lifetime. The ticket is modeled on the session resumption ticket, but is
distinct from it. Specifically, the ticket acts as a "second factor" for proving
the server's identity;
the ticket does not authenticate the client. The
committed lifetime indicates for how long the server promises to
retain the server-side ticket-encryption key, which allows it to complete
the protocol exchange correctly and prove its identity. The committed lifetime is
typically on the order of weeks or months. We follow the Trust On First Use (TOFU)
model, in that the first server authentication is only based on PKI certificate
validation, but for any follow-on sessions, the client is further ensuring the server's
identity based on the server's ability to decrypt the ticket and complete the handshake
correctly.

This version of the draft only discusses TLS 1.3. We believe that the idea can also
be back-fitted into earlier versions of the protocol.

The main advantages of this protocol over earlier pinning solutions are:

* The protocol is at the TLS level, and as a result is not restricted to HTTP at the
application level.
* Once a single parameter is configured (the ticket secret's lifetime), operation
is fully automated. The server administrator need not bother with the
management of backup certificates or explicit pins.
* For server clusters, we reuse the existing {{RFC5077}} infrastructure where
it exists.
* Pinning errors, presumably resulting from MITM attacks, can be detected both by the
client and the server. This allows for server-side detection of MITM attacks using
large-scale analytics.

A note on terminology: unlike other solutions in this space, we do not
do "certificate pinning" (or "public key pinning"), since the protocol is oblivious to the server's
certificate. We prefer the term "server identity pinning" for this new solution.

## Conventions used in this document
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in <xref target="RFC2119"/>.

# Protocol Overview
The protocol consists of two phases: the first time a particular client connects to a
server, and subsequent connections.

This protocol supports full TLS handshakes, as well as 0-RTT handshakes.
Below we present it in the context of a full handshake, but behavior in 0-RTT
handshakes should be identical.

The preshared key (PSK) variant of TLS 1.3 is orthogonal to this protocol. A TLS session can
be established using PKI and a pinning ticket, and later resumed with PSK.
The PSK handshake MUST NOT include the extension defined here.

## Initial Connection

When a client first connects to a server, it requests a pinning ticket by sending an
empty PinningTicket extension, and receives it
as part of the
server's first response, in the returned PinningTicket extension.

     Client                                               Server

     ClientHello
       + key_share
       + PinningTicket         -------->
                                                     ServerHello
                                                     + key_share
                                           {EncryptedExtensions
                                                + PinningTicket}
                                          {ServerConfiguration*}
                                                  {Certificate*}
                                           {CertificateRequest*}
                                            {CertificateVerify*}
                               <--------              {Finished}
     {Certificate*}
     {CertificateVerify*}
     {Finished}                -------->
     [Application Data]        <------->      [Application Data]

            *  Indicates optional or situation-dependent
               messages that are not always sent.

            {} Indicates messages protected using keys
               derived from the ephemeral secret.

            [] Indicates messages protected using keys
               derived from the master secret.

The server computes a pinning_secret value ({{pinning-secret}})
in order to generate the ticket.
When the connection setup is complete, the client computes
the same pinning\_secret value and saves it locally, together with the received
ticket.

The client SHOULD cache the ticket and the pinning_secret for the lifetime received from
the server. The client MUST forget these values
at the end of this duration.

The returned ticket is sent as part of the ServerHello encrypted extensions, and MUST NOT be sent
as part of a HelloRetryRequest.

## Subsequent Connections

When the client initiates a connection to a server it has previously seen (see
{{indexing}}
on identifying servers), it SHOULD send the pinning ticket for that server.

The server MUST extract the original pinning_secret from the ticket
and MUST respond with a PinningTicket extension, which includes:

* A proof that the server can understand
the ticket that was sent by the client; this proof also binds the pinning ticket to
the server's (current) public key. The proof is MANDATORY if a ticket was sent by
the client.
* A fresh pinning ticket. The main reason for refreshing the ticket on each connection
is privacy: to avoid the ticket serving as a fixed client identifier. It is RECOMMENDED
to include a fresh ticket with each response.

If the server cannot validate the ticket, that might indicate an earlier MITM attack
on this client. The server MUST then abort the connection with a
handshake_failure alert, and SHOULD log this failure.

The client MUST verify the proof, and if it fails to do so,
MUST issue a handshake\_failure alert
and abort the connection (see also {{client_error}}). When the connection is successfully
set up, i.e. after the Finished message is verified, the
client SHOULD store the new ticket along with the corresponding pinning\_secret.

Although this is an extension, if the client already has a ticket for a server,
the client MUST interpret a missing PinningTicket extension in the
server's response as an attack, because of the server's prior commitment to
respect the ticket. The client MUST abort the connection in this case.
See also {{ramp_down}} on ramping down support for this extension.

## Indexing the Pins {#indexing}

Each pin is associated with a host name, protocol (TLS or DTLS) and port number.
In other words, the pin for port TCP/443 may be different from that for DTLS or from
the pin for port TCP/8443.
The host name MUST be
the value sent inside the Server Name Indication (SNI) extension.
This definition is similar to
a Web Origin {{RFC6454}}, but does not assume the existence of a URL.

IP addresses are ephemeral and forbidden in SNI and therefore Pins MUST NOT be associated
with IP addresses.

# Message Definitions

This section defines the format of the PinningTicket extension.
We follow the message notation of {{I-D.ietf-tls-tls13}}.

     opaque pinning_ticket<0..2^16-1>;

     opaque pinning_proof<0..2^8-1>;

     struct {
       select (Role) {
       case client:
	       uint16 ticket_len; //zero if no ticket
           pinning_ticket ticket<0..2^16-1>; //omitted on 1st connection

         case server:
	       uint16 proof_len; //zero if no proof
           pinning_proof proof<0..2^8-1>; //no proof on 1st connection
	       uint16 ticket_len; //zero if no ticket
           pinning_ticket ticket<0..2^16-1>; //omitted on ramp down
           uint32 lifetime;
       }
    } PinningTicketExtension;

ticket
: a pinning ticket sent by the client or returned by the server. The ticket is opaque
to the client. The extension MUST contain exactly 0 or 1 tickets.

proof
: a demonstration by the server that it understands the ticket and therefore that
it is in possession of the secret that was used to generate it originally. The
proof is further bound to the server's public key to prevent some MITM attacks.
The extension MUST contain exactly 0 or 1 proofs.

ticket\_len, proof\_len
: the length in octet of the ticket or, respectively, the proof. The length values are each
2 bytes, in network order.

lifetime
: the duration (in seconds) that the server commits to accept the newly offered
ticket in the future.

# Cryptographic Operations {#crypto}

This section provides details on the cryptographic operations performed
by the protocol peers.

## Pinning Secret {#pinning-secret}

On each connection that includes the PinningTicket extension, both peers
derive the the value pinning_secret from the shared Diffie Hellman secret. They compute:

    pinning_secret = HKDF(xSS + xES, "pinning secret", L)

using the notation of {{I-D.ietf-tls-tls13}}, sec. Key Schedule. This secret
is used by the server to generate the new ticket that it returns to the client.

The length of the secret L is determined by the server, and MUST be between 16 and 63 octets, inclusive.

## Pinning Ticket {#pinning-ticket}
The pinning ticket's format is not specified by this document, but it MUST be
encrypted and integrity-protected using a long-term pinning-ticket protection key.
The server MUST rotate the protection key periodically, and therefore the ticket
MUST contain a protection key ID or serial number.
The ticket MUST allow the server to recover the pinning_secret value, and MAY
include additional information.

As noted in {{cluster}}, if the server is actually a cluster of machines,
the protection key MUST
be synchronized between them. An easy way to do it is to derive it from the
session-ticket protection key, which is already synchronized. For example:

    pinning_protection_key = HKDF(resumption_protection_key,
                                  "pinning protection", L)

## Pinning Proof

The proof sent by the server consists of this value:

    proof = HMAC(original_pinning_secret, "pinning proof" + crlen +
	        client.random + srlen + server.random +
	        Hash(server_public_key))

where HMAC {{RFC2104}} uses the Hash algorithm for the handshake,
and the same hash is also used over the server's public key. The server\_public\_key value
is the DER representation of the public key, specifically
the SubjectPublicKeyInfo structure as-is. The nonce lengths crlen and srlen
are a single octet each.

# Operational Considerations

The main motivation behind the current protocol is to enable identity
pinning without the need for manual operations. Manual operations are susceptible
to human error and in the case of public key pinning, can easily result in
"server bricking": the server becoming inaccessible to some or all of its users.

## Protection Key Synchronization {#cluster}

The only operational requirement when deploying this protocol is that if the server
is part of a cluster, protection keys (the keys used to encrypt tickets)
MUST be synchronized between all cluster members. The protocol is designed so
that if resumption ticket protection keys {{RFC5077}} are already synchronized between
cluster members, nothing more needs to be done.

Moreover, synchronization does not need
to be instantaneous, e.g. protection keys can be distributed a few minutes
or hours in advance of their rollover.

Misconfiguration can lead to the server's clock being off by a large amount of time. Therefore we recommend
never to automatically delete protection keys, even when they are long expired.

## Ticket Lifetime

The lifetime of the ticket is a commitment by the server to retain the ticket's corresponding
protection key for this duration, so that the server can prove to the client that it
knows the secret embedded in the ticket. For production systems, the lifetime SHOULD be between 7 and 30 days.

## Certificate Renewal

The protocol ensures that the client will continue speaking to the correct server
even when the server's certificate is renewed. In this sense, we are not "pinning
certificates" and the protocol should more precisely be called "server identity pinning".

## Certificate Revocation

The protocol is orthogonal to certificate validation, in the sense that, if the
server's certificate has been revoked or is invalid for some other reason,
the client MUST refuse to connect to it.

## Disabling Pinning {#ramp_down}

A server implementing this protocol MUST have a "ramp down" mode of operation where:

  * The server continues to accept valid pinning tickets and responds
  correctly with a proof.
  * The server does not send back a new pinning ticket.

After a while no clients will hold valid tickets any more and the feature may be
disabled.

## Server Compromise

If a server compromise is detected, the pinning secret MUST be rotated immediately,
but the server MUST still accept valid tickets that use the old, compromised key.
Clients that still hold old pinning tickets will remain vulnerable to MITM attacks,
but those that connect to the correct server will immediately receive new tickets.

## Disaster Recovery

All web servers in production need to be backed up, so that they can be recovered if a
disaster (including a malicious activity) ever wipes them out. Backup typically includes
the certificate and its private key, which must be backed up securely. The pinning secret, including
earlier versions that are still being accepted, must be backed up regularly. However
since it is only used as an authentication second factor, it does not require the same level
of confidentiality as the server's private key.

Readers should note that {{RFC5077}} session resumption keys are more security sensitive, and
should normally not be backed up but rather treated as ephemeral keys. Even when servers derive
pinning secrets from resumption keys ({{pinning-secret}}), they MUST NOT back up resumption keys.

# Previous Work

This section compares ticket pinning to two earlier proposals, HPKP and TACK.

## Comparison: HPKP Deployment {#hpkp}

The current IETF standard for pinning the identity of web servers
is the Public Key Pinning Extension
for HTTP, or HPKP {{RFC7469}}. Unfortunately HPKP has not seen wide deployment yet.
As of March 2016, the number of servers using HPKP was less than 3000 {{Netcraft}}.
This may simply be
due to inertia, but we believe the main reason is the onerous manual certificate
management which is needed to implement HPKP for enterprise servers. The penalty
for making mistakes (e.g. being too early or too late to deploy new pins) is
having the server become unusable for some of the clients.

To demonstrate
this point, we present a list of the steps involved in deploying HPKP on
a security-sensitive Web server.

1. Generate two public/private key-pairs on a computer that is not the Live server. The second one is
the "backup1" key-pair.

    `openssl genrsa -out "example.com.key" 2048;`

    `openssl genrsa -out "example.com.backup1.key" 2048;`


2. Generate hashes for both of the public keys. These will be used in the HPKP header:

    `openssl rsa -in "example.com.key" -outform der -pubout | openssl dgst -sha256 -binary | openssl enc -base64`

    `openssl rsa -in "example.com.backup1.key" -outform der -pubout | openssl dgst -sha256 -binary | openssl enc -base64`

3. Generate a single CSR (Certificate Signing Request) for the first key-pair, where you
include the domain name in the CN (Common Name) field:

    `openssl req -new -subj "/C=GB/ST=Area/L=Town/O=Company/CN=example.com"
            -key "example.com.key" -out "example.com.csr";`

4. Send this CSR to the CA (Certificate Authority), and go though the dance to prove you own the domain.
The CA will give you back a single certificate that will typically expire within a year or two.

5. On the Live server, upload and setup the first key-pair (and its certificate).
At this point you can add the "Public-Key-Pins" header, using the two hashes you created in step 2.

    Note that only the first key-pair has been uploaded to the server so far.

6. Store the second (backup1) key-pair somewhere safe, probably somewhere encrypted like a password manager.
It won't expire, as it's just a key-pair, it just needs to be ready for when you need to get your next certificate.

7. Time passes... probably just under a year (if waiting for a certificate to expire), or maybe sooner if you find
that your server has been compromised and you need to replace the key-pair and certificate.

8. Create a new CSR (Certificate Signing Request) using the "backup1" key-pair, and get a new certificate
from your CA.

9. Generate a new backup key-pair (backup2), get its hash, and store it in a safe place (again,
not on the Live server).

10. Replace your old certificate and old key-pair, and update the "Public-Key-Pins" header to remove
the old hash, and add the new "backup2" key-pair.

Note that in the above steps, both the certificate issuance as well as the storage of the backup key pair
involve manual steps. Even with an automated CA that runs the ACME protocol, key backup would be a challenge
to automate.

## Comparison: TACK {#tack}

Compared with HPKP, TACK {{I-D.perrin-tls-tack}} is a lot more similar to the current draft.
It can even be argued that this document is a symmetric-cryptography variant of TACK.
That said, there are still a few significant differences:

- Probably the most important difference is that with TACK, validation of the server
certificate is no longer required, and in fact TACK specifies it as a "MAY" requirement
(Sec. 5.3).
With ticket pinning, certificate validation by the client remains a MUST requirement, and the
ticket acts only as a second factor. If the pinning secret is compromised, the server's
security is not immediately at risk.
- Both TACK and the current draft are mostly orthogonal to the server certificate as far as
their life cycle, and so both can be deployed with no manual steps.
- TACK uses ECDSA to sign the server's public key. This allows cooperating clients
to share server assertions between themselves. This is an optional TACK feature,
and one that cannot be done with pinning tickets.
- TACK allows multiple servers to share its public keys. Such sharing is disallowed
by the current document.
- TACK does not allow the server to track a particular client, and so has better
privacy properties than the current draft.
- TACK has an interesting way to determine the pin's lifetime, setting it
to the time period since the pin was first observed, with a hard upper bound of 30 days.
The current draft makes the lifetime explicit, which may be more flexible to deploy.
For example, Web sites which are only visited rarely by users may opt for a longer
period than other sites that expect users to visit on a daily basis.

# Implementation Status

[Note to RFC Editor: please remove this section before publication.]

This section records the status of known implementations of the
protocol defined by this specification at the time of posting of
this Internet-Draft, and is based on a proposal described in [RFC6982].
The description of implementations in this section is intended to assist the IETF
in its decision processes in
progressing drafts to RFCs.  Please note that the listing of any individual
implementation here does not imply endorsement by the
IETF.  Furthermore, no effort has been spent to verify the
information presented here that was supplied by IETF contributors.
This is not intended as, and must not be construed to be, a
catalog of available implementations or their features.  Readers
are advised to note that other implementations may exist.

According to RFC 6982, "this will allow reviewers and working
groups to assign due consideration to documents that have the
benefit of running code, which may serve as evidence of valuable
experimentation and feedback that have made the implemented
protocols more mature.  It is up to the individual working groups
to use this information as they see fit".

## Mint Fork

### Overview
A fork of the Mint TLS 1.3 implementation, developed by Yaron Sheffer
and available at https://github.com/yaronf/mint.

### Description
This is a fork of the TLS 1.3 implementation, and includes client and server code.
In addition to the actual protocol, several utilities are provided allowing
to manage protection keys on the server side, and pinning tickets on the client side.

### Level of Maturity
This is a prototype.

### Coverage
The entire protocol is implemented.

### Licensing
Mint itself and this fork are available under an MIT license.

### Contact Information
See author details below.

# Security Considerations

This section reviews several security aspects related to the proposed extension.

## Trust on First Use (TOFU) and MITM Attacks

This protocol is a "trust on first use" protocol. If a client initially connects to
the "right" server, it will be protected against MITM attackers for the lifetime
of each received ticket. If it connects regularly (depending of course on the
server-selected lifetime), it will stay constantly protected against fake
certificates.

However if it initially connects to an attacker, subsequent connections to
the "right" server will fail. Server operators might want to advise clients
on how to remove corrupted pins, once such large scale attacks are detected and remediated.

The protocol is designed so that it is not vulnerable to an active MITM attacker
who has real-time access
to the original server. The pinning proof includes a hash of the server's
public key, to ensure the client that the proof was in fact generated by
the server with which it is initiating the connection.

## Pervasive Monitoring

Some organizations, and even some countries perform pervasive monitoring on their
constituents {{RFC7258}}. This often takes the form of always-active SSL proxies. Because of
the TOFU property, this protocol does not provide any security in such cases.

## Server-Side Error Detection {#server_error}

Uniquely, this protocol allows the server to detect clients that present incorrect
tickets and therefore can be assumed to be victims of a MITM attack. Server operators
can use such cases as indications of ongoing attacks, similarly to fake certificate
attacks that took place in a few countries in the past.

## Client Policy and SSL Proxies {#client_policy}

Like it or not, some clients are normally deployed behind an SSL proxy.
Similarly to {{RFC7469}}, it is acceptable to allow pinning to be disabled for some hosts
according to local policy. For example, a UA MAY disable pinning for hosts whose
validated certificate chain terminates at a user-defined trust anchor, rather than
a trust anchor built-in to the UA (or underlying platform). Moreover, a client MAY accept
an empty PinningTicket extension from such hosts as a valid response.

## Client-Side Error Behavior {#client_error}

When a client receives a malformed or empty PinningTicket extension from a pinned server, it MUST
abort the handshake and MUST NOT retry with no PinningTicket in the request. Doing
otherwise would expose the client to trivial fallback attacks, similar to
those described in {{RFC7507}}.

This rule can however have negative affects on clients that move from behind SSL proxies into
the open Internet and vice versa, if the advice in {{client_policy}} is not followed.
Therefore, we RECOMMEND that browser and library vendors provide a documented way to
remove stored pins.

## Client Privacy

This protocol is designed so that an external attacker cannot correlate between different requests
of a single client, provided the client requests and receives a fresh ticket upon each connection.

On the other hand, the server to which the client is connecting can easily track the client.
This may be an issue when the client expects to connect to the server (e.g., a mail server)
with multiple identities. Implementations SHOULD allow the user to opt out of pinning, either
in general or for particular servers.

# IANA Considerations

IANA is requested to allocate a TicketPinning extension value in the
TLS ExtensionType Registry.

No registries are defined by this document.

# Acknowledgements

The original idea behind this proposal was published in {{Oreo}} by Moty Yung,
Benny Pinkas
and Omer Berkman. The current protocol is but a
distant relative of the original Oreo protocol, and any errors are the
draft author's alone.

I would like to thank Dave Garrett, Daniel Kahn Gillmor and Yoav Nir for their comments on this draft.
Special thanks to Craig Francis for contributing the HPKP deployment script, and to Ralph Holz
for several fruitful discussions.

--- back

# Document History

## draft-sheffer-tls-pinning-ticket-02

- Added an Implementation Status section.
- Added lengths into the extension structure.
- Changed the computation of the pinning proof to be more robust.
- Clarified requirements on the length of the pinning_secret.
- Revamped the HPKP section to be more in line with current practices, and added recent
statistics on HPKP deployment.

## draft-sheffer-tls-pinning-ticket-01

- Corrected the notation for variable-sized vectors.
- Added a section on disaster recovery and backup.
- Added a section on privacy.
- Clarified the assumptions behind the HPKP procedure in the comparison section.
- Added a definition of pin indexing (origin).
- Adjusted to the latest TLS 1.3 notation.

## draft-sheffer-tls-pinning-ticket-00

Initial version.
