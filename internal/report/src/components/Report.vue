<script setup lang="ts">
  import { onMounted, reactive, computed } from 'vue'
  import TestComponent from './Test.vue'
  import { Test, testName } from '../data/Test'
  import { ClockIcon } from '@heroicons/vue/24/solid';
  import formatDuration from 'date-fns/formatDuration'
  import intervalToDuration from 'date-fns/intervalToDuration'

  const state = reactive<{
    title: string,
    test: Test | undefined,
    isLoading: boolean,
    filter: {
      testName: string
    }
  }>({
    title: "Go Test Report",
    test: undefined,
    isLoading: true,
    filter: {
      testName: ""
    }
  });

  const stats = computed(() => {
    const shownTests = tests.value;
    return {
      total: shownTests.length,
      passed: shownTests.filter((t) => t.done && t.passed).length,
      failed: shownTests.filter((t) => t.done && !t.passed).length
    }
  })

  const elapsed = computed(() => {
    if (!state.test) {
      return "?"
    }

    const formatDistanceLocale: Record<string, string>
      = { xSeconds: '{{count}}s', xMinutes: '{{count}}m', xHours: '{{count}}h' }
    const shortEnLocale = {
      formatDistance: (token: string, count: string) => {
        return formatDistanceLocale[token].replace('{{count}}', count)
      }
    }

    return formatDuration(intervalToDuration({
      start: 0,
      end: state.test.elapsed * 1000
    }), {
      format: ["hours", "minutes", "seconds"],
      locale: shortEnLocale,
      delimiter: ""
    })
  })

  const tests = computed(() => {
    if (!state.test) {
      return []
    }

    return Object.values(state.test.tests)
      .filter((t) => {
        return testName(t).toLowerCase().includes(state.filter.testName.toLowerCase())
      })
      .sort((t1, t2) => t1.index - t2.index)
  })

  onMounted(() => {
    loadReport()
  })

  function setTitle(title: string) {
    document.title = title
    state.title = title
  }

  function setParentTest(test: Test) {
    setTitle(test.full_name)
    state.test = test
  }

  function readData(cb: (text: string) => void) {
    function decompress(byteArray: ArrayBuffer, encoding: CompressionFormat) {
      const dcs = new DecompressionStream(encoding);
      const writer = dcs.writable.getWriter();
      writer.write(byteArray);
      writer.close();
      return new Response(dcs.readable).arrayBuffer().then(function (ab) {
        return new TextDecoder().decode(ab);
      });
    }

    const raw = document.getElementById("report-data")!.innerText;
    fetch(new URL(raw))
      .then(r => r.arrayBuffer())
      .then(ab => decompress(ab, "gzip"))
      .then(text => cb(text))
      .catch(e => console.error(e))
  }

  function loadReport() {
    state.isLoading = true
    readData(text => {
      let rootTest: Test = JSON.parse(text);
      if (rootTest) {
        // Go deep until there's more than one child test
        let test = rootTest;
        while (test.tests && Object.keys(test.tests).length == 1) {
          test = Object.values(test.tests)[0];
        }
        setParentTest(test);
        state.isLoading = false
      }
    })
  }

  function openJSON() {
    readData(text => {
      const blob = new Blob([text], { type: 'application/json' });
      var url = window.URL.createObjectURL(blob);
      window.open(url, "_blank")
    })
  }
</script>

<template>
  <template v-if="!state.isLoading">
    <div class="flex items-center text-3xl font-bold mb-5">
      <h1>{{ state.title }}</h1>
      <div class="flex items-center font-normal text-xl text-gray-500 ms-5">
        <ClockIcon class="h-6 w-6" />
        <span class="ms-1">{{ elapsed }}</span>
      </div>
      <div class="ms-auto">
        <span class="text-green-700">{{ stats?.passed }}</span> / <span class="text-red-700">{{ stats?.failed }}</span> /
        <span>{{ stats?.total }}</span>
      </div>
    </div>
    <div class="flex items-center">
      <input v-model="state.filter.testName" type="text">
    </div>
    <table v-if="tests.length > 0" class="w-full table-fixed border-collapse">
      <tr v-for="test in tests" :key="test.index" class="border-solid border border-neutral-800 rounded-md">
        <TestComponent :test="test" />
      </tr>
    </table>
    <p v-else>Empty report!</p>
    <button @click="openJSON()">JSON</button>
  </template>
</template>
