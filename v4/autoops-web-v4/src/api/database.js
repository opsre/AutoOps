import request from '@/utils/request'

/**
 * Unified database API bundle covering MySQL/PostgreSQL, Redis, Elasticsearch, and MongoDB.
 *
 * Exports:
 * - default: dbAPI methods flattened, plus redis/es/mongo namespaces
 * - dbAPI: MySQL/PostgreSQL API
 * - redisAPI: Redis API
 * - esAPI: Elasticsearch API
 * - mongoAPI: MongoDB API
 */
/**
 * Database Instance Management API
 * Support MySQL and PostgreSQL
 */
export const dbAPI = {
  // ==================== Database Instance Management ====================

  /**
   * Get database instance list (unified interface)
   */
  getInstanceList(params) {
    return request({
      url: '/api/v1/db/unified/instances',
      method: 'get',
      params
    })
  },

  /**
   * Get database instance by type and ID (unified interface)
   */
  getInstanceByTypeAndId(dbType, id) {
    return request({
      url: `/api/v1/db/unified/instances/${dbType}/${id}`,
      method: 'get'
    })
  },

  /**
   * Get instance password/secret (unified interface)
   * For mongodb, returns full URI containing credentials
   */
  getInstancePassword(dbType, id) {
    return request({
      url: `/api/v1/db/unified/instances/${dbType}/${id}/password`,
      method: 'get'
    })
  },

  /**
   * Get database statistics (unified interface)
   */
  getInstanceStats() {
    return request({
      url: '/api/v1/db/unified/instances/stats',
      method: 'get'
    })
  },

  /**
   * Get database instance list (legacy)
   * @deprecated Use getInstanceList instead
   */
  getInstanceListLegacy(params) {
    return request({
      url: '/api/v1/db/instances',
      method: 'get',
      params
    })
  },

  /**
   * Create database (DDL)
   */
  createDatabaseDDL(dbId, data) {
    return request({
      url: `/api/v1/db/databases/${dbId}/create-database`,
      method: 'post',
      data
    })
  },

  /**
   * Create table (DDL)
   */
  createTableDDL(dbId, data) {
    return request({
      url: `/api/v1/db/databases/${dbId}/create-table`,
      method: 'post',
      data
    })
  },

  /**
   * Get database instance by ID
   */
  getInstanceById(id) {
    return request({
      url: `/api/v1/db/instances/${id}`,
      method: 'get'
    })
  },

  /**
   * Create database instance
   */
  createInstance(data) {
    return request({
      url: '/api/v1/db/instances',
      method: 'post',
      data
    })
  },

  /**
   * Update database instance
   */
  updateInstance(data) {
    return request({
      url: '/api/v1/db/instances',
      method: 'put',
      data
    })
  },

  /**
   * Delete database instance
   */
  deleteInstance(id) {
    return request({
      url: `/api/v1/db/instances/${id}`,
      method: 'delete'
    })
  },

  /**
   * Test database connection
   */
  testConnection(data) {
    return request({
      url: '/api/v1/db/instances/test-conn',
      method: 'post',
      data
    })
  },

  /**
   * Get databases under instance
   */
  getInstanceDatabases(instanceId, params) {
    return request({
      url: `/api/v1/db/instances/${instanceId}/databases`,
      method: 'get',
      params
    })
  },

  // ==================== Database Management ====================

  /**
   * Get database list
   */
  getDatabaseList(params) {
    return request({
      url: '/api/v1/db/databases',
      method: 'get',
      params
    })
  },

  /**
   * Get database by ID
   */
  getDatabaseById(id) {
    return request({
      url: `/api/v1/db/databases/${id}`,
      method: 'get'
    })
  },

  /**
   * Create database config
   */
  createDatabase(data) {
    return request({
      url: '/api/v1/db/databases',
      method: 'post',
      data
    })
  },

  /**
   * Update database config
   */
  updateDatabase(data) {
    return request({
      url: '/api/v1/db/databases',
      method: 'put',
      data
    })
  },

  /**
   * Delete database config
   */
  deleteDatabase(id) {
    return request({
      url: `/api/v1/db/databases/${id}`,
      method: 'delete'
    })
  },

  // ==================== Table Management ====================

  /**
   * Get table list
   */
  getTables(dbId, dbName) {
    return request({
      url: `/api/v1/db/databases/${dbId}/tables`,
      method: 'get',
      params: { dbName }
    })
  },

  /**
   * Get table column info
   */
  getTableColumns(dbId, dbName, tableName) {
    return request({
      url: `/api/v1/db/databases/${dbId}/columns`,
      method: 'get',
      params: { db: dbName, table: tableName }
    })
  },

  /**
   * Get table index info
   */
  getTableIndexes(dbId, dbName, tableName) {
    return request({
      url: `/api/v1/db/databases/${dbId}/indexes`,
      method: 'get',
      params: { db: dbName, table: tableName }
    })
  },

  /**
   * Get table DDL
   */
  getTableDDL(dbId, dbName, tableName) {
    return request({
      url: `/api/v1/db/databases/${dbId}/ddl`,
      method: 'get',
      params: { db: dbName, table: tableName }
    })
  },

  /**
   * Get approximate row count for a table
   */
  getApproxCount(dbId, dbName, tableName) {
    return request({
      url: `/api/v1/db/databases/${dbId}/approx-count`,
      method: 'get',
      params: { db: dbName, table: tableName }
    })
  },

  /**
   * Copy table
   */
  copyTable(data) {
    return request({
      url: '/api/v1/db/databases/copy-table',
      method: 'post',
      data
    })
  },

  // ==================== SQL Execution ====================

  /**
   * Execute SQL
   */
  execSQL(data) {
    return request({
      url: '/api/v1/db/databases/exec-sql',
      method: 'post',
      data
    })
  },

  /**
   * Format SQL
   */
  formatSQL(data) {
    return request({
      url: '/api/v1/db/format-sql',
      method: 'post',
      data
    })
  },

  /**
   * Pre-check SQL before execution
   */
  preCheckSQL(data) {
    return request({
      url: '/api/v1/db/pre-check-sql',
      method: 'post',
      data
    })
  },

  /**
   * Execute SQL file upload (multipart)
   */
  execSqlFile(dbId, dbName, file) {
    const form = new FormData()
    form.append('db', dbName)
    form.append('file', file)
    return request({
      url: `/api/v1/db/databases/${dbId}/exec-sql-file`,
      method: 'post',
      data: form,
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  /**
   * Query SQL file execution job status
   */
  getExecSqlFileStatus(jobId) {
    return request({
      url: `/api/v1/db/databases/exec-sql-file/${jobId}/status`,
      method: 'get'
    })
  },

  /**
   * Get hint tables for SQL editor autocomplete
   */
  getHintTables(dbId, dbName) {
    return request({
      url: `/api/v1/db/databases/${dbId}/hint-tables`,
      method: 'get',
      params: { db: dbName }
    })
  },

  /**
   * Get single row by primary key
   */
  getRow(dbId, dbName, tableName, pk, pkv) {
    return request({
      url: `/api/v1/db/databases/${dbId}/row`,
      method: 'get',
      params: { db: dbName, table: tableName, pk, pkv }
    })
  },

  /**
   * Delete single row by primary key
   */
  deleteRow(dbId, dbName, tableName, pk, pkv) {
    return request({
      url: `/api/v1/db/databases/${dbId}/row`,
      method: 'delete',
      params: { db: dbName, table: tableName, pk, pkv }
    })
  },

  /**
   * Build export CSV URL (browser download)
   */
  buildExportCsvUrl(dbId, query) {
    const p = new URLSearchParams(query)
    return `/api/v1/db/databases/${dbId}/export-csv?${p.toString()}`
  },

  /**
   * Build export SQL URL (browser download)
   */
  buildExportSqlUrl(dbId, query) {
    const p = new URLSearchParams(query)
    return `/api/v1/db/databases/${dbId}/export-sql?${p.toString()}`
  },

  /**
   * Export CSV via POST (Blob)
   */
  exportCsv(dbId, data) {
    return request({
      url: `/api/v1/db/databases/${dbId}/export-csv`,
      method: 'post',
      data,
      responseType: 'blob'
    })
  },

  /**
   * Export SQL via POST (Blob)
   */
  exportSql(dbId, data) {
    return request({
      url: `/api/v1/db/databases/${dbId}/export-sql`,
      method: 'post',
      data,
      responseType: 'blob'
    })
  },

  /**
   * Export entire database structure (DDL only)
   */
  exportDatabaseStructure(dbId, params) {
    return request({
      url: `/api/v1/db/databases/${dbId}/export-database-structure`,
      method: 'get',
      params,
      responseType: 'blob'
    })
  },

  /**
   * Export entire database structure + data (DDL + INSERT)
   */
  exportDatabaseFull(dbId, params) {
    return request({
      url: `/api/v1/db/databases/${dbId}/export-database-full`,
      method: 'get',
      params,
      responseType: 'blob'
    })
  },

  /**
   * Start async full database export (DDL + data)
   * Returns 200 with { taskId, status, message }
   */
  exportDatabaseFullAsync(dbId, dbName, exportType = 'full') {
    return request({
      url: `/api/v1/db/databases/${dbId}/export-database-full-async`,
      method: 'post',
      params: { db: dbName, exportType }
    })
  },

  /**
   * Get export task status
   */
  getExportJobStatus(taskId) {
    return request({
      url: `/api/v1/db/export-tasks/${taskId}`,
      method: 'get'
    })
  },

  /**
   * Download export task file (Blob)
   */
  downloadExportJobFile(taskId) {
    return request({
      url: `/api/v1/db/export-tasks/${taskId}/download`,
      method: 'get',
      responseType: 'blob'
    })
  },

  // ==================== SQL Template Management ====================

  /**
   * Get SQL template list
   */
  getSqlTemplates(params) {
    return request({
      url: '/api/v1/db/sql-templates',
      method: 'get',
      params
    })
  },

  /**
   * Get SQL template by ID
   */
  getSqlTemplateById(id) {
    return request({
      url: `/api/v1/db/sql-templates/${id}`,
      method: 'get'
    })
  },

  /**
   * Create SQL template
   */
  createSqlTemplate(data) {
    return request({
      url: '/api/v1/db/sql-templates',
      method: 'post',
      data
    })
  },

  /**
   * Update SQL template
   */
  updateSqlTemplate(data) {
    return request({
      url: '/api/v1/db/sql-templates',
      method: 'put',
      data
    })
  },

  /**
   * Delete SQL template
   */
  deleteSqlTemplate(id) {
    return request({
      url: `/api/v1/db/sql-templates/${id}`,
      method: 'delete'
    })
  },

  /**
   * Get SQL template name list
   */
  getSqlTemplateNames(dbId, db) {
    return request({
      url: '/api/v1/db/sql-templates/names',
      method: 'get',
      params: { dbId, db }
    })
  },

  // ==================== SQL Execution History ====================

  /**
   * Get SQL execution history
   */
  getSqlExecHistory(params) {
    return request({
      url: '/api/v1/db/sql-execs',
      method: 'get',
      params
    })
  },

  /**
   * Get SQL execution detail by ID
   */
  getSqlExecById(id) {
    return request({
      url: `/api/v1/db/sql-execs/${id}`,
      method: 'get'
    })
  },

  /**
   * Get SQL execution statistics
   */
  getSqlExecStatistics(params) {
    return request({
      url: '/api/v1/db/sql-execs/statistics',
      method: 'get',
      params
    })
  }
}

// Redis Management API
export const redisAPI = {
  // Instance list
  getInstanceList(params) {
    return request({
      url: '/api/v1/redis',
      method: 'get',
      params
    })
  },

  // Instance detail
  getInstanceById(id) {
    return request({
      url: `/api/v1/redis/${id}`,
      method: 'get'
    })
  },

  // Create instance
  createInstance(data) {
    return request({
      url: '/api/v1/redis',
      method: 'post',
      data
    })
  },

  // Update instance
  updateInstance(data) {
    return request({
      url: '/api/v1/redis',
      method: 'put',
      data
    })
  },

  // Delete instance
  deleteInstance(id) {
    return request({
      url: `/api/v1/redis/${id}`,
      method: 'delete'
    })
  },

  // Test connection
  testConnection(data) {
    return request({
      url: '/api/v1/redis/test-conn',
      method: 'post',
      data
    })
  },

  // Info
  getInfo(id, params) {
    return request({
      url: `/api/v1/redis/${id}/info`,
      method: 'get',
      params
    })
  },

  // Cluster info
  getClusterInfo(id) {
    return request({
      url: `/api/v1/redis/${id}/cluster-info`,
      method: 'get'
    })
  },

  // Scan keys
  scanKeys(id, db, data) {
    return request({
      url: `/api/v1/redis/${id}/${db}/scan`,
      method: 'post',
      data
    })
  },

  // Scan all DBs (standalone/sentinel)
  scanAll(id, data) {
    return request({
      url: `/api/v1/redis/${id}/scan`,
      method: 'post',
      data
    })
  },

  // List available DBs of an instance
  getDBs(id) {
    return request({
      url: `/api/v1/redis/${id}/dbs`,
      method: 'get'
    })
  },

  // Key info
  getKeyInfo(id, db, key) {
    return request({
      url: `/api/v1/redis/${id}/${db}/key-info`,
      method: 'get',
      params: { key }
    })
  },

  // Key TTL
  getKeyTTL(id, db, key) {
    return request({
      url: `/api/v1/redis/${id}/${db}/key-ttl`,
      method: 'get',
      params: { key }
    })
  },

  // Key memory usage
  getKeyMemUse(id, db, key) {
    return request({
      url: `/api/v1/redis/${id}/${db}/key-memuse`,
      method: 'get',
      params: { key }
    })
  },

  // Run command
  runCmd(id, db, data) {
    return request({
      url: `/api/v1/redis/${id}/${db}/run-cmd`,
      method: 'post',
      data
    })
  }
}

// Elasticsearch Management API
export const esAPI = {
  // List instances
  getInstanceList(params) {
    return request({
      url: '/api/v1/elasticsearch',
      method: 'get',
      params
    })
  },

  // Get instance detail
  getInstanceById(id) {
    return request({
      url: `/api/v1/elasticsearch/${id}`,
      method: 'get'
    })
  },

  // Create
  createInstance(data) {
    return request({
      url: '/api/v1/elasticsearch',
      method: 'post',
      data
    })
  },

  // Update
  updateInstance(data) {
    return request({
      url: '/api/v1/elasticsearch',
      method: 'put',
      data
    })
  },

  // Delete
  deleteInstance(id) {
    return request({
      url: `/api/v1/elasticsearch/${id}`,
      method: 'delete'
    })
  },

  // Test connection
  testConnection(data) {
    return request({
      url: '/api/v1/elasticsearch/test-conn',
      method: 'post',
      data
    })
  },

  // Proxy ES API (optional)
  proxy(id, path, method = 'get', data, params, headers) {
    return request({
      url: `/api/v1/elasticsearch/${id}/proxy/${path}`,
      method,
      data,
      params,
      headers
    })
  }
}

// MongoDB Management API
export const mongoAPI = {
  // List instances
  getInstanceList(params) {
    return request({
      url: '/api/v1/mongodb',
      method: 'get',
      params
    })
  },

  // Get instance detail
  getInstanceById(id) {
    return request({
      url: `/api/v1/mongodb/${id}`,
      method: 'get'
    })
  },

  // Create
  createInstance(data) {
    return request({
      url: '/api/v1/mongodb',
      method: 'post',
      data
    })
  },

  // Update
  updateInstance(data) {
    return request({
      url: '/api/v1/mongodb',
      method: 'put',
      data
    })
  },

  // Delete
  deleteInstance(id) {
    return request({
      url: `/api/v1/mongodb/${id}`,
      method: 'delete'
    })
  },

  // Test connection
  testConnection(data) {
    return request({
      url: '/api/v1/mongodb/test-conn',
      method: 'post',
      data
    })
  },

  // Get databases list
  getDatabases(id) {
    return request({
      url: `/api/v1/mongodb/${id}/databases`,
      method: 'get'
    })
  },

  // Get collections list
  getCollections(id, params) {
    return request({
      url: `/api/v1/mongodb/${id}/collections`,
      method: 'get',
      params // { database: 'dbname' }
    })
  },

  // Find documents
  findDocuments(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/command/find`,
      method: 'post',
      data // { database, collection, filter, limit, skip, sort }
    })
  },

  // Insert document
  insertDocument(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/command/insert`,
      method: 'post',
      data // { database, collection, document }
    })
  },

  // Update document
  updateDocument(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/command/update-by-id`,
      method: 'post',
      data // { database, collection, filter, update }
    })
  },

  // Delete document
  deleteDocument(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/command/delete-by-id`,
      method: 'post',
      data // { database, collection, filter }
    })
  },

  // Collection stats - 使用 run-command 接口
  collectionStats(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/run-command`,
      method: 'post',
      data: {
        database: data.database,
        command: {
          collStats: data.collection
        }
      }
    })
  },

  // Drop collection - 使用 run-command 接口
  dropCollection(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/run-command`,
      method: 'post',
      data: {
        database: data.database,
        command: {
          drop: data.collection
        }
      }
    })
  },

  // Drop database - 使用 run-command 接口
  dropDatabase(id, data) {
    return request({
      url: `/api/v1/mongodb/${id}/run-command`,
      method: 'post',
      data: {
        database: data.database,
        command: {
          dropDatabase: 1
        }
      }
    })
  },

  // Proxy MongoDB API (optional, for custom commands)
  proxy(id, path, method = 'get', data, params, headers) {
    return request({
      url: `/api/v1/mongodb/${id}/proxy/${path}`,
      method,
      data,
      params,
      headers
    })
  }
}

// Unified export for convenience (import from '@/api/database')
export default {
  ...dbAPI,
  redis: redisAPI,
  es: esAPI,
  mongo: mongoAPI
}
