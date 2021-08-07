export function New(time) {
  return {
    toJSON() {
      return time
    },
    toString() {
      return time.toString()
    },
  }
}
