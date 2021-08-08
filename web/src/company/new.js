import { entity } from '@justindfuller/entity';

const defaults = entity.defaults({
    Name: "",
    EmailDomain: "",
  })

export function New(input) {
  const data = defaults(input)

  return entity.New(data, New)
}

