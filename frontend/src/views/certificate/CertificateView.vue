<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">📜 电子合格证</h3>
          <p class="page-desc">查看带有二维码、防伪印章和检测结果的电子合格证</p>
        </div>
        <div>
          <el-select v-model="selectedBatchId" placeholder="选择批次" style="width:260px;margin-right:12px" @change="loadCert">
            <el-option v-for="b in batches" :key="b.id" :label="`${b.batch_no} - ${b.product_name}`" :value="b.id" />
          </el-select>
          <el-button type="primary" icon="Printer" @click="printCert">打印合格证</el-button>
        </div>
      </div>
    </el-card>

    <div style="margin-top:16px" v-if="cert">
      <!-- 电子合格证主体 -->
      <div class="cert-wrapper" id="cert-print-area">
        <div class="cert-card">
          <!-- 顶部标题 -->
          <div class="cert-header">
            <div class="cert-title-area">
              <div class="cert-gov-logo">🏛️</div>
              <div>
                <h2 class="cert-main-title">农产品质量安全合格证</h2>
                <p class="cert-sub-title">Certificate of Agricultural Product Quality and Safety</p>
              </div>
            </div>
            <div class="cert-no-area">
              <div class="cert-no">证书编号：{{ cert.cert_no }}</div>
              <div class="cert-blockchain-badge">
                <el-icon><Connection /></el-icon>
                区块链存证
              </div>
            </div>
          </div>

          <el-divider />

          <el-row :gutter="24">
            <!-- 左侧信息 -->
            <el-col :span="16">
              <el-descriptions :column="2" border size="small" class="cert-desc">
                <el-descriptions-item label="产品名称" :span="2">
                  <strong style="font-size:16px;color:#1a1a1a">{{ cert.product_name }}</strong>
                </el-descriptions-item>
                <el-descriptions-item label="生产主体">{{ cert.producer_name }}</el-descriptions-item>
                <el-descriptions-item label="联系电话">{{ cert.producer_phone }}</el-descriptions-item>
                <el-descriptions-item label="生产地址" :span="2">{{ cert.producer_addr }}</el-descriptions-item>
                <el-descriptions-item label="产品数量">{{ cert.quantity }}</el-descriptions-item>
                <el-descriptions-item label="签发日期">{{ cert.issue_date }}</el-descriptions-item>
                <el-descriptions-item label="有效期至">{{ cert.valid_until }}</el-descriptions-item>
                <el-descriptions-item label="签发机构">{{ cert.issue_org }}</el-descriptions-item>
              </el-descriptions>

              <!-- 检测结果 -->
              <div class="check-results">
                <div class="check-title">📋 检测结果</div>
                <div class="check-items">
                  <div class="check-item" :class="{ pass: cert.pesticide_ok }">
                    <el-icon><component :is="cert.pesticide_ok ? 'CircleCheck' : 'CircleClose'" /></el-icon>
                    <span>农药残留检测</span>
                    <el-tag :type="cert.pesticide_ok ? 'success' : 'danger'" size="small">{{ cert.pesticide_ok ? '合格' : '不合格' }}</el-tag>
                  </div>
                  <div class="check-item" :class="{ pass: cert.heavy_metal_ok }">
                    <el-icon><component :is="cert.heavy_metal_ok ? 'CircleCheck' : 'CircleClose'" /></el-icon>
                    <span>重金属含量检测</span>
                    <el-tag :type="cert.heavy_metal_ok ? 'success' : 'danger'" size="small">{{ cert.heavy_metal_ok ? '合格' : '不合格' }}</el-tag>
                  </div>
                  <div class="check-item" :class="{ pass: cert.microbe_ok }">
                    <el-icon><component :is="cert.microbe_ok ? 'CircleCheck' : 'CircleClose'" /></el-icon>
                    <span>微生物指标检测</span>
                    <el-tag :type="cert.microbe_ok ? 'success' : 'danger'" size="small">{{ cert.microbe_ok ? '合格' : '不合格' }}</el-tag>
                  </div>
                </div>
              </div>

              <!-- 区块链存证信息 -->
              <div class="blockchain-info">
                <div class="bi-title">🔗 区块链存证信息</div>
                <div class="bi-hash">交易哈希：<span class="hash-text">{{ cert.tx_hash }}</span></div>
              </div>
            </el-col>

            <!-- 右侧二维码+印章 -->
            <el-col :span="8" class="cert-right">
              <div class="qr-area">
                <canvas ref="qrCanvas" class="qr-canvas"></canvas>
                <div class="qr-label">扫码溯源</div>
                <div class="qr-code-text">{{ cert.cert_no }}</div>
              </div>
              <!-- 防伪印章 -->
              <div class="seal-area">
                <div class="seal">
                  <div class="seal-inner">
                    <div class="seal-star">★</div>
                    <div class="seal-text-main">新疆阿克苏</div>
                    <div class="seal-text-sub">农业农村局</div>
                    <div class="seal-text-bottom">质量认证</div>
                  </div>
                </div>
                <div class="seal-label">官方防伪印章</div>
              </div>
            </el-col>
          </el-row>

          <!-- 底部声明 -->
          <div class="cert-footer">
            <p>本证书由 <strong>{{ cert.issue_org }}</strong> 依据国家农产品质量安全相关标准签发，并已在 Hyperledger Fabric 国密区块链上完成不可篡改存证。</p>
            <p>本证书电子版与纸质版具有同等法律效力。</p>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-else description="请选择批次查看电子合格证" style="margin-top:60px" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import QRCode from 'qrcode'
import { Connection } from '@element-plus/icons-vue'
import { getCertificate, listBatches } from '../../api/index'

const batches = ref<any[]>([])
const selectedBatchId = ref<number | null>(null)
const cert = ref<any>(null)
const qrCanvas = ref<HTMLCanvasElement>()

async function loadCert(batchId: number) {
  try {
    const res: any = await getCertificate(batchId)
    cert.value = res.data || res
    // 生成二维码
    setTimeout(() => {
      if (qrCanvas.value && cert.value) {
        const traceUrl = `${window.location.origin}/trace?code=${cert.value.cert_no}`
        QRCode.toCanvas(qrCanvas.value, traceUrl, { width: 140, margin: 1 })
      }
    }, 100)
  } catch {
    cert.value = null
  }
}

function printCert() {
  window.print()
}

onMounted(async () => {
  try {
    const res: any = await listBatches()
    batches.value = res.data || []
    if (batches.value.length > 0) {
      selectedBatchId.value = batches.value[0].id
      loadCert(batches.value[0].id)
    }
  } catch {}
})
</script>

<style scoped>
.page-header-card { margin-bottom: 0; }
.page-header { display: flex; justify-content: space-between; align-items: center; }
.page-title { margin: 0; font-size: 18px; font-weight: 600; }
.page-desc { margin: 4px 0 0; color: #888; font-size: 13px; }

.cert-wrapper { display: flex; justify-content: center; }
.cert-card {
  width: 900px; background: #fff; border-radius: 12px;
  padding: 32px; box-shadow: 0 4px 20px rgba(0,0,0,.1);
  border: 2px solid #e8f4e8;
}
.cert-header { display: flex; justify-content: space-between; align-items: center; }
.cert-title-area { display: flex; align-items: center; gap: 12px; }
.cert-gov-logo { font-size: 40px; }
.cert-main-title { font-size: 22px; font-weight: 700; color: #1a1a1a; margin: 0; }
.cert-sub-title { font-size: 12px; color: #888; margin: 2px 0 0; }
.cert-no-area { text-align: right; }
.cert-no { font-size: 14px; color: #555; }
.cert-blockchain-badge {
  display: inline-flex; align-items: center; gap: 4px;
  background: #e6f4ff; color: #1677ff; padding: 2px 8px;
  border-radius: 4px; font-size: 12px; margin-top: 4px;
}
.cert-desc { margin-bottom: 16px; }
.check-results { background: #f9fafb; border-radius: 8px; padding: 12px; margin-bottom: 12px; }
.check-title { font-size: 13px; font-weight: 600; color: #333; margin-bottom: 8px; }
.check-items { display: flex; flex-direction: column; gap: 6px; }
.check-item { display: flex; align-items: center; gap: 8px; font-size: 13px; color: #555; }
.check-item.pass { color: #52c41a; }
.blockchain-info { background: #f0f7ff; border-radius: 8px; padding: 10px 12px; }
.bi-title { font-size: 13px; font-weight: 600; color: #1677ff; margin-bottom: 4px; }
.bi-hash { font-size: 11px; color: #666; word-break: break-all; }
.hash-text { font-family: monospace; color: #1677ff; }

.cert-right { display: flex; flex-direction: column; align-items: center; gap: 20px; }
.qr-area { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.qr-canvas { border: 2px solid #e8e8e8; border-radius: 4px; }
.qr-label { font-size: 12px; color: #888; }
.qr-code-text { font-size: 10px; color: #aaa; font-family: monospace; }

.seal-area { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.seal {
  width: 100px; height: 100px; border-radius: 50%;
  border: 3px solid #cc0000; display: flex; align-items: center; justify-content: center;
  background: rgba(204,0,0,0.05);
}
.seal-inner { text-align: center; }
.seal-star { color: #cc0000; font-size: 20px; }
.seal-text-main { color: #cc0000; font-size: 11px; font-weight: 700; }
.seal-text-sub { color: #cc0000; font-size: 10px; }
.seal-text-bottom { color: #cc0000; font-size: 9px; }
.seal-label { font-size: 11px; color: #888; }

.cert-footer {
  margin-top: 20px; padding-top: 16px; border-top: 1px dashed #e8e8e8;
  font-size: 12px; color: #888; line-height: 1.8;
}
.cert-footer p { margin: 0; }

@media print {
  .page-header-card { display: none; }
  .cert-card { box-shadow: none; border: 1px solid #ccc; }
}
</style>
