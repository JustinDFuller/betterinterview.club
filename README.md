<img src="https://repository-images.githubusercontent.com/393375822/08a8ec9f-798b-4924-8ba6-1b666ac1de34" />

Let's make it easy to manage interviews.

## How it works

- A `recruiter` creates an `interview type`. (ex. Software Engineer Technical Phone Screen)
- A `recruiter` creates yes or no `feedback questions` for that interview type.

- A `manager` creates their `team`.
- A `manager` opens a `role` for their `team`.
- A `recruiter` creates an interview `schedule` for a `candidate`. 
- A `recruiter` adds one or more `interview`s to a `schedule`.
- A `recruiter` adds one or more `interviewer`s to an `interview`.
- The `candidate` can review their `schedule`, `role`s and `team`s they are interviewing for.
- During the `interview`, `interviewer`s can add `notes` about the `interview`.
- After the `interview`, an `interviewer` can answer one or more `feedback question`s.

## Data Types

```
All types have an UUIDv4 ID. All types have CreatedAt and UpdatedAt timestamps.

Person:
  - Name
  - Email
  - Type (manager, recruiter, candidate, interviewer)
  
Candidate: Person
  - Status (interviewing, declined, rejected, offered, hired)
  
Interview Type:
  - Name
  - Description
  
Feedback Questions:
  - Interview Type ID
  - Question
  
Team:
  - Title
  - Description
  - Parent Team ID
  
Role:
  - Title
  - Description
  - Team ID
  - Status (Open, Filled, Closed)
  
Link:
  - URL
  - Description (google hangout, coderpad)

Interview:
  - StartTime
  - EndTime
  - Interview Type ID
  - Array of Interviewer IDs
  - Array of Feedback Question IDs
  - Links IDs

Schedule:
  - Candidate ID
  - Array of Role IDs
  - Array of Interview IDs
  
Notes:
  - Text
  - Interview ID
  - Interviewer ID

Feedback Answer:
  - Answer (true, false)
  - Notes
  - Interview ID
  - Feedback Question ID
  - Interviewer ID
```

## Views

### Recruiter

- Interview Types
    - List (a table that links to edit and delete, links to feedback questions for this interview type)
    - Edit
    - Create

- Feedback Questions
    - List (a table that links to edit and delete)
    - Edit
    - Create

- Candidates
    - List (a table that links to edit and delete, link to view candidate schedule)
    - Edit
    - Create

- Schedules
    - List (by all, by candidate ID, or by role ID) (links to edit and delete)
    - Edit (should it be able to add interviews in-line? probably)
    - Create

- Interviews
  - List by schedule ID (links to edit and delete)
  - Create interview for schedule ID
  - Edit interview

- Interviewers
  - List (links to edit and delete)
  - Create (should this link to interviewers as well? not so sure)
  - Edit

### Manager

- Teams
    - List my teams (links to edit and delete, links to roles for this team)
    - Edit
    - Create

- Role
    - List my roles (or roles for a given team) (links to edit and delete, links to interviews for this role)
    - Edit
    - Create

### Candidate

- Schedule
    - See schedule in timeline form
        - Shows interview type, description, times, links
        - Maybe shows interviewers
        - Link to learn more about team(s), maybe show inline

    - See my status (open, pending, etc.) Use nice langugage for this, not just "REJECTED".

### Interviewer

- Interviews
    - List (show only future interviews by default, add view to see previous interviews)
        - Links to take notes
        - Links to answer feedback questions

## Components

There are some clear patterns that can become components.

1. Many of the datatypes have a LIST, a CREATE, and an EDIT view.
2. The LIST view is a table.
3. The CREATE and EDIT views are forms.
4. The DELETE action is a button.
5. It may also make sense to have a VIEW but I am not sure yet. LIST may suffice.
6. The table will contain links to other data types. There can be a common `link` that links an ID to a URL and displays a TITLE.
7. The interviewer should see a really nice vertical timeline component.
8. The interviews list may make more sense as a set of cards than a table. It will give more details that way.

## Client Considerations

- All logged-in data, no huge SEO concerns here.
    - SSR may not be necessary
- Will likely need to be themed based on Company's theme (white label) 
- Each company may want to have a landing page to onboard new candidates

## Backend Considerations

- Backend will mostly be a CRUD system
- The backend will need to store and validate updates to this data.
    - Will need to make sure candidates can't see all open roles or other people's schedules
- Will need to support SSO for each company
    - Not sure if Magic.link will be compliant? 
    - Google SSO is another option that will probably hit a lot of companies.
