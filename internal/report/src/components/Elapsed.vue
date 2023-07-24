<script setup lang="ts">
  import { computed } from 'vue'
  import { ClockIcon } from '@heroicons/vue/24/solid';
  import formatDuration from 'date-fns/formatDuration'
  import intervalToDuration from 'date-fns/intervalToDuration'

  const props = defineProps<{
    elapsed: number | undefined
    class: string
    showIcon: boolean
  }>()

  const elapsed = computed(() => {
    if (!props.elapsed) {
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
      end: props.elapsed * 1000
    }), {
      format: ["hours", "minutes", "seconds"],
      locale: shortEnLocale,
      delimiter: ""
    })
  })
</script>

<template>
  <div class="flex items-center" :class="props.class">
    <ClockIcon v-if="props.showIcon" class="h-6 w-6" />
    <span :class="{ 'ms-1': props.showIcon }">{{ elapsed }}</span>
  </div>
</template>
