import { entity } from "@justindfuller/entity";

const defaults = entity.defaults({
  Title: entity.String(),
  Description: entity.String(),
  ParentTeamID: entity.ID(),
  CompanyID: entity.ID(),
});

export function New(input) {
  const data = defaults(input);

  return entity.New(data, New);
}
