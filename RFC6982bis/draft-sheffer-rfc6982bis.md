---
title: "Improving Awareness of Running Code: The Implementation Status Section"
abbrev: Running Code
docname: draft-sheffer-rfc6982bis-latest
category: bcp
updates:
obsoletes: 6982

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
    ins: A. Farrel
    name: Adrian Farrel
    organization: Juniper Networks
    email: adrian@olddog.co.uk
    

normative:

informative:
  RFC1264:
  RFC1369:
  RFC4794:
  RFC5080:
  RFC6982:
  Tao:
    author:
      name: Paul Hoffman
      ins: P. Hoffman
    title: >
      "The Tao of IETF: A Novice's Guide to the Internet Engineering Task Force"
    target: "http://www.ietf.org/tao.html"
    date: 2012

--- abstract

   This document describes a simple process that allows authors of
   Internet-Drafts to record the status of known implementations by
   including an Implementation Status section.  This will allow
   reviewers and working groups to assign due consideration to documents
   that have the benefit of running code, which may serve as evidence of
   valuable experimentation and feedback that have made the implemented
   protocols more mature.

   This process is not mandatory. Authors of
   Internet-Drafts are encouraged to consider using the process for
   their documents, and working groups are invited to think about
   applying the process to all of their protocol specifications.
   This document obsoletes RFC 6982, advancing it to a
   Best Current Practice. 

--- middle

# Introduction

   Most IETF participants are familiar with the saying "rough consensus
   and running code" [Tao] and can identify with its pragmatic approach.
   However, implementation is not a requirement for publication as an
   RFC.  There are many examples of Internet-Drafts containing protocol
   specification that have gone through to publication as Proposed
   Standard RFCs without implementation.  Some of them may never get
   implemented.

   Over time, a variety of policies have been applied within the IETF to
   consider running code.  In the Routing Area, it used to be a
   requirement that one or more implementations must exist before an
   Internet-Draft could be published as a Proposed Standard RFC
   [RFC1264].  That RFC was later obsoleted and the requirement for
   implementation was lifted, but each working group was given the
   authority to impose its own implementation requirements [RFC4794] and
   at least one working group, Inter-Domain Routing (IDR), continues to
   require two independent implementations.

   The hypothesis behind the current document is that there are benefits
   to the IETF standardization process of producing implementations of
   protocol specifications before publication as RFCs.  These benefits,
   which include determining that the specification is comprehensible
   and that there is sufficient interest to implement, are further
   discussed in {{benefits}}.

   This document describes a simple mechanism that allows authors of
   Internet-Drafts to record and publicize the status of known
   implementations by including an Implementation Status section.  The
   document defines (quite informally) the contents of this section to
   ensure that the relevant information is included.  This will allow
   reviewers and working groups to assign due consideration to documents
   that have the benefit of running code, which may serve as evidence of
   valuable experimentation and feedback that have made the implemented
   protocols more mature.

   It is up to the individual working groups to use this information as
   they see fit, but one result might be the preferential treatment of
   documents, resulting in them being processed more rapidly. We
   recommend that the Implementation Status section should be removed
   from Internet-Drafts before they are published as RFCs.  As a result,
   we do not envisage changes to this section after approval of the
   document for publication, while the document sits in the RFC-editor 
   queue, e.g., the RFC errata process does not apply.


This process is not mandatory.  Authors of Internet-Drafts
are encouraged to consider using the process for their
   documents, and working groups are invited to think about applying the
   process to all of their protocol specifications.

   The scope of this process is all Internet-Drafts (I-Ds)
   that contain implementable specifications, whether produced within
   IETF working groups or outside working groups but intended for IETF
   consensus.  I-Ds published on the Independent Stream are explicitly
   out of scope.  It is expected that the greatest benefit
   will be seen with Standards Track documents developed
   within working groups.

   This process was initially proposed as an experiment in
   [RFC6982]. That document is now obsoleted, and the
   process advanced to Best Current Practice.

Historically there have been other ways for experience based on protocol implementations
to feed back into the IETF process. Many "implementation reports" have been published,
in some cases several years after the protocol was
originally published. Providing feedback to published protocols is a related goal,
but different from the current document's
focus. Two notable examples of published implementation reports
are {{RFC1369}} and {{RFC5080}}.

# The "Implementation Status" Section

   Each Internet-Draft may contain a section entitled "Implementation
   Status".  This section, if it appears, should be located just before
   the "Security Considerations" section and contain, for each existing
   implementation, some or all of the following:

   -  The organization responsible for the implementation, if any.

   - The implementation's name and/or a link to a web page where the implementation
	 or a description of it can be found.

   -  A brief general description.

   -  The implementation's level of maturity: research, prototype,
      alpha, beta, production, widely used, etc.

   -  Coverage: which parts of the protocol specification are
      implemented and which versions of the Internet-Draft were
      implemented.

   -  Licensing: the terms under which the implementation can be used.
      For example: proprietary, royalty licensing, freely distributable
      with acknowledgement (BSD style), freely distributable with
      requirement to redistribute source (General Public License (GPL)
      style), and other (specify).

   -  Implementation experience: any useful information the implementers
	  want to share with the community.
   
   -  Contact information: ideally a person's name and email address,
	  but possibly just a URL or mailing list.

   - The date when information about this particular implementation was
	 last updated.

   In addition, this section can contain information about the
   interoperability of any or all of the implementations, including
   references to test-case descriptions and interoperability reports,
   when such exist.

   Working group chairs and area directors (ADs) are requested to ensure
   that this section is not used as a marketing venue for specific
   implementations.

   Since this information is necessarily time dependent, it is
   inappropriate for inclusion in a published RFC.  The authors should
   include a note to the RFC Editor requesting that the section be
   removed before publication.

## Introductory Text

   The following boilerplate text is proposed to head the Implementation
   Status section:

      This section records the status of known implementations of the
      protocol defined by this specification at the time of posting of
      this Internet-Draft, and is based on a proposal described in RFC
      [[RFC Editor: replace by RFC number]].
      The description of implementations in this section is
      intended to assist the IETF in its decision processes in
      progressing drafts to RFCs.  Please note that the listing of any
      individual implementation here does not imply endorsement by the
      IETF.  Furthermore, no effort has been spent to verify the
      information presented here that was supplied by IETF contributors.
      This is not intended as, and must not be construed to be, a
      catalog of available implementations or their features.  Readers
      are advised to note that other implementations may exist.

      According to RFC [[RFC Editor: replace by RFC number]],
      "this will allow reviewers and working
      groups to assign due consideration to documents that have the
      benefit of running code, which may serve as evidence of valuable
      experimentation and feedback that have made the implemented
      protocols more mature.  It is up to the individual working groups
      to use this information as they see fit".

   Authors are requested to add a note to the RFC Editor at the top of
   this section, advising the Editor to remove the entire section before
   publication, as well as the reference to
   RFC [[RFC Editor: replace by RFC number]].

# Alternative Formats

   Sometimes it can be advantageous to publish the implementation status
   separately from the base Internet-Draft, e.g., on the IETF wiki:

   -  When the Implementation Status section becomes too large to be
      conveniently managed within the document.

   -  When a working group decides to have implementors, rather than
      authors, keep the status of their implementations current.

   -  When a working group already maintains an active wiki and prefers
      to use it for this purpose.

   -  If the working group decides that the information is still
      valuable (and needs to be kept current) after the I-D is published
      as an RFC, and the Implementation Status section had been removed
      from it.

   It is highly desirable for all readers of the Internet-Draft to be
   made aware of this information.  Initially, this can be done by
   replacing the Implementation Status section's contents with a URL
   pointing to the wiki.  Later, the IETF Tools may support this
   functionality, e.g., by including such a link in the HTML file of the
   document, similar to the IPR link.

   If the implementation status is published separately from the I-D,
   then this information needs to be openly available without requiring
   authentication, registration, or access controls if it is to have any
   useful effects.

# Benefits {#benefits}

   Publishing the information about implementations provides the working
   group with several benefits:

   -  Working group members, chairs, and ADs may use the information
      provided to help prioritize the progress of I-Ds, e.g., when there
      are several competing proposals to solve a particular problem.

   -  Similarly, the information is useful when deciding whether the
      document should be progressed on a different track (individual
      submission, Experimental, etc.).

   -  Making this information public and an explicit part of WG
      deliberations will motivate participants to implement protocol
      proposals, which in turn helps in discovering protocol flaws at an
      early stage.

   -  Other participants can use the software to evaluate the usefulness
      of protocol features, its correctness (to some degree), and other
      properties, such as resilience and scalability.

   -  WG members may choose to perform interoperability testing with
      known implementations, especially when they are publicly
      available.

   -  In the case of open source, people may want to study the code to
      better understand the protocol and its limitations, determine if
      the implementation matches the protocol specification, and whether
      the protocol specification has omissions or ambiguities.

   -  And lastly, some protocol features may be hard to understand, and
      for such features, the mere assurance that they can be implemented
      is beneficial.  We note though that code should never be used in
      lieu of a clear specification.

   We do not specify here whether and to what degree working groups are
   expected to prefer proposals that have "running code" associated with
   them, over others that do not.

   Working group chairs are invited to suggest this
   mechanism to document editors in their working groups, and to draw
   the attention of their working group participants to Implementation
   Status sections where they exist.

# Security Considerations

   This is a process document; therefore, it does not have a direct
   effect on the security of any particular IETF protocol.  However,
   better-reviewed protocols are likely to also be more secure.

# Acknowledgements

   We would like to thank Stephen Farrell, who reawakened community
   interest in this topic.  Several reviewers provided important input,
   including Loa Andersson, Dave Crocker, Ned Freed, Joel M. Halpern,
   Christer Holmberg,
   Denis Ovsienko, and Curtis Villamizar.

   This document was originally prepared using the lyx2rfc tool, and we
   would like to thank Nico Williams, its author.

--- back

# Document History

## draft-sheffer-rfc6982bis-02

* Implemented IESG review comments.

## draft-sheffer-rfc6982bis-01

* Removed the IANA Considerations section.
* Implemented GENART comments.

## draft-sheffer-rfc6982bis-00

Initial version. RFC 6982 as-is, without the experiment sections.
