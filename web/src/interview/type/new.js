import { entity } from "@justindfuller/entity";

const defaults = entity.defaults({
  Name: entity.String(),
  Description: entity.String(),
  CompanyID: entity.ID(),
});

export function New(input) {
  const data = defaults(input);

  return entity.New(data, New);
}
