const CLUSTER_STORAGE_KEY = 'k8s_selected_cluster_id'

export const getStoredClusterId = () => {
  if (typeof window === 'undefined') return ''

  try {
    return localStorage.getItem(CLUSTER_STORAGE_KEY) || ''
  } catch (error) {
    console.error('读取存储的集群失败:', error)
    return ''
  }
}

export const saveStoredClusterId = (clusterId) => {
  if (typeof window === 'undefined') return

  try {
    if (clusterId) {
      localStorage.setItem(CLUSTER_STORAGE_KEY, String(clusterId))
      return
    }

    localStorage.removeItem(CLUSTER_STORAGE_KEY)
  } catch (error) {
    console.error('保存集群选择失败:', error)
  }
}

export const resolvePreferredClusterId = (clusters, preferredClusterId = '') => {
  if (!Array.isArray(clusters) || clusters.length === 0) return ''

  const normalizedPreferredId = preferredClusterId ? String(preferredClusterId) : ''
  const matchedOnlineCluster = normalizedPreferredId
    ? clusters.find(cluster => String(cluster.id) === normalizedPreferredId && cluster.status === 2)
    : null

  if (matchedOnlineCluster) {
    return matchedOnlineCluster.id
  }

  const onlineCluster = clusters.find(cluster => cluster.status === 2)
  if (onlineCluster) {
    return onlineCluster.id
  }

  const matchedCluster = normalizedPreferredId
    ? clusters.find(cluster => String(cluster.id) === normalizedPreferredId)
    : null

  return matchedCluster ? matchedCluster.id : clusters[0].id
}
