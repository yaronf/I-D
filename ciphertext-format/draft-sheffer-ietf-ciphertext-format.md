---
title: "A Generic Ciphertext Format"
abbrev: Ciphertext Format
docname: draft-sheffer-ciphertext-format-latest
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
  -
    ins: G. Keselman
    name: Gleb Keselman
    organization: Intuit
    email: gleb.keselman@gmail.com
  -
    ins: Y. Nir
    name: Yoav Nir
    organization: Dell Technologies
    email: ynir.ietf@gmail.com

normative:

informative:

    NISTSP800-38D: DOI.10.6028/NIST.SP.800-38D

--- abstract

This document defines a format for encrypted data, that allows to detect such data in large data stores, and associate it back to the system where it was created and the key with which it was encrypted. This allows organizations to extend the concept of data governance to encrypted data, and to manage such data even when encrypted by multiple systems and cloud providers.

We intend to standardize this format through the IETF, and the document is organized as an Internet Draft.

--- middle

# Introduction and Design Principles

Organizations that manage sensitive data often employ application-level encryption to protect data at rest. When this solution is used, it is common that very large numbers of encrypted data items are stored, potentially for a long time. Security best practices, complicated organizational structures, as well as the existence of modern key management systems, lead to the proliferation of large numbers of encryption keys. After a while it becomes difficult to identify the encryption key that was used for a particular piece of data, with the situation becoming even more complicated when multiple key management systems are used by the same organization.

Tagging encrypted data with metadata supports a number of important use cases: it allows the organization to better catalog the data (a.k.a. “data governance”), to discover the owner of each piece of encrypted data, to detect data encrypted with outdated keys.

## Terminology

{::boilerplate bcp14}

# Motivation

Our main goal in defining a common ciphertext format is to allow organizations to manage large scale data, encrypted at rest using multiple key management and encryption services. Additional motivations for an enterprise to use a common format are:

*   Cross-KMS-provider interoperability, to simplify automated management of data sourced from multiple origins.

*   Proprietary data encryption formats mean that the data remains tied to a single vendor.

*   Standardisation around key management best practices.

## Design Goals

Some of the goals behind this design include:

*   The format should allow simple and efficient detection of encrypted data, in support of automated data governance and key lifecycle management.

*   The format should be space-efficient, since it may be used for very large numbers of small encrypted items. As a result, important information is associated with the (stored) key, rather than the ciphertext.

*   Specifically, following security best practices, a given key material should be used with only a single cryptographic algorithm. Therefore, the algorithm identifier should be stored with the key (or the key version), rather than with the ciphertext.

*   The format defined here only covers the ciphertext header, and not the cipher text itself (referred to as "body" in this document). The body is defined elsewhere, such as {{NISTSP800-38D}} for AES-GCM.

*   The header is not encrypted. Integrity-protection is optional. See {{integrity-protection}} for details.

*   The format should support key versioning, i.e. automated, periodic rotation of keys.

*   The format should support granular key management by allowing for key derivation and key wrapping.

*   The format should allow for generic tools to perform partial attribution of ciphertext, i.e. to associate it with a specific key provider. More specific, possibly provider-specific tools are required for full attribution.

## Previous Work

A few notable formats are:

* The Amazon Web Services SDK message format, documented [here](https://docs.aws.amazon.com/encryption-sdk/latest/developer-guide/message-format.html). This format is specific to the AWS library, and aimed at users of the AWS Key Management System (KMS).
* The [wire format](https://github.com/google/tink/blob/master/docs/WIRE-FORMAT.md) defined by Google's Tink library.
* The format defined by the [KMIP 2.1](https://docs.oasis-open.org/kmip/kmip-profiles/v2.1/csprd01/kmip-profiles-v2.1-csprd01.html) specification, which is targeted at data transmittal, rather than storage.

# The Ciphertext Format

## Format Overview

The ciphertext is prefixed by a header, which in turn, consists of a fixed header, a 2-octet header length, and a sequence of Type-Length-Value (TLV) structures. The TLV structures may appear in any order.

Following the header is the body of the ciphertext. The format (including length) of the body is out of scope for this document.

### Fixed Header

The fixed header consists of:

*   A single constant octet 0x08 (see {{fixed-header-rationale}}).

*   1 octet denoting the format version. The version is 0x01 for the format defined in this document.

*   2 octets for the Key Provider (the organization responsible for the key management system).

### Header Length

The header length field consists of 2 octets of length (covering the complete header including the fixed header), in network order.

### TLV Structures

All other fields consist of:

*   1 octet type

*   2 octet length of the Value field, where the length is encoded in network order. For TV structures, the length field is omitted and the length of the value is implied by the "type" field.

*   Variable-length value.

Issue: we could define “short” vs. “long” TLVs, where for short TLVs the type and length are packed into a single byte (MS nibble for type, LS nibble for length).

The fields are defined by the following table.

| Field Name | Type Tag | Meaning                                                  | Mandatory |
| -------------- | ------------ | ------------------------------------------------------------ | ------------- |
| Key ID         | 1            | An encryption key, as stored in a key  management system. This must denote a unique key, even if the Provider  supports multiple tenants. Encoding of this field is Provider-specific. The  field must appear once. | Y             |
| Key Version    | 2            | A version of a key, where the key is  rotated on a periodic basis. Encoding of this field is Provider-specific. The  field must appear at most once. | N             |
| Auxiliary Data | 3            | Additional data required to derive a specific key from the referenced key  (and key version, if any), see also {{deriving-a-specific-key}}. The field must  appear at most once. | N             |

### Deriving a Specific Key
{: #deriving-a-specific-key}

The Auxiliary Data field is used to support derivation of a key, specific to the ciphertext being managed. There are two common ways to obtain this specific key:

*   Using a key derivation function: SK = KDF(key, aux-data)

*   Decryption of a wrapped key: SK = Decrypt(key, aux-data)

The exact algorithm is implementation dependent, and should be uniquely defined by the combination of Key Provider, Key ID and (if given) Key Version.

## Receiving Ciphertext

Correct interpretation of the format may have security implications, making it important to define the exact semantics even when the entity that receives a ciphertext may not understand parts of the header.

*   A recipient MUST reject a malformed header, e.g. if the total length is larger than the physical length allocated to it based on higher-level network protocols.

*   A recipient MUST reject a ciphertext if it does not recognize the format version.

*   A recipient MUST accept a ciphertext even if it does not recognize some of the TLVs, it MUST ignore the unknown TLVs and MUST interpret all known ones. In other words, the only way to introduce new mandatory TLVs is by incrementing the format version.

*   If ciphertext integrity protection coverage includes the header, a recipient MUST reject the header as well as the ciphertext if the integrity protection fails to validate.

## Fixed Header Rationale
{: #fixed-header-rationale}

We chose the initial byte 0x08, since strings are very unlikely to start with it, as we explain below. Automated tools can detect encrypted data in structured contexts (e.g., a SQL database column) by sampling a number of data items and if all start with this byte, determining that they are encrypted with a high probability.

The byte 0x08 encodes the ASCII control character “backspace”. It has the same meaning in UTF-8, and the 08 block of UTF-16 characters is only populated by two very small languages and rarely-used extended [Arabic characters](https://en.wikipedia.org/wiki/Arabic_Extended-A).

# IANA Considerations

TBD: establish a registry for Types, with 128-255 as private use.

TBD: establish a registry of Key Providers.

# Security Considerations

## Integrity Protection
{: #integrity-protection}

The format defined here does not include integrity protection for the header, and neither does it mandate that the encrypted item's integrity protection should include the header.

Data encrypted at rest is typically vulnerable to denial of service attacks, since (assuming the data is integrity protected) an attacker that can change the ciphertext can trivially cause it to fail validation.

There are cases where it is convenient to manipulate the ciphertext header, even if the data itself remains encrypted and unmodified. For example, when migrating between formats or when bulk-changing metadata associated with the ciphertext. On the other hand, it is a best practice to protect cryptographic metadata against malicious modification. We are currently not aware of a specific threat vector associated with malicious changes to the proposed format, at least assuming the use of AEAD ciphers.

--- back

# Document History

## draft-sheffer-ietf-ciphertext-format-00

Initial version.
