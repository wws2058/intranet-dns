<template>
    <div>
        <pre class="line-numbers"><code class="language-javascript">{{ jsonStr }}</code></pre>
    </div>
</template>

<script setup>
import { watch, nextTick } from 'vue';
import Prism from 'prismjs';
import 'prismjs/themes/prism.css';
import 'prism-themes/themes/prism-atom-dark.css';
import 'prismjs/plugins/line-numbers/prism-line-numbers.css';
import 'prismjs/plugins/line-numbers/prism-line-numbers';
import 'prismjs/components/prism-json';

const props = defineProps({
    jsonStr: {
        type: String,
        required: true
    }
});

// 此处需要借助nextTick在DOM更新完成以后执行Prism, 防止watch回调时, DOM还未更新, 会导致Prism处理到旧DOM
watch(() => props.jsonStr, async () => {
    await nextTick();
    Prism.highlightAll();
}, {
    immediate: true
});
</script>

<style scoped></style>