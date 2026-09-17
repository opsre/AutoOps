<template>
  <div class="db-terminal-page">
    <el-container class="terminal-container">
      <!-- 左侧树形结构 -->
      <el-aside :width="asideWidth + 'px'" class="tree-aside">
        <el-card shadow="never" class="tree-card">
          <template #header>
            <div class="tree-header">
              <span>数据库资源</span>
              <el-button text size="small" @click="loadInstances" :loading="loadingTree">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>

          <el-input v-model="treeFilter" size="small" placeholder="过滤..." clearable class="mb10">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <el-scrollbar height="calc(100vh - 200px)">
            <el-tree
              ref="treeRef"
              :data="treeData"
              :props="treeProps"
              :filter-node-method="filterNode"
              node-key="id"
              highlight-current
              :current-node-key="selectedNode && selectedNode.id"
              :expand-on-click-node="false"
              @node-click="handleNodeClick"
            >
              <template #default="{ node, data }">
                <span class="tree-node" @contextmenu.prevent.stop="onTreeContextMenu($event, data)">
                  <el-icon v-if="data.type === 'instance'">
                    <img :src="getDbIcon(data.dbType)" class="node-icon" />
                  </el-icon>
                  <el-icon v-else-if="data.type === 'database'"><img src="@/assets/image/数据库.svg" class="node-icon" /></el-icon>
                  <el-icon v-else-if="data.type === 'table'"><img src="@/assets/image/数据表.svg" class="node-icon" /></el-icon>
                  <span class="node-label">{{ node.label }}</span>
                  <span v-if="data.type === 'table' && data.size" class="node-size">{{ data.size }}</span>
                  <span v-if="data.comment" class="node-comment">{{ data.comment }}</span>
                </span>
              </template>
            </el-tree>
          </el-scrollbar>
        </el-card>
      </el-aside>
      <!-- Draggable resizer between aside and main -->
      <div class="aside-resizer" @mousedown="startAsideResize">
        <el-icon class="resizer-icon"><DArrowRight /></el-icon>
      </div>

            <!-- 右侧内容区域 -->
            <el-main class="main-content" ref="mainContent">
              <div v-if="openTabs.length > 0">
                <el-tabs v-model="activeTab" type="card" closable @tab-remove="handleTabRemove">
                  <el-tab-pane
                    v-for="tab in openTabs"
                    :key="tab.name"
                    :label="tab.title"
                    :name="tab.name"
                  >
                    <!-- Conditional rendering based on tab type -->
                    <div v-if="!tab.type || tab.type === 'table'">
                      <el-card shadow="never" class="content-card">
                        <template #header>
                          <div class="content-header">
                            <div class="header-info">
                              <el-icon><img src="@/assets/image/数据表.svg" class="node-icon" /></el-icon>
                              <span class="header-title">{{ tab.title }}</span>
                            </div>
                            <div v-if="tab.connection" class="header-connection-info">
                              <span style="margin-right: 15px;"><strong>数据库名称:</strong> <span class="info-value">{{ tab.connection.instanceName }}</span></span>
                              <span v-if="tab.connection.connectionAddress" style="margin-right: 15px;"><strong>连接地址:</strong> <span class="info-value">{{ tab.connection.connectionAddress }}</span></span>
                              <span style="margin-right: 15px;"><strong>库名:</strong> <span class="info-value">{{ tab.connection.dbName }}</span></span>
                              <span style="margin-right: 15px;"><strong>表名:</strong> <span class="info-value">{{ tab.title }}</span></span>
                            </div>
                            <div class="header-actions">
                              <el-button type="primary" circle @click="openInsertTab(tab)" title="新增数据">
                                <el-icon><Plus /></el-icon>
                              </el-button>
                              <el-button type="success" circle @click="refreshTableData(tab)" :loading="tab.loadingData" title="刷新">
                                <el-icon><Refresh /></el-icon>
                              </el-button>
                              <el-button color="#FFA500" circle @click="openQueryTab(tab)" title="查询数据">
                                <el-icon><Search /></el-icon>
                              </el-button>
                              <el-button type="warning" circle @click="saveTableData(tab)" :disabled="!tab.modifiedRows || tab.modifiedRows.size === 0" title="保存修改">
                                <el-icon><Check /></el-icon>
                              </el-button>
                              <el-upload
                                :show-file-list="false"
                                accept=".sql"
                                :before-upload="beforeUploadSql"
                                :http-request="(p) => uploadSqlFile(tab, p.file)"
                                style="display: inline-block;"
                              >
                                <el-button type="primary" circle title="导入SQL脚本">
                                  <el-icon><Upload /></el-icon>
                                </el-button>
                              </el-upload>
                            </div>
                          </div>
                        </template>
                        <div class="table-data-section">
                          <el-table
                            :data="tab.tableData"
                            v-loading="tab.loadingData"
                            border
                            stripe
                            size="small"
                            max-height="calc(100vh - 360px)"
                            class="data-table"
                            @cell-dblclick="(row, column) => handleCellDoubleClick(row, column, tab)"
                            @row-contextmenu="(row, column, e) => onRowContextMenu(row, column, e, tab)"
                          >
                            <!-- Table columns -->
                            <el-table-column type="index" width="60" label="序号" />
                            <el-table-column
                              v-for="col in tab.tableColumns"
                              :key="col.columnName"
                              :prop="col.columnName"
                              :label="col.columnName"
                              :min-width="col.columnName === 'id' ? 90 : 130"
                              show-overflow-tooltip
                            >
                              <template #header>
                                <div class="column-header">
                                  <div class="column-title">
                                    <span>{{ col.columnName }}</span>
                                    <div v-if="col.columnName === 'id' || col.isPrimaryKey" class="pk-hint">(主键)</div>
                                  </div>
                                  <el-tooltip
                                    v-if="(col.columnComment || col.comment) && !((col.columnName === 'id' || col.isPrimaryKey) && sanitizeComment(col.columnComment || col.comment) === '主键')"
                                    :content="sanitizeComment(col.columnComment || col.comment)"
                                    placement="top"
                                  >
                                    <el-icon color="#909399"><InfoFilled /></el-icon>
                                  </el-tooltip>
                                </div>
                              </template>
                              <template #default="scope">
                                <div v-if="editingCell && editingCell.row === scope.row && editingCell.column === col.columnName" class="edit-cell">
                                  <el-input
                                    v-model="scope.row[col.columnName]"
                                    @blur="handleInputBlur(scope.row, col.columnName, tab)"
                                    @keyup.enter.native="$event.target.blur()"
                                    size="small"
                                    autofocus
                                  ></el-input>
                                </div>
                                <div v-else>
                                  <span :class="{ 'null-cell': isNullish(scope.row[col.columnName]) }">
                                    {{ isNullish(scope.row[col.columnName]) ? 'null' : scope.row[col.columnName] }}
                                  </span>
                                </div>
                              </template>
                            </el-table-column>
                          </el-table>
                          <div class="pagination-section">
                            <el-pagination
                              v-model:current-page="tab.pagination.currentPage"
                              v-model:page-size="tab.pagination.pageSize"
                              :page-sizes="[10, 20, 50, 100]"
                              :total="tab.pagination.total"
                              layout="total, sizes, prev, pager, next"
                              @size-change="size => handlePageSizeChange(size, tab)"
                              @current-change="page => handlePageChange(page, tab)"
                            />
                          </div>
                          </div>
                      </el-card>
                    </div>
                    
                    <div v-else-if="tab.type === 'query'">
                      <el-card shadow="never" class="content-card">
                        <template #header>
                          <div class="content-header">
                            <div class="header-info">
                              <el-icon><img src="@/assets/image/数据库.svg" class="node-icon" /></el-icon>
                              <span class="header-title">{{ tab.title }}</span>
                            </div>
                            <div v-if="tab.connection" class="header-connection-info">
                              <span style="margin-right: 15px;"><strong>数据库名称:</strong> <span class="info-value">{{ tab.connection.instanceName }}</span></span>
                              <span v-if="tab.connection.connectionAddress" style="margin-right: 15px;"><strong>连接地址:</strong> <span class="info-value">{{ tab.connection.connectionAddress }}</span></span>
                              <span style="margin-right: 15px;"><strong>库名:</strong> <span class="info-value">{{ tab.connection.dbName }}</span></span>
                              <span style="margin-right: 15px;"><strong>表名:</strong> <span class="info-value">Query</span></span>
                            </div>
                          </div>
                        </template>
                        <div class="query-tab-container">
                          <el-form :inline="true" size="small" label-width="80px">
                            <el-form-item label="使用缓存">
                              <el-checkbox v-model="tab.useCache">启用</el-checkbox>
                            </el-form-item>
                            <el-form-item label="缓存时长(秒)">
                              <el-input-number v-model="tab.cacheSeconds" :min="10" :max="600" />
                            </el-form-item>
                            <el-form-item label="最大行数">
                              <el-input-number v-model="tab.maxRows" :min="100" :max="5000" :step="100" />
                            </el-form-item>
                            <el-form-item label="超时(秒)">
                              <el-input-number v-model="tab.timeoutSeconds" :min="5" :max="300" />
                            </el-form-item>
                          </el-form>
                          <CodeEditor v-model="tab.sqlQuery" :completion-items="tab.sqlCompletionItems" language="sql" height="250px" />
                                                  <div class="sql-actions">
                                                    <el-button type="primary" @click="executeSql(tab)" :loading="tab.executing">
                                                      <el-icon><VideoPlay /></el-icon>
                                                      执行
                                                    </el-button>
                                                    <el-button type="info" @click="preCheckSql(tab)" :loading="tab.preChecking">
                                                      <el-icon><View /></el-icon>
                                                      预检查
                                                    </el-button>
                                                    <el-button type="warning" @click="formatSql(tab)" :loading="tab.formatting">
                                                      <el-icon><MagicStick /></el-icon>
                                                      格式化
                                                    </el-button>
                                                    <el-button
                                                      type="success"
                                                      plain
                                                      :disabled="!tab.sqlResult || !tab.sqlResult.success || !(tab.sqlResult.rows && tab.sqlResult.rows.length)"
                                                      @click="exportQueryResultCsv(tab)"
                                                    >
                                                      <el-icon><Download /></el-icon>
                                                      导出CSV
                                                    </el-button>
                                                    <el-text v-if="tab.sqlResult && tab.sqlResult.success" type="success" style="margin-left: 15px;">
                                                      执行成功，耗时: {{ tab.sqlResult.execTime }}ms
                                                    </el-text>
                                                    <el-text v-if="tab.preCheckResult" type="info" style="margin-left: 15px;">
                                                      预检查: {{ tab.preCheckResult.sqlTypeText }} 语句，预计影响 {{ tab.preCheckResult.affectedRows }} 行
                                                    </el-text>
                                                  </div>
                                                  <div v-if="tab.sqlResult && (!tab.sqlResult.success || (tab.sqlResult.rows && tab.sqlResult.rows.length > 0))" class="sql-result-container">
                                                    <el-divider>执行结果</el-divider>
                                                    <div v-if="tab.sqlResult.success">
                                                        <el-table
                                                          v-if="tab.sqlResult.rows && tab.sqlResult.rows.length > 0"
                                                          :data="tab.sqlResult.rows"
                                                          border
                                                          stripe
                                                          size="small"
                                                          class="sql-result-table"
                                                          max-height="320"
                                                        >
                                                          <el-table-column
                                                            v-for="col in tab.sqlResult.columns"
                                                            :key="col"
                                                            :prop="col"
                                                            :label="col"
                                                            min-width="100"
                                                            show-overflow-tooltip
                                                          />
                                                        </el-table>
                                                        <el-text v-else>查询成功，但没有返回任何行。</el-text>
                                                    </div>
                                                    <el-result v-else icon="error" title="执行失败" :sub-title="tab.sqlResult.error" />
                                                  </div>                        </div>
                      </el-card>
                    </div>
                  </el-tab-pane>
                </el-tabs>
              </div>
              <el-empty v-else description="请从左侧选择一张表进行查看" />

              <!-- Insert Dialog -->
              <el-dialog v-model="insertDialogVisible" :title="insertFormTab && insertFormTab.node ? `新增: ${insertFormTab.node.tableName}` : '新增'" width="600px" append-to-body>
                <div class="insert-form-container">
                  <el-form :model="insertFormData" label-width="120px" size="small">
                    <el-form-item v-for="col in insertFormColumns" :key="col.columnName" :label="col.columnName">
                      <el-input :placeholder="sanitizeComment(col.columnComment || col.comment) ? '(' + sanitizeComment(col.columnComment || col.comment) + ')' : ''" v-model="insertFormData[col.columnName]" />
                      <el-text v-if="sanitizeComment(col.columnComment || col.comment)" type="info" size="small" class="field-hint">
                        ({{ sanitizeComment(col.columnComment || col.comment) }})
                      </el-text>
                    </el-form-item>
                  </el-form>
                </div>
                <template #footer>
                  <span class="dialog-footer">
                    <el-button @click="cancelInsertForm">取消</el-button>
                    <el-button type="primary" @click="submitInsertForm" :loading="inserting">提交</el-button>
                  </span>
                </template>
              </el-dialog>

              <!-- Export Job Progress Dialog -->
              <el-dialog v-model="exportJobDialog.visible" title="整库导出进度" width="640px" append-to-body @close="onExportDialogClose">
                <div>
                  <div style="margin-bottom: 10px;">
                    任务ID: <el-text type="info">{{ exportJobDialog.jobId }}</el-text>
                    <span style="margin-left: 12px;">状态:
                      <el-tag :type="exportJobDialog.status === 'completed' ? 'success' : (exportJobDialog.status === 'failed' ? 'danger' : 'warning')" size="small">
                        {{ exportJobDialog.status || '-' }}
                      </el-tag>
                    </span>
                  </div>
                  <el-progress :percentage="computeExportPercent()" :status="exportJobDialog.status === 'failed' ? 'exception' : (exportJobDialog.status === 'completed' ? 'success' : undefined)" />
                  <div style="margin-top: 10px; display:flex; gap:16px;">
                    <el-text>总表数: {{ exportJobDialog.tables }}</el-text>
                    <el-text type="primary">已完成: {{ exportJobDialog.doneTables }}</el-text>
                    <el-text>文件大小: {{ formatFileSize(exportJobDialog.fileSize) }}</el-text>
                  </div>
                </div>
                <template #footer>
                  <el-button @click="exportJobDialog.visible = false">关闭</el-button>
                  <el-button type="primary" :disabled="!canDownloadExport()" @click="downloadExportResult">下载结果</el-button>
                </template>
              </el-dialog>

              <!-- SQL Import Progress Dialog -->
              <el-dialog v-model="sqlImportDialog.visible" title="SQL 导入进度" width="640px" append-to-body @close="onSqlImportDialogClose">
                <div>
                  <div style="margin-bottom: 10px;">
                    任务ID: <el-text type="info">{{ sqlImportDialog.jobId }}</el-text>
                    <span style="margin-left: 12px;">状态: 
                      <el-tag :type="sqlImportDialog.status === 'completed' ? 'success' : (sqlImportDialog.status === 'failed' ? 'danger' : 'warning')" size="small">
                        {{ sqlImportDialog.status || '-' }}
                      </el-tag>
                    </span>
                  </div>
                  <el-progress :percentage="computeImportPercent()" :status="sqlImportDialog.status === 'failed' ? 'exception' : (sqlImportDialog.status === 'completed' ? 'success' : undefined)" />
                  <div style="margin-top: 10px; display:flex; gap:16px;">
                    <el-text>总数: {{ sqlImportDialog.total }}</el-text>
                    <el-text type="success">成功: {{ sqlImportDialog.success }}</el-text>
                    <el-text type="danger">失败: {{ sqlImportDialog.failed }}</el-text>
                  </div>
                  <div v-if="sqlImportDialog.errors && sqlImportDialog.errors.length" style="margin-top: 12px;">
                    <el-alert title="错误列表" type="error" :closable="false" />
                    <el-scrollbar height="160px" style="margin-top:6px;">
                      <ul class="sql-import-errors">
                        <li v-for="(err, idx) in sqlImportDialog.errors" :key="idx">
                          <el-text type="danger">{{ typeof err === 'string' ? err : (err.message || JSON.stringify(err)) }}</el-text>
                        </li>
                      </ul>
                    </el-scrollbar>
                  </div>
                </div>
                <template #footer>
                  <el-button @click="sqlImportDialog.visible = false">关闭</el-button>
                </template>
              </el-dialog>

              <!-- Row Context Menu -->
              <div v-show="contextMenu.visible" :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }" class="row-context-menu" @click.stop>
                <div class="menu-item" @click="handleViewJson">
                  <el-icon><Document /></el-icon>
                  查看 JSON
                </div>
                <div class="menu-item" @click="handleViewRowForm">
                  <el-icon><Tickets /></el-icon>
                  表单视图
                </div>
                <div class="menu-item" @click="handleExportCsv">
                  <el-icon><Download /></el-icon>
                  导出 CSV
                </div>
                <div class="menu-item" @click="handleExportSql">
                  <el-icon><Document /></el-icon>
                  导出 SQL
                </div>
                <div class="menu-item danger" @click="handleDeleteRow">
                  <el-icon><Delete /></el-icon>
                  删除该行
                </div>
              </div>

              <!-- Tree Context Menu (teleport to body, fixed near cursor) -->
              <teleport to="body">
                <div v-show="treeMenu.visible" :style="{ top: treeMenu.y + 'px', left: treeMenu.x + 'px' }" class="tree-context-menu" @click.stop>
                  <div class="menu-item" @click="handleCopyTableStructure">
                    <el-icon><DocumentCopy /></el-icon>
                    复制表结构
                  </div>
                  <div class="menu-item" @click="handleRenameTable">
                    <el-icon><EditPen /></el-icon>
                    重命名
                  </div>
                  <div class="menu-item" @click="handleEditTable">
                    <el-icon><EditPen /></el-icon>
                    编辑表
                  </div>
                  <div class="menu-item danger" @click="handleDropTable">
                    <el-icon><Delete /></el-icon>
                    删除表
                  </div>
                </div>
              <div v-show="dbMenu.visible" :style="{ top: dbMenu.y + 'px', left: dbMenu.x + 'px' }" class="tree-context-menu" @click.stop>
                <div class="menu-item" @click="handleCreateSiblingDatabaseSql">
                  <el-icon><Folder /></el-icon>
                  创建数据库
                </div>
                <div class="menu-item" @click="handleCreateTableInDbSql">
                  <el-icon><DocumentAdd /></el-icon>
                  创建数据表
                </div>
                <div class="menu-item" @click="onDbMenuImportSql">
                  <el-icon><Upload /></el-icon>
                  导入SQL文件
                </div>
                <div class="menu-item" @click="handleExportDatabaseSqlWithDataAsync">
                  <el-icon><Document /></el-icon>
                  导出SQL (结构+数据)
                </div>
                <div class="menu-item" @click="handleExportDatabaseSqlStructure">
                  <el-icon><Document /></el-icon>
                  导出SQL (结构)
                </div>
              </div>
              </teleport>

              <!-- JSON Dialog -->
              <el-dialog v-model="jsonDialog.visible" title="当前行 JSON" width="600px" append-to-body>
                <el-input type="textarea" :rows="12" v-model="jsonDialog.text" readonly />
                <template #footer>
                  <el-button @click="jsonDialog.visible = false">关闭</el-button>
                  <el-button type="primary" @click="copyJson">
                    <el-icon><DocumentCopy /></el-icon>
                    复制 JSON
                  </el-button>
                </template>
              </el-dialog>

              <!-- Hidden file input for DB-level SQL import -->
              <input ref="hiddenSqlFileInput" type="file" accept=".sql" @change="handleHiddenSqlFileChange" style="display:none" />

              <!-- Row Form Dialog (read-only) -->
              <el-dialog v-model="rowFormDialog.visible" :title="rowFormDialog.tab && rowFormDialog.tab.node ? `表单视图: ${rowFormDialog.tab.node.tableName}` : '表单视图'" width="600px" append-to-body>
                <div class="insert-form-container">
                  <el-form :model="rowFormDialog.data" label-width="120px" size="small">
                    <el-form-item v-for="col in rowFormDialog.columns" :key="col.columnName" :label="col.columnName">
                      <el-input v-model="rowFormDialog.data[col.columnName]" readonly />
                      <el-text v-if="sanitizeComment(col.columnComment || col.comment)" type="info" size="small" class="field-hint">
                        ({{ sanitizeComment(col.columnComment || col.comment) }})
                      </el-text>
                    </el-form-item>
                  </el-form>
                </div>
                <template #footer>
                  <el-button @click="rowFormDialog.visible = false">关闭</el-button>
                  <el-button @click="copyRowInsertSql" type="primary" plain>
                    <el-icon><DocumentAdd /></el-icon>
                    复制 INSERT
                  </el-button>
                  <el-button @click="copyRowUpdateSql" type="primary" :disabled="!hasPrimaryKey(rowFormDialog.columns)">
                    <el-icon><EditPen /></el-icon>
                    复制 UPDATE
                  </el-button>
                </template>
              </el-dialog>
            </el-main>
          </el-container>

    <!-- Table Structure Dialog -->
    <teleport to="body">
      <el-dialog
        v-model="tableStructureDialog.visible"
        :title="`表结构: ${tableStructureDialog.tableName}`"
        width="900px"
        :z-index="30000"
        class="table-structure-dialog"
      >
        <div class="sql-display-container">
          <pre class="sql-code">{{ tableStructureDialog.sql }}</pre>
        </div>
        <template #footer>
          <el-button @click="tableStructureDialog.visible = false">关闭</el-button>
          <el-button type="primary" @click="copyTableStructureSql">
            <el-icon><DocumentCopy /></el-icon>
            复制
          </el-button>
        </template>
      </el-dialog>
    </teleport>

    <!-- Create Table SQL Dialog -->
    <teleport to="body">
      <el-dialog
        v-model="createTableSqlDialog.visible"
        title="创建数据表 (SQL)"
        width="900px"
        :z-index="30000"
        class="table-structure-dialog"
      >
        <div class="sql-display-container">
          <CodeEditor v-model="createTableSqlDialog.sql" language="sql" height="450px" />
        </div>
        <template #footer>
          <el-button @click="createTableSqlDialog.visible = false">取消</el-button>
          <el-button type="primary" :loading="createTableSqlDialog.loading" @click="confirmCreateTableSql">执行</el-button>
        </template>
      </el-dialog>
    </teleport>

    <!-- Create Database SQL Dialog -->
    <teleport to="body">
      <el-dialog
        v-model="createDbSqlDialog.visible"
        title="创建数据库 (SQL)"
        width="900px"
        :z-index="30000"
        class="table-structure-dialog"
      >
        <div class="sql-display-container">
          <CodeEditor v-model="createDbSqlDialog.sql" language="sql" height="300px" />
        </div>
        <template #footer>
          <el-button @click="createDbSqlDialog.visible = false">取消</el-button>
          <el-button type="primary" :loading="createDbSqlDialog.loading" @click="confirmCreateDatabaseSql">执行</el-button>
        </template>
      </el-dialog>
    </teleport>

        </div>
</template>

<script>
import { dbAPI } from '@/api/database'
import { Refresh, Search, Folder, Document, EditPen, VideoPlay, InfoFilled, Plus, Check, Upload, DArrowRight, DocumentCopy, Delete, Download, Tickets, DocumentAdd, MagicStick, View } from '@element-plus/icons-vue'
import CodeEditor from '@/components/CodeEditor.vue'

export default {
  components: { Refresh, Search, Folder, Document, EditPen, VideoPlay, InfoFilled, CodeEditor, Plus, Check, Upload, DArrowRight, DocumentCopy, Delete, Download, Tickets, DocumentAdd, MagicStick, View },
  data() {
    return {
      treeFilter: '',
      loadingTree: false,
      treeData: [],
      treeProps: { children: 'children', label: 'label', isLeaf: 'isLeaf' },
      selectedNode: null,
      tableData: [],
      tableColumns: [],
      loadingData: false,
      pagination: { currentPage: 1, pageSize: 20, total: 0 },
      editingCell: null,
      tableKeys: [],
      currentInstance: null,
      openTabs: [],
      activeTab: '',
      // Left aside width (resizable)
      asideWidth: 300,
      asideResizing: false,
      _asideResizeStart: { x: 0, w: 300 },
      // Insert dialog state
      insertDialogVisible: false,
      insertFormData: {},
      insertFormColumns: [],
      insertFormTab: null,
      inserting: false,
      uploadingScript: false,
      // When importing from DB context menu
      sqlImportTargetNode: null,
      // SQL import (async job) dialog state
      sqlImportDialog: {
        visible: false,
        jobId: '',
        status: '', // queued|running|completed|failed
        total: 0,
        success: 0,
        failed: 0,
        errors: [],
        startedAt: '',
        endedAt: '',
        timer: null,
        relatedTab: null
      },
      // Export (async job) dialog state
      exportJobDialog: {
        visible: false,
        jobId: '',
        status: '', // queued|running|completed|failed
        tables: 0,
        doneTables: 0,
        fileSize: 0,
        timer: null,
        node: null
      },
      // Context menu
      contextMenu: { visible: false, x: 0, y: 0, row: null, tab: null },
      // JSON dialog
      jsonDialog: { visible: false, text: '' },
      // Row form dialog
      rowFormDialog: { visible: false, data: {}, columns: [], tab: null, saving: false },
      // Tree (table) context menu
      treeMenu: { visible: false, x: 0, y: 0, node: null },
      // DB (database) context menu
      dbMenu: { visible: false, x: 0, y: 0, node: null },
      // Copy table dialog
      copyTableDialog: { visible: false, oldName: '', newName: '', copyData: true, loading: false, node: null },
      // Rename table dialog
      renameDialog: { visible: false, oldName: '', newName: '', loading: false, node: null },
      // Edit table dialog (DDL)
      editTableDialog: { visible: false, tableName: '', ddl: '', loading: false, node: null, columns: [], activeTab: 'columns' },
      // Create database dialog
      createDbDialog: { visible: false, dbName: '', charset: '', collate: '', ifNotExists: true, loading: false, node: null },
      // Create table SQL dialog
      createTableSqlDialog: { visible: false, sql: '', loading: false, node: null },
      // Create database SQL dialog
      createDbSqlDialog: { visible: false, sql: '', loading: false, node: null },
      // Table structure dialog
      tableStructureDialog: { visible: false, tableName: '', sql: '' },
    }
  },
  watch: {
    treeFilter(val) {
      this.$refs.treeRef.filter(val)
    }
  },
  created() {
    this.loadInstances()
  },
  beforeUnmount() {
    this.clearSqlImportPolling()
    this.clearExportJobPolling()
  },
  methods: {
    getDbIcon(type) {
      const iconMap = {
        'mysql': require('@/assets/image/mysqldb.svg'),
        'postgres': require('@/assets/image/PostgreSQL.svg')
      }
      return iconMap[type] || iconMap['mysql']
    },
    filterNode(value, data) {
      if (!value) return true
      return data.label.toLowerCase().includes(value.toLowerCase())
    },
    async loadInstances() {
      this.loadingTree = true
      try {
        const res = await dbAPI.getInstanceList({ pageNum: 1, pageSize: 100 })
        if (res.data?.code === 200) {
          let instances = res.data?.data
          if (!Array.isArray(instances)) {
            if (instances && Array.isArray(instances.list)) instances = instances.list
            else if (instances && Array.isArray(instances.data)) instances = instances.data
            else instances = []
          }
          // 仅保留支持在终端操作的数据库类型（MySQL/PostgreSQL）
          const supported = ['mysql', 'postgres']
          this.treeData = instances
            .filter(inst => supported.includes((inst.dbType || inst.type || '').toLowerCase()))
            .map(inst => ({
            id: `inst-${inst.id}`,
            label: inst.name,
            type: 'instance',
            dbType: (inst.dbType || inst.type || 'mysql').toLowerCase(),
            instanceId: inst.id,
            children: [],
            isLeaf: false
          }))
        }
      } catch (error) {
        this.$message.error('加载实例列表失败')
      } finally {
        this.loadingTree = false
      }
    },
    async handleNodeClick(data) {
      this.selectedNode = data
      if (data.type === 'instance') {
        this.currentInstance = data
        if (data.children.length === 0) {
          await this.loadDatabases(data)
        }
      } else if (data.type === 'database') {
        this.currentInstance = this.treeData.find(inst => inst.instanceId === data.instanceId)
        if (data.children.length === 0) {
          await this.loadTables(data)
        }
      } else if (data.type === 'table') {
        this.currentInstance = this.treeData.find(inst => inst.instanceId === data.instanceId)
        this.openTableTab(data)
      }
    },

    async openTableTab(node) {
      const tabName = `table-${node.instanceId}-${node.databaseName}-${node.tableName}`
      const existingTab = this.openTabs.find(tab => tab.name === tabName)
      if (existingTab) {
        this.activeTab = tabName
      } else {
        let instanceName = this.currentInstance.label
        let connectionAddress = ''
        try {
          const instanceDetailRes = await dbAPI.getInstanceById(node.instanceId)
          if (instanceDetailRes.data?.code === 200 && instanceDetailRes.data?.data) {
            instanceName = instanceDetailRes.data.data.name || instanceName
            connectionAddress = instanceDetailRes.data.data.host || ''
          }
        } catch (error) {
          console.error('Failed to fetch instance details:', error)
          this.$message.error('加载实例详情失败')
        }

        const newTab = {
          name: tabName,
          title: node.tableName,
          node: node,
          connection: {
            instanceName: instanceName,
            connectionAddress: connectionAddress,
            dbName: node.databaseName,
          },
          tableData: [],
          tableColumns: [],
          tableKeys: [],
          loadingData: false,
          pagination: { currentPage: 1, pageSize: 20, total: 0 },
          modifiedRows: new Set(), // Track modified rows
          renderTable: false, // Flag to control el-table rendering
        }
        this.openTabs.push(newTab)
        this.activeTab = tabName
        
        // Ensure we pass the reactive proxy to the load function
        const reactiveTab = this.openTabs[this.openTabs.length - 1];
        this.loadTableData(reactiveTab);
      }
    },

    handleTabRemove(targetName) {
      const tabs = this.openTabs
      let activeName = this.activeTab
      if (activeName === targetName) {
        tabs.forEach((tab, index) => {
          if (tab.name === targetName) {
            const nextTab = tabs[index + 1] || tabs[index - 1]
            if (nextTab) {
              activeName = nextTab.name
            }
          }
        })
      }
      this.activeTab = activeName
      this.openTabs = tabs.filter(tab => tab.name !== targetName)
    },

    async loadDatabases(node) {
      try {
        // Request with stats and a global timeout to avoid long blocking
        const res = await dbAPI.getInstanceDatabases(node.instanceId, { withStats: true, statsTimeoutSeconds: 30 })
        if (res.data?.code === 200) {
          const list = res.data.data || []
          node.children = list.map(item => {
            const isObj = item && typeof item === 'object'
            const name = isObj ? (item.name || item.db || item.database || '') : String(item)
            const tables = isObj ? (item.tables ?? item.tableCount) : undefined
            const sizeStr = isObj ? (item.size || (typeof item.sizeBytes === 'number' ? `${(item.sizeBytes / (1024*1024)).toFixed(1)} MB` : undefined)) : undefined
            const comment = isObj ? [
              typeof tables === 'number' ? `表: ${tables}` : null,
              sizeStr ? `大小: ${sizeStr}` : null
            ].filter(Boolean).join(' | ') : ''
            return {
              id: `db-${node.instanceId}-${name}`,
              label: name,
              type: 'database',
              instanceId: node.instanceId,
              databaseName: name,
              comment,
              children: [],
              isLeaf: false
            }
          })
        }
      } catch (error) {
        this.$message.error('加载数据库失败')
      }
    },
    async loadTables(node) {
      try {
        const res = await dbAPI.getTables(node.instanceId, node.databaseName)
        if (res.data?.code === 200) {
          node.children = (res.data.data || []).map(table => ({
            id: `tbl-${node.instanceId}-${node.databaseName}-${table.tableName}`,
            label: table.tableName,
            type: 'table',
            comment: table.tableComment,
            size: table.size,
            instanceId: node.instanceId,
            databaseName: node.databaseName,
            tableName: table.tableName,
            isLeaf: true
          }))
        }
      } catch (error) {
        this.$message.error('加载表列表失败')
      }
    },
    async loadTableData(tab) {
      // Guard clause to ensure this is only called for valid table tabs
      if (!tab || !tab.node || tab.node.type !== 'table') {
        console.error("loadTableData called with an invalid tab object:", tab)
        return
      }

      tab.loadingData = true
      const node = tab.node
      try {
        const offset = (tab.pagination.currentPage - 1) * tab.pagination.pageSize
        const sql = `SELECT * FROM ${node.tableName} LIMIT ${tab.pagination.pageSize} OFFSET ${offset}`
        const res = await dbAPI.execSQL({ dbId: node.instanceId, db: node.databaseName, sql })

    if (res.data?.code === 200 && res.data?.data) {
      const { columns, rows, columnMetas } = res.data.data
      if (Array.isArray(columnMetas) && columnMetas.length > 0) {
        tab.tableColumns = columnMetas.map(m => ({
          columnName: m.name,
          columnComment: m.comment,
          type: m.type,
          // Try best-effort PK detection; fallback to 'id'
          isPrimaryKey: !!(m.isPrimaryKey || m.pk || m.primary || m.primaryKey || m.columnKey === 'PRI' || m.name === 'id')
        }))
      } else {
        tab.tableColumns = columns ? columns.map(col => ({ columnName: col, isPrimaryKey: col === 'id' })) : []
      }
      tab.tableData = rows || []
    }

    // Fetch total count for pagination using approx-count API
    const approxCountRes = await dbAPI.getApproxCount(node.instanceId, node.databaseName, node.tableName)
    if (approxCountRes.data?.code === 200 && approxCountRes.data?.data) {
      tab.pagination.total = approxCountRes.data.data.rows
    } else {
      // Fallback if approx-count fails or is not available
      tab.pagination.total = tab.tableData.length
    }
  } catch (error) {
    console.error('Failed to load table data in try block:', error);
    this.$message.error('加载表数据失败')
  } finally {
    tab.loadingData = false
  }
    },
    refreshTableData(tab) {
      this.loadTableData(tab)
    },
    // Resize aside (left tree panel)
    startAsideResize(e) {
      e.preventDefault()
      this.asideResizing = true
      this._asideResizeStart = { x: e.clientX, w: this.asideWidth }
      const onMove = (evt) => {
        if (!this.asideResizing) return
        const dx = evt.clientX - this._asideResizeStart.x
        let w = this._asideResizeStart.w + dx
        const min = 200
        const max = 600
        if (w < min) w = min
        if (w > max) w = max
        this.asideWidth = w
      }
      const onUp = () => {
        this.asideResizing = false
        window.removeEventListener('mousemove', onMove)
        window.removeEventListener('mouseup', onUp)
      }
      window.addEventListener('mousemove', onMove)
      window.addEventListener('mouseup', onUp)
    },
    isNullish(val) {
      return val === null || val === undefined || val === ''
    },
    onRowContextMenu(row, column, event, tab) {
      event.preventDefault()
      // Position relative to main content container to avoid transform/scroll issues
      const menuW = 200
      const menuH = 100
      const gap = 2
      const containerEl = this.$refs.mainContent?.$el || this.$refs.mainContent
      const rect = containerEl?.getBoundingClientRect?.() || { left: 0, top: 0, width: window.innerWidth, height: window.innerHeight }
      const clientX = event.clientX
      const clientY = event.clientY
      let x = clientX - rect.left + gap
      let y = clientY - rect.top + gap
      if (x + menuW > rect.width) x = clientX - rect.left - menuW - gap
      if (y + menuH > rect.height) y = clientY - rect.top - menuH - gap
      this.contextMenu.visible = true
      this.contextMenu.x = Math.max(0, x)
      this.contextMenu.y = Math.max(0, y)
      this.contextMenu.row = row
      this.contextMenu.tab = tab
      const hide = () => {
        this.contextMenu.visible = false
        window.removeEventListener('click', hide)
        window.removeEventListener('scroll', hide, true)
        window.removeEventListener('resize', hide)
      }
      window.addEventListener('click', hide)
      window.addEventListener('scroll', hide, true)
      window.addEventListener('resize', hide)
    },
    onTreeContextMenu(event, data) {
      try { console.log('[DB-Terminal] onTreeContextMenu', { type: data?.type, db: data?.databaseName, table: data?.tableName }); } catch (_) {}
      try { console.log('[DB-Terminal] onTreeContextMenu', { type: data?.type, label: data?.label, table: data?.tableName, db: data?.databaseName }); } catch (_) {}
      if (!data) return
      const gap = 2
      const menuW = 220
      const menuH = 160
      const vw = window.innerWidth
      const vh = window.innerHeight
      let x = (event.pageX ?? event.clientX) + gap
      let y = (event.pageY ?? event.clientY) + gap
      if (x + menuW > vw) x = (event.pageX ?? event.clientX) - menuW - gap
      if (y + menuH > vh) y = (event.pageY ?? event.clientY) - menuH - gap
      const showAt = (menu) => {
        this.treeMenu.visible = false
        this.dbMenu.visible = false
        menu.visible = true
        menu.x = Math.max(0, x)
        menu.y = Math.max(0, y)
        menu.node = data
        const hide = () => {
          menu.visible = false
          window.removeEventListener('click', hide)
          window.removeEventListener('scroll', hide, true)
          window.removeEventListener('resize', hide)
        }
        window.addEventListener('click', hide)
        window.addEventListener('scroll', hide, true)
        window.addEventListener('resize', hide)
      }
      if (data.type === 'table') {
        showAt(this.treeMenu)
      } else if (data.type === 'database') {
        showAt(this.dbMenu)
      }
    },
    async handleCopyTable() {
      const n = this.treeMenu.node
      if (!n) return
      this.treeMenu.visible = false
      // 打开对话框，默认复制数据
      this.copyTableDialog.visible = true
      this.copyTableDialog.node = n
      this.copyTableDialog.oldName = n.tableName
      this.copyTableDialog.newName = `${n.tableName}_copy`
      this.copyTableDialog.copyData = true
      this.copyTableDialog.structureOnly = false
      this.copyTableDialog.loading = false
    },
    async handleCopyTableStructure() {
      const n = this.treeMenu.node
      if (!n) return
      this.treeMenu.visible = false

      const loading = this.$loading({
        lock: true,
        text: '正在获取表结构...',
        background: 'rgba(0, 0, 0, 0.7)'
      })

      try {
        const res = await dbAPI.getTableDDL(n.instanceId, n.databaseName, n.tableName)
        console.log('getTableDDL response:', res)

        if (res.data?.code === 200) {
          const ddl = res.data?.data?.ddl || res.data?.data?.DDL || ''
          console.log('DDL:', ddl)

          if (ddl) {
            this.tableStructureDialog.tableName = n.tableName
            this.tableStructureDialog.sql = ddl
            this.tableStructureDialog.visible = true
            console.log('tableStructureDialog:', this.tableStructureDialog)

            this.$nextTick(() => {
              console.log('After nextTick, visible:', this.tableStructureDialog.visible)
            })
          } else {
            this.$message.error('未获取到表结构SQL')
          }
        } else {
          this.$message.error(res.data?.message || '获取表结构失败')
        }
      } catch (e) {
        console.error('获取表结构失败:', e)
        this.$message.error('获取表结构失败: ' + (e.message || ''))
      } finally {
        loading.close()
      }
    },
    copyTableStructureSql() {
      this.copyToClipboard(this.tableStructureDialog.sql)
    },
    async handleDropTable() {
      const n = this.treeMenu.node
      if (!n) return
      try {
        await this.$confirm(`确认删除表 ${n.tableName} 吗？此操作不可恢复！`, '警告', { type: 'warning' })
      } catch (_) {
        this.treeMenu.visible = false
        return
      }
      try {
        const sql = `DROP TABLE ${n.tableName}`
        const res = await dbAPI.execSQL({ dbId: n.instanceId, db: n.databaseName, sql })
        if (res.data?.code === 200) {
          this.$message.success('删除成功')
          await this.reloadDatabaseTables(n.instanceId, n.databaseName)
        } else {
          this.$message.error(res.data?.message || '删除失败')
        }
      } catch (e) {
        this.$message.error('删除失败')
      } finally {
        this.treeMenu.visible = false
      }
    },
    async handleRenameTable() {
      const n = this.treeMenu.node
      if (!n) return
      this.treeMenu.visible = false
      try {
        const { value: newName } = await this.$prompt('请输入新表名', '重命名表', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          inputValue: `${n.tableName}_new`,
          inputValidator: (val) => (/^(?:[A-Za-z0-9_]+)$/.test(val) && val !== n.tableName) || '需不同且仅允许字母/数字/下划线'
        })
        if (!newName) return
        const dbType = this.getInstanceDbType(n.instanceId)
        const sql = dbType === 'postgres'
          ? `ALTER TABLE ${n.tableName} RENAME TO ${newName}`
          : `RENAME TABLE ${n.tableName} TO ${newName}`
        const res = await dbAPI.execSQL({ dbId: n.instanceId, db: n.databaseName, sql })
        if (res.data?.code === 200) {
          this.$message.success('重命名成功')
          await this.reloadDatabaseTables(n.instanceId, n.databaseName)
        } else {
          this.$message.error(res.data?.message || '重命名失败')
        }
      } catch (_) {
        // cancelled
      }
    },
    async handleEditTable() {
      try { console.log('[DB-Terminal] handleEditTable'); } catch (_) {}
      try { console.log('[DB-Terminal] handleEditTable click'); } catch (_) {}
      const n = this.treeMenu.node
      if (!n) return
      this.treeMenu.visible = false
      this.editTableDialog.visible = true
      this.editTableDialog.node = n
      this.editTableDialog.tableName = n.tableName
      this.editTableDialog.loading = true
      this.editTableDialog.activeTab = 'columns'
      try {
        // Fetch DDL
        const ddlRes = await dbAPI.getTableDDL(n.instanceId, n.databaseName, n.tableName)
        try { console.log('[DB-Terminal] getTableDDL resp', { code: ddlRes?.data?.code }); } catch (_) {}
        if (ddlRes.data?.code === 200) {
          const d = ddlRes.data?.data
          this.editTableDialog.ddl = (d && (d.ddl || d.DDL)) || ''
        } else {
          this.editTableDialog.ddl = ''
        }
      } catch (e) {
        this.editTableDialog.ddl = ''
      }
      try {
        // Fetch columns
        const colRes = await dbAPI.getTableColumns(n.instanceId, n.databaseName, n.tableName)
        try { console.log('[DB-Terminal] getTableColumns resp', { code: colRes?.data?.code, count: (colRes?.data?.data || []).length }); } catch (_) {}
        if (colRes.data?.code === 200) {
          const cols = colRes.data?.data || []
          const mapped = (cols || []).map(x => ({
            name: x.name || x.columnName,
            columnName: x.columnName || x.name,
            type: x.columnType || x.type || x.dataType,
            nullable: typeof x.isNullable === 'string' ? x.isNullable === 'YES' : !!(x.nullable ?? x.isNullable),
            default: x.defaultValue ?? x.default ?? x.columnDefault ?? '',
            comment: x.columnComment ?? x.comment ?? '',
            autoIncrement: (x.extra && String(x.extra).toLowerCase().includes('auto_increment')) || !!(x.autoIncrement ?? x.isAutoIncrement),
            pk: (x.columnKey === 'PRI') || !!(x.pk ?? x.isPrimaryKey),
          }))
          this.editTableDialog.columns = mapped
          // snapshot for diff
          this.editTableDialog.originalColumns = JSON.parse(JSON.stringify(mapped))
        } else {
          this.editTableDialog.columns = []
        }
      } catch (e) {
        this.editTableDialog.columns = []
      } finally {
        this.editTableDialog.loading = false
      }
      // menu already hidden
    },
    async applyEditTable() {
      const d = this.editTableDialog
      if (!d.node) return
      // If user is on columns tab and mysql, generate ALTERs automatically
      const dbType = this.getInstanceDbType(d.node.instanceId)
      if (d.activeTab === 'columns' && dbType !== 'postgres') {
        const sql = this.buildAlterFromColumns(d)
        if (!sql) {
          this.$message.info('没有检测到结构变化')
          return
        }
        d.ddl = sql
        d.activeTab = 'ddl'
      }
      if (!d.ddl || !d.ddl.trim()) {
        this.$message.warning('DDL 内容为空')
        return
      }
      d.loading = true
      try {
        const res = await dbAPI.execSQL({ dbId: d.node.instanceId, db: d.node.databaseName, sql: d.ddl })
        if (res.data?.code === 200) {
          this.$message.success('执行成功')
          d.visible = false
          await this.reloadDatabaseTables(d.node.instanceId, d.node.databaseName)
        } else {
          this.$message.error(res.data?.message || '执行失败')
        }
      } catch (e) {
        this.$message.error('执行失败')
      } finally {
        d.loading = false
      }
    },
    generateAlterSQL() {
      const d = this.editTableDialog
      if (!d.node) return
      const dbType = this.getInstanceDbType(d.node.instanceId)
      if (dbType === 'postgres') {
        this.$message.warning('PostgreSQL 请直接在 DDL 页签编辑')
        d.activeTab = 'ddl'
        return
      }
      const sql = this.buildAlterFromColumns(d)
      d.ddl = sql || '-- 无结构变化'
      d.activeTab = 'ddl'
    },
    addEditColumn() {
      this.editTableDialog.columns.push({ columnName: '', type: '', nullable: true, default: '', pk: false, autoIncrement: false, comment: '' })
    },
    removeEditColumn(idx) {
      this.editTableDialog.columns.splice(idx, 1)
    },
    buildAlterFromColumns(dlg) {
      const table = dlg.tableName
      const cur = dlg.columns || []
      const orig = dlg.originalColumns || []
      const byName = (arr) => {
        const m = new Map()
        arr.forEach(c => m.set(c.columnName, c))
        return m
      }
      const mOrig = byName(orig)
      const mCur = byName(cur)
      const parts = []
      // Added columns
      for (const c of cur) {
        if (!mOrig.has(c.columnName)) {
          parts.push(`ADD COLUMN ${this.columnDef(c)}`)
        }
      }
      // Dropped columns
      for (const c of orig) {
        if (!mCur.has(c.columnName)) {
          parts.push(`DROP COLUMN ${c.columnName}`)
        }
      }
      // Modified columns
      for (const c of cur) {
        const o = mOrig.get(c.columnName)
        if (!o) continue
        if (this.colChanged(o, c)) {
          parts.push(`MODIFY COLUMN ${this.columnDef(c)}`)
        }
      }
      // Primary key changes
      const names = (arr) => arr.filter(c => c.pk).map(c => c.columnName).join(', ')
      const pkOrig = names(orig)
      const pkCur = names(cur)
      if (pkOrig !== pkCur) {
        if (pkOrig) parts.push('DROP PRIMARY KEY')
        if (pkCur) parts.push(`ADD PRIMARY KEY (${pkCur})`)
      }
      if (!parts.length) return ''
      return `ALTER TABLE ${table} \n  ${parts.join(',\n  ')};`
    },
    colChanged(a, b) {
      return (
        (a.type || '') !== (b.type || '') ||
        !!a.nullable !== !!b.nullable ||
        (String(a.default ?? '') !== String(b.default ?? '')) ||
        !!a.autoIncrement !== !!b.autoIncrement ||
        !!a.pk !== !!b.pk ||
        (a.comment || '') !== (b.comment || '')
      )
    },
    columnDef(c) {
      let s = `${c.columnName} ${c.type || 'varchar(255)'} `
      const isPk = !!c.pk
      const nullable = isPk ? false : !!c.nullable
      s += nullable ? 'NULL ' : 'NOT NULL '
      if (c.autoIncrement) s += 'AUTO_INCREMENT '
      if (c.default !== undefined && c.default !== '') {
        const val = this.formatDefault(c.default)
        s += `DEFAULT ${val} `
      }
      if (c.comment) {
        const esc = String(c.comment).replace(/'/g, "''")
        s += `COMMENT '${esc}' `
      }
      return s.trim()
    },
    formatDefault(v) {
      if (v === null || v === 'NULL') return 'NULL'
      if (typeof v === 'number') return String(v)
      const low = String(v).toLowerCase()
      if (low === 'null') return 'NULL'
      if (low === 'current_timestamp' || low === 'now()') return low.toUpperCase()
      if (/^\d+(\.\d+)?$/.test(String(v))) return String(v)
      return `'${String(v).replace(/'/g, "''")}'`
    },
    handleCreateSiblingDatabase() {
      const n = this.dbMenu.node
      if (!n) return
      try { console.log('[DB-Terminal] handleCreateSiblingDatabase', { inst: n?.instanceId, db: n?.databaseName }); } catch (_) {}
      this.dbMenu.visible = false
      this.createDbSqlDialog.node = n
      const dbType = this.getInstanceDbType(n.instanceId)
      const tmplMysql = `CREATE DATABASE IF NOT EXISTS your_database_name \n  DEFAULT CHARACTER SET utf8mb4;`
      const tmplPg = `CREATE DATABASE your_database_name;`
      this.createDbSqlDialog.sql = dbType === 'postgres' ? tmplPg : tmplMysql
      this.$nextTick(() => { this.createDbSqlDialog.visible = true })
    },
    async confirmCreateDatabaseSql() {
      const d = this.createDbSqlDialog
      try { console.log('[DB-Terminal] confirmCreateDatabaseSql'); } catch (_) {}
      if (!d.node) return
      if (!d.sql || !d.sql.trim()) {
        this.$message.warning('请输入创建数据库的 SQL')
        return
      }
      d.loading = true
      try {
        // Execute on current selected database connection; for MySQL this is acceptable
        const res = await dbAPI.execSQL({ dbId: d.node.instanceId, db: d.node.databaseName, sql: d.sql })
        if (res.data?.code === 200) {
          this.$message.success('执行成功')
          d.visible = false
          // reload databases under instance
          const instNode = this.treeData.find(i => i.instanceId === d.node.instanceId)
          if (instNode) await this.loadDatabases(instNode)
        } else {
          this.$message.error(res.data?.message || '执行失败')
        }
      } catch (e) {
        this.$message.error('执行失败')
      } finally {
        d.loading = false
      }
    },
    // Prompt-based database creation (SQL) - 改用弹框
    handleCreateSiblingDatabaseSql() {
      const n = this.dbMenu.node
      if (!n) return
      console.log('[DB-Terminal] handleCreateSiblingDatabaseSql', { node: n })
      this.dbMenu.visible = false
      this.createDbSqlDialog.node = n
      const dbType = this.getInstanceDbType(n.instanceId)
      const tmplMysql = `CREATE DATABASE IF NOT EXISTS your_database_name
  DEFAULT CHARACTER SET utf8mb4;`
      const tmplPg = `CREATE DATABASE your_database_name;`
      this.createDbSqlDialog.sql = dbType === 'postgres' ? tmplPg : tmplMysql
      console.log('[DB-Terminal] Opening createDbSqlDialog')
      this.$nextTick(() => {
        this.createDbSqlDialog.visible = true
        console.log('[DB-Terminal] createDbSqlDialog.visible set to true')
      })
    },
    handleCreateTableInDb() {
      const n = this.dbMenu.node
      if (!n) return
      try { console.log('[DB-Terminal] handleCreateTableInDb', { inst: n?.instanceId, db: n?.databaseName }); } catch (_) {}
      this.dbMenu.visible = false
      this.createTableSqlDialog.node = n
      const dbType = this.getInstanceDbType(n.instanceId)
      const tmplMysql = `CREATE TABLE ${n.databaseName}.your_table_name (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  -- your columns here
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
      const tmplPg = `CREATE TABLE ${n.databaseName}.your_table_name (
  id BIGSERIAL PRIMARY KEY,
  -- your columns here
);`
      this.createTableSqlDialog.sql = dbType === 'postgres' ? tmplPg : tmplMysql
      this.$nextTick(() => { this.createTableSqlDialog.visible = true })
    },
    async confirmCreateTableSql() {
      const d = this.createTableSqlDialog
      try { console.log('[DB-Terminal] confirmCreateTableSql'); } catch (_) {}
      if (!d.node) return
      if (!d.sql || !d.sql.trim()) {
        this.$message.warning('请输入创建表的 SQL')
        return
      }
      d.loading = true
      try {
        const res = await dbAPI.execSQL({ dbId: d.node.instanceId, db: d.node.databaseName, sql: d.sql })
        if (res.data?.code === 200) {
          this.$message.success('执行成功')
          d.visible = false
          await this.reloadDatabaseTables(d.node.instanceId, d.node.databaseName)
        } else {
          this.$message.error(res.data?.message || '执行失败')
        }
      } catch (e) {
        this.$message.error('执行失败')
      } finally {
        d.loading = false
      }
    },
    // Prompt-based table creation (SQL) - 改用弹框
    handleCreateTableInDbSql() {
      const n = this.dbMenu.node
      if (!n) return
      console.log('[DB-Terminal] handleCreateTableInDbSql', { node: n })
      this.dbMenu.visible = false
      this.createTableSqlDialog.node = n
      const dbType = this.getInstanceDbType(n.instanceId)
      const tmplMysql = `CREATE TABLE ${n.databaseName}.your_table_name (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  -- your columns here
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
      const tmplPg = `CREATE TABLE ${n.databaseName}.your_table_name (
  id BIGSERIAL PRIMARY KEY,
  -- your columns here
);`
      this.createTableSqlDialog.sql = dbType === 'postgres' ? tmplPg : tmplMysql
      console.log('[DB-Terminal] Opening createTableSqlDialog')
      this.$nextTick(() => {
        this.createTableSqlDialog.visible = true
        console.log('[DB-Terminal] createTableSqlDialog.visible set to true')
      })
    },
    getInstanceDbType(instanceId) {
      const inst = (this.treeData || []).find(i => i.instanceId === instanceId)
      return inst?.dbType || 'mysql'
    },
    async reloadDatabaseTables(instanceId, dbName) {
      // find database node and reload tables
      const inst = (this.treeData || []).find(i => i.instanceId === instanceId)
      if (!inst) return
      const dbNode = (inst.children || []).find(c => c.type === 'database' && c.databaseName === dbName)
      if (dbNode) {
        await this.loadTables(dbNode)
      }
    },
    handleViewJson() {
      const { row } = this.contextMenu
      this.jsonDialog.text = JSON.stringify(row || {}, null, 2)
      this.jsonDialog.visible = true
      this.contextMenu.visible = false
    },
    copyJson() {
      this.copyToClipboard(this.jsonDialog.text || '')
    },
    async handleViewRowForm() {
      const { row, tab } = this.contextMenu
      this.rowFormDialog.tab = tab
      // Use already loaded columns
      const columns = (tab.tableColumns || [])
      this.rowFormDialog.columns = columns
      this.rowFormDialog.data = { ...row }
      this.rowFormDialog.visible = true
      this.contextMenu.visible = false
    },
    handleExportCsv() {
      const { tab } = this.contextMenu
      const payload = { db: tab.node.databaseName, table: tab.node.tableName }
      dbAPI.exportCsv(tab.node.instanceId, payload).then(res => {
        const blob = new Blob([res.data], { type: 'text/csv;charset=utf-8;' })
        const file = `${tab.node.tableName}.csv`
        this.downloadBlob(blob, file)
      }).catch(() => {
        this.$message.error('导出CSV失败')
      }).finally(() => {
        this.contextMenu.visible = false
      })
    },
    handleExportSql() {
      const { tab } = this.contextMenu
      const payload = { db: tab.node.databaseName, table: tab.node.tableName }
      dbAPI.exportSql(tab.node.instanceId, payload).then(res => {
        const blob = new Blob([res.data], { type: 'application/sql;charset=utf-8;' })
        const file = `${tab.node.tableName}.sql`
        this.downloadBlob(blob, file)
      }).catch(() => {
        this.$message.error('导出SQL失败')
      }).finally(() => {
        this.contextMenu.visible = false
      })
    },
    async confirmCopyTable() {
      const d = this.copyTableDialog
      if (!d.node) return
      const newName = (d.newName || '').trim()
      if (!/^(?:[A-Za-z0-9_]+)$/.test(newName)) {
        this.$message.warning('新表名仅允许字母/数字/下划线')
        return
      }
      d.loading = true
      try {
        const payload = {
          dbId: d.node.instanceId,
          db: d.node.databaseName,
          tableName: d.node.tableName,
          newTable: newName,
          copyData: !!d.copyData && !d.structureOnly
        }
        const res = await dbAPI.copyTable(payload)
        if (res.data?.code === 200) {
          this.$message.success(d.structureOnly ? '复制表结构成功' : '复制表成功')
          d.visible = false
          await this.reloadDatabaseTables(d.node.instanceId, d.node.databaseName)
        } else {
          this.$message.error(res.data?.message || '复制失败')
        }
      } catch (e) {
        this.$message.error('复制失败')
      } finally {
        d.loading = false
      }
    },
    handleExportDatabaseSqlStructure() {
      const n = this.dbMenu.node
      if (!n) return
      this.dbMenu.visible = false
      dbAPI.exportDatabaseStructure(n.instanceId, { db: n.databaseName }).then(res => {
        const blob = new Blob([res.data], { type: 'application/sql;charset=utf-8;' })
        const ts = new Date().toISOString().replace(/[:T]/g, '-').slice(0,19)
        const file = `${n.databaseName}_structure_${ts}.sql`
        this.downloadBlob(blob, file)
      }).catch(() => {
        this.$message.error('导出SQL失败')
      })
    },
    handleExportDatabaseSqlWithData() {
      const n = this.dbMenu.node
      if (!n) return
      this.dbMenu.visible = false
      dbAPI.exportDatabaseFull(n.instanceId, { db: n.databaseName }).then(res => {
        const blob = new Blob([res.data], { type: 'application/sql;charset=utf-8;' })
        const ts = new Date().toISOString().replace(/[:T]/g, '-').slice(0,19)
        const file = `${n.databaseName}_full_${ts}.sql`
        this.downloadBlob(blob, file)
      }).catch(() => {
        this.$message.error('导出SQL失败')
      })
    },
    async handleDeleteRow() {
      const { row, tab } = this.contextMenu
      const pk = (tab.tableColumns.find(c => c.isPrimaryKey) || { columnName: 'id' }).columnName || 'id'
      const pkv = row[pk]
      try {
        await this.$confirm(`确认删除该行数据？${pk}=${pkv}`, '提示', { type: 'warning' })
      } catch (_) {
        this.contextMenu.visible = false
        return
      }
      try {
        const res = await dbAPI.deleteRow(tab.node.instanceId, tab.node.databaseName, tab.node.tableName, pk, pkv)
        if (res.data?.code === 200) {
          this.$message.success('删除成功')
          this.refreshTableData(tab)
        } else {
          this.$message.error(res.data?.message || '删除失败')
        }
      } catch (e) {
        this.$message.error('删除失败')
      } finally {
        this.contextMenu.visible = false
      }
    },
    hasPrimaryKey(columns) {
      return !!(columns || []).find(c => c.isPrimaryKey || c.columnName === 'id')
    },
    // Copy helpers
    quoteValue(val) {
      if (val === null) return 'NULL'
      if (val === undefined) return 'NULL'
      if (typeof val === 'number' || typeof val === 'bigint') return String(val)
      if (typeof val === 'boolean') return val ? 'TRUE' : 'FALSE'
      // default string
      return `'${String(val).replace(/'/g, "''")}'`
    },
    buildInsertSqlFromRow(tab, row, columns) {
      const cols = (columns || []).map(c => c.columnName)
      const values = cols.map(c => this.quoteValue(row[c]))
      return `INSERT INTO ${tab.node.tableName} (${cols.join(', ')}) VALUES (${values.join(', ')});`
    },
    buildUpdateSqlFromRow(tab, row, columns) {
      const pkCol = (columns.find(c => c.isPrimaryKey) || { columnName: 'id' }).columnName || 'id'
      const pairs = (columns || [])
        .filter(c => c.columnName !== pkCol)
        .map(c => `${c.columnName} = ${this.quoteValue(row[c.columnName])}`)
      const where = `${pkCol} = ${this.quoteValue(row[pkCol])}`
      return `UPDATE ${tab.node.tableName} SET ${pairs.join(', ')} WHERE ${where};`
    },
    async copyToClipboard(text) {
      try {
        if (navigator.clipboard && navigator.clipboard.writeText) {
          await navigator.clipboard.writeText(text)
        } else {
          const ta = document.createElement('textarea')
          ta.value = text
          ta.style.position = 'fixed'
          ta.style.opacity = '0'
          document.body.appendChild(ta)
          ta.focus()
          ta.select()
          document.execCommand('copy')
          document.body.removeChild(ta)
        }
        this.$message.success('已复制到剪贴板')
      } catch (e) {
        this.$message.error('复制失败')
      }
    },
    downloadBlob(blob, filename) {
      try {
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = filename
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        URL.revokeObjectURL(url)
      } catch (e) {
        this.$message.error('下载失败')
      }
    },
    copyRowInsertSql() {
      const sql = this.buildInsertSqlFromRow(this.rowFormDialog.tab, this.rowFormDialog.data, this.rowFormDialog.columns)
      this.copyToClipboard(sql)
    },
    copyRowUpdateSql() {
      const sql = this.buildUpdateSqlFromRow(this.rowFormDialog.tab, this.rowFormDialog.data, this.rowFormDialog.columns)
      this.copyToClipboard(sql)
    },
    sanitizeComment(comment) {
      if (!comment) return ''
      try {
        const text = String(comment).replace(/^'+|'+$/g, '')
        return text
      } catch (e) {
        return ''
      }
    },
    handleCellDoubleClick(row, column, tab) {
      // Check if the column is a primary key
      const primaryKey = (tab.tableColumns || []).find(c => c.isPrimaryKey)
      if (primaryKey && primaryKey.columnName === column.property) {
        this.$message.warning('主键列不能修改')
        return
      }
      this.editingCell = { row, column: column.property }
    },
    handleInputBlur(row, column, tab) {
      // Mark the row as modified, don't save immediately
      tab.modifiedRows.add(row);
      this.editingCell = null
    },

    async saveTableData(tab) {
      try {
        if (!tab || !tab.modifiedRows || tab.modifiedRows.size === 0) {
          this.$message.info('没有需要保存的修改')
          return
        }

        const pkCol = (tab.tableColumns.find(c => c.isPrimaryKey) || { columnName: 'id' }).columnName || 'id'
        const toSave = Array.from(tab.modifiedRows)

        if (!pkCol) {
          this.$message.error('无法识别主键列，无法保存')
          return
        }

        let success = 0
        let skipped = 0
        let failed = 0

        for (const row of toSave) {
          const pkv = row[pkCol]
          if (pkv === undefined || pkv === null) {
            skipped++
            continue
          }
          try {
            const sql = this.buildUpdateSqlFromRow(tab, row, tab.tableColumns)
            const res = await dbAPI.execSQL({ dbId: tab.node.instanceId, db: tab.node.databaseName, sql })
            if (res.data?.code === 200) success++
            else failed++
          } catch (e) {
            failed++
          }
        }

        tab.modifiedRows.clear()

        if (failed === 0) {
          const msg = skipped > 0 ? `保存成功 ${success} 条，跳过 ${skipped} 条（无主键）` : `保存成功 ${success} 条`
          this.$message.success(msg)
        } else {
          this.$message.error(`部分保存失败：成功 ${success}，失败 ${failed}，跳过 ${skipped}`)
        }

        this.refreshTableData(tab)
      } catch (e) {
        console.error('saveTableData error:', e)
        this.$message.error('保存失败')
      }
    },

    async submitNewRow(tab) {
      const data = tab.newRowData
      const columns = Object.keys(data)
      const values = Object.values(data).map(val => {
        // Simple quoting, might need improvement for different data types
        return typeof val === 'string' ? `'${val.replace(/'/g, "''")}'` : val
      })

      if (columns.length === 0) {
        return this.$message.warning('没有要插入的数据')
      }
      
      const sql = `INSERT INTO ${tab.node.tableName} (${columns.join(', ')}) VALUES (${values.join(', ')})`

      try {
        const res = await dbAPI.execSQL({ dbId: tab.node.instanceId, db: tab.node.databaseName, sql })
        if (res.data?.code === 200) {
          this.$message.success('新增成功！')
          // Close the insert tab
          this.handleTabRemove(tab.name)
          // Find and refresh the corresponding table tab
          const tableTabName = `table-${tab.node.instanceId}-${tab.node.databaseName}-${tab.node.tableName}`
          const tableTab = this.openTabs.find(t => t.name === tableTabName)
          if (tableTab) {
            this.refreshTableData(tableTab)
            this.activeTab = tableTabName // Switch back to the table tab
          }
        } else {
          this.$message.error(res.data?.message || '新增失败')
        }
      } catch (error) {
        this.$message.error('新增失败，请检查控制台')
        console.error("Insert failed:", error)
      }
    },

    openInsertTab(tab) {
      // Open dialog instead of a new tab
      this.insertFormTab = tab
      this.insertFormColumns = (tab.tableColumns || []).filter(c => c.columnName !== 'id' && !(c.columnName || '').endsWith('_at'))
      // Initialize form data
      const data = {}
      this.insertFormColumns.forEach(c => { data[c.columnName] = '' })
      this.insertFormData = data
      this.insertDialogVisible = true
    },

    cancelInsertForm() {
      this.insertDialogVisible = false
      this.inserting = false
    },

    async submitInsertForm() {
      const tab = this.insertFormTab
      if (!tab || !tab.node) {
        this.$message.error('无效的插入上下文')
        return
      }
      const data = this.insertFormData
      const columns = Object.keys(data).filter(k => data[k] !== undefined)
      const values = columns.map(col => {
        const val = data[col]
        return typeof val === 'string' ? `'${val.replace(/'/g, "''")}'` : val
      })
      if (columns.length === 0) {
        return this.$message.warning('没有要插入的数据')
      }
      const sql = `INSERT INTO ${tab.node.tableName} (${columns.join(', ')}) VALUES (${values.join(', ')})`
      this.inserting = true
      try {
        const res = await dbAPI.execSQL({ dbId: tab.node.instanceId, db: tab.node.databaseName, sql })
        if (res.data?.code === 200) {
          this.$message.success('新增成功！')
          this.insertDialogVisible = false
          this.inserting = false
          this.refreshTableData(tab)
        } else {
          this.$message.error(res.data?.message || '新增失败')
          this.inserting = false
        }
      } catch (error) {
        this.$message.error('新增失败，请检查控制台')
        console.error('Insert failed:', error)
        this.inserting = false
      }
    },

    beforeUploadSql(file) {
      const isSql = file.name.toLowerCase().endsWith('.sql')
      if (!isSql) {
        this.$message.warning('仅支持 .sql 文件')
        return false
      }
      const isLt10M = file.size / 1024 / 1024 < 10
      if (!isLt10M) {
        this.$message.warning('文件不能超过 10MB')
        return false
      }
      return true
    },
    async uploadSqlFile(tab, file) {
      if (this.uploadingScript) return
      this.uploadingScript = true
      try {
        const res = await dbAPI.execSqlFile(tab.node.instanceId, tab.node.databaseName, file)
        const data = res?.data || {}
        // Backend now returns 202 with job info
        if (data.code === 202 && data.data?.jobId) {
          const { jobId, status } = data.data
          this.startSqlImportJob(jobId, status, tab)
        } else if (data.code === 200) {
          // Backward compatibility (synchronous response)
          const d = data.data || {}
          this.$message.success(`执行完成：总数 ${d.total || 0}，成功 ${d.success || 0}，失败 ${d.failed || 0}`)
          this.refreshTableData(tab)
        } else {
          this.$message.error(data.message || '提交执行失败')
        }
      } catch (e) {
        this.$message.error('上传或提交失败')
      } finally {
        // Keep uploadingScript true if we entered polling, it will be cleared when job finishes
        if (!this.sqlImportDialog.visible) {
          this.uploadingScript = false
        }
      }
    },

    startSqlImportJob(jobId, status, tab) {
      // Initialize dialog and start polling
      this.clearSqlImportPolling()
      this.sqlImportDialog.visible = true
      this.sqlImportDialog.jobId = jobId
      this.sqlImportDialog.status = status || 'queued'
      this.sqlImportDialog.total = 0
      this.sqlImportDialog.success = 0
      this.sqlImportDialog.failed = 0
      this.sqlImportDialog.errors = []
      this.sqlImportDialog.startedAt = ''
      this.sqlImportDialog.endedAt = ''
      this.sqlImportDialog.relatedTab = tab
      // Begin polling every 2 seconds
      this.sqlImportDialog.timer = setInterval(() => {
        this.pollSqlImportStatus()
      }, 2000)
      // Trigger an immediate fetch
      this.pollSqlImportStatus()
    },

    async pollSqlImportStatus() {
      const jobId = this.sqlImportDialog.jobId
      if (!jobId) return
      try {
        const res = await dbAPI.getExecSqlFileStatus(jobId)
        const d = res?.data?.data || {}
        // Update dialog state
        this.sqlImportDialog.status = d.status || this.sqlImportDialog.status
        this.sqlImportDialog.total = Number(d.total || 0)
        this.sqlImportDialog.success = Number(d.success || 0)
        this.sqlImportDialog.failed = Number(d.failed || 0)
        this.sqlImportDialog.errors = Array.isArray(d.errors) ? d.errors : []
        this.sqlImportDialog.startedAt = d.startedAt || ''
        this.sqlImportDialog.endedAt = d.endedAt || ''

        // Terminal detection
        const total = this.sqlImportDialog.total
        const done = this.sqlImportDialog.success + this.sqlImportDialog.failed
        const finishedByCount = total > 0 && done >= total
        const ended = Number(this.sqlImportDialog.endedAt || 0) > 0
        // If finished by count or endedAt is set but status isn't terminal, coerce to completed (do not infer failed)
        if ((finishedByCount || ended) && !['completed', 'failed'].includes(this.sqlImportDialog.status)) {
          this.sqlImportDialog.status = 'completed'
        }

        if (['completed', 'failed'].includes(this.sqlImportDialog.status)) {
          // Stop polling
          this.clearSqlImportPolling()
          this.uploadingScript = false
          const msg = `执行完成：总数 ${this.sqlImportDialog.total}，成功 ${this.sqlImportDialog.success}，失败 ${this.sqlImportDialog.failed}`
          // Always show a green success toast for completion summary per requirement
          this.$message.success(msg)
          // Refresh table data if tab still present
          if (this.sqlImportDialog.relatedTab) {
            this.refreshTableData(this.sqlImportDialog.relatedTab)
          }
        }
      } catch (e) {
        // On polling error, stop polling
        this.$message.error('查询导入进度失败')
        this.clearSqlImportPolling()
        this.uploadingScript = false
      }
    },

    clearSqlImportPolling() {
      if (this.sqlImportDialog && this.sqlImportDialog.timer) {
        clearInterval(this.sqlImportDialog.timer)
        this.sqlImportDialog.timer = null
      }
    },

    computeImportPercent() {
      const t = this.sqlImportDialog.total || 0
      const done = (this.sqlImportDialog.success || 0) + (this.sqlImportDialog.failed || 0)
      if (t > 0) return Math.min(100, Math.round((done * 100) / t))
      // If finished but total unknown, mark 100
      if (['completed', 'failed'].includes(this.sqlImportDialog.status)) return 100
      return 0
    },

    onSqlImportDialogClose() {
      // Stop polling when user closes the dialog
      this.clearSqlImportPolling()
      this.uploadingScript = false
    },

    // ===== Database async export =====
    handleExportDatabaseSqlWithDataAsync() {
      const n = this.dbMenu.node
      if (!n) return
      this.dbMenu.visible = false
      dbAPI.exportDatabaseFullAsync(n.instanceId, n.databaseName, 'full').then(res => {
        console.log('[Export] API Response:', res)
        const data = res?.data || {}
        console.log('[Export] Response data:', data)

        // 兼容多种响应格式，优先使用 taskId
        let taskId, status

        // 格式1: { code: 200, data: { taskId, status } }
        if (data.code === 200 && data.data?.taskId) {
          taskId = data.data.taskId
          status = data.data.status
        }
        // 格式2: { code: 200, taskId, status }
        else if (data.code === 200 && data.taskId) {
          taskId = data.taskId
          status = data.status
        }
        // 格式3: 直接返回 { taskId, status }
        else if (data.taskId) {
          taskId = data.taskId
          status = data.status
        }
        // 兼容旧格式 jobId
        else if (data.data?.jobId) {
          taskId = data.data.jobId
          status = data.data.status
        }
        else if (data.jobId) {
          taskId = data.jobId
          status = data.status
        }

        if (taskId) {
          console.log('[Export] Starting export task:', taskId, status)
          this.startExportJob(taskId, status, n)
        } else {
          console.error('[Export] No taskId/jobId found in response:', data)
          this.$message.error(data.message || '启动导出任务失败')
        }
      }).catch(err => {
        console.error('[Export] API Error:', err)
        this.$message.error('启动导出任务失败: ' + (err.message || ''))
      })
    },
    startExportJob(jobId, status, node) {
      this.clearExportJobPolling()
      this.exportJobDialog.visible = true
      this.exportJobDialog.jobId = jobId
      this.exportJobDialog.status = status || 'queued'
      this.exportJobDialog.tables = 0
      this.exportJobDialog.doneTables = 0
      this.exportJobDialog.fileSize = 0
      this.exportJobDialog.node = node || null
      this.exportJobDialog.timer = setInterval(() => {
        this.pollExportJobStatus()
      }, 2000)
      this.pollExportJobStatus()
    },
    async pollExportJobStatus() {
      const jobId = this.exportJobDialog.jobId
      if (!jobId) return
      try {
        const res = await dbAPI.getExportJobStatus(jobId)
        const d = res?.data?.data || {}
        this.exportJobDialog.status = d.status || this.exportJobDialog.status
        this.exportJobDialog.tables = Number(d.tables || d.total || 0)
        this.exportJobDialog.doneTables = Number(d.doneTables || d.done || 0)
        this.exportJobDialog.fileSize = Number(d.fileSize || 0)

        const finishedByCount = this.exportJobDialog.tables > 0 && this.exportJobDialog.doneTables >= this.exportJobDialog.tables
        const ended = Number(d.endedAt || 0) > 0
        if ((finishedByCount || ended) && !['completed', 'failed'].includes(this.exportJobDialog.status)) {
          this.exportJobDialog.status = 'completed'
        }
        if (['completed', 'failed'].includes(this.exportJobDialog.status)) {
          this.clearExportJobPolling()
        }
      } catch (e) {
        this.clearExportJobPolling()
        this.$message.error('查询导出进度失败')
      }
    },
    clearExportJobPolling() {
      if (this.exportJobDialog && this.exportJobDialog.timer) {
        clearInterval(this.exportJobDialog.timer)
        this.exportJobDialog.timer = null
      }
    },
    onExportDialogClose() {
      this.clearExportJobPolling()
    },
    canDownloadExport() {
      return this.exportJobDialog.status === 'completed'
    },
    async downloadExportResult() {
      const jobId = this.exportJobDialog.jobId
      if (!jobId) return
      try {
        const res = await dbAPI.downloadExportJobFile(jobId)
        const node = this.exportJobDialog.node
        const ts = new Date().toISOString().replace(/[:T]/g, '-').slice(0,19)
        const filename = `${(node && node.databaseName) || 'database'}_full_${ts}.sql`
        this.downloadBlob(res.data, filename)
      } catch (e) {
        this.$message.error('下载失败')
      }
    },
    computeExportPercent() {
      const t = this.exportJobDialog.tables || 0
      const d = this.exportJobDialog.doneTables || 0
      if (t > 0) return Math.min(100, Math.round((d * 100) / t))
      if (['completed', 'failed'].includes(this.exportJobDialog.status)) return 100
      return 0
    },
    formatFileSize(size) {
      const s = Number(size || 0)
      if (s <= 0) return '0B'
      const units = ['B','KB','MB','GB','TB']
      let idx = 0
      let v = s
      while (v >= 1024 && idx < units.length - 1) { v /= 1024; idx++ }
      return `${v.toFixed(1)}${units[idx]}`
    },

    onDbMenuImportSql() {
      const n = this.dbMenu.node
      if (!n) return
      this.dbMenu.visible = false
      this.sqlImportTargetNode = n
      const inp = this.$refs.hiddenSqlFileInput
      if (inp) {
        try { inp.value = '' } catch (_) {}
        inp.click()
      }
    },

    handleHiddenSqlFileChange(e) {
      const file = e?.target?.files?.[0]
      if (!file) return
      if (!this.beforeUploadSql(file)) {
        try { e.target.value = '' } catch (_) {}
        return
      }
      const n = this.sqlImportTargetNode
      if (!n) return
      this.uploadSqlFileForNode(n, file)
      try { e.target.value = '' } catch (_) {}
    },

    async uploadSqlFileForNode(node, file) {
      if (this.uploadingScript) return
      this.uploadingScript = true
      try {
        const res = await dbAPI.execSqlFile(node.instanceId, node.databaseName, file)
        const data = res?.data || {}
        if (data.code === 202 && data.data?.jobId) {
          const { jobId, status } = data.data
          // Start job without binding to a table tab
          this.startSqlImportJob(jobId, status, null)
        } else if (data.code === 200) {
          const d = data.data || {}
          this.$message.success(`执行完成：总数 ${d.total || 0}，成功 ${d.success || 0}，失败 ${d.failed || 0}`)
        } else {
          this.$message.error(data.message || '提交执行失败')
        }
      } catch (e) {
        this.$message.error('上传或提交失败')
      } finally {
        if (!this.sqlImportDialog.visible) {
          this.uploadingScript = false
        }
      }
    },

    async openQueryTab(sourceTab) {
      // Each database gets its own query tab.
      const dbName = sourceTab.node.databaseName
      const instanceId = sourceTab.node.instanceId
      const tabName = `query-${instanceId}-${dbName}`

      const existingTab = this.openTabs.find(t => t.name === tabName)
      if (existingTab) {
        this.activeTab = tabName
        return
      }

      // Create a new query tab
              const newQueryTab = {
                name: tabName,
                title: `Query: ${dbName}`,
                type: 'query', // Differentiate from 'table' tabs
                node: sourceTab.node, // For context
                connection: sourceTab.connection,
                sqlQuery: `SELECT * FROM ${sourceTab.node.tableName} LIMIT 100;`,
                sqlResult: null,
                sqlCompletionItems: [],
                executing: false,
                formatting: false,
                preChecking: false,
                preCheckResult: null,
                // New performance optimization parameters
                useCache: false,
                cacheSeconds: 60,
                maxRows: 1000,
                                  timeoutSeconds: 60,
                                                };
                                        this.openTabs.push(newQueryTab);
                                        this.activeTab = tabName;
      // Fetch hints for autocompletion
      try {
        const hintRes = await dbAPI.getHintTables(instanceId, dbName)
        if (hintRes.data?.code === 200) {
          newQueryTab.sqlCompletionItems = hintRes.data.data
        }
      } catch (error) {
        console.error('Failed to load SQL hints:', error)
        this.$message.warning('加载SQL自动补全提示失败')
      }
    },

        async executeSql(tab) {

          if (!tab || tab.type !== 'query') {

            this.$message.error('无效的查询标签页')

            return

          }

          if (!tab.sqlQuery.trim()) {

            this.$message.warning('请输入SQL')

            return

          }



          tab.executing = true

          const start = Date.now()

                    try {

                      const res = await dbAPI.execSQL({

                        dbId: tab.node.instanceId,

                        db: tab.node.databaseName,

                        sql: tab.sqlQuery.trim(),

                        // New performance parameters

                        ...(tab.useCache && { useCache: tab.useCache }),

                        ...(tab.useCache && tab.cacheSeconds && { cacheSeconds: tab.cacheSeconds }),

                        ...(tab.maxRows && { maxRows: tab.maxRows }),

                        ...(tab.timeoutSeconds && { timeoutSeconds: tab.timeoutSeconds }),

                      })

                      if (res.data?.code === 200) {

              const data = res.data.data

              tab.sqlResult = {

                success: true,

                execTime: data.execTime || (Date.now() - start),

                affectedRows: data.affectedRows || 0,

                columns: data.columns || [],

                rows: data.rows || []

              }

            } else {

              tab.sqlResult = { success: false, error: res.data?.message || '执行失败' }

            }

          } catch (error) {

            tab.sqlResult = { success: false, error: error.message || 'SQL执行出错' }

          } finally {

            tab.executing = false

          }

        },

    async formatSql(tab) {
      if (!tab || tab.type !== 'query') {
        this.$message.error('无效的查询标签页')
        return
      }

      if (!tab.sqlQuery.trim()) {
        this.$message.warning('请输入SQL')
        return
      }

      tab.formatting = true
      try {
        const res = await dbAPI.formatSQL({
          sql: tab.sqlQuery.trim()
        })

        if (res.data?.code === 200 && res.data?.data?.formattedSql) {
          tab.sqlQuery = res.data.data.formattedSql
          this.$message.success('格式化成功')
        } else {
          this.$message.error(res.data?.message || '格式化失败')
        }
      } catch (error) {
        this.$message.error('格式化失败: ' + (error.message || ''))
      } finally {
        tab.formatting = false
      }
    },

    async preCheckSql(tab) {
      if (!tab || tab.type !== 'query') {
        this.$message.error('无效的查询标签页')
        return
      }

      if (!tab.sqlQuery.trim()) {
        this.$message.warning('请输入SQL')
        return
      }

      tab.preChecking = true
      tab.preCheckResult = null
      try {
        const res = await dbAPI.preCheckSQL({
          instanceId: tab.node.instanceId,
          db: tab.node.databaseName,
          sql: tab.sqlQuery.trim()
        })

        if (res.data?.code === 200 && res.data?.data) {
          tab.preCheckResult = res.data.data
          const { sqlTypeText, affectedRows, estimateSQL } = res.data.data

          this.$alert(
            `SQL类型: ${sqlTypeText}\n预计影响行数: ${affectedRows}\n\n预估SQL:\n${estimateSQL}`,
            '预检查结果',
            {
              confirmButtonText: '确定',
              type: 'info',
              dangerouslyUseHTMLString: false
            }
          )
        } else {
          this.$message.error(res.data?.message || '预检查失败')
        }
      } catch (error) {
        this.$message.error('预检查失败: ' + (error.message || ''))
      } finally {
        tab.preChecking = false
      }
    },
    exportQueryResultCsv(tab) {
      if (!tab || !tab.sqlResult || !tab.sqlResult.success) return
      const cols = tab.sqlResult.columns || []
      const rows = tab.sqlResult.rows || []
      if (!cols.length || !rows.length) return
      const esc = (v) => {
        if (v === null || v === undefined) return ''
        let s = typeof v === 'string' ? v : (typeof v === 'object' ? JSON.stringify(v) : String(v))
        s = s.replace(/"/g, '""')
        return `"${s}"`
      }
      const header = cols.map(esc).join(',')
      const lines = rows.map(r => cols.map(c => esc(r[c])).join(','))
      const csv = '\ufeff' + [header, ...lines].join('\n')
      const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
      const db = tab.node?.databaseName || 'db'
      const name = tab.node?.tableName ? `query_${db}_${tab.node.tableName}` : `query_${db}`
      const ts = new Date().toISOString().replace(/[:T]/g, '-').slice(0,19)
      this.downloadBlob(blob, `${name}_${ts}.csv`)
    },
    handlePageChange(page, tab) {
      tab.pagination.currentPage = page
      this.refreshTableData(tab)
    },
    handlePageSizeChange(size, tab) {
      tab.pagination.pageSize = size
      tab.pagination.currentPage = 1
      this.refreshTableData(tab)
    }
  }
}
</script>

<style scoped>
.db-terminal-page {
  padding: 20px;
  min-height: 100vh;
  background: #eef1f6;
}

.terminal-container {
  height: calc(100vh - 40px);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  }

.tree-aside {
  background: rgba(255, 255, 255, 0.5);
  border-right: 1px solid rgba(103, 126, 234, 0.2);
  padding: 15px;
}

.aside-resizer {
  width: 8px;
  height: 100%;
  cursor: col-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  user-select: none;
}
.aside-resizer:hover {
  background: rgba(103, 126, 234, 0.08);
}
.aside-resizer .resizer-icon {
  color: #909399;
}

.tree-card {
  height: 100%;
  border-radius: 8px;
}

/* Highlight current selected node in tree */
:deep(.el-tree--highlight-current .el-tree-node.is-current > .el-tree-node__content) {
  background-color: #eef5ff !important;
  color: #409eff;
  font-weight: 600;
  border-radius: 6px;
}

:deep(.el-tree--highlight-current .el-tree-node.is-current > .el-tree-node__content .node-label) {
  color: #409eff;
}

.tree-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  padding: 4px 0;
}

.node-icon {
  width: 16px;
  height: 16px;
}

.node-label {
  flex: 1;
  font-weight: 500;
}

.node-comment {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.node-size {
  font-size: 12px;
  color: #909399;
  margin-left: 6px;
}

.main-content {
  padding: 15px;
  position: relative;
}

.content-card {
  height: 100%;
  border-radius: 8px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  width: 24px;
  height: 24px;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.header-connection-info {
  margin-left: 20px;
  font-size: 14px;
  color: #606266;
  display: flex;
  align-items: center;
}

.header-connection-info .info-value {
  color: #667eea;
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.table-data-section {
  margin-top: 15px;
}

.data-table {
  margin-bottom: 20px;
}

/* Compact table cells */
:deep(.el-table__cell) {
  padding: 6px 8px;
}

.null-cell {
  color: #b1b4ba;
  font-style: italic;
}

.column-header {
  display: flex;
  align-items: center;
  gap: 4px;
}

.column-title {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.pk-hint {
  font-size: 11px;
  color: #909399;
}

.column-comment {
  font-size: 11px;
  color: #909399;
}

/* Make tooltip icon show hand cursor */
.column-header .el-icon { cursor: pointer; }
.column-header { cursor: pointer; }
.column-title { cursor: pointer; }
.pk-hint { cursor: pointer; }
.column-comment { cursor: pointer; }

.pagination-section {
  display: flex;
  justify-content: flex-end;
  padding: 15px 0;
}

.insert-form-container {
  padding: 20px;
  max-width: 700px;
  margin: 0 auto;
}

.field-hint {
  margin-top: 4px;
}

/* Compact form spacing */
:deep(.el-form-item) {
  margin-bottom: 10px;
}

.sql-actions {
  margin-top: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.sql-result-container {
  margin-top: 20px;
  padding: 12px;
  background-color: #f9fafb;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

/* Compact SQL result table */
.sql-result-table :deep(.el-table__cell) {
  padding: 4px 6px;
  font-size: 12px;
}
.sql-result-table :deep(.el-table__header .el-table__cell) {
  padding: 4px 6px;
  font-weight: 600;
}
.sql-result-table :deep(.el-table__inner-wrapper::before) {
  /* thinner header separator */
  height: 1px;
}

.mb10 {
  margin-bottom: 10px;
}

.database-info,
.instance-info {
  padding: 20px;
}

/* Row context menu */
.row-context-menu {
  position: absolute;
  z-index: 3000;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.12);
  width: 200px;
  height: auto;
  overflow: visible;
}
.tree-context-menu {
  position: fixed;
  z-index: 10000;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  box-shadow: 0 8px 20px rgba(0,0,0,0.12);
  width: 200px;
  height: auto;
  overflow: visible;
}
.tree-context-menu .menu-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.tree-context-menu .menu-item:hover { background: #f5f7fa; }
.tree-context-menu .menu-item.danger { color: #f56c6c; }
.tree-context-menu .menu-item,
.tree-context-menu .menu-item:hover,
.tree-context-menu .menu-item.danger {
  /* inherit same styles as row-context-menu items */
}
.row-context-menu .menu-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.row-context-menu .menu-item:hover {
  background: #f5f7fa;
}
.row-context-menu .menu-item.danger {
  color: #f56c6c;
}
.tree-context-menu .menu-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.tree-context-menu .menu-item:hover {
  background: #f5f7fa;
}
.tree-context-menu .menu-item.danger {
  color: #f56c6c;
}

/* Edit table columns layout */
.edit-columns-container { padding: 8px 4px; }
.edit-col-row { margin-bottom: 6px; }
.edit-col-grid {
  display: grid;
  grid-template-columns: 1fr auto auto auto 1fr 1fr;
  gap: 8px;
  align-items: center;
}

/* Table Structure Dialog */
.table-structure-dialog .sql-display-container {
  background: #f5f7fa;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 0;
  min-height: 400px;
}

.table-structure-dialog .sql-code {
  margin: 0;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #2c3e50;
  background: #f5f7fa;
  white-space: pre;
  overflow-x: auto;
  tab-size: 2;
  max-height: 70vh;
  overflow-y: auto;
}

/* CodeEditor 容器高度控制 */
.table-structure-dialog .sql-display-container .monaco-editor-container {
  min-height: 450px !important;
  height: 450px !important;
}

.table-structure-dialog .el-dialog__body {
  padding: 20px;
  min-height: 450px;
}
</style>
              <!-- Copy Table Dialog -->
              <el-dialog v-model="copyTableDialog.visible" title="复制表" width="500px" append-to-body>
                <el-form label-width="120px" size="small">
                  <el-form-item label="源表">
                    <el-input v-model="copyTableDialog.oldName" readonly />
                  </el-form-item>
                  <el-form-item label="新表名">
                    <el-input v-model="copyTableDialog.newName" placeholder="输入新表名" />
                  </el-form-item>
                  <el-form-item label="复制数据">
                    <el-switch v-model="copyTableDialog.copyData" />
                  </el-form-item>
                </el-form>
                <template #footer>
                  <el-button @click="copyTableDialog.visible = false">取消</el-button>
                  <el-button type="primary" :loading="copyTableDialog.loading" @click="confirmCopyTable">复制</el-button>
                </template>
              </el-dialog>

              <!-- Rename Table Dialog -->
              <el-dialog v-model="renameDialog.visible" title="重命名表" width="500px" append-to-body>
                <el-form label-width="120px" size="small">
                  <el-form-item label="原表名">
                    <el-input v-model="renameDialog.oldName" readonly />
                  </el-form-item>
                  <el-form-item label="新表名">
                    <el-input v-model="renameDialog.newName" placeholder="输入新表名" />
                  </el-form-item>
                </el-form>
                <template #footer>
                  <el-button @click="renameDialog.visible = false">取消</el-button>
                  <el-button type="primary" :loading="renameDialog.loading" @click="confirmRenameTable">重命名</el-button>
                </template>
              </el-dialog>

              <!-- Edit Table Dialog -->
              <teleport to="body">
              <el-dialog v-model="editTableDialog.visible" :title="`编辑表: ${editTableDialog.tableName || ''}`" width="860px" :z-index="20000" :teleported="true">
                <el-tabs v-model="editTableDialog.activeTab" v-loading="editTableDialog.loading">
                  <el-tab-pane label="字段" name="columns">
                    <div class="edit-columns-container">
                      <div style="display:flex;gap:8px;margin-bottom:8px;">
                        <el-button type="primary" size="small" @click="addEditColumn">新增列</el-button>
                      </div>
                      <el-table :data="editTableDialog.columns" size="small" border>
                        <el-table-column label="列名" min-width="140">
                          <template #default="scope">
                            <el-input v-model="scope.row.columnName" placeholder="列名" size="small" />
                          </template>
                        </el-table-column>
                        <el-table-column label="类型" min-width="160">
                          <template #default="scope">
                            <el-input v-model="scope.row.type" placeholder="int,varchar(255),datetime..." size="small" />
                          </template>
                        </el-table-column>
                        <el-table-column label="可空" width="80" align="center">
                          <template #default="scope">
                            <el-switch v-model="scope.row.nullable" />
                          </template>
                        </el-table-column>
                        <el-table-column label="默认值" min-width="140">
                          <template #default="scope">
                            <el-input v-model="scope.row.default" placeholder="NULL/数字/函数/文本" size="small" />
                          </template>
                        </el-table-column>
                        <el-table-column label="主键" width="80" align="center">
                          <template #default="scope">
                            <el-switch v-model="scope.row.pk" />
                          </template>
                        </el-table-column>
                        <el-table-column label="自增" width="80" align="center">
                          <template #default="scope">
                            <el-switch v-model="scope.row.autoIncrement" />
                          </template>
                        </el-table-column>
                        <el-table-column label="备注" min-width="180">
                          <template #default="scope">
                            <el-input v-model="scope.row.comment" placeholder="备注" size="small" />
                          </template>
                        </el-table-column>
                        <el-table-column label="操作" width="90" align="center">
                          <template #default="scope">
                            <el-button text type="danger" size="small" @click="removeEditColumn(scope.$index)">删除</el-button>
                          </template>
                        </el-table-column>
                      </el-table>
                      <div style="margin-top:10px;display:flex;gap:8px;">
                        <el-button size="small" @click="generateAlterSQL">生成变更SQL</el-button>
                        <el-button size="small" type="primary" @click="applyEditTable">应用变更</el-button>
                      </div>
                    </div>
                  </el-tab-pane>
                  <el-tab-pane label="DDL" name="ddl">
                    <CodeEditor v-model="editTableDialog.ddl" language="sql" height="360px" />
                  </el-tab-pane>
                </el-tabs>
                <template #footer>
                  <el-button @click="editTableDialog.visible = false">关闭</el-button>
                  <el-button type="primary" :loading="editTableDialog.loading" @click="applyEditTable">应用 DDL</el-button>
                </template>
              </el-dialog>
              </teleport>

              <!-- Create Database Dialog -->
              <el-dialog v-model="createDbDialog.visible" title="创建数据库" width="520px" append-to-body>
                <el-form label-width="120px" size="small">
                  <el-form-item label="库名">
                    <el-input v-model="createDbDialog.dbName" placeholder="字母/数字/下划线" />
                  </el-form-item>
                  <el-form-item label="字符集">
                    <el-input v-model="createDbDialog.charset" placeholder="可选，例如 utf8mb4" />
                  </el-form-item>
                  <el-form-item label="排序规则">
                    <el-input v-model="createDbDialog.collate" placeholder="可选，例如 utf8mb4_general_ci" />
                  </el-form-item>
                  <el-form-item label="不存在则创建">
                    <el-switch v-model="createDbDialog.ifNotExists" />
                  </el-form-item>
                </el-form>
                <template #footer>
                  <el-button @click="createDbDialog.visible = false">取消</el-button>
                  <el-button type="primary" :loading="createDbDialog.loading" @click="confirmCreateDatabase">创建</el-button>
                </template>
              </el-dialog>

              <!-- Create Table SQL Dialog -->
              <el-dialog
                v-model="createTableSqlDialog.visible"
                title="创建数据表 (SQL)"
                width="900px"
                append-to-body
                :teleported="true"
                :z-index="20000"
                class="table-structure-dialog"
              >
                <div class="sql-display-container">
                  <CodeEditor v-model="createTableSqlDialog.sql" language="sql" height="450px" />
                </div>
                <template #footer>
                  <el-button @click="createTableSqlDialog.visible = false">取消</el-button>
                  <el-button type="primary" :loading="createTableSqlDialog.loading" @click="confirmCreateTableSql">执行</el-button>
                </template>
              </el-dialog>

              <!-- Create Database SQL Dialog -->
              <el-dialog
                v-model="createDbSqlDialog.visible"
                title="创建数据库 (SQL)"
                width="900px"
                append-to-body
                :teleported="true"
                :z-index="20000"
                class="table-structure-dialog"
              >
                <div class="sql-display-container">
                  <CodeEditor v-model="createDbSqlDialog.sql" language="sql" height="300px" />
                </div>
                <template #footer>
                  <el-button @click="createDbSqlDialog.visible = false">取消</el-button>
                  <el-button type="primary" :loading="createDbSqlDialog.loading" @click="confirmCreateDatabaseSql">执行</el-button>
                </template>
              </el-dialog>
