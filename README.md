<img src="https://repository-images.githubusercontent.com/393375822/f196bfe1-d84f-4ede-9df9-1d33dc095b02" />

[www.betterinterview.club](https://www.betterinterview.club) makes it easy to manage interviews.

## Workflow

1. I create a new organization for my company by signing up with `justin@betterinterviews.com`.
1. I invite the first hiring manager: `max@betterinterviews.com`.
1. Max receives a registration link in her email. Clicking it brings her to the organization overview page.
1. Max is hiring for a `Staff Software Engineer` role at our company. She clicks `Open a new role`.
1. She picks 5 `yes` or `no` questions that she wants her interviewers to answer after interviewing the candidate.
1. She types in the emmails of the two interviewers to request their feedback: `kit@betterinterviews.com` and `thea@betterinterviews.com`.
1. Kit and Thea receive an email requesting their feedback. Clicking the link brings them to the `Give Feedback` page.
1. On the `Give Feedback` page, they simply select `Yes` or `No` for each question.
1. After feedback is received, Max, the person who opened the role, receives an email linking her to the `Feedback Given` page, where she can see the feedback.
1. Once feedback is complete, Max closes the role.

## TODO

### Features

* [x] Send email after invite.
* [x] Send email to request feedback.
* [x] Log in from home page.
* [x] Log out.
* [x] Send email to creator after feedback received.
* [x] Close a role once feedback is complete.
* [x] Send all emails async (don't wait to render page)
* [x] Submit yes/no to recommend candidate.
* [x] Open role and request feedback separately.
* [x] Section for managers, "My Open Roles"
* [x] Section for interviewers, "My Requested Feedback"
* [x] Explain feedback answers
* [ ] Section for interviewers, "My Given Feedback"
* [ ] Section for managers, "My Closed Roles"
* [ ] Section for Admins, "All Open Roles"

### Marketing

* [x] Add explanation to landing page
  * What problem am I solving?
  * Why would someone want to use this tool?
  * Why only yes/no questions?
  * Why the additional recommend yes/no?
* [x] Favicon

### Technical

* [x] Gzip responses
* [ ] Use correct time zones instead of UTC.
* [ ] Minify responses
* [x] Persist data locally
* [x] Persist data in production
* [x] HTML Lang attribute for accessibility
* [x] Content Security Policy
* [ ] Reusable components

### Security

* [x] Block common domains (gmail, yahoo, etc.)
* [x] Can't send emails to other domains
* [ ] Rate limit email sends
* [ ] Do not allow double feedback responses.

## Live Screenshot

![Screenshot](https://image.thum.io/get/maxAge/12/width/1000/https://www.betterinterview.club "Screenshot")

