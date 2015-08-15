---
title: TLS Certificate Pinning with Tickets
abbrev: Pinning Tickets
docname: draft-sheffer-tls-pinning-ticket-latest
category: std
updates: 
obsoletes:

ipr: pre5378Trust200902
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
  I-D.ietf-tls-tls13:
  RFC5077:

informative:
  RFC6962:
  RFC7469:
  Oreo:
    title: "Firm Grip Handshakes: A Tool for Bidirectional Vouching"
    date: 2012
    author: 
      - ins: O. Berkman
      - ins: B. Pinkas
      - ins: M. Yung
    seriesinfo: Cryptology and Network Security, pp. 142-157

--- abstract

Session resumption tickets have been part of TLS for years, and are even more extensively used in TLS 1.3.
This document proposes to use similar tickets as a way to pin server certificates.
 that is, to ensure the client that it is connecting to the right server.

--- middle

# Introduction
The weaknesses of the global PKI system are by now widely known. There are many attempts
to resolve them, including Certificate Transparency (CT) {{RFC6962}} and HTTP Public Key
Pinning (HPKP) {{RFC7469}}. CT requires cooperation of a large number of the hundreds of
extant certificate authorities (CAs) before it can be used "for real", in enforcing
mode. HPKP has so far proven hard to deploy. This proposal augments these mechanisms
with a much easier to implement and deploy solution for certificate pinning, making use
of the session resumption mechanism.

When a client first connects to a server, the server responds with a ticket and a
committed lifetime. The ticket is modeled on the session resumption ticket, but is
distinct from it. Specifically, the ticket does not prove the client's identity. The
committed lifetime indicates for how long the server promises to accept the ticket,
i.e. to retain the server-side ticket-encryption key. The committed lifetime is
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
can be fully automated. The server administrator need not bother with the
management of backup certificates or explicit pins.
* For larger deployments, we reuse the existing {{RFC5077}} infrastructure where
it exists.
* Pinning errors, presumably resulting from MITM attacks, can be detected both by the
client and the server. This allows for server-side detection of MITM attacks using
large-scale analytics.

## Conventions used in this document
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in <xref target="RFC2119"/>.

# Protocol Overview
We review the two major use cases: the first time a particular client connects to a
server, and subsequent connections.

## Initial Connection
When a client first connects to a server, it receives a pinning ticket as part of the
server's first response, in a XXX message.
The pinning ticket is generated exactly like a session resumption
ticket ({{I-D.ietf-tls-tls13}}, Sec. Key Schedule), except that a separate secret is
computed as follows:

    pinning_secret = HKDF-Expand-Label(master_secret,
                                       "pinning master secret",
                                        session_hash, L)

When the connection setup is complete, the client computes 
pinning_secret and saves it locally, together with the received
ticket.

The client SHOULD cache the ticket and the secret for the pinning duration received from
the server, and MUST forget these values
at the end of this duration.

## Subsequent Connections

When the client initiates a connection to a server it has previously seen (see
{{indexing}}
on identifying servers and origins), it MUST send the pinning ticket for that server.

The server MUST extract the pinning_secret from the ticket
and MUST respond with a PinningAck message, which binds the pinning ticket to
the server's (current) public key.

## Indexing the Pins {#indexing}

# Message Definitions

# Operational Considerations

The main motivation behind the current protocol is to enable certificate
pinning without the need for manual operations. Manual operations are susceptible
to human error and in the case of certificate pinning, can easily result in
"server bricking": the server becoming inaccessible to some or all of its users.

## Ticket Secret Synchronization

The only operational requirement when deploying this protocol is that if the server
is part of a cluster, ticket
secrets MUST be synchronized between all cluster members. The protocol is designed so
that if resumption tickets {{RFC5077}} are already synchronized between
cluster members, nothing more needs to be done.

Moreover, synchronization does not need
to be instanteneous, e.g. tickets can be distributed a few minutes or hours in advance
of their rollover.

## Certificate Renewal

The protocol ensures that the client will continue speaking to the correct server
even when the server's certificate is renewed. In fact the protocol should more
precisely be called "server identity pinning".

## Certificate Revocation

The protocol is orthogonal to certificate validation, in the sense that, if the
server's certificate has been revoked or is invalid for some other reason,
the client MUST refuse to connect to it.

## Disabling Pinning

A server implementing this protocol MUST have a "ramp down" mode of operation where:

  * The server continues to accept valid pinning tickets and responds correctly with
a TicketAck.
  * The server does not send back a new PinningTicket.

After a while no clients will hold valid tickets any more and the feature may be
disabled.

## Server Compromise

If a server compromise is detected, the pinning secret MUST be rotated immediately,
but the server MUST still accept valid tickets that use the old, compromised key.
Clients who still hold old pinning tickets will remain vulnerable to MITM attacks,
but those that connect to the correct server will immediately receive new tickets.

# Comparison: HPKP Deployment

The current IETF standard for certificate pinning is the Public Key Pinning Extension
for HTTP, or HPKP. Unfortunately HPKP has not seen wide deployment yet. This may
simply be
due to inertia, but we believe the main reason is the onerous manual certificate
management which is needed to implement HPKP for enterprise servers. The penalty
for making mistakes (e.g. being too early or too late to deploy new pins) is
often bricking the server for some clients.

To demonstrate
this point, we present an analysis of what it would take to deploy HPKP for
a security-sensitive Web server.

1. Pin only end-entity certs. Pinning an intermediate cert means that the enterprise is
at risk if the CA makes sudden operational changes. Pinning the root cert is useless:
it still allows every "brand" (sub-CA) to issue a fake certificate for the servers.
1. Make sure the default reminder period from the certificate management system is
long, e.g. 3 months. I'm assuming a pin period ("max age") of 1 month.
1. Issue 2 certs with the same validity period, the main and a backup one.
1. Once we get the expiration reminder, issue 2 new certs and install the new "main"
cert on servers. Change the HPKP header to send the old main cert as the main pin,
the new main cert as the backup, and the new backup cert as a secondary backup
(in case the new main cert gets compromised). This transition period must be at least
1 month, so as not to break clients who still pin to the old main cert.
1. Shortly before expiration, change the HPKP header again to send the new main cert as the main pin and the new backup cert as the backup pin.

To summarize:

|---
| Period | Main server cert | Backup pin | Secondary backup pin
| Regular operation: before rotation | Old main cert | Old backup cert
| >1 month before expiration of old certs |	Old main cert |	New main cert |	New backup cert
| Shortly before expiration but not earlier than the previous change + 1 month | New main cert | New backup cert
| Regular operation: after rotation | New main cert | New backup cert
	

# Security Considerations

## Server-Side Error Detection

## Client Privacy

# IANA Considerations

[[TODO]]

* PinningTicket message
* PinningAck message
* Alerts?

# Acknowledgements

The original idea behind this proposal is due to Moty Yung, Benny Pinkas 
and Omer Berkman, and was published in {{Oreo}}. The current protocol is but a
distant relative of the original Oreo protocol, and any errors are the
draft author's alone.
