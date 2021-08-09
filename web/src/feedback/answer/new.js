import { entity } from "@justindfuller/entity";

const defaults = entity.defaults({
  Answer: entity.Boolean(),
  Notes: entity.String(),
  InterviewID: entity.ID(),
  FeedbackQuestionID: entity.ID(),
  InterviewerID: entity.ID(),
});

export function New(input) {
  const data = defaults(input);

  return entity.New(data, New);
}
