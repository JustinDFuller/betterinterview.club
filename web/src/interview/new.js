import { entity } from 'entity';
import { time } from 'time';

const defaults = entity.defaults({
  StartTime: time.New(),
  EndTime: time.New(),
  InterviewTypeID: entity.ID(),
  ScheduleID: entity.ID(),
  InterviewerIDs: entity.IDs(),
  FeedbackQuestionIDs: entity.IDs(),
})

export function New(input) {
  const data = defaults(input)

  return entity.New(data, New)
}

