export interface TestOutput {
  index: number
  text: string
}

export interface TestResult {
  index: number
  started_at: string | null
  ended_at: string | null
  full_name: string
  package: string
  output: Array<TestOutput>
  done: boolean
  passed: boolean
  elapsed: number
  tests: Record<string, TestResult>;
  // Capture timestamps are only set for the root test
  capture_started_at: string | undefined
  capture_ended_at: string | undefined
}

export function testName(t: TestResult): string {
  const parts = t.full_name.split("/")
  return parts[parts.length - 1]
}
