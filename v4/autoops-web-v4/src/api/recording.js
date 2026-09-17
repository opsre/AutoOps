import request from "@/utils/request"

export default {
    // ==================== 会话录制管理 ====================

    /**
     * 获取会话录制列表
     * @param {Object} params - 查询参数
     * @param {number} params.page - 页码
     * @param {number} params.pageSize - 每页数量
     * @param {number} params.adminId - 用户ID过滤
     * @param {number} params.hostId - 主机ID过滤
     * @param {number} params.status - 状态过滤
     * @param {number} params.riskLevel - 风险等级过滤
     * @param {string} params.keyword - 关键词搜索
     * @param {string} params.startTime - 开始时间
     * @param {string} params.endTime - 结束时间
     */
    getRecordingList(params) {
        return request({
            url: '/recordings',
            method: 'get',
            params: params
        })
    },

    /**
     * 获取会话录制详情
     * @param {number} id - 录制记录ID
     */
    getRecordingDetail(id) {
        return request({
            url: `/recordings/${id}`,
            method: 'get'
        })
    },

    /**
     * 获取录制cast文件内容（用于播放器）
     * @param {number} id - 录制记录ID
     */
    getRecordingCast(id) {
        return request({
            url: `/recordings/${id}/cast`,
            method: 'get',
            responseType: 'text'
        })
    },

    /**
     * 下载录制文件
     * @param {number} id - 录制记录ID
     */
    downloadRecording(id) {
        return request({
            url: `/recordings/${id}/download`,
            method: 'get',
            responseType: 'blob'
        })
    },

    /**
     * 删除录制记录
     * @param {number} id - 录制记录ID
     */
    deleteRecording(id) {
        return request({
            url: `/recordings/${id}`,
            method: 'delete'
        })
    },

    /**
     * 批量删除录制记录
     * @param {Array<number>} ids - 录制记录ID列表
     */
    batchDeleteRecording(ids) {
        return request({
            url: '/recordings/batch',
            method: 'delete',
            data: { ids }
        })
    },

    /**
     * 获取命令审计列表
     * @param {Object} params - 查询参数
     */
    getCommandAuditList(params) {
        return request({
            url: '/recordings/commands',
            method: 'get',
            params: params
        })
    }
}
