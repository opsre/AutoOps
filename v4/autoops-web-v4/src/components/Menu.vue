<template>
  <el-menu
      :default-active="activePath || $route.path"
      background-color="#304156"
      text-color="#fff"
      active-text-color="#ffd04b"
      router
      mode="vertical"
      :collapse="isCollapse"
      :collapse-transition="false"
  >
    <!-- 递归渲染菜单 -->
    <menu-item
        v-for="item in processedMenuList"
        :key="item.id"
        :menu-item="item"
    />
  </el-menu>
</template>

<script>
import storage from '@/utils/storage'

export default {
  name: 'Menu',
  props: {
    menuList: {
      type: Array,
      default: () => [],
      required: true
    },
    activePath: {
      type: String,
      default: '/dashboard'
    },
    isCollapse: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      fallbackMenuList: storage.getItem('leftMenuList') || []
    }
  },
  computed: {
    processedMenuList() {
      const sourceData = this.menuList.length > 0 ? this.menuList : this.fallbackMenuList
      const formatted = this.formatMenu(sourceData)
      console.log('【格式化后菜单】', formatted)
      return formatted
    }
  },
  methods: {
    formatMenu(menuItems) {
      let nextId = 10000; // 起始 ID（可自定义）
      return menuItems.map(item => {
        const menuItem = {
          id: item.id !== undefined ? item.id : nextId++, // ✅ 如果 id 不存在，自动分配
          menuName: item.menuName,
          icon: item.icon,
          url: item.url,
          children: item.menuSvoList ? this.formatMenu(item.menuSvoList) : []
        }

        console.log(`【处理菜单项 ${menuItem.menuName}】`, {
          hasChildren: menuItem.children.length > 0,
          url: menuItem.url,
          id: menuItem.id // ✅ 确保 id 不为 undefined
        })

        return menuItem
      })
    }
  },
  components: {
    MenuItem: {
      name: 'MenuItem',
      props: {
        menuItem: {
          type: Object,
          required: true
        }
      },
      template: `
        <div>
          <!-- 有子菜单 -->
          <el-submenu
            v-if="hasChildren"
            :index="menuItem.id + ''" <!-- ✅ 使用 id 作为 index -->
            :key="menuItem.id" <!-- ✅ 使用 id 作为 key -->
          >
            <template slot="title">
              <i :class="menuItem.icon"></i>
              <span>{{ menuItem.menuName }}</span>
            </template>
            <!-- 递归渲染子菜单 -->
            <menu-item
              v-for="child in menuItem.children"
              :key="child.id" <!-- ✅ 使用子菜单的 id -->
              :menu-item="child"
            />
          </el-submenu>
          <!-- 无子菜单 -->
          <el-menu-item
            v-else
            :index="\`/\${menuItem.url}\`" <!-- ✅ 使用完整路径 -->
            :key="menuItem.id" <!-- ✅ 使用 id 作为 key -->
          >
            <i :class="menuItem.icon"></i>
            <span slot="title">{{ menuItem.menuName }}</span>
          </el-menu-item>
        </div>
      `,
      computed: {
        hasChildren() {
          return this.menuItem.children && this.menuItem.children.length > 0
        }
      }
    }
  }
}
</script>

<style scoped>
.el-menu {
  border-right: none;
}
</style>
