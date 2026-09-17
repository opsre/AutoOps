<template>
  <div class="markdown-table-container">
    <el-table
      :data="tableData.rows"
      stripe
      border
      size="small"
      :max-height="maxHeight"
      style="width: 100%"
      class="ai-markdown-table"
    >
      <el-table-column
        v-for="col in tableData.columns"
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :min-width="80"
        show-overflow-tooltip
      />
    </el-table>
  </div>
</template>

<script>
import { parseMarkdownTable } from '@/utils/parseMarkdownTable'

export default {
  name: 'MarkdownTable',
  props: {
    source: {
      type: [String, Object],
      required: true
    },
    maxHeight: {
      type: Number,
      default: 400
    }
  },
  computed: {
    tableData() {
      return parseMarkdownTable(this.source)
    }
  }
}
</script>

<style lang="less" scoped>
.markdown-table-container {
  margin: 8px 0;
  border-radius: 6px;
  overflow: hidden;

  .ai-markdown-table {
    background: transparent;

    :deep(.el-table__header-wrapper) {
      background: rgba(26, 31, 46, 0.9);

      th {
        background: rgba(26, 31, 46, 0.9);
        color: rgba(255, 255, 255, 0.95);
        font-weight: 700;
        font-size: 11px;
        border-color: rgba(102, 126, 234, 0.25);
        padding: 6px 8px;
        line-height: 1.4;
      }
    }

    :deep(.el-table__body-wrapper) {
      td {
        background: rgba(26, 31, 46, 0.4);
        color: rgba(255, 255, 255, 0.85);
        font-size: 11px;
        border-color: rgba(102, 126, 234, 0.15);
        padding: 5px 8px;
        line-height: 1.4;
      }

      .el-table__row--striped td {
        background: rgba(26, 31, 46, 0.6);
      }

      .el-table__row:hover td {
        background: rgba(102, 126, 234, 0.2) !important;
      }
    }

    :deep(.el-table__border-left-patch),
    :deep(.el-table__border-bottom-patch) {
      background: rgba(102, 126, 234, 0.2);
    }

    // 紧凑型表格样式
    :deep(.el-table--small) {
      .el-table__cell {
        padding: 5px 0;
      }
    }

    // 滚动条样式
    :deep(.el-table__body-wrapper)::-webkit-scrollbar {
      width: 6px;
      height: 6px;
    }

    :deep(.el-table__body-wrapper)::-webkit-scrollbar-track {
      background: rgba(102, 126, 234, 0.05);
      border-radius: 3px;
    }

    :deep(.el-table__body-wrapper)::-webkit-scrollbar-thumb {
      background: rgba(102, 126, 234, 0.3);
      border-radius: 3px;
      transition: all 0.3s ease;

      &:hover {
        background: rgba(102, 126, 234, 0.5);
      }
    }
  }
}
</style>
