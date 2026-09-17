// Parse a compact Markdown-like table string into columns and rows
// Input can be either a raw string or an object { type: 'text', content: '...' }
// Returns: { columns: [{ label, prop }], rows: Array<Record<string,string>> }

function slugify(label) {
  return String(label)
    .trim()
    .replace(/\s+/g, '_')
    .replace(/[%]/g, 'pct')
    .replace(/[^a-zA-Z0-9_\u4e00-\u9fa5]/g, '')
    .replace(/^_+|_+$/g, '') || 'col'
}

function normalizeToLines(input) {
  const text = typeof input === 'string' ? input : (input && input.content) || ''
  if (!text) return []

  // Ensure row separation. Backend may concatenate rows like "| h | |----| | r |"
  let s = text
  // Insert newlines between adjacent pipe-bounded rows
  s = s.replace(/\|\s*\|/g, '|\n|')
  // Also split on explicit newlines just in case
  const lines = s
    .split(/\r?\n/)
    .map((l) => l.trim())
    .filter((l) => l.length > 0)
  return lines
}

function isSeparatorLine(line) {
  // Matches markdown header separator: |-----|-----| or -----|-----
  const cleaned = line.replace(/\s|:|\|/g, '')
  return cleaned.length > 0 && /^-+$/.test(cleaned)
}

export function parseMarkdownTable(input) {
  const lines = normalizeToLines(input)
  if (!lines.length) return { columns: [], rows: [] }

  // Remove any leading/trailing non-table noise until we hit a line starting with '|'
  const tableLines = lines.filter((l) => l.startsWith('|') || l.includes('|'))
  if (!tableLines.length) return { columns: [], rows: [] }

  // First non-separator is header
  const headerLine = tableLines.find((l) => !isSeparatorLine(l))
  if (!headerLine) return { columns: [], rows: [] }

  const headerCells = headerLine
    .split('|')
    .map((c) => c.trim())
    .filter((c) => c.length > 0)

  const columns = headerCells.map((label) => ({ label, prop: slugify(label) }))

  const headerIndex = tableLines.indexOf(headerLine)
  const dataLines = tableLines.slice(headerIndex + 1).filter((l) => !isSeparatorLine(l))

  const rows = dataLines
    .map((line) => {
      const cells = line
        .split('|')
        .map((c) => c.trim())
        .filter((c) => c.length > 0)
      if (!cells.length) return null
      const row = {}
      columns.forEach((col, idx) => {
        row[col.prop] = cells[idx] !== undefined ? cells[idx] : ''
      })
      return row
    })
    .filter(Boolean)

  return { columns, rows }
}

export default parseMarkdownTable

