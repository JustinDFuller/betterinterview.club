import { v4 } from "uuid";

export function defaults(defaultObj) {
  return function (data) {
    return Object.freeze(Object.assign({ ID: v4() }, defaultObj, data));
  };
}
