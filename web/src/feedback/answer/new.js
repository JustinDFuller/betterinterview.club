import { entity } from "@justindfuller/entity";

const defaults = entity.defaults({
  Answer: false,
  Notes: "",
  InterviewID: entity.ID(),
  FeedbackQuestionID: entity.ID(),
  InterviewerID: entity.ID(),
});

export function New(input) {
  const data = defaults(input);

  return entity.New(data, New);
}
