<template>
  <!-- 标签组件 -->
  <div class="tags" @contextmenu.prevent="showContextMenu">
    <el-tag class="tag" size="medium" :effect="item.title == $route.meta.tTitle ? 'dark' : 'plain'"
            v-for="(item, index) in tags" :key="item.path" @click="goTo(item.path)" @close="close(index)"
            :closable="index > 0">
            <i class="circular" v-show="item.title == $route.meta.tTitle"></i>
            {{item.title}}
    </el-tag>

    <!-- 右键菜单 -->
    <div v-if="contextMenuVisible" class="context-menu" :style="{ left: menuLeft + 'px', top: menuTop + 'px' }">
      <div class="menu-item" @click="closeAllTags">
        <el-icon><Close /></el-icon>
        <span>清空所有标签</span>
      </div>
    </div>
  </div>
</template>

<script>
import index from "vuex";
import { Close } from '@element-plus/icons-vue';

export default {
  name: "Tags",
  components: {
    Close
  },
  computed: {
    index() {
      return index
    }
  },
  data() {
    return {
      tags: [{
        path: "/dashboard",
        title: "仪表盘",
      }],
      contextMenuVisible: false,
      menuLeft: 0,
      menuTop: 0
      }
        },
    watch: {
      $route: {
        immediate: true,
        handler(val) {
          const boolean = this.tags.find(item => {
            return val.path == item.path
          })
          if (!boolean) {
            this.tags.push({
              title: val.meta.tTitle,
              path: val.path
            })
            // 限制最多显示6个标签，超过则移除最旧的标签（保留仪表盘）
            if (this.tags.length > 6) {
              this.tags.splice(1, 1)
            }
          }
        }
      }
    },
  methods: {
    // 路由跳转到指定位置
    goTo(path) {
      this.$router.push(path)
    },
    // 点击关闭标签
    close(i) {
      // 判断是否关闭的是当前页面的标签
      const isCurrentTag = this.tags[i].path === this.$route.path

      // 删除标签
      this.tags.splice(i, 1)

      // 如果关闭的是当前标签,跳转到前一个标签
      if (isCurrentTag && this.tags.length > 0) {
        const targetTag = this.tags[i - 1] || this.tags[0]
        this.$router.push(targetTag.path)
      }
    },
    // 显示右键菜单
    showContextMenu(e) {
      this.contextMenuVisible = true
      this.menuLeft = e.clientX
      this.menuTop = e.clientY
      // 点击其他区域关闭菜单
      document.addEventListener('click', this.hideContextMenu)
    },
    // 隐藏右键菜单
    hideContextMenu() {
      this.contextMenuVisible = false
      document.removeEventListener('click', this.hideContextMenu)
    },
    // 清空所有标签
    closeAllTags() {
      // 只保留仪表盘标签
      this.tags = [{
        path: "/dashboard",
        title: "仪表盘",
      }]
      // 跳转到仪表盘
      this.$router.push('/dashboard')
      // 关闭右键菜单
      this.hideContextMenu()
    }
  },
  beforeUnmount() {
    // 组件销毁前移除事件监听
    document.removeEventListener('click', this.hideContextMenu)
  }
}
</script>

<style>
.tags {
    padding-left: 20px;
    padding-top: 2px;
    padding-bottom: 2px;
    position: relative;
}
.tag {
  cursor: pointer;
  margin-right: 3px;
}
.circular {
  width: 8px;
  height: 8px;
  margin-right: 4px;
  background-color: #fff;
  border-radius: 50%;
  display: inline-block;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 9999;
  padding: 4px 0;
  min-width: 140px;
}

.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
  transition: all 0.2s;
}

.menu-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.menu-item .el-icon {
  font-size: 16px;
}
</style>
