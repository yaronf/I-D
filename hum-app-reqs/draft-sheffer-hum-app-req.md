---
title: "Virtual Hum Application: Requirements"
abbrev: Hum Application Requirements
docname: draft-sheffer-hum-app-req-latest
category: info

ipr: trust200902
area: Security
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

informative:
  Coldewey:
      -: ta
      target: "https://techcrunch.com/2018/06/11/the-messy-musical-process-behind-the-webs-new-security-standard/"
      title: "The messy, musical process behind the web’s new security standard"
      date: June 2018
      author:
        -
          name: Devin Coldewey


--- abstract

The IETF has been having virtual meetings for a long time. Until recently these have been interim meetings,
and following the all-virtual IETF-107 we expect to see more and more WG meetings take the virtual route. 
A common practice at the IETF is to use room "humming" to gauge working group consensus,
though the final consensus is determined by the working group chairs and typically
requires a mailing list poll as well. We do not
have a technique equivalent to the hum for virtual meetings, and arguably this
reduces the effectiveness of these meetings.

This document defines the requirements from a web application whose
goal is to most faithfully replicate the "feel" of hums, through the use
of a traditional web user interface.


--- middle

# Introduction

The IETF has been having virtual meetings for a long time. Until recently these have been interim meetings,
and following the all-virtual IETF-107 we expect to see more and more WG meetings take the virtual route. 
A common practice at the IETF is to use room "humming" to gauge working group consensus,
though the final consensus is determined by the working group chairs and typically
requires a mailing list poll as well. We do not
have a technique equivalent to the hum for virtual meetings, and arguably, this
reduces the effectiveness of these meetings.

This document defines the requirements from a web application whose
goal is to most faithfully replicate the "feel" of hums, through the use
of a traditional web user interface.

The document's scope is strictly on the web application, and not
on the process implications of humming or of replacing it by
a different (though hopefully similar) human protocol.
Please refer to {{?RFC7282}}
for the authoritative discussion of what IETF consensus means, how
it can be reached, and the role - as well as the limitations - of humming
in achieving consensus.

## Conventions used in this document

{::boilerplate bcp14}

# Background

Note: the intended audience for this section is application developers who are not familiar
with the IETF process.

IETF, the Internet Engineering Task Force, is the most important standards body for the Internet.
Its main product is RFC documents that define protocols. For example, the IP protocol is defined by RFC 791,
the HTTP protocol is defined by a series of RFCs, TLS 1.3 is defined by RFC 8446.
The IETF has a very long history and very detailed processes associated with its operation.
It has been holding 3 annual face-to-face meetings for a very long time,
and is only now moving more fully into virtual meetings. In fact the first fully virtual IETF meeting
is the upcoming IETF 107, taking place next week.
The IETF consists of dozens of working groups, and they come to decisions using a process
called “rough consensus” which means that most participants are in favor of a certain decision
and there is no large faction against or an even smaller faction but with strongly held opinions.
Quoting “the Tao or the IETF”:

4.2 Getting Things Done in a Working Group

One fact that confuses many newcomers is that the face-to-face WG meetings are much less
important in the IETF than they are in most other organizations. Any decision made at
a face-to-face meeting must also gain consensus on the WG mailing list.

There are numerous examples of important decisions made in WG meetings that are later overturned
on the mailing list, often because someone who couldn't attend the meeting pointed out a serious
flaw in the logic used to come to the decision. Finally, WG meetings aren't "drafting sessions",
as they are in some other standards bodies: in the IETF, drafting is done elsewhere.

Another aspect of Working Groups that confounds many people is the fact that there is no formal voting.
The general rule on disputed topics is that the Working Group has to come to "rough consensus",
meaning that a very large majority of those who care must agree. The exact method of determining
rough consensus varies from Working Group to Working Group. Sometimes consensus is determined by "humming"
— if you agree with a proposal, you hum when prompted by the chair. Most "hum" questions come in two parts:
you hum to the first part if you agree with the proposal, or you hum to the second part if you
disagree with the proposal. Newcomers find it quite peculiar, but it works.
It is up to the chair to decide when the Working Group has reached rough consensus.

The lack of formal voting has caused some very long delays for some proposals, but most IETF participants
who have witnessed rough consensus after acrimonious debates feel that the delays often result
in better protocols. (And, if you think about it, how could you have "voting" in a group that
invites all interested individuals to participate, and when it's impossible to count the participants?)
Rough consensus has been defined in many ways; a simple version is that it means that strongly
held objections must be debated until most people are satisfied that these objections are wrong.

See also this article [Coldewey] for another view on the humming practice.

With the move to virtual meetings, real audio-based humming is no longer an option.
We would like to develop a replacement that preserves as much as possible of the spirit
behind this practice but is workable for widely distributed virtual working group meetings.

# General Requirements

This is a relatively simple web application. It needs to be usable by people who are seeing it for the first time (meeting participants) or people who have had very minimal practice and need to operate it under pressure (working group chairs).

* Administration requires a desktop browser.
* Nice to have: participation from a mobile browser.
* Nice to have: OpenID authentication for admins.

# Hum Rooms
Anybody can open a “hum room”. The room is available for a period of time (default: 6 hours) and then
it is archived. The room consists of:

* A name, defined by the room admin.
* A secret, random management URL, which may be shared with 2-3 additional admins, and visible to the IETF Secretariat.
* A secret, random participation URL, which will be shared with all participants (should allow up to 500). After the room is archived, the archive view will be returned, as the URL will wind up in minutes.
* Configuration: the expected number of participants. Entered by admins, visible to all participants. (Note: this value may not be needed, this depends on the exact rules we define for gauging consensus).
* A list of questions, see below.
Some analytics, including the total number of participants seen, the total number of hums taken, number of unique IPs etc. Analytics are open to all participants.
* Expiry time of the room.
* A “get summary” button that enables downloading (e.g. as JSON) a summary of all analytics, questions and results, so they can be used in the meeting minutes. This button is available to everybody.
* A way to close the room even before it has expired. Available only to admins.

# Hum Questions
Questions are typically entered on-the-fly by the admin, during the meeting. So the interface must be very minimal/simple.

* Introductory text (up to 2-3 lines of text).
* Between 2-4 options, with text and a button for each. “I don’t understand the question” shall be suggested as one of the options, however it should be possible to delete it.
* A link or popup with the detailed rules for consensus for this question, visible to all participants.

For example:

~~~
Should we require encryption of all HTTP traffic, as a MUST?

Yes [button]

No [button]

Don’t have enough information [button]  (This is not the same as “I don’t understand the question”)
~~~

* Questions are visible to all participants as soon as they are entered (even keystroke by keystroke), similarly to Etherpad/hack-MD.
* Buttons become available only when the admin enables the question.
* Buttons are available for a short duration, e.g. 3 minutes
* Buttons are treated as toggles, i.e. a second press disables the selection.
* People are allowed to press more than one button (this is weird but it replicates the hum experience).
* All participants see an indicator (e.g. progress bar) of how much time remains.

# Gauging Consensus
When the time expires for a question, each option is evaluated separately:
* Zero responses: silence.
* Less than 20% of the total number of people who responded: weak hum.
* 20-70%: medium hum.
* 70-100%: strong hum.

Only the summary (e.g. “medium hum”) is displayed/retained, not the exact numbers. In addition, we display the total number of responses.

Admins (working group chairs) are expected to announce the results to the protocol, and decide whether consensus has been reached, before moving on to the next question.

# Graphics/UI

Please include the IETF logo, https://www.ietf.org/logo/.

# Transport Security

* HTTPS, and redirection from HTTP to HTTPS.
* Please use Let’s Encrypt for certificates. You should probably use the Certbot client.
* The server should have scheduled code that fetches a new certificate automatically 2 weeks before the cert expires.
* To authorize the server to Let’s Encrypt, use the HTTP-01 challenge.

# Security, Access Control, Fraud
Basically none. Counting unique IPs allows for detection of simple-minded fraud.


# IANA Considerations

This is a process document, with no IANA implications.

# Security Considerations

The process described here may have operational security considerations
related to the IETF process, but none that apply directly to any IETF deliverables.


# Privacy Considerations

IETF processes are not expected to ensure anonymity of participants.
The process described
here does not make any changes to the existing privacy guarantees.

# Acknowledgements

I would like to thank Michael Richardson for his extensive comments to a previous version of this document.


--- back

# Document History

[[Note to RFC Editor: please remove before publication.]]

## draft-sheffer-hum-app-req-00

Initial version.
