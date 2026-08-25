# Alice Suite

Alice Suite is a companion to a physical classic-literature book. It gives a Reader immediate word help, contextual AI assistance, and access to a human Consultant.

## Reading assistance

**Reading Context**:
The specific part of the physical-book reading experience that an AI request or Help Request concerns: a section, a page, or a reader-selected passage. It includes the source text and its book/page/section location.
_Avoid_: Context, prompt context, page state

**Selected Passage**:
The exact contiguous text a Reader deliberately selects from the rendered book page. It is a Reading Context with selection scope, not merely transient browser selection.
_Avoid_: Selected text, browser selection, highlight

**AI Assistance Request**:
A Reader's requested explanation, simplification, observation, or question paired with one Reading Context. It is distinct from the stored AI interaction, which is the record of the completed request and response.
_Avoid_: AI chat, AI interaction, prompt

**Help Request**:
A Reader's request for a human Consultant, paired with the relevant Reading Context. It has a lifecycle of pending, assigned, then resolved.
_Avoid_: Support ticket, consultant message

**Consultant Assignment**:
The active relationship that gives one Consultant responsibility for a Reader and book. It is distinct from assignment of a specific Help Request.
_Avoid_: Help-request assignment, primary consultant
