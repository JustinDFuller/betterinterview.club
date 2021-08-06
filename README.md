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

Interview:
  - StartTime
  - EndTime
  - Interview Type ID
  - Array of Interviewer IDs
  - Array of Feedback Question IDs

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
