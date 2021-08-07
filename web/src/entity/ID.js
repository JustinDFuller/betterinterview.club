export function ID(id = "") {
  return {
    toString() {
      return id;
    },
    toJSON() {
      return id;
    },
    validate() {
      return id !== "";
    },
  }
}
