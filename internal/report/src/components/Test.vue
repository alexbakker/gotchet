<script setup lang="ts">
  import { computed, reactive } from 'vue'
  import { TestOutput, testName } from '../data/Test'
  import Elapsed from './Elapsed.vue'
  import {
    ChevronDownIcon,
    CheckCircleIcon,
    ExclamationCircleIcon,
    QuestionMarkCircleIcon,
    XCircleIcon,
  } from '@heroicons/vue/24/solid'
  import {
    ClipboardIcon,
    DocumentTextIcon
  } from '@heroicons/vue/24/outline'
  import { TestResult } from '../stores/report.ts'
  import { useReportStore } from '../stores/report.ts'

  const store = useReportStore()

  const props = defineProps<{
    test: TestResult
    depth: number
  }>()

  const hasChildTests = computed(() => {
    return props.test.tests.length > 0
  })

  const state = reactive({
    logsCollapsed: true,
    showClipboardButton: false
  })

  const name = computed(() => testName(props.test.data))

  const tests = computed(() =>
    props.test.tests
      .filter((t) => isTestShown(t))
      .sort((t1, t2) => t1.data.index - t2.data.index))

  const showLogsToggle = computed(() => {
    return !props.test.collapsed && hasChildTests.value
  })

  const selfOutput = computed(() => {
    const outputs: Array<TestOutput> = Object.assign([], props.test.data.output)
    outputs.sort((o1, o2) => o1.index - o2.index)
    return outputs
  })

  const indentStyle = reactive({
    'margin-left': (props.depth * 20) + 'px'
  })

  function toggleCollapse() {
    const collapse = !props.test.collapsed
    props.test.collapsed = collapse
    if (collapse) {
      state.logsCollapsed = true
    } else if (hasChildTests.value) {
      state.logsCollapsed = true
    } else {
      state.logsCollapsed = false
    }
  }

  function toggleLogsCollapse() {
    const collapse = !state.logsCollapsed
    if (!collapse && props.test.collapsed) {
      props.test.collapsed = false
    }
    if (collapse && !hasChildTests.value) {
      props.test.collapsed = true
    }
    state.logsCollapsed = collapse
  }

  async function copyToClipboard() {
    let log = ""
    for (const line of selfOutput.value) {
      log += line.text
    }
    await navigator.clipboard.writeText(log)
  }

  function testMatchesFilter(t: TestResult): boolean {
    const matches = testName(t.data).toLowerCase().includes(store.filter.testName.toLowerCase())
    if (!matches) {
      for (const st of t.tests) {
        if (testMatchesFilter(st)) {
          return true
        }
      }
    }

    return matches
  }

  function isTestShown(t: TestResult): boolean {
    if (t.data.done) {
      if (t.data.passed && !store.filter.showPassed) {
        return false
      }

      if (!t.data.passed && !store.filter.showFailed) {
        return false
      }

      if (t.data.skipped && !store.filter.showSkipped) {
        return false
      }
    }

    return testMatchesFilter(t)
  }
</script>

<template>
  <div class="-mb-px -ms-px -me-px border-solid border border-neutral-800">
    <div class="flex flex-col">
      <div class="flex items-center" :style="indentStyle">
        <div class="flex items-center cursor-pointer p-2 grow" @click="toggleCollapse()">
          <div class="me-1">
            <CheckCircleIcon v-if="test.data.passed" class="h-6 w-6 text-green-700" />
            <QuestionMarkCircleIcon v-else-if="!test.data.done" class="h-6 w-6 text-orange-700" />
            <ExclamationCircleIcon v-else-if="test.data.skipped" class="h-6 w-6 text-yellow-700" />
            <XCircleIcon v-else class="h-6 w-6 text-red-700" />
          </div>
          <span>{{ name }}</span>
        </div>
        <button v-show="showLogsToggle" @click="toggleLogsCollapse()" class="hover:bg-gray-300 rounded ms-auto p-1">
          <DocumentTextIcon class="h-6 w-6 text-gray-500" />
        </button>
        <Elapsed :showIcon="false" :elapsed="test.data.elapsed" class="font-normal text-gray-400 ms-auto cursor-pointer"
          @click="toggleCollapse()" />
        <ChevronDownIcon class="h-6 w-6 ms-2 me-2 cursor-pointer" :class="{ 'rotate-180': !props.test.collapsed }"
          @click="toggleCollapse()" />
      </div>
      <div v-if="!props.test.collapsed && !state.logsCollapsed && selfOutput.length > 0"
        class="log-container px-2 py-2 mt-2 self-stretch bg-gr" :style="indentStyle"
        @mouseenter="state.showClipboardButton = true" @mouseleave="state.showClipboardButton = false">
        <button :style="{ 'visibility': state.showClipboardButton ? 'visible' : 'hidden' }"
          class="float-right me-1 mt-1 border-solid border border-neutral-500 rounded p-1 hover:bg-gray-300"
          @click="copyToClipboard()">
          <ClipboardIcon class="h-6 w-6 text-gray-500" />
        </button>
        <code class="log" v-for="output in selfOutput">{{ output.text }}</code>
      </div>
      <Test v-if="!props.test.collapsed" v-for="test in tests" :key="test.data.index" :test="test"
        :depth="props.depth + 1" />
    </div>
  </div>
</template>

<style scoped>
  .log {
    white-space: pre-wrap;
    word-wrap: break-word;
    font-size: 1rem;
  }

  .log-container {
    background-color: rgb(245, 245, 245);
  }
</style>
