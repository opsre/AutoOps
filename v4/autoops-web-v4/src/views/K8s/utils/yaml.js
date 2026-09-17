import yaml from 'js-yaml'

export function normalizeYamlContent(rawData, fallbackContent = '') {
  if (rawData === null || rawData === undefined) {
    return String(fallbackContent || '')
  }

  if (typeof rawData === 'string') {
    return rawData
  }

  if (typeof rawData === 'object') {
    if (typeof rawData.yamlContent === 'string') {
      return rawData.yamlContent
    }

    if (rawData.yaml !== undefined && rawData.yaml !== null) {
      if (typeof rawData.yaml === 'string') {
        return rawData.yaml
      }

      try {
        return yaml.dump(rawData.yaml, { indent: 2, lineWidth: -1 })
      } catch (error) {
        console.error('YAML转换失败:', error)
        return JSON.stringify(rawData.yaml, null, 2)
      }
    }

    try {
      return yaml.dump(rawData, { indent: 2, lineWidth: -1 })
    } catch (error) {
      console.error('对象转YAML失败:', error)
      return JSON.stringify(rawData, null, 2)
    }
  }

  return String(rawData)
}

export function getYamlContentFromApiResponse(responseData, fallbackContent = '') {
  return normalizeYamlContent(responseData?.data, fallbackContent)
}

export async function copyYamlToClipboard(content, emptyMessage = '没有内容可以复制') {
  if (!content || typeof content !== 'string' || content.trim() === '') {
    return {
      success: false,
      message: emptyMessage
    }
  }

  try {
    await navigator.clipboard.writeText(content)
    return { success: true }
  } catch (error) {
    console.error('复制失败:', error)

    try {
      const textArea = document.createElement('textarea')
      textArea.value = content
      textArea.setAttribute('readonly', 'readonly')
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      textArea.style.pointerEvents = 'none'
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      return { success: true }
    } catch (fallbackError) {
      console.error('备用复制方案也失败:', fallbackError)
      return {
        success: false,
        message: '复制失败，请手动复制'
      }
    }
  }
}
