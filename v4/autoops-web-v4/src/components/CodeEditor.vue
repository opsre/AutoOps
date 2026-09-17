<template>
  <div class="monaco-editor-container" ref="editorContainer" :style="{ height: height }"></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'sql'
  },
  readonly: {
    type: Boolean,
    default: false
  },
  completionItems: {
    type: Array,
    default: () => []
  },
  height: {
    type: String,
    default: '400px'
  },
  theme: {
    type: String,
    default: 'vs-dark'
  },
  wordWrap: {
    type: String,
    default: 'off'
  },
  searchText: {
    type: String,
    default: ''
  },
  searchResults: {
    type: Array,
    default: () => []
  },
  currentSearchIndex: {
    type: Number,
    default: -1
  }
})

const emit = defineEmits(['update:modelValue'])

const editorContainer = ref(null)
let editor = null
let completionProvider = null
let searchDecorations = []

const ensureThemes = () => {
  monaco.editor.defineTheme('terminal-yellow', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      { token: '', foreground: 'C9A227' }
    ],
    colors: {
      'editor.background': '#000000',
      'editor.foreground': '#C9A227',
      'editorLineNumber.foreground': '#6F5A16',
      'editorCursor.foreground': '#C9A227',
      'editor.selectionBackground': '#3D320F',
      'editor.inactiveSelectionBackground': '#241D08'
    }
  })
}

const updateSearchDecorations = () => {
  if (!editor) {
    return
  }

  const hasSearch = props.searchText && props.searchResults.length > 0
  if (!hasSearch) {
    searchDecorations = editor.deltaDecorations(searchDecorations, [])
    return
  }

  const decorations = props.searchResults.map((result, index) => {
    const startLineNumber = Number(result.lineIndex || 0) + 1
    const startColumn = Number(result.startIndex || 0) + 1
    const endColumn = Number(result.endIndex || 0) + 1
    const isCurrent = index === props.currentSearchIndex

    return {
      range: new monaco.Range(startLineNumber, startColumn, startLineNumber, endColumn),
      options: {
        inlineClassName: isCurrent ? 'editor-search-current-match' : 'editor-search-match'
      }
    }
  })

  searchDecorations = editor.deltaDecorations(searchDecorations, decorations)
}

const revealCurrentSearchResult = () => {
  if (!editor || props.currentSearchIndex < 0 || props.currentSearchIndex >= props.searchResults.length) {
    return
  }

  const currentResult = props.searchResults[props.currentSearchIndex]
  const lineNumber = Number(currentResult.lineIndex || 0) + 1
  const startColumn = Number(currentResult.startIndex || 0) + 1
  const endColumn = Number(currentResult.endIndex || 0) + 1

  editor.setSelection(new monaco.Range(lineNumber, startColumn, lineNumber, endColumn))
  editor.revealLineInCenter(lineNumber)
}

// --- Auto-completion ---
const setupAutoCompletion = (items) => {
  // If a provider already exists, dispose of it first
  if (completionProvider) {
    completionProvider.dispose()
  }

  if (!items || items.length === 0) {
    return
  }

  completionProvider = monaco.languages.registerCompletionItemProvider(props.language, {
    provideCompletionItems: (model, position) => {
      const suggestions = []

      // Add tables and columns to suggestions
      items.forEach(table => {
        // Suggest table name
        suggestions.push({
          label: table.tableName,
          kind: monaco.languages.CompletionItemKind.Folder,
          insertText: table.tableName,
          detail: 'Table'
        })

        // Suggest columns for that table
        table.columns.forEach(column => {
          suggestions.push({
            label: column.columnName,
            kind: monaco.languages.CompletionItemKind.Field,
            insertText: column.columnName,
            detail: `Column in ${table.tableName} (${column.columnType})`
          })
        })
      })
      
      // Add basic SQL keywords
      const keywords = ['SELECT', 'FROM', 'WHERE', 'INSERT', 'INTO', 'VALUES', 'UPDATE', 'SET', 'DELETE', 'GROUP BY', 'ORDER BY', 'AND', 'OR', 'LIMIT', 'OFFSET', 'COUNT', 'AVG', 'SUM', 'MIN', 'MAX', 'JOIN', 'LEFT JOIN', 'RIGHT JOIN', 'INNER JOIN', 'ON']
      keywords.forEach(keyword => {
        suggestions.push({
          label: keyword,
          kind: monaco.languages.CompletionItemKind.Keyword,
          insertText: keyword
        })
      })

      return {
        suggestions: suggestions
      }
    }
  })
}

watch(() => props.completionItems, (newItems) => {
  if (editor) {
    setupAutoCompletion(newItems)
  }
}, { deep: true, immediate: true })

// --- Editor Lifecycle ---
onMounted(() => {
  if (editorContainer.value) {
    ensureThemes()

    editor = monaco.editor.create(editorContainer.value, {
      value: props.modelValue,
      language: props.language,
      readOnly: props.readonly,
      theme: props.theme,
      automaticLayout: true,
      fontSize: 14,
      minimap: {
        enabled: false // Disable minimap
      },
      scrollBeyondLastLine: false,
      wordWrap: props.wordWrap,
    })

    // Sync content back to parent
    editor.onDidChangeModelContent(() => {
      const currentValue = editor.getValue()
      if (currentValue !== props.modelValue) {
        emit('update:modelValue', currentValue)
      }
    })

    // Setup initial completion items if they exist
    if (props.completionItems.length > 0) {
      setupAutoCompletion(props.completionItems)
    }

    updateSearchDecorations()
    revealCurrentSearchResult()
  }
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
  if (completionProvider) {
    completionProvider.dispose()
  }
})

// Watch for external changes to modelValue
watch(() => props.modelValue, (newValue) => {
  if (editor && editor.getValue() !== newValue) {
    editor.setValue(newValue)
  }
})

// Watch for external changes to readonly state
watch(() => props.readonly, (newVal) => {
  if (editor) {
    editor.updateOptions({ readOnly: newVal })
  }
})

watch(() => props.theme, (newVal) => {
  if (editor) {
    monaco.editor.setTheme(newVal)
  }
})

watch(() => props.wordWrap, (newVal) => {
  if (editor) {
    editor.updateOptions({ wordWrap: newVal })
  }
})

watch(
  () => [props.searchText, props.searchResults, props.currentSearchIndex],
  () => {
    updateSearchDecorations()
    revealCurrentSearchResult()
  },
  { deep: true }
)

// Expose focus method
const focus = () => {
  if (editor) {
    editor.focus()
  }
}

const scrollToBottom = () => {
  if (!editor) {
    return
  }

  const model = editor.getModel()
  const lineCount = model?.getLineCount() || 1
  editor.revealLine(lineCount)

  // 再直接设置一次滚动位置，避免某些情况下 revealLine 未贴底
  const scrollHeight = editor.getScrollHeight()
  editor.setScrollTop(scrollHeight)
}

const scrollToLine = (lineNumber) => {
  if (!editor) {
    return
  }

  const model = editor.getModel()
  const maxLine = model?.getLineCount() || 1
  const targetLine = Math.min(Math.max(1, lineNumber), maxLine)
  editor.revealLineInCenter(targetLine)
}

defineExpose({ focus, scrollToBottom, scrollToLine })

</script>

<style scoped>
.monaco-editor-container {
  width: 100%;
  height: 100%;
  min-height: 200px; /* Ensure a minimum height */
  border: 1px solid #ccc;
  border-radius: 4px;
  overflow: hidden;
}
</style>

<style>
.monaco-editor .editor-search-match {
  background: rgba(135, 206, 250, 0.42) !important;
  border-bottom: 1px solid rgba(135, 206, 250, 0.9) !important;
  color: #0b0b0b !important;
}

.monaco-editor .editor-search-current-match {
  background: rgba(135, 206, 250, 0.72) !important;
  border-bottom: 1px solid #87cefa !important;
  color: #000000 !important;
}
</style>
