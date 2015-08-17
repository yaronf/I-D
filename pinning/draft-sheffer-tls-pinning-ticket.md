﻿---
title: TLS Certificate Pinning with Tickets
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
  RFC6962:
  RFC7469:
  I-D.perrin-tls-tack:
  Oreo:
    title: "Firm Grip Handshakes: A Tool for Bidirectional Vouching"
    date: 2012
    author: 
      - ins: O. Berkman
      - ins: B. Pinkas
      - ins: M. Yung
    seriesinfo: Cryptology and Network Security, pp. 142-157

--- abstract

Fake public-key certificates are an ongoing problem for users of TLS. Several solutions
have been proposed, but none is currently in wide use.
This document proposes to extend TLS with opaque tickets, similar to the ones being
used for TLS session resumption, as a way to pin server certificates.
That is, to ensure the client that it is connecting to the right server even
in the presence of corrupt certificate authorities. The main advantage
of this solution is that no manual management actions are required.

--- middle

# Introduction
The weaknesses of the global PKI system are by now widely known. There are many attempts
to resolve them, including Certificate Transparency (CT) {{RFC6962}}, HTTP Public Key
Pinning (HPKP) {{RFC7469}}, and TACK {{I-D.perrin-tls-tack}}.
CT requires cooperation of a large
portion of the hundreds of
extant certificate authorities (CAs) before it can be used "for real", in enforcing
mode. TACK has not been standardized. HPKP is a standard,
but so far has proven hard to deploy (see {{hpkp}}).
This proposal augments these mechanisms
with a much easier to implement and deploy solution for certificate pinning, by
reusing some of the mechanisms behind TLS session resumption.

When a client first connects to a server, the server responds with a ticket and a
committed lifetime. The ticket is modeled on the session resumption ticket, but is
distinct from it. Specifically, the ticket acts as a "second factor" for proving
the server's identity;
the ticket does not authenticate the client. The
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
is fully automated. The server administrator need not bother with the
management of backup certificates or explicit pins.
* For server clusters, we reuse the existing {{RFC5077}} infrastructure where
it exists.
* Pinning errors, presumably resulting from MITM attacks, can be detected both by the
client and the server. This allows for server-side detection of MITM attacks using
large-scale analytics.

## Conventions used in this document
The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in <xref target="RFC2119"/>.

# Protocol Overview
We review the two major use cases: the first time a particular client connects to a
server, and subsequent connections.

This protocol supports full TLS handshakes, as well as 0-RTT handshakes.
Below we present it in the context of a full handshake, but behavior in 0-RTT
handshakes should be identical.

The PSK variant of TLS 1.3 is orthogonal to this protocol. A TLS session can
be established using PKI and a pinning ticket, and later resumed with PSK.
The PSK handshake MUST NOT include the extension defined here.

## Initial Connection

When a client first connects to a server, it requests a pinning ticket and receives it
as part of the
server's first response, in a PinningTicket extension.

     Client                                               Server

     ClientHello
       + ClientKeyShare
       + PinningTicket         -------->
                                                     ServerHello
                                                 ServerKeyShare*
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

The server computes a pinning_secret value in order to generate the ticket.
When the connection setup is complete, the client computes 
the same pinning\_secret value and saves it locally, together with the received
ticket.

The client SHOULD cache the ticket and the pinning_secret for the lifetime received from
the server. The client MUST forget these values
at the end of this duration.

The returned ticket is sent as a ServerHello protected extension, and MUST NOT be sent
as part of a HelloRetryRequest.

## Subsequent Connections

When the client initiates a connection to a server it has previously seen (see
{{indexing}}
on identifying servers and origins), it SHOULD send the pinning ticket for that server.

The server MUST extract the pinning_secret from the ticket
and MUST respond with a PinningTicket message, which includes

* A proof it can understand
the ticket that was sent by the client; this proof also binds the pinning ticket to
the server's (current) public key. The proof is MANDATORY if a ticket was sent by
the client.
* A fresh pinning ticket. The main reason for refreshing the ticket on each connection
is privacy: to avoid the ticket serving as a fixed client identifier. It is RECOMMENDED
to include a fresh ticket with each response.

If the server cannot validate the ticket, that might point to an earlier MITM attack
on this client. The server MUST then abort the connection with a
handshake_failure alert, and SHOULD log this failure.

The client MUST verify the proof, and if it fails to do so,
MUST issue a handshake_failure alert
and abort the connection. When the connection is successfully set up, the
client SHOULD store the new ticket along with the corresponding pinning\_secret.

Although this is an extension, if the client already has a ticket for a server,
the client MUST interpret a missing PinningTicket extension in the
server's response as an attack, because of the server's prior commitment to
respect the ticket. See also {{ramp_down}} on ramping down support for this extension.

## Indexing the Pins {#indexing}

[[TODO]]

# Message Definitions

This section defines the format of the PinningTicket extension.
We follow the message notation of {{I-D.ietf-tls-tls13}}.

     opaque pinning_ticket<0..2^16-1>;

     opaque pinning_proof<0..2^8-1>;

     struct {
         select (Role) {
             case client:
		         pinning_ticket ticket<0..1>; // 0 tickets on 1st connection

             case server:
                 pinning_proof proof<0..1>; // 0 proofs on 1st connection
                 pinning_ticket ticket<0..1>; // omitted only on ramp down
                 uint32 lifetime;
       }
	} PinningTicketExtension;

ticket
: a pinning ticket sent by the client or returned by the server. The ticket is opaque
to the client.

proof
: a demonstration by the server that it understands the ticket and therefore that
it is in possession of the secret that was used to generate it originally. The
proof is further bound to the server's public key to prevent some MITM attacks.

lifetime
: the duration (in seconds) that the server commits to accept the newly offered
ticket in the future. This period MUST be at least one week.

# Cryptographic Operations {#crypto}

This section provides details on the cryptographic operations performed
by the protocol peers.

## Pinning Secret

On each connection that includes the PinningTicket extension, both peers
compute:

    pinning_secret = HKDF(xSS, xES, "pinning secret", L)

using the notation of {{I-D.ietf-tls-tls13}}, sec. Key Schedule. This secret
is used by the server to generate the new ticket that it returns to the client.

## Pinning Ticket
The pinning ticket's format is not specified by this document, but it MUST be
encrypted and integrity-protected using a long-term pinning-ticket protection key.
The ticket MUST allow the server to recover the pinning_secret value, and MAY
include additional information.

As noted in {{cluster}}, if the server is actually a cluster of machines,
the protection key MUST
be synchronized between them. An easy way to do it is to derive it from the
session-ticket protection key, which is already synchronized. For example:

    pinning_protection_key = HKDF(0, resumption_protection_key,
                                  "pinning protection", L)

## Pinning Proof

The proof sent by the server consists of this value:

    proof = HMAC(old_pinning_secret, "pinning proof" + '\0' +
                 client.random + server.random + Hash(server-public-key))

where HMAC {{RFC2104}} uses the Hash algorithm for the handshake.

# Operational Considerations

The main motivation behind the current protocol is to enable certificate
pinning without the need for manual operations. Manual operations are susceptible
to human error and in the case of certificate pinning, can easily result in
"server bricking": the server becoming inaccessible to some or all of its users.

## Ticket Secret Synchronization {#cluster}

The only operational requirement when deploying this protocol is that if the server
is part of a cluster, ticket
secrets MUST be synchronized between all cluster members. The protocol is designed so
that if resumption tickets {{RFC5077}} are already synchronized between
cluster members, nothing more needs to be done.

Moreover, synchronization does not need
to be instantaneous, e.g. tickets can be distributed a few minutes or hours in advance
of their rollover.

## Certificate Renewal

The protocol ensures that the client will continue speaking to the correct server
even when the server's certificate is renewed. In fact the protocol should more
precisely be called "server identity pinning".

## Certificate Revocation

The protocol is orthogonal to certificate validation, in the sense that, if the
server's certificate has been revoked or is invalid for some other reason,
the client MUST refuse to connect to it.

## Disabling Pinning {#ramp_down}

A server implementing this protocol MUST have a "ramp down" mode of operation where:

  * The server continues to accept valid pinning tickets and responds correctly with
a proof.
  * The server does not send back a new PinningTicket.

After a while no clients will hold valid tickets any more and the feature may be
disabled.

## Server Compromise

If a server compromise is detected, the pinning secret MUST be rotated immediately,
but the server MUST still accept valid tickets that use the old, compromised key.
Clients who still hold old pinning tickets will remain vulnerable to MITM attacks,
but those that connect to the correct server will immediately receive new tickets.

# Comparison: HPKP Deployment {#hpkp}

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
1. Shortly before expiration, change the HPKP header again to send the new main cert as
the main pin and the new backup cert as the backup pin.

To summarize:

|---
| Period | Main server cert | Backup pin | Secondary backup pin
| Regular operation: before rotation | Old main cert | Old backup cert
| >1 month before expiration of old certs |	Old main cert |	New main cert |	New backup cert
| Shortly before expiration but not earlier than the previous change + 1 month | New main cert | New backup cert
| Regular operation: after rotation | New main cert | New backup cert
	

# Security Considerations

## Trust on First Use (TOFU) and MITM Attacks

This protocol is a "trust on first use" protocol. If a client initially connects to
the "right" server, it will be protected against MITM attackers for the lifetime
of each received ticket. If it connects regularly (depending of course on the
server-selected lifetime), it will stay constantly protected against fake
certificates.

However if it initially connects to an attacker, subsequent connections to
the "right" server will fail. Server operators might want to advise clients
on how to remove corrupted pins, if such large scale attacks are detected.

The protocol is designed to deal with an active MITM attacker that has real-time access
to the original server. The pinning proof includes a hash of the server's
public key, to ensure the client that the proof was in fact generated by
the server with which it is initiating the connection.

## Server-Side Error Detection

Uniquely, this protocol the server to detect clients that present incorrect
tickets and therefore can be assumed to be victims of a MITM attack. Server operators
can use such cases as indications of ongoing attacks, similarly to attacks that
were observed in a few countries in the past.

## Client Privacy

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