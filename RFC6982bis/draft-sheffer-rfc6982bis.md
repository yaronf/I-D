---
title: "Improving Awareness of Running Code: The Implementation Status Section"
abbrev: Running Code
docname: draft-sheffer-rfc6982bis-latest
category: bcp
updates:
obsoletes: RFC6982

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
  RFC3933:
  RFC4794:
  Stats:
    author:
      name: "Jari Arkko"
    title: Distribution of Processing Times
    target: "http://www.arkko.com/tools/lifecycle/wgdistr.html"
    date: "2012-12"
  Tao:
    author:
      name: Paul Hoffman
    title: "The Tao of IETF: A Novice's Guide to the Internet Engineering Task Force"
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

   The process in this document is offered as an experiment.  Authors of
   Internet-Drafts are encouraged to consider using the process for
   their documents, and working groups are invited to think about
   applying the process to all of their protocol specifications.  The
   authors of this document intend to collate experiences with this
   experiment and to report them to the community.

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
   documents, resulting in them being processed more rapidly.  We
   recommend that the Implementation Status section should be removed
   from Internet-Drafts before they are published as RFCs.  As a result,
   we do not envisage changes to this section after approval of the
   document for publication, e.g., the RFC errata process does not
   apply.

   The process in this document is offered as an experiment (though not
   as an [RFC3933] experiment; see {{experiment}}).  Authors of Internet-
   Drafts are encouraged to consider using the process for their
   documents, and working groups are invited to think about applying the
   process to all of their protocol specifications.

   The scope of the intended experiment is all Internet-Drafts (I-Ds)
   that contain implementable specifications, whether produced within
   IETF working groups or outside working groups but intended for IETF
   consensus.  I-Ds published on the Independent Stream are explicitly
   out of scope.  It is expected that the greatest benefit in the
   experiment will be seen with Standards Track documents developed
   within working groups.

   The authors of this document intend to collate experiences with this
   experiment and to report them to the community.

# The "Implementation Status" Section

   Each Internet-Draft may contain a section entitled "Implementation
   Status".  This section, if it appears, should be located just before
   the "Security Considerations" section and contain, for each existing
   implementation, some or all of the following:

   -  The organization responsible for the implementation, if any.

   -  The implementation's name and/or a link to a web page describing
      the implementation.

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
      6982.  The description of implementations in this section is
      intended to assist the IETF in its decision processes in
      progressing drafts to RFCs.  Please note that the listing of any
      individual implementation here does not imply endorsement by the
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

   Authors are requested to add a note to the RFC Editor at the top of
   this section, advising the Editor to remove the entire section before
   publication, as well as the reference to RFC 6982.

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

# Process Experiment {#experiment}

   The current proposal is proposed as an experiment.  The inclusion of
   Implementation Status sections in Internet-Drafts is not mandatory,
   but the authors of this document wish to encourage authors of other
   Internet-Drafts to try out this simple mechanism to discover whether
   it is useful.  Working group chairs are invited to suggest this
   mechanism to document editors in their working groups, and to draw
   the attention of their working group participants to Implementation
   Status sections where they exist.

   Following a community discussion, it was concluded that [RFC3933] is
   not an appropriate framework for this experiment, primarily because
   no change is required to any existing process.

## Duration {#duration}

   Given the typical time to produce an RFC (see [Stats]), we propose a
   duration of 18 months for the experiment.  Thus, 18 months after the
   date of publication of this document as an RFC, the authors will
   report on the experiment as described in the next section.

   I-D authors are obviously free to include Implementation Status
   sections in their documents even after the experiment has concluded.

##  Summary Report

   The authors will summarize the results of the experiment at the end
   of the period assigned to the experiment (see {{duration}}).  If
   nothing happens (no I-Ds or only a handful include an Implementation
   Status section), an email to the IETF list will be sufficient.  This
   would obviously constitute a failure to adopt the idea and the
   authors will abandon the experiment.

   If this idea is adopted by document authors, a summary I-D will be
   written containing the statistics of such adoption, as well as
   (necessarily subjective) reports by working group members, chairs,
   and area directors who have used this mechanism.

   The authors may then propose more wide-scale use of the process and
   might suggest more formal adoption of the process by the IETF.

## Success Criteria

   The goal of this experiment is to improve the quality of IETF
   specifications.  This is impossible to quantify, of course.  We
   suggest that generally positive answers to the following questions
   would indicate that the experiment was successful:

   -  Did the working group make decisions that were more informed when
      comparing multiple competing solutions for the same work item?

   -  Did authors significantly modify proposed protocols based on
      implementation experience?

   -  Did disclosure of implementations encourage more interoperability
      testing than previously?

   -  Did non-authors review documents based on interactions with
      running code and/or inspection of the code itself?

# Security Considerations

   This is a process document; therefore, it does not have a direct
   effect on the security of any particular IETF protocol.  However,
   better-reviewed protocols are likely to also be more secure.

# Acknowledgements

   We would like to thank Stephen Farrell, who reawakened community
   interest in this topic.  Several reviewers provided important input,
   including Loa Andersson, Dave Crocker, Ned Freed, Christer Holmberg,
   Denis Ovsienko, and Curtis Villamizar.

   This document was originally prepared using the lyx2rfc tool, and we
   would like to thank Nico Williams, its author.

--- back

# Document History

## draft-sheffer-rfc6982bis-00

Initial version.
