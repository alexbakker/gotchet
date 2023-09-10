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
  skipped: boolean
  passed: boolean
  elapsed: number
  tests: Record<string, TestResult>;
}

export interface TestCapture {
  tests: Record<string, TestResult>;
  title: string
  started_at: string | null
  ended_at: string | null
  capture_started_at: string
  capture_ended_at: string
}

export function testName(t: TestResult): string {
  if (t.full_name == "") {
    return t.package
  }
  const parts = t.full_name.split("/")
  return parts[parts.length - 1]
}
