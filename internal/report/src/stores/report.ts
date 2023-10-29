import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  TestCapture as TestCaptureData,
  TestResult as TestResultData,
  readReportData
} from '../data/Test'

export interface TestResult {
  data: TestResultData
  tests: Record<string, TestResult>;
  collapsed: boolean
}

export interface TestCapture {
  data: TestCaptureData
  tests: Record<string, TestResult>;
}

export const useReportStore = defineStore('report', () => {
  const isLoading = ref(true)
  const testCapture = ref<TestCapture | null>(null)

  function wrapTestResults(parentTestResultData: TestResultData | null, testResultsData: Record<string, TestResultData>): Record<string, TestResult> {
    const tests: Record<string, TestResult> = {}
    for (const [testName, testResultData] of Object.entries(testResultsData)) {
      const subTestCount = Object.keys(testResultData.tests).length
      const testResult: TestResult = {
        data: testResultData,
        tests: {},
        collapsed: true
      }
      if (parentTestResultData == null || (Object.keys(parentTestResultData.tests).length == 1 && subTestCount > 0)) {
        testResult.collapsed = false
      }
      if (!testResultData.skipped && !testResultData.passed && subTestCount > 0) {
        testResult.collapsed = false
      }
      testResult.tests = wrapTestResults(testResultData, testResultData.tests)
      tests[testName] = testResult
    }
    return tests
  }

  function loadReport() {
    isLoading.value = true
    readReportData(text => {
      let testCaptureData: TestCaptureData = JSON.parse(text)
      if (testCaptureData) {
        testCapture.value = {
          data: testCaptureData,
          tests: wrapTestResults(null, testCaptureData.tests)
        }
        isLoading.value = false
      }
    })
  }

  loadReport()
  return { isLoading, testCapture }
})
