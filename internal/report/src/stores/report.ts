import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import {
  TestCapture as TestCaptureData,
  TestResult as TestResultData,
  readReportData
} from '../data/Test'

export interface TestResult {
  data: TestResultData
  tests: Array<TestResult>
  collapsed: boolean
}

export interface TestCapture {
  data: TestCaptureData
  tests: Array<TestResult>
}

export const useReportStore = defineStore('report', () => {
  const isLoading = ref(true)
  const testCapture = ref<TestCapture | null>(null)
  const filter = reactive<{
    testName: string
    showPassed: boolean
    showFailed: boolean
    showSkipped: boolean
  }>({
    testName: "",
    showPassed: true,
    showFailed: true,
    showSkipped: true
  })

  function wrapTestResults(parentTestResultData: TestResultData | null, testResultsData: Array<TestResultData>): Array<TestResult> {
    const tests: Array<TestResult> = []
    for (const testResultData of testResultsData) {
      const subTestCount = testResultData.tests.length
      const testResult: TestResult = {
        data: testResultData,
        tests: [],
        collapsed: true
      }
      if (parentTestResultData == null || (parentTestResultData.tests.length == 1 && subTestCount > 0)) {
        testResult.collapsed = false
      }
      if (!testResultData.skipped && !testResultData.passed && subTestCount > 0) {
        testResult.collapsed = false
      }
      testResult.tests = wrapTestResults(testResultData, testResultData.tests)
      tests.push(testResult)
    }
    return tests
  }

  function loadReport() {
    isLoading.value = true
    const startTime = Date.now()
    readReportData(text => {
      let testCaptureData: TestCaptureData = JSON.parse(text)
      if (testCaptureData) {
        testCapture.value = {
          data: testCaptureData,
          tests: wrapTestResults(null, testCaptureData.tests)
            .sort((t1, t2) => t1.data.index - t2.data.index)
        }
        isLoading.value = false
        console.log(`Report load took: ${Date.now() - startTime}ms`)
      }
    })
  }

  loadReport()
  return { isLoading, testCapture, filter }
})
