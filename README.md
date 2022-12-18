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
* [ ] Ensure one data query / lock per endpoint
* [ ] Load templates on application startup

### Security

* [x] Block common domains (gmail, yahoo, etc.)
* [x] Can't send emails to other domains
* [ ] Rate limit email sends
* [ ] Do not allow double feedback responses.

## Live Screenshot

![Screenshot](https://image.thum.io/get/maxAge/12/width/1000/https://www.betterinterview.club "Screenshot")

## Workflow

### 1. Sign Up with your email

![betterinterviews1](https://user-images.githubusercontent.com/11914897/208248226-0369c8f7-cd2e-4b5f-a97f-e3b59ef2ba6a.png)

### 2 Head over to your email
![betterinterviews2](https://user-images.githubusercontent.com/11914897/208248225-1214a7fe-2863-4df2-b48b-b1a9a4350809.png)

### 3. Click the login link
![betterinterviews3](https://user-images.githubusercontent.com/11914897/208248223-7810cced-127e-4c6c-8691-5eaa198e5c14.png)

### 4. Open a new role
![betterinterviews4](https://user-images.githubusercontent.com/11914897/208248222-2ffd4b88-72ae-44fa-be64-b2395bc178d5.png)

### 5. Choose your interview signals
![betterinterviews5](https://user-images.githubusercontent.com/11914897/208248221-f35d4494-6616-40fd-a763-a929f90349a8.png)

### 6. Request Feedback
![betterinterviews6](https://user-images.githubusercontent.com/11914897/208248220-1091b9c8-014e-44a5-bf2b-38b55c9b6d2c.png)

### 7. Enter the candidate name and interviewer emails
![betterinterviews7](https://user-images.githubusercontent.com/11914897/208248219-3525ec1d-d4c1-4308-b5ac-a8c9a8cfd5c3.png)

### 8. Interviewers receive feedback requests
![betterinterviews8](https://user-images.githubusercontent.com/11914897/208248218-edad11f9-f614-4c9a-815c-070edcc922d8.png)

### 9. Interviewers answer questions
![betterinterviews9](https://user-images.githubusercontent.com/11914897/208248217-1b2e5e6f-1ffc-4a4a-8418-374d864f0b3e.png)

### 10. Hiring manager receives feedback via email
![betterinterviews10](https://user-images.githubusercontent.com/11914897/208248216-5e21a88b-0df8-499b-97e9-1a12cd042a19.png)

### 11. Feedback is also viewable on the website
![betterinterviews11](https://user-images.githubusercontent.com/11914897/208248213-b52cf95a-55fa-4973-bb44-c1abc40d4032.png)

### 12. Close the role once you make a hire!
![betterinterviews12](https://user-images.githubusercontent.com/11914897/208248211-345df5f1-ed41-4aeb-9cb2-9238b905d1b1.png)
