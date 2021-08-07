import { entity } from 'entity';

const defaults = entity.defaults({
  Question: "",
  InterviewTypeID: entity.ID(),
})

export function New(input) {
  const data = defaults(input)

  return entity.New(data, New)
}

