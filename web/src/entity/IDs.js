export function IDs(ids = []) {
  return {
    toJSON() {
      return ids
    },
  }
}
