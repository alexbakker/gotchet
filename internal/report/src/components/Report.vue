<script setup lang="ts">
  import { reactive, computed } from 'vue'
  import Test from './Test.vue'
  import Elapsed from './Elapsed.vue'
  import { testName, readReportData } from '../data/Test'
  import { TestResult, useReportStore } from '../stores/report.ts'

  const store = useReportStore()

  const state = reactive<{
    filter: {
      testName: string
      showPassed: boolean
      showFailed: boolean
      showSkipped: boolean
    }
  }>({
    filter: {
      testName: "",
      showPassed: true,
      showFailed: true,
      showSkipped: true
    }
  })

  const stats = computed(() => {
    const shownTests = tests.value
    return {
      total: shownTests.length,
      passed: shownTests.filter((t) => t.data.done && t.data.passed).length,
      failed: shownTests.filter((t) => t.data.done && !t.data.passed).length
    }
  })

  const tests = computed(() => {
    if (!store.testCapture) {
      return []
    }

    return Object.values(store.testCapture.tests)
      .filter((t) => isTestShown(t))
      .sort((t1, t2) => t1.data.index - t2.data.index)
  })

  const totalElapsed = computed(() => {
    if (!store.testCapture) {
      return 0
    }

    return Object.values(store.testCapture.tests)
      .reduce((sum, t) => sum + t.data.elapsed, 0);
  })

  function isTestShown(t: TestResult): boolean {
    if (t.data.done) {
      if (t.data.passed && !state.filter.showPassed) {
        return false
      }

      if (!t.data.passed && !state.filter.showFailed) {
        return false
      }

      if (t.data.skipped && !state.filter.showSkipped) {
        return false
      }
    }

    return testName(t.data).toLowerCase().includes(state.filter.testName.toLowerCase())
  }

  function openJSON() {
    readReportData(text => {
      const blob = new Blob([text], { type: 'application/json' })
      var url = window.URL.createObjectURL(blob)
      window.open(url, "_blank")
    })
  }
</script>

<template>
  <template v-if="!store.isLoading">
    <div class="flex items-center text-3xl font-bold mb-5">
      <h1>{{ store.testCapture?.data.title }}</h1>
      <Elapsed :showIcon="true" :elapsed="totalElapsed" class="font-normal text-xl text-gray-500 ms-5" />
      <div class="ms-auto">
        <span class="text-green-700">{{ stats?.passed }}</span> / <span class="text-red-700">{{ stats?.failed }}</span> /
        <span>{{ stats?.total }}</span>
      </div>
    </div>
    <div class="flex flex-row items-center mb-3">
      <input v-model="state.filter.showPassed" id="check-show-passed" type="checkbox"
        class="border-solid border border-neutral-800 p-1">
      <label for="check-show-passed" class="ms-1">Passed</label>
      <input v-model="state.filter.showFailed" id="check-show-failed" type="checkbox"
        class="border-solid border border-neutral-800 p-1 ms-2">
      <label for="check-show-failed" class="ms-1">Failed</label>
      <input v-model="state.filter.showSkipped" id="check-show-skipped" type="checkbox"
        class="border-solid border border-neutral-800 p-1 ms-2">
      <label for="check-show-skipped" class="ms-1">Skipped</label>
    </div>
    <div class="flex flex-row items-center mb-3">
      <input v-model="state.filter.testName" type="text" class="border-solid border border-neutral-800 p-1 grow"
        placeholder="filter">
    </div>
    <div v-if="tests.length > 0">
      <div class="w-full mb-2">
        <Test v-for="test in tests" :key="test.data.index" :test="test" :depth="0" />
      </div>
    </div>
    <p v-else class="mb-3">Empty report!</p>
    <div class="flex items-start">
      <button @click="openJSON()" class="border-solid border border-neutral-800 rounded p-1">JSON</button>
      <div class="ms-auto">
        <p class="text-gray-500" v-if="store.testCapture">Test run started: {{
          store.testCapture.data.started_at }}</p>
        <p class="text-gray-500 mb-5" v-if="store.testCapture && store.testCapture.data.capture_started_at">Report
          generated:
          {{ store.testCapture.data.capture_started_at }}</p>
      </div>
    </div>
  </template>
</template>
