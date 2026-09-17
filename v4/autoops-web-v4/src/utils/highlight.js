import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/javascript';
import python from 'highlight.js/lib/languages/python';
import go from 'highlight.js/lib/languages/go';
import yaml from 'highlight.js/lib/languages/yaml';
import bash from 'highlight.js/lib/languages/bash';
import jsonlang from 'highlight.js/lib/languages/json';
import 'highlight.js/styles/github-dark.css';

// 注册语言
hljs.registerLanguage('javascript', javascript);
hljs.registerLanguage('python', python);
hljs.registerLanguage('go', go);
hljs.registerLanguage('yaml', yaml);
hljs.registerLanguage('bash', bash);
hljs.registerLanguage('json', jsonlang);

/**
 * 代码高亮函数
 * @param {string} code - 需要高亮的代码
 * @param {string} [language] - 代码语言(可选，自动检测)
 * @returns {string} 高亮后的HTML字符串
 */
export function highlight(code, language) {
  // 确保处理所有可能的换行符情况
  let formattedCode = '';
  
  if (typeof code === 'string') {
    // 处理字符串中的转义换行符和实际换行符
    formattedCode = code
      .replace(/\\r\\n|\\n\\r|\\r/g, '\n')  // 替换所有转义换行符为实际换行符
      .replace(/\r\n|\n\r|\r/g, '\n');     // 统一换行符为\n
  } else {
    // 如果是对象，转换为JSON并保留格式
    formattedCode = JSON.stringify(code, null, 2)
      .replace(/\\r\\n|\\n\\r|\\r/g, '\n')
      .replace(/\r\n|\n\r|\r/g, '\n');
  }
  
  // 确保highlight.js处理时保留换行符
  let result;
  try {
    result = language ? hljs.highlight(formattedCode, { language }).value : hljs.highlightAuto(formattedCode).value;
  } catch (e) {
    // Fallback to auto if specified language missing
    try {
      result = hljs.highlightAuto(formattedCode).value;
    } catch (_) {
      // Last resort: escape HTML
      result = formattedCode
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
    }
  }
    
  // 确保结果中的换行符被保留
  return result.replace(/\n/g, '<br>');
}

/**
 * 高亮DOM元素内的代码
 * @param {HTMLElement} element - 包含代码的DOM元素
 */
export function highlightElement(element) {
  hljs.highlightElement(element);
}

/**
 * 初始化高亮，处理整个页面
 */
export function initHighlighting() {
  hljs.highlightAll();
}

// 默认导出常用方法
export default {
  highlight,
  highlightElement,
  initHighlighting
};
