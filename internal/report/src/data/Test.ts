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
  tests: Array<TestResult>
}

export interface TestCapture {
  tests: Array<TestResult>
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

export function readReportData(cb: (text: string) => void) {
  function decompress(byteArray: ArrayBuffer, encoding: CompressionFormat) {
    const dcs = new DecompressionStream(encoding)
    const writer = dcs.writable.getWriter()
    writer.write(byteArray)
    writer.close()
    return new Response(dcs.readable).arrayBuffer().then(function (ab) {
      return new TextDecoder().decode(ab)
    })
  }

  const raw = document.getElementById("report-data")!.innerText
  fetch(new URL(raw))
    .then(r => r.arrayBuffer())
    .then(ab => decompress(ab, "gzip"))
    .then(text => cb(text))
    .catch(e => console.error(e))
}
