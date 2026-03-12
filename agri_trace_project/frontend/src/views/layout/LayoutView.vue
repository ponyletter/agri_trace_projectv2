<template>
  <el-container class="layout-container">
    <!-- 左侧侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '220px'" class="sidebar">
      <!-- Logo -->
      <div class="logo-area">
        <img src="/vite.svg" alt="logo" class="logo-img" />
        <span v-if="!isCollapse" class="logo-text">阿克苏苹果溯源</span>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :collapse-transition="false"
        router
        background-color="#001529"
        text-color="#ffffffa0"
        active-text-color="#ffffff"
        class="sidebar-menu"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>数据大屏</template>
        </el-menu-item>

        <el-sub-menu index="trace-group">
          <template #title>
            <el-icon><List /></el-icon>
            <span>溯源管理</span>
          </template>
          <el-menu-item index="/batches">批次管理</el-menu-item>
          <el-menu-item index="/planting">种植管理</el-menu-item>
          <el-menu-item index="/processing">加工管理</el-menu-item>
          <el-menu-item index="/inspection">质检管理</el-menu-item>
          <el-menu-item index="/logistics">物流追踪</el-menu-item>
          <el-menu-item index="/trace-records">溯源记录</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/certificate">
          <el-icon><Medal /></el-icon>
          <template #title>电子合格证</template>
        </el-menu-item>

        <el-menu-item index="/blockchain">
          <el-icon><Connection /></el-icon>
          <template #title>区块链浏览器</template>
        </el-menu-item>

        <el-menu-item v-if="userStore.role === 'admin'" index="/admin/users">
          <el-icon><User /></el-icon>
          <template #title>用户管理</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="isCollapse = !isCollapse">
            <Fold v-if="!isCollapse" /><Expand v-else />
          </el-icon>
          <!-- 面包屑 -->
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-tag type="success" size="small" style="margin-right:12px">
            <el-icon><CircleCheck /></el-icon> 区块链已连接
          </el-tag>
          <el-dropdown @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32" style="background:#1677ff">{{ userStore.realName?.charAt(0) }}</el-avatar>
              <span style="margin-left:8px">{{ userStore.realName }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../../store/user'
import {
  DataAnalysis, List, Medal, Connection, User,
  Fold, Expand, CircleCheck, ArrowDown
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isCollapse = ref(false)

const activeMenu = computed(() => route.path)
const titleMap: Record<string, string> = {
  '/dashboard': '数据大屏',
  '/batches': '批次管理',
  '/planting': '种植管理',
  '/processing': '加工管理',
  '/inspection': '质检管理',
  '/logistics': '物流追踪',
  '/trace-records': '溯源记录',
  '/certificate': '电子合格证',
  '/blockchain': '区块链浏览器',
  '/admin/users': '用户管理',
}
const currentTitle = computed(() => titleMap[route.path] || '首页')

function handleCommand(cmd: string) {
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.layout-container { height: 100vh; overflow: hidden; }
.sidebar {
  background: #001529;
  transition: width 0.3s;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.logo-area {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  border-bottom: 1px solid #ffffff15;
  overflow: hidden;
  white-space: nowrap;
}
.logo-img { width: 32px; height: 32px; flex-shrink: 0; }
.logo-text { color: #fff; font-size: 15px; font-weight: 600; margin-left: 10px; }
.sidebar-menu { border-right: none; flex: 1; overflow-y: auto; }
.header {
  height: 64px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0,0,0,.08);
}
.header-left { display: flex; align-items: center; gap: 16px; }
.collapse-btn { font-size: 20px; cursor: pointer; color: #666; }
.header-right { display: flex; align-items: center; }
.user-info { display: flex; align-items: center; cursor: pointer; color: #333; }
.main-content { background: #f5f7fa; overflow-y: auto; padding: 20px; }
</style>
