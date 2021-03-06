---
title: "A Generic Ciphertext Format"
abbrev: Ciphertext Format
docname: draft-sheffer-ietf-ciphertext-format-latest
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

    RFC8949:

informative:

    NISTSP800-38D: DOI.10.6028/NIST.SP.800-38D

--- abstract

This document defines a set of structured headers for encrypted data. The main goal of this format is to enable detection of encrypted data in large data stores, and associating it back to the system where it was created and the key with which it was encrypted. This allows organizations to extend the concept of data governance to encrypted data, and to manage such data even when encrypted by multiple different systems and cloud providers.

--- middle

# Introduction and Design Principles

Organizations that manage sensitive data often employ application-level encryption to protect data at rest. When this solution is used, it is common that very large numbers of encrypted data items are stored, potentially for a long time. Security best practices, complicated organizational structures, as well as the existence of modern key management systems, lead to the proliferation of large numbers of encryption keys. After a while it becomes difficult to identify the encryption key that was used for a particular piece of data, with the situation becoming even more complicated when multiple key management systems are used by the same organization.

Application-level encryption can be deployed at different scales: in some cases a multi-megabyte file may be encrypted with a single key. In other cases, we may want to deploy encryption for specific database fields, which can easily manifest itself as millions of keys for a single database table.

Tagging encrypted data with metadata supports a number of important use cases: it allows the organization to better catalog the data (a.k.a. “data governance”), to discover the owner of each piece of encrypted data, to detect data encrypted with outdated keys.

## Terminology

{::boilerplate bcp14}

# Motivation

Our main goal in defining a common ciphertext format is to allow organizations to manage large scale data, encrypted at rest using multiple key management and encryption services. Additional motivations for an enterprise to use a common format are:

*   Cross-KMS-provider interoperability, to simplify automated management of data sourced from multiple origins.

*   Proprietary data encryption formats mean that the data remains tied to a single vendor.

*   Standardization around key management best practices.

## Design Goals

Some of the goals behind this design include:

*   The format should allow simple and efficient detection of encrypted data, in support of automated data governance and key lifecycle management.

*   The format should be space-efficient, since it may be used for very large numbers of small encrypted items. As a result, important information is associated with the (stored) key, rather than the ciphertext.

*   Specifically, following security best practices, a given key material should be used with only a single cryptographic algorithm. Therefore, the algorithm identifier should be stored with the key (or the key version), rather than with the ciphertext.

*   The format defined here only covers the ciphertext header, and not the ciphertext itself (referred to as "body" in this document). The body is defined elsewhere, such as {{NISTSP800-38D}} for AES-GCM.

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

The ciphertext is prefixed by a header, which in turn, consists of a short fixed header and variable header. The variable header is a CBOR {{RFC8949}} map.

Following the header is the body of the ciphertext. The format (including length) of the body is out of scope for this document.

### Fixed Header

The fixed header consists of:

*   A single constant octet 0x08 (see {{fixed-header-rationale}}).
*   A single octet denoting the format version. The version is 0x01 for the format defined in this document.

### Variable Header

The variable header is a CBOR map consisting of elements from the following table.

| Field Name | Map Key | Value Type       | Meaning                                                      | Mandatory |
| -------------- | ------------ | ------------------------------------------------------------ | ------------- | -------------- |
| Key Provider   | 1       | Unsigned integer | The organization responsible for the key management system.  | Y |
| Key ID         | 2       | Byte string      | An encryption key identifier, where the key is stored in a key management system. This must denote a unique key, even if the Provider supports multiple tenants. Encoding of this field is Provider-specific. The field must appear once. | Y |
| Key Version    | 3       | Unsigned integer | A version of a key, where the key is rotated on a periodic basis. Encoding of this field is Provider-specific. The field must appear at most once. | N |
| Auxiliary Data | 4 | Byte string      | Additional data required to derive a specific key from the referenced key (and key version, if any), see also {{deriving-a-specific-key}}. The field must appear at most once. | N |
| Nonce | 5 | Byte string | A nonce or initialization vector (IV), if required by the cipher algorithm. We note that an implementation may prefer to store the nonce and authentication tag in-line with the ciphertext. | N |
| Authentication Tag | 6 | Byte string | An authentication tag or integrity check value (ICV), if required by the cipher algorithm. | N |
| Additional Authenticated Data | 7 | Byte string | Additional authenticated data (AAD), which is integrity-protected but not encrypted by the cipher. | N |
| Structured Additional Authenticated Data | 8 | Any | Additional authenticated data (AAD), encoded as an arbitrary CBOR structure. See {{aad-multiple}}. | N |

### Deriving a Specific Key
{: #deriving-a-specific-key}

The Auxiliary Data field is used to support derivation of a key, specific to the ciphertext being managed. There are two common ways to obtain this specific key:

*   Using a key derivation function: SK = KDF(key, aux-data)

*   Decryption of a wrapped key: SK = Decrypt(key, aux-data)

The exact algorithm is implementation dependent, and should be uniquely defined by the combination of Key Provider, Key ID and (if given) Key Version.

### Multiple Values of Additional Authenticated Data
{: #aad-multiple}

When multiple values of Additional Authenticated Data (AAD) are needed, they must be carefully encoded to prevent security issues as a result of ambiguities.

For example, if we were to simply concatenate the AAD values, then the following two sequences of byte strings would result in the same AAD used as input to the cipher:

`aabbcc, dd, eeff`

`aabb, ccddee, ff`

This might cause a data recipient to be confused about the internal structure of protected data. A common solution is to prepend the length of each string to the string itself: `03aabbcc01dd02eeff`. However, we would like to avoid a custom encoding scheme.

Therefore, when the application requires support for multiple AAD values, it MUST use the Structured Additional Authenticated Data field type, whose content consists of a well-formed CBOR structure, such as an array of byte strings. The cipher algorithm takes the entire structure (including the initial array marker, if any) as its AAD input.

The internal structure of this field type (e.g. CBOR array or map) is implementation dependent.

## Receiving Ciphertext

Correct interpretation of the format may have security implications, making it important to define the exact semantics even when the entity that receives a ciphertext may not understand parts of the header.

*   A recipient MUST reject a malformed header, e.g. if the total length is larger than the physical length allocated to it based on higher-level network protocols or storage formats.
*   A recipient MUST reject a ciphertext if it does not recognize the format version.
*   A recipient MUST reject a ciphertext if the variable header is not valid CBOR, as per {{RFC8949}} Sec. 5.3.1. In particular, it MUST reject duplicate map keys.
*   A recipient MUST accept a ciphertext even if it does not recognize some of the map keys. It MUST ignore the unknown map keys and MUST interpret all known ones. In other words, the only way to introduce new mandatory map keys is by incrementing the format version.
*   If ciphertext integrity protection coverage includes the header, a recipient MUST reject the header as well as the ciphertext if the integrity protection fails to validate.

## Fixed Header Rationale
{: #fixed-header-rationale}

We chose the initial byte 0x08, since strings are very unlikely to start with it, as we explain below. Automated tools can detect encrypted data in structured contexts (e.g., a SQL database column) by sampling a number of data items and if all start with this byte, determining that they are encrypted with a high probability.

The byte 0x08 encodes the ASCII control character “backspace”. It has the same meaning in UTF-8, and the 08 block of UTF-16 characters is only populated by two very small languages and rarely-used extended [Arabic characters](https://en.wikipedia.org/wiki/Arabic_Extended-A).

# Example

## Fixed Header

``08 01``

## Variable Header: CBOR Diagnostic Notation

```
{1: 65535, 2: h'1122334455', 3: 6, }
```

## Variable Header: Binary

```
a3 01 19 ff ff 02 45 11 22 33 44 55 03 06
```

## Complete Header

```
08 01 a3 01 19 ff ff 02 45 11 22 33 44 55 03 06
```

## CDDL

The following non-normative snippet defines the format of the variable header using CDDL {{?RFC8610}}.

~~~ 
{::include cddl/var-header.cddl}
~~~

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

## draft-sheffer-ietf-ciphertext-format-02

* Addressed secdispatch feedback from Russ Housley and Carsten Bormann.

## draft-sheffer-ietf-ciphertext-format-01

* SAAG feedback: the variable header is now CBOR.
* Binary example.
* Non-normative CDDL.
* Additional types for non-inline AEAD.

## draft-sheffer-ietf-ciphertext-format-00

* Initial version.