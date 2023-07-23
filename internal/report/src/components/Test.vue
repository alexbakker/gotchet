<script setup lang="ts">
  import { computed, reactive } from 'vue'
  import { Test, TestOutput } from '../data/Test'
  import {
    ChevronDownIcon,
    CheckCircleIcon,
    QuestionMarkCircleIcon,
    XCircleIcon
  } from '@heroicons/vue/24/solid'

  const props = defineProps<{
    test: Test
  }>()

  const state = reactive({
    collapsed: true
  })

  const name = computed(() => {
    const parts = props.test.full_name.split("/")
    return parts[parts.length - 1]
  })

  /*const tests = computed(() =>
    Object.values(props.test.tests).sort((t1, t2) => t1.index - t2.index))*/

  const fullOutput = computed(() => {
    const outputs: Array<TestOutput> = Object.assign([], props.test.output)

    const merge = (t: Test) => {
      for (const test of Object.values(t.tests)) {
        outputs.push(...test.output)
        merge(test)
      }
    }
    merge(props.test)

    outputs.sort((o1, o2) => o1.index - o2.index)
    return outputs
  })
</script>

<template>
  <div class="flex flex-col p-2" :class="{ 'expanded': !state.collapsed }">
    <div @click="state.collapsed = !state.collapsed" class="flex items-center cursor-pointer">
      <div class="me-1">
        <CheckCircleIcon v-if="test.passed" class="h-6 w-6 text-green-700" />
        <QuestionMarkCircleIcon v-else-if="!test.done" class="h-6 w-6 text-yellow-700" />
        <XCircleIcon v-else="test.done" class="h-6 w-6 text-red-700" />
      </div>
      <span>{{ name }}</span>
      <ChevronDownIcon class="h-6 w-6 ms-auto self-end" :class="{ 'rotate-180': !state.collapsed }" />
    </div>
    <div v-if="!state.collapsed" class="log-container px-2 py-2 mt-2 overflow-y-scroll self-stretch bg-gr">
      <code class="log" v-for="output in fullOutput">{{ output.text }}</code>
    </div>
  </div>
</template>

<style scoped>
  .expanded {
    height: 50vh;
  }

  .log {
    white-space: pre-wrap;
    word-wrap: break-word;
    font-size: 1rem;
  }

  .log-container {
    background-color: rgb(245, 245, 245)
  }
</style>
