import { entity } from "@justindfuller/entity";

const defaults = entity.defaults({
  StartTime: entity.Time(),
  EndTime: entity.Time(),
  InterviewTypeID: entity.ID(),
  ScheduleID: entity.ID(),
  InterviewerIDs: entity.IDs(),
  FeedbackQuestionIDs: entity.IDs(),
});

export function New(input) {
  const data = defaults(input);

  return entity.New(data, New);
}
