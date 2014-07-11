namespace go crystal.thrift

typedef binary Id

enum Error {
  NONE = 0,
  CLOCK_RUNNING_BACKWARDS = 1,
}

struct IdGenerationResult {
  1: Error err
  2: Id id
}

service IdGenerationService {
  IdGenerationResult generate()
}
