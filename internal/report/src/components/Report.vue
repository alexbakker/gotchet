<script setup lang="ts">
  import { onMounted, reactive, computed } from 'vue'
  import Test from './Test.vue'
  import Elapsed from './Elapsed.vue'
  import { TestResult, testName } from '../data/Test'

  const state = reactive<{
    rootTest: TestResult | null,
    isLoading: boolean,
    filter: {
      testName: string
      showPassed: boolean
      showFailed: boolean
    }
  }>({
    rootTest: null,
    isLoading: true,
    filter: {
      testName: "",
      showPassed: true,
      showFailed: true
    }
  })

  const stats = computed(() => {
    const shownTests = tests.value
    return {
      total: shownTests.length,
      passed: shownTests.filter((t) => t.done && t.passed).length,
      failed: shownTests.filter((t) => t.done && !t.passed).length
    }
  })

  const tests = computed(() => {
    if (!state.rootTest) {
      return []
    }

    return Object.values(state.rootTest.tests)
      .filter((t) => isTestShown(t))
      .sort((t1, t2) => t1.index - t2.index)
  })

  onMounted(() => {
    loadReport()
  })

  function isTestShown(t: TestResult): boolean {
    if (t.done) {
      if (t.passed && !state.filter.showPassed) {
        return false
      }

      if (!t.passed && !state.filter.showFailed) {
        return false
      }
    }

    return testName(t).toLowerCase().includes(state.filter.testName.toLowerCase())
  }

  function readData(cb: (text: string) => void) {
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

  function loadReport() {
    state.isLoading = true
    readData(text => {
      let rootTest: TestResult = JSON.parse(text)
      if (rootTest) {
        state.rootTest = rootTest
        state.isLoading = false
      }
    })
  }

  function openJSON() {
    readData(text => {
      const blob = new Blob([text], { type: 'application/json' })
      var url = window.URL.createObjectURL(blob)
      window.open(url, "_blank")
    })
  }
</script>

<template>
  <template v-if="!state.isLoading">
    <div class="flex items-center text-3xl font-bold mb-5">
      <h1>{{ state.rootTest?.title }}</h1>
      <Elapsed :showIcon="true" :elapsed="state.rootTest?.elapsed" class="font-normal text-xl text-gray-500 ms-5" />
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
    </div>
    <div class="flex flex-row items-center mb-3">
      <input v-model="state.filter.testName" type="text" class="border-solid border border-neutral-800 p-1 grow"
        placeholder="filter">
    </div>
    <div v-if="tests.length > 0">
      <div class="w-full mb-2">
        <Test v-for="test in tests" :key="test.index" :test="test" :depth="0" />
      </div>
    </div>
    <p v-else class="mb-3">Empty report!</p>
    <div class="flex items-start">
      <button @click="openJSON()" class="border-solid border border-neutral-800 rounded p-1">JSON</button>
      <div class="ms-auto">
        <p class="text-gray-500" v-if="state.rootTest">Test run started: {{
          state.rootTest.started_at }}</p>
        <p class="text-gray-500 mb-5" v-if="state.rootTest && state.rootTest.capture_started_at">Report generated: {{
          state.rootTest.capture_started_at }}</p>
      </div>
    </div>
  </template>
</template>
