export interface TestOutput {
  index: number
  text: string
}

export interface Test {
  index: number
  full_name: string
  package: string
  output: Array<TestOutput>
  done: boolean
  passed: boolean
  elapsed: number
  tests: Record<string, Test>;
}

export function testName(t: Test): string {
  const parts = t.full_name.split("/")
  return parts[parts.length - 1]
}
