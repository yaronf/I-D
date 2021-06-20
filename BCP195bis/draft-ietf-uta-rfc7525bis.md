---
title: Recommendations for Secure Use of Transport Layer Security (TLS) and Datagram Transport Layer Security (DTLS)
abbrev: TLS Recommendations
docname: draft-ietf-uta-rfc7525bis-latest
category: bcp
obsoletes: 7525

ipr: trust200902
area: Applications
workgroup: UTA Working Group
keyword: Internet-Draft
consensus: false

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
  subcompact: no

author:
 -
    ins: Y. Sheffer
    name: Yaron Sheffer
    organization: Intuit
    email: yaronf.ietf@gmail.com
 -
    ins: R. Holz
    name: Ralph Holz
    organization: University of Twente
    email: ralph.ietf@gmail.com
 -
    ins: P. Saint-Andre
    name: Peter Saint-Andre
    organization: Mozilla
    email: stpeter@mozilla.com

normative:

informative:

  DegabrieleP07: DOI.10.1109/SP.2007.8

  triple-handshake: DOI.10.1109/SP.2014.14

  Soghoian2011: DOI.10.2139/ssrn.1591033

  SESSION-HASH: RFC7627

  POODLE:
    author:
    - org: US-CERT
    date: October 2014
    target: https://www.us-cert.gov/ncas/alerts/TA14-290A
    title: SSL 3.0 Protocol Vulnerability and POODLE Attack

  TLS-XMPP: RFC7590

  CAB-Baseline:
    author:
    - org: CA/Browser Forum
    date: '2013'
    target: https://www.cabforum.org/documents.html
    title: Baseline Requirements for the Issuance and Management of Publicly-Trusted Certificates Version 1.1.6

  Heninger2012:
    author:
    - ins: N. Heninger
      name: Nadia Heninger
    - ins: Z. Durumeric
      name: Zakir Durumeric
    - ins: E. Wustrow
      name: Eric Wustrow
    - ins: J.A. Halderman
      name: J. Alex Halderman
    date: '2012'
    title: 'Mining Your Ps and Qs: Detection of Widespread Weak Keys in Network Devices'
    seriesinfo:
      Usenix Security Symposium: '2012'

  DANE-SMTP: RFC7672

  PatersonRS11: DOI.10.1007/978-3-642-25385-0_20

  DANE-SRV: RFC7673

  Kleinjung2010: DOI.10.1007/978-3-642-14623-7_18

  IANA_TLS:
    author:
    - org: IANA
    target: http://www.iana.org/assignments/tls-parameters
    title: Transport Layer Security (TLS) Parameters
    date: false

  Smith2013:
    author:
    - ins: B. Smith
      name: Brian Smith
    date: '2013'
    target: https://briansmith.org/browser-ciphersuites-01.html
    title: Proposal to Change the Default TLS Ciphersuites Offered by Browsers.

  ECRYPT-II:
    author:
    - ins: N. Smart
      name: Nigel Smart (ed.)
    date: '2012'
    target: http://www.ecrypt.eu.org/ecrypt2/
    title: ECRYPT II Yearly Report on Algorithms and Keysizes (2011-2012)

  Krawczyk2001:
    author:
    - ins: H. Krawczyk
      name: Hugo Krawczyk
    date: '2001'
    seriesinfo:
      CRYPTO: '01'
    target: https://www.iacr.org/archive/crypto2001/21390309.pdf
    title: 'The Order of Encryption and Authentication for Protecting Communications (Or: How Secure is SSL?)'

  Multiple-Encryption: DOI.10.1145/358699.358718 

  BETTERCRYPTO:
    author:
    - org: bettercrypto.org
    date: April 2015
    target: https://bettercrypto.org/
    title: Applied Crypto Hardening

  NIST.SP.800-56A: DOI.10.6028/NIST.SP.800-56Ar3

  DEP-SSLv3: RFC7568

--- abstract
Transport Layer Security (TLS) and Datagram Transport Layer Security (DTLS) are widely used to protect data exchanged over application protocols such as HTTP, SMTP, IMAP, POP, SIP, and XMPP.  Over the last few years, several serious attacks on TLS have emerged, including attacks on its most commonly used cipher suites and their modes of operation.  This document provides recommendations for improving the security of deployed services that use TLS and DTLS. The recommendations are applicable to the majority of use cases.

This document was published as RFC 7525 when the industry was in the midst of its transition to TLS 1.2. Years later this transition is largely complete and TLS 1.3 is widely available. Given the new environment, we believe new guidance is needed.

--- middle

# Introduction

Transport Layer Security (TLS) {{!RFC5246}} and Datagram Transport Security Layer (DTLS) {{!RFC6347}} are widely used to protect data exchanged over application protocols such as HTTP, SMTP, IMAP, POP, SIP, and XMPP.  Over the years leading to 2015, several serious attacks on TLS have emerged, including attacks on its most commonly used cipher suites and their modes of operation.  For instance, both the AES-CBC {{?RFC3602}} and RC4 {{!RFC7465}} encryption algorithms, which together have been the most widely deployed ciphers, have been attacked in the context of TLS.  A companion document {{?RFC7457}} provides detailed information about these attacks and will help the reader understand the rationale behind the recommendations provided here.

The TLS community reacted to these attacks in two ways:

- Detailed guidance was published on the use of TLS 1.2 and earlier protocol versions. This guidance is included in the original {{?RFC7525}} and mostly retained in this revised version.
- A new protocol version was released, TLS 1.3 {{!RFC8446}}, which largely mitigates or resolves these attacks.

Those who implement and deploy TLS and DTLS, in particular versions 1.2 or earlier of these protocols, need guidance on how TLS can be used securely.  This document provides guidance for deployed services as well as for software implementations, assuming the implementer expects his or her code to be deployed in environments defined in {{applicability}}. Concerning deployment, this document targets a wide audience -- namely, all deployers who wish to add authentication (be it one-way only or mutual), confidentiality, and data integrity protection to their communications.

The recommendations herein take into consideration the security of various mechanisms, their technical maturity and interoperability, and their prevalence in implementations at the time of writing.  Unless it is explicitly called out that a recommendation applies to TLS alone or to DTLS alone, each recommendation applies to both TLS and DTLS.

As noted, the TLS 1.3 specification resolves many of the vulnerabilities listed in this document. A system that deploys TLS 1.3 should have fewer vulnerabilities than TLS 1.2 or below. This document is being republished with this in mind, and with an explicit goal to migrate most uses of TLS 1.2 into TLS 1.3.


These are minimum recommendations for the use of TLS in the vast majority of implementation and deployment scenarios, with the exception of unauthenticated TLS (see {{applicability}}). Other specifications that reference this document can have stricter requirements related to one or more aspects of the protocol, based on their particular circumstances (e.g., for use with a particular application protocol); when that is the case, implementers are advised to adhere to those stricter requirements. Furthermore, this document provides a floor, not a ceiling, so stronger options are always allowed (e.g., depending on differing evaluations of the importance of cryptographic strength vs. computational load).


Community knowledge about the strength of various algorithms and feasible attacks can change quickly, and experience shows that a Best Current Practice (BCP) document about security is a point-in-time statement.  Readers are advised to seek out any errata or updates that apply to this document.
    
# Terminology

A number of security-related terms in this document are used in the sense defined in {{!RFC4949}}.

{::boilerplate bcp14}

# General Recommendations
{: #rec}

This section provides general recommendations on the secure use of TLS. Recommendations related to cipher suites are discussed in the following section.

## Protocol Versions

### SSL/TLS Protocol Versions
{: #rec-versions}

It is important both to stop using old, less secure versions of SSL/TLS and to start using modern, more secure versions; therefore, the following are the recommendations concerning TLS/SSL protocol versions:


* Implementations MUST NOT negotiate SSL version 2.

  Rationale: Today, SSLv2 is considered insecure {{!RFC6176}}.

* Implementations MUST NOT negotiate SSL version 3.

  Rationale: SSLv3 {{?RFC6101}} was an improvement over SSLv2 and plugged some significant security holes but did not support strong cipher suites. SSLv3 does not support TLS extensions, some of which (e.g., renegotiation_info {{!RFC5746}}) are security-critical.  In addition, with the emergence of the POODLE attack {{POODLE}}, SSLv3 is now widely recognized as fundamentally insecure.  See {{DEP-SSLv3}} for further details.

* Implementations MUST NOT negotiate TLS version 1.0 {{?RFC2246}}.

  Rationale: TLS 1.0 (published in 1999) does not support many modern, strong cipher suites. In addition, TLS 1.0 lacks a per-record Initialization Vector (IV) for CBC-based cipher suites and does not warn against common padding errors. This and other recommendations in this section are in line with {{!RFC8996}}.

* Implementations MUST NOT negotiate TLS version 1.1 {{?RFC4346}}.


  Rationale: TLS 1.1 (published in 2006) is a security improvement over TLS 1.0 but still does not support certain stronger cipher suites.

  NOTE: This recommendation has been changed from SHOULD NOT to MUST NOT on the assumption that {{!I-D.ietf-tls-oldversions-deprecate}} will be published as an RFC before this document.
            
* Implementations MUST support TLS 1.2 {{!RFC5246}} and MUST prefer to negotiate TLS version 1.2 over earlier versions of TLS.

  Rationale: Several stronger cipher suites are available only with TLS 1.2 (published in 2008). In fact, the cipher suites recommended by this document for TLS 1.2 ({{rec-cipher}} below) are only available in this version.

* Implementations SHOULD support TLS 1.3 {{!RFC8446}} and if implemented, MUST prefer to negotiate TLS 1.3 over earlier versions of TLS.
               
  Rationale: TLS 1.3 is a major overhaul to the protocol and resolves many of the security issues with TLS 1.2. We note that as long as TLS 1.2 is still allowed by a particular implementation, even if it defaults to TLS 1.3, implementers MUST still follow all the recommendations in this document.

* Implementations of "greenfield" protocols or deployments, where there is no need to support legacy endpoints, SHOULD support TLS 1.3, with no negotiation of earlier versions. Similarly, we RECOMMEND that new protocol designs that embed the TLS mechanisms (such as QUIC has done {{?RFC9001}}) include TLS 1.3.

  Rationale: secure deployment of TLS 1.3 is significantly easier and less error prone than the secure deployment of TLS 1.2.

This BCP applies to TLS 1.2, 1.3 and to earlier versions. It is not safe for readers to assume that the recommendations in this BCP apply to any future version of TLS.

### DTLS Protocol Versions

DTLS, an adaptation of TLS for UDP datagrams, was introduced when TLS 1.1 was published.  The following are the recommendations with respect to DTLS:

* Implementations MUST NOT negotiate DTLS version 1.0 {{?RFC4347}}.

  Version 1.0 of DTLS correlates to version 1.1 of TLS (see above).

* Implementations MUST support and (unless a higher version is available) MUST prefer to negotiate DTLS version 1.2 {{!RFC6347}} 


  Version 1.2 of DTLS correlates to version 1.2 of TLS (see above).
  (There is no version 1.1 of DTLS.)

* Implementations SHOULD support and, if available, MUST prefer to negotiate DTLS version 1.3 as specified in {{!I-D.ietf-tls-dtls13}}.

  Version 1.3 of DTLS correlates to version 1.3 of TLS (see above).

### Fallback to Lower Versions
{: #rec-fallback}

Clients that "fall back" to lower versions of the protocol after the server rejects higher versions of the protocol MUST NOT fall back to SSLv3 or earlier. Implementations of TLS/DTLS 1.2 or earlier MUST implement the Fallback SCSV mechanism {{!RFC7507}} to prevent
such fallback being forced by an attacker.

Rationale: Some client implementations revert to lower versions of TLS or even to SSLv3 if the server rejected higher versions of the protocol. This fallback can be forced by a man-in-the-middle (MITM) attacker. TLS 1.0 and SSLv3 are significantly less secure than TLS 1.2 but at least TLS 1.0 is still allowed by many web servers. As of this writing, the Fallback SCSV solution is widely deployed and proven as a robust solution to this problem.

## Strict TLS


The following recommendations are provided to help prevent SSL Stripping (an attack that is summarized in Section 2.1 of {{?RFC7457}}):
      
* In cases where an application protocol allows implementations or deployments a choice between strict TLS configuration and dynamic upgrade from unencrypted to TLS-protected traffic (such as STARTTLS), clients and servers SHOULD prefer strict TLS configuration.



* Application protocols typically provide a way for the server to offer TLS during an initial protocol exchange, and sometimes also provide a way for the server to advertise support for TLS (e.g., through a flag indicating that TLS is required); unfortunately, these indications are sent before the communication channel is encrypted. A client SHOULD attempt to negotiate TLS even if these indications are not communicated by the server.



* HTTP client and server implementations MUST support the HTTP Strict Transport
      Security (HSTS) header {{?RFC6797}}, in order to allow Web servers to 
      advertise that they are
      willing to accept TLS-only clients.



* Web servers SHOULD use HSTS to indicate that they are willing to accept TLS-only clients, unless they are deployed in such a way that using HSTS would in fact weaken overall security (e.g., it can be problematic to use HSTS with self-signed certificates, as described in Section 11.3 of {{?RFC6797}}).
      

Rationale: Combining unprotected and TLS-protected communication opens the way to SSL Stripping and similar attacks, since an initial part of the communication is not integrity protected and therefore can be manipulated by an attacker whose goal is to keep the communication in the clear. 


## Compression
{: #rec-compress}

In order to help prevent compression-related attacks (summarized in Section 2.6 of {{?RFC7457}}), when using TLS 1.2 implementations and deployments SHOULD disable TLS-level compression (Section 6.2.2 of {{!RFC5246}}), unless the application protocol in question has been shown not to be open to such attacks. Note: this recommendation applies to TLS 1.2 only, because compression has been removed from TLS 1.3.


Rationale: TLS compression has been subject to security attacks, such as the CRIME attack.

Implementers should note that compression at higher protocol levels can allow an active attacker to extract cleartext information from the connection. The BREACH attack is one such case. These issues can only be mitigated outside of TLS and are thus outside the scope of this document. See Section 2.6 of {{?RFC7457}} for further details.
      

## TLS Session Resumption
{: #rec-resume}

If TLS session resumption is used in the context of TLS 1.2, care ought to be taken to do so safely. In particular, when using session tickets {{?RFC5077}}, the resumption information MUST be authenticated and encrypted to prevent modification or eavesdropping by an attacker. Further recommendations apply to session tickets:

* A strong cipher suite MUST be used when encrypting the ticket (as least as strong as the main TLS cipher suite).

* Ticket keys MUST be changed regularly, e.g., once every week, so as not to negate the benefits of forward secrecy (see {{sec-pfs}} for details on forward secrecy).

* For similar reasons, session ticket validity SHOULD be limited to a reasonable duration (e.g., half as long as ticket key validity).



Rationale: session resumption is another kind of TLS handshake, and therefore must be as secure as the initial handshake. This document ({{detail}}) recommends the use of cipher suites that provide forward secrecy, i.e. that prevent an attacker who gains momentary access to the TLS endpoint (either client or server) and its secrets from reading either past or future communication. The tickets must be managed so as not to negate this security property.


## TLS Renegotiation

Where handshake renegotiation is implemented, both clients and servers MUST implement the renegotiation_info extension, as defined in {{!RFC5746}}. Note: this recommendation applies to TLS 1.2 only, because renegotiation has been removed from TLS 1.3.

The most secure option for countering the Triple Handshake attack is to refuse any change of certificates during renegotiation.  In addition, TLS clients SHOULD apply the same validation policy for all certificates received over a connection.  The {{triple-handshake}} document suggests several other possible countermeasures, such as binding the master secret to the full handshake (see {{SESSION-HASH}}) and binding the abbreviated session resumption handshake to the original full handshake.  Although the latter two techniques are still under development and thus do not qualify as current practices, those who implement and deploy TLS are advised to watch for further development of appropriate countermeasures.
      
      

## Server Name Indication

TLS implementations MUST support the Server Name Indication (SNI) extension defined in Section 3 of {{!RFC6066}} for those higher-level protocols that would benefit from it, including HTTPS. However, the actual use of SNI in particular circumstances is a matter of local policy.  Implementers are strongly encouraged to support TLS Encrypted Client Hello (formerly called Encrypted SNI) once {{?I-D.ietf-tls-esni}} has been standardized.



Rationale: SNI supports deployment of multiple TLS-protected virtual servers on a single
      address, and therefore enables fine-grained security for these virtual servers,
      by allowing each one to have its own certificate. However, SNI also leaks the 
      target domain for a given connection; this information leak will be plugged by 
      use of TLS Encrypted Client Hello.

## Zero Round Trip Time (0-RTT) Data in TLS 1.3

The 0-RTT early data feature is new in TLS 1.3. It provides improved latency
when TLS connections are resumed, at the potential cost of security.
As a result, it requires special attention from implementers on both
the server and the client side. Typically this extends to both the
TLS library as well as protocol layers above it.

For use in HTTP-over-TLS, readers are referred to {{?RFC8470}} for guidance.

For QUIC-on-TLS, refer to Sec. 9.2 of {{?RFC9001}}.

For other protocols, generic guidance is given in Sec. 8 and Appendix E.5
of {{RFC8446}}.
Given the complexity, we RECOMMEND to avoid this feature altogether unless
an explicit specification exists for the application protocol in question to clarify
when 0-RTT is appropriate and secure. This can take the form of an IETF RFC,
a non-IETF standard, or even documentation associated with a non-standard protocol.

# Recommendations: Cipher Suites
{: #detail}

TLS and its implementations provide considerable flexibility in the
         selection of cipher suites. Unfortunately, some available cipher
         suites are insecure, some do not provide the targeted security
         services, and some no longer provide enough security.  Incorrectly
         configuring a server leads to no or reduced security.  This section
         includes recommendations on the selection and negotiation of
         cipher suites.
         

## General Guidelines
{: #rec-cipher-guidelines}


Cryptographic algorithms weaken over time as cryptanalysis improves: algorithms that were once considered strong become weak. Such algorithms need to be phased out over time and replaced with more secure cipher suites. This helps to ensure that the desired security properties still hold. SSL/TLS has been in existence for almost 20 years and many of the cipher suites that have been recommended in various versions of SSL/TLS are now considered weak or at least not as strong as desired. Therefore, this section modernizes the recommendations concerning cipher suite selection.



* Implementations MUST NOT negotiate the cipher suites with NULL encryption.
               
  Rationale: The NULL cipher suites do not encrypt traffic and 
               so provide no confidentiality services. Any entity in the 
               network with access to the connection can view the plaintext 
               of contents being exchanged by the client and server.  
               Nevertheless, this document does not discourage software from
               implementing NULL cipher suites, since they can be useful for 
               testing and debugging. 
            

* Implementations MUST NOT negotiate RC4 cipher suites. 
               
  Rationale: The RC4 stream cipher has a variety of cryptographic 
               weaknesses, as documented in {{!RFC7465}}.
       Note that DTLS specifically forbids the use of RC4 already.
            

* Implementations MUST NOT negotiate cipher suites offering less 
               than 112 bits of security, including so-called "export-level" 
               encryption (which provide 40 or 56 bits of security).

  Rationale: Based on {{!RFC3766}}, at least 112 bits 
               of security is needed.  40-bit and 56-bit security are considered 
               insecure today.  TLS 1.1 and 1.2 never negotiate 40-bit or 56-bit 
               export ciphers.
            

* Implementations SHOULD NOT negotiate cipher suites that use 
               algorithms offering less than 128 bits of security.

  Rationale: Cipher suites that offer between 112-bits and 128-bits 
               of security are not considered weak at this time; however, it is 
               expected that their useful lifespan is short enough to justify 
               supporting stronger cipher suites at this time.  128-bit ciphers 
               are expected to remain secure for at least several years, and 
               256-bit ciphers until the next fundamental technology 
               breakthrough.  Note that, because of so-called 
               "meet-in-the-middle" attacks {{Multiple-Encryption}},
               some legacy cipher suites (e.g., 168-bit 3DES) have an effective 
               key length that is smaller than their nominal key length (112 
               bits in the case of 3DES).  Such cipher suites should be 
               evaluated according to their effective key length.
            

* Implementations SHOULD NOT negotiate cipher suites based on 
               RSA key transport, a.k.a. "static RSA".

  Rationale: These cipher suites, which have assigned values starting 
               with the string "TLS_RSA_WITH_*", have several drawbacks, especially
               the fact that they do not support forward secrecy.
            

* Implementations MUST support and prefer to negotiate cipher suites 
               offering forward secrecy, such as those in the Ephemeral 
               Diffie-Hellman and Elliptic Curve Ephemeral Diffie-Hellman ("DHE" 
               and "ECDHE") families.

  Rationale: Forward secrecy (sometimes called "perfect forward 
               secrecy") prevents the recovery of information that was encrypted 
               with older session keys, thus limiting the amount of time during 
               which attacks can be successful. See {{sec-pfs}} for 
               a detailed discussion.
            
## Recommended Cipher Suites
{: #rec-cipher}

Given the foregoing considerations, implementation and deployment of the following cipher suites is RECOMMENDED:


* TLS_DHE_RSA_WITH_AES_128_GCM_SHA256

* TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256

* TLS_DHE_RSA_WITH_AES_256_GCM_SHA384

* TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
        

These cipher suites are supported only in TLS 1.2 and not in earlier protocol versions, because they are authenticated encryption (AEAD) algorithms {{?RFC5116}}.

Typically, in order to prefer these suites, the order of suites needs to be explicitly configured in server software. (See {{BETTERCRYPTO}} for helpful deployment guidelines, but note that its recommendations differ from the current document in some details.)  It would be ideal if server software implementations were to prefer these suites by default.

Some devices have hardware support for AES-CCM but not AES-GCM, so they are unable to follow the foregoing recommendations regarding cipher suites.  There are even devices that do not support public key cryptography at all, but they are out of scope entirely.

### Implementation Details
{: #detail-neg}

Clients SHOULD include TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 as the first proposal to any server, unless they have prior knowledge that the server cannot respond to a TLS 1.2 client_hello message.

Servers MUST prefer this cipher suite over weaker cipher suites whenever it is proposed, even if it is not the first proposal.

Clients are of course free to offer stronger cipher suites, e.g., using AES-256; when they do, the server SHOULD prefer the stronger cipher suite unless there are compelling reasons (e.g., seriously degraded performance) to choose otherwise.

This document does not change the mandatory-to-implement TLS cipher suite(s) prescribed by TLS. To maximize interoperability, RFC 5246 mandates implementation of the TLS_RSA_WITH_AES_128_CBC_SHA cipher suite, which is significantly weaker than the cipher suites recommended here. (The GCM mode does not suffer from the same weakness, caused by the order of MAC-then-Encrypt in TLS {{Krawczyk2001}}, since it uses an AEAD mode of operation.) Implementers should consider the interoperability gain against the loss in security when deploying the TLS_RSA_WITH_AES_128_CBC_SHA cipher suite. Other application protocols specify other cipher suites as mandatory to implement (MTI).

Note that some profiles of TLS 1.2 use different cipher suites. For example, {{?RFC6460}} defines a profile that uses the TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 and TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384 cipher suites.

{{!RFC4492}} allows clients and servers to negotiate ECDH parameters (curves).  Both clients and servers SHOULD include the "Supported Elliptic Curves" extension {{!RFC4492}}.  For interoperability, clients and servers SHOULD support the NIST P-256 (secp256r1) curve {{!RFC4492}}. In addition, clients SHOULD send an ec_point_formats extension with a single element, "uncompressed".

## Cipher Suites for TLS 1.3

This document does not specify any cipher suites for TLS 1.3. Readers
are referred to Sec. 9.1 of {{RFC8446}} for cipher suite recommendations.

## Public Key Length
{: #rec-keylength}

When using the cipher suites recommended in this document, two public keys are 
      normally used in the TLS handshake: one for the Diffie-Hellman key agreement
      and one for server authentication. Where a client certificate is used, a third 
      public key is added.

With a key exchange based on modular exponential (MODP) Diffie-Hellman groups ("DHE" cipher suites), DH key lengths of at least 2048 bits are RECOMMENDED.

 Rationale: For various reasons, in practice, DH keys are typically generated in lengths
 that are powers of two (e.g., 2^10 = 1024 bits, 2^11 = 2048 bits, 2^12 = 4096 bits).
 Because a DH key of 1228 bits would be roughly equivalent to only an 80-bit symmetric key
{{?RFC3766}}, it is better to use keys longer than that for the "DHE" family of cipher suites.
A DH key of 1926 bits would be roughly equivalent to a 100-bit symmetric key {{?RFC3766}}.
A DH key of 2048 bits (equivalent to a 112-bit symmetric key) 
is allowed by the latest revision of {{NIST.SP.800-56A}} (see in particular Appendix D).
See {{detail-alt}} for additional
information on the use of MODP Diffie-Hellman in TLS.

As noted in {{?RFC3766}}, correcting for the emergence of a TWIRL machine would imply that 1024-bit DH keys yield about 65 bits of equivalent strength and that a 2048-bit DH key would yield about 92 bits of equivalent strength.

With regard to ECDH keys, the IANA "EC Named Curve Registry" (within the
   "Transport Layer Security (TLS) Parameters" registry {{IANA_TLS}}) contains 160-bit
elliptic curves that are considered to be roughly equivalent to only an 80-bit
symmetric key {{ECRYPT-II}}.   Curves of less than 192 bits SHOULD NOT be used.

When using RSA, servers SHOULD authenticate using certificates with at least a 2048-bit modulus for the public key.  In addition, the use of the SHA-256 hash algorithm is RECOMMENDED (see {{CAB-Baseline}} for more details). Clients SHOULD indicate to servers that they request SHA-256, by using the "Signature Algorithms" extension defined in TLS 1.2. 

## Modular Exponential vs. Elliptic Curve DH Cipher Suites
{: #detail-alt}

[[TODO: revise after Logjam?]]

Not all TLS implementations support both modular exponential (MODP) and elliptic curve (EC) Diffie-Hellman groups, as required by {{rec-cipher}}. Some implementations are severely limited in the length of DH values. When such implementations need to be accommodated, the following are RECOMMENDED (in priority order):

1. Elliptic Curve DHE with appropriately negotiated parameters (e.g., the curve to be used) and a Message Authentication Code (MAC) algorithm stronger than HMAC-SHA1 {{!RFC5289}}

1. TLS_DHE_RSA_WITH_AES_128_GCM_SHA256 {{!RFC5288}}, with 2048-bit Diffie-Hellman parameters

1. TLS_DHE_RSA_WITH_AES_128_GCM_SHA256, with 1024-bit parameters
        

Rationale: Although Elliptic Curve Cryptography is widely deployed, there are some communities where its adoption has been limited for several reasons, including its complexity compared to modular arithmetic and longstanding perceptions of IPR concerns (which, for the most part, have now been resolved {{?RFC6090}}).  Note that ECDHE cipher suites exist for both RSA and ECDSA certificates, so moving to ECDHE cipher suites does not require moving away from RSA-based certificates.  On the other hand, there are two related issues hindering effective use of MODP Diffie-Hellman cipher suites in TLS:

* There are no standardized, widely implemented protocol mechanisms to negotiate the DH groups or parameter lengths supported by client and server.

* Many servers choose DH parameters of 1024 bits or fewer.

* There are widely deployed client implementations that reject received DH parameters if they are longer than 1024 bits.  In addition, several implementations do not perform appropriate validation of group parameters and are vulnerable to attacks referenced in Section 2.9 of {{?RFC7457}}.
        

Note that with DHE and ECDHE cipher suites, the TLS master key only depends on the Diffie-Hellman parameters and not on the strength of the RSA certificate; moreover, 1024 bit MODP DH parameters are generally considered insufficient at this time.

With MODP ephemeral DH, deployers ought to carefully evaluate interoperability vs. security considerations when configuring their TLS endpoints.
      
      
## Truncated HMAC

Implementations MUST NOT use the Truncated HMAC extension, defined in Section 7 of {{!RFC6066}}.



Rationale: the extension does not apply to the AEAD
      cipher suites recommended above. However it does apply to most other TLS cipher suites. Its use
      has been shown to be insecure in {{PatersonRS11}}.

# Applicability Statement
{: #applicability}


The recommendations of this document primarily apply to the implementation and deployment of application protocols that are most commonly used with TLS and DTLS on the Internet today.  Examples include, but are not limited to:


* Web software and services that wish to protect HTTP traffic with TLS.

* Email software and services that wish to protect IMAP, POP3, or SMTP traffic with TLS.

* Instant-messaging software and services that wish to protect Extensible Messaging and Presence Protocol (XMPP) or Internet Relay Chat (IRC) traffic with TLS.

* Realtime media software and services that wish to protect Secure Realtime Transport Protocol (SRTP) traffic with DTLS.
        

This document does not modify the implementation and deployment recommendations (e.g., mandatory-to-implement cipher suites) prescribed by existing application protocols that employ TLS or DTLS. If the community that uses such an application protocol wishes to modernize its usage of TLS or DTLS to be consistent with the best practices recommended here, it needs to explicitly update the existing application protocol definition (one example is {{TLS-XMPP}}, which updates {{?RFC6120}}). 



  Designers of new application protocols developed through the Internet
  Standards Process {{?RFC2026}} are expected at minimum to conform to the best
  practices recommended here, unless they provide documentation of
  compelling reasons that would prevent such conformance (e.g.,
  widespread deployment on constrained devices that lack support for
  the necessary algorithms).


## Security Services

This document provides recommendations for an audience that wishes to secure their communication with TLS to achieve the following:


* Confidentiality: all application-layer communication is encrypted with the goal that no party should be able to decrypt it except the intended receiver.

* Data integrity: any changes made to the communication in transit are detectable by the receiver.

* Authentication: an endpoint of the TLS communication is authenticated as the intended entity to communicate with.
          

With regard to authentication, TLS enables authentication of one or both endpoints in the communication.  In the context of opportunistic security {{?RFC7435}}, TLS is sometimes used without authentication. As discussed in {{oppsec}}, considerations for opportunistic security are not in scope for this document.

If deployers deviate from the recommendations given in this document, they need to be aware that they might lose access to one of the foregoing security services.

This document applies only to environments where confidentiality is required. It recommends algorithms and configuration options that enforce secrecy of the data in transit.

This document also assumes that data integrity protection is always one of the goals of a deployment. In cases where integrity is not required, it does not make sense to employ TLS in the first place. There are attacks against confidentiality-only protection that utilize the lack of integrity to also break confidentiality (see, for instance, {{DegabrieleP07}} in the context of IPsec).

This document addresses itself to application protocols that are most commonly used on the Internet with TLS and DTLS. Typically, all communication between TLS clients and TLS servers requires all three of the above security services. This is particularly true where TLS clients are user agents like Web browsers or email software.

This document does not address the rarer deployment scenarios where one of the above three properties is not desired, such as the use case described in {{oppsec}} below.  As another scenario where confidentiality is not needed, consider a monitored network where the authorities in charge of the respective traffic domain require full access to unencrypted (plaintext) traffic, and where users collaborate and send their traffic in the clear.
        
## Opportunistic Security
{: #oppsec}

There are several important scenarios in which the use of TLS is optional, i.e., the client decides dynamically ("opportunistically") whether to use TLS with a particular server or to connect in the clear.  This practice, often called "opportunistic security", is described at length in {{?RFC7435}} and is often motivated by a desire for backward compatibility with legacy deployments.

In these scenarios, some of the recommendations in this document might be too strict, since adhering to them could cause fallback to cleartext, a worse outcome than using TLS with an outdated protocol version or cipher suite.

# Security Considerations
{: #sec}

This entire document discusses the security practices directly affecting applications
    using the TLS protocol. This section contains broader security considerations related
    to technologies used in conjunction with or by TLS.

## Host Name Validation

  Application authors should take note that some TLS implementations
  do not validate host names.  If the TLS implementation they are
  using does not validate host names, authors might need to write their
  own validation code or consider using a different TLS implementation.


It is noted that the requirements regarding host name validation (and, in general, binding between the TLS layer and the protocol that runs above it) vary between different protocols. For HTTPS, these requirements are defined by Section 3 of {{!RFC2818}}.

Readers are referred to {{!RFC6125}} for further details regarding generic host name validation in the TLS context. In addition, that RFC contains a long list of example protocols, some of which implement a policy very different from HTTPS.

If the host name is discovered indirectly and in an insecure manner (e.g., by an insecure DNS query for an MX or SRV record), it SHOULD NOT be used as a reference identifier {{!RFC6125}} even when it matches the presented certificate.  This proviso does not apply if the host name is discovered securely (for further discussion, see {{DANE-SRV}} and {{DANE-SMTP}}).

Host name validation typically applies only to the leaf "end entity" certificate. Naturally, in order to ensure proper authentication in the context of the PKI, application clients need to verify the entire certification path in accordance with {{?RFC5280}} (see also 
        {{!RFC6125}}).

## AES-GCM
{: #sec-aes}

{{rec-cipher}} above recommends the use of the AES-GCM authenticated encryption algorithm. Please refer to Section 11 of {{RFC5246}} for general security considerations when using TLS 1.2, and to Section 6 of {{!RFC5288}} for security considerations that apply specifically to AES-GCM when used with TLS.
      
## Forward Secrecy
{: #sec-pfs}

Forward secrecy (also called "perfect forward secrecy" or "PFS" and defined in {{?RFC4949}}) is a defense against an attacker who records encrypted conversations where the session keys are only encrypted with the communicating parties' long-term keys. 

Should the attacker be able to obtain these long-term keys at some point later in time, the session keys and thus the entire conversation could be decrypted.

In the context of TLS and DTLS, such compromise of long-term keys is not entirely implausible. It can happen, for example, due to:

* A client or server being attacked by some other attack vector, and the private key retrieved.

* A long-term key retrieved from a device that has been sold or otherwise decommissioned without prior wiping.

* A long-term key used on a device as a default key {{Heninger2012}}.

* A key generated by a trusted third party like a CA, and later retrieved from it either by extortion or compromise {{Soghoian2011}}.

* A cryptographic break-through, or the use of asymmetric keys with insufficient length {{Kleinjung2010}}.

* Social engineering attacks against system administrators.

* Collection of private keys from inadequately protected backups.
        


Forward secrecy ensures in such cases that it is not feasible for an attacker to determine the session keys even if the attacker has obtained the long-term keys some time after the conversation. It also protects against an attacker who is in possession of the long-term keys but remains passive during the conversation.


Forward secrecy is generally achieved by using the Diffie-Hellman scheme to derive session keys. The Diffie-Hellman scheme has both parties maintain private secrets and send parameters over the network as modular powers over certain cyclic groups. The properties of the so-called Discrete Logarithm Problem (DLP) allow the parties to derive the session keys without an eavesdropper being able to do so. There is currently no known attack against DLP if sufficiently large parameters are chosen. A variant of the Diffie-Hellman scheme uses Elliptic Curves instead of the originally proposed modular arithmetics.

Unfortunately, many TLS/DTLS cipher suites were defined that do not feature forward secrecy, e.g., TLS_RSA_WITH_AES_256_CBC_SHA256.  This document therefore advocates strict use of forward-secrecy-only ciphers.
      
      
## Diffie-Hellman Exponent Reuse

For performance reasons, many TLS implementations reuse Diffie-Hellman and Elliptic Curve Diffie-Hellman exponents across multiple connections. Such reuse can result in major security issues:

* If exponents are reused for too long (e.g., even more than a few hours), an attacker      who gains access to the host can decrypt previous connections. In other words, exponent reuse negates the effects of forward secrecy.

* TLS implementations that reuse exponents should test the DH public key they receive for group membership, in order to avoid some known attacks. These tests are not standardized in TLS at the time of writing. See {{?RFC6989}} for recipient tests required of IKEv2 implementations that reuse DH exponents.
      
      
## Certificate Revocation

The following considerations and recommendations represent the current state of the art regarding certificate revocation, even though no complete and efficient solution exists for the problem of checking the revocation status of common public key certificates {{RFC5280}}:


* Although Certificate Revocation Lists (CRLs) are the most widely supported mechanism for distributing revocation information, they have known scaling challenges that limit their usefulness (despite workarounds such as partitioned CRLs and delta CRLs).

* Proprietary mechanisms that embed revocation lists in the Web browser's configuration database cannot scale beyond a small number of the most heavily used Web servers.

* The On-Line Certification Status Protocol (OCSP) {{?RFC6960}} presents both scaling and privacy issues. In addition, clients typically "soft-fail", meaning that they do not abort the TLS connection if the OCSP server does not respond. (However, this might be a workaround to avoid denial-of-service attacks if an OCSP responder is taken offline.)

* The TLS Certificate Status Request extension (Section 8 of {{!RFC6066}}), commonly called "OCSP stapling", resolves the operational issues with OCSP. However, it is still ineffective in the presence of a MITM attacker because the attacker can simply ignore the client's request for a stapled OCSP response.

* OCSP stapling as defined in {{!RFC6066}} does not extend to intermediate certificates used in a certificate chain. Although the Multiple Certificate Status extension {{?RFC6961}} addresses this shortcoming, it is a recent addition without much deployment.

* Both CRLs and OCSP depend on relatively reliable connectivity to the Internet, which might not be available to certain kinds of nodes (such as newly provisioned devices that need to establish a secure connection in order to boot up for the first time).
      

With regard to common public key certificates, servers SHOULD support the following as a best practice given the current state of the art and as a foundation for a possible future solution:

1. OCSP {{?RFC6960}}

1. Both the status\_request extension defined in {{!RFC6066}} and the status_request_v2 extension defined in {{?RFC6961}} (This might enable interoperability with the widest range of clients.)

1. The OCSP stapling extension defined in {{?RFC6961}}
      

The considerations in this section do not apply to scenarios where the DANE-TLSA resource record {{?RFC6698}} is used to signal to a client which certificate a server considers valid and good to use for TLS connections.


# Acknowledgments
{: #d1e1127}

The following acknowledgements are inherited from {{?RFC7525}}.

Thanks to RJ Atkinson, Uri Blumenthal, Viktor Dukhovni, Stephen Farrell, Daniel Kahn Gillmor, Paul Hoffman, Simon Josefsson, Watson Ladd, Orit Levin, Ilari Liusvaara, Johannes Merkle, Bodo Moeller, Yoav Nir, Massimiliano Pala, Kenny Paterson, Patrick Pelletier, Tom Ritter, Joe St. Sauver, Joe Salowey, Rich Salz, Brian Smith, Sean Turner, and Aaron Zauner for their feedback and suggested improvements. Thanks also to Brian Smith, who has provided a great resource in his "Proposal to Change the Default TLS Ciphersuites Offered by Browsers" {{Smith2013}}. Finally, thanks to all others who commented on the TLS, UTA, and other discussion lists but who are not mentioned here by name.

Robert Sparks and Dave Waltermire provided helpful reviews on behalf of the General Area Review Team and the Security Directorate, respectively.

During IESG review, Richard Barnes, Alissa Cooper, Spencer Dawkins, Stephen Farrell, Barry Leiba, Kathleen Moriarty, and Pete Resnick provided comments that led to further improvements.

Ralph Holz gratefully acknowledges the support by Technische Universitaet Muenchen.

The authors gratefully acknowledge the assistance of Leif Johansson and Orit Levin as the working group chairs and Pete Resnick as the sponsoring Area Director.

--- back

# Differences from RFC 7525

* Clarified some items (e.g. renegotiation) that only apply to TLS 1.2 - many more TBD.
* Changed status of TLS 1.0 and 1.1 from SHOULD NOT to MUST NOT.
* Added TLS 1.3 at a SHOULD level.
* Similar changes to DTLS, pending publication of DTLS 1.3.
* Fallback SCSV as a MUST for TLS 1.2.
* Added mention of TLS Encrypted Client Hello, but no recommendation to use yet.
* New TLS 1.3 capabilities: 0-RTT.

# Document History

[[Note to RFC Editor: please remove before publication.]]

## draft-ietf-uta-rfc7525bis-00

- Renamed: WG document.
- Started populating list of changes from RFC 7525.
- General rewording of abstract and intro for revised version.
- Protocol versions, fallback.
- Reference to ECHO.

## draft-sheffer-uta-rfc7525bis-00

- Renamed, since the BCP number does not change.
- Added an empty "Differences from RFC 7525" section.

## draft-sheffer-uta-bcp195bis-00

- Initial release, the RFC 7525 text as-is, with some minor editorial
changes to the references.

