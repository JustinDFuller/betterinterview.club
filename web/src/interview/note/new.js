import { entity } from "@justindfuller/entity";

const defaults = entity.defaults({
  Text: entity.String(),
  InterviewID: entity.ID(),
  InterviewerID: entity.ID(),
});

export function New(input) {
  const data = defaults(input);

  return entity.New(data, New);
}
