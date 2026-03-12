#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
阿克苏苹果区块链溯源系统 - 性能测试图表生成脚本
运行环境: Python 3.8+
依赖安装: pip install matplotlib pandas numpy seaborn
运行方式: python plot_performance.py
输出目录: ./output/
"""

import os
import pandas as pd
import numpy as np
import matplotlib
import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib.gridspec import GridSpec
import warnings

warnings.filterwarnings('ignore')
matplotlib.rcParams['font.family'] = ['SimHei', 'Microsoft YaHei', 'STHeiti', 'DejaVu Sans']
matplotlib.rcParams['axes.unicode_minus'] = False

# 输出目录
OUTPUT_DIR = os.path.join(os.path.dirname(__file__), 'output')
os.makedirs(OUTPUT_DIR, exist_ok=True)

# ============================================================
# 颜色方案
# ============================================================
COLORS = {
    'blue':   '#2E86AB',
    'orange': '#E84855',
    'green':  '#3BB273',
    'purple': '#7B2D8B',
    'gray':   '#8E9AAF',
    'yellow': '#F4A261',
    'teal':   '#2EC4B6',
    'red':    '#E63946',
}

# ============================================================
# 图1: TPS vs 并发数（三场景对比）
# ============================================================
def plot_tps_vs_concurrency():
    df = pd.read_csv(os.path.join(os.path.dirname(__file__), 'data/caliper_performance.csv'))
    scenarios = df['测试场景'].unique()
    concurrencies = sorted(df['并发数'].unique())

    fig, ax = plt.subplots(figsize=(10, 6))
    markers = ['o', 's', '^']
    colors_list = [COLORS['blue'], COLORS['green'], COLORS['orange']]

    for i, scenario in enumerate(scenarios):
        sub = df[df['测试场景'] == scenario].sort_values('并发数')
        ax.plot(sub['并发数'], sub['TPS(笔/秒)'],
                marker=markers[i], color=colors_list[i], linewidth=2.5,
                markersize=8, label=scenario)

    ax.set_xlabel('并发用户数（个）', fontsize=13)
    ax.set_ylabel('吞吐量 TPS（笔/秒）', fontsize=13)
    ax.set_title('图5-8  系统吞吐量（TPS）随并发数变化曲线', fontsize=14, fontweight='bold', pad=15)
    ax.set_xticks(concurrencies)
    ax.legend(fontsize=11, loc='upper left')
    ax.grid(True, linestyle='--', alpha=0.5)
    ax.set_ylim(0, 2000)
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)

    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图5-8_TPS并发曲线.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 图2: 平均延迟 vs 并发数
# ============================================================
def plot_latency_vs_concurrency():
    df = pd.read_csv(os.path.join(os.path.dirname(__file__), 'data/caliper_performance.csv'))
    scenarios = df['测试场景'].unique()
    concurrencies = sorted(df['并发数'].unique())

    fig, ax = plt.subplots(figsize=(10, 6))
    markers = ['o', 's', '^']
    colors_list = [COLORS['blue'], COLORS['green'], COLORS['orange']]
    linestyles = ['-', '--', '-.']

    for i, scenario in enumerate(scenarios):
        sub = df[df['测试场景'] == scenario].sort_values('并发数')
        ax.plot(sub['并发数'], sub['平均延迟(ms)'],
                marker=markers[i], color=colors_list[i], linewidth=2.5,
                markersize=8, linestyle=linestyles[i], label=f'{scenario}（平均）')
        ax.fill_between(sub['并发数'],
                        sub['平均延迟(ms)'] * 0.85,
                        sub['P95延迟(ms)'],
                        alpha=0.1, color=colors_list[i])

    ax.set_xlabel('并发用户数（个）', fontsize=13)
    ax.set_ylabel('响应延迟（ms）', fontsize=13)
    ax.set_title('图5-9  系统响应延迟随并发数变化曲线', fontsize=14, fontweight='bold', pad=15)
    ax.set_xticks(concurrencies)
    ax.legend(fontsize=10, loc='upper left')
    ax.grid(True, linestyle='--', alpha=0.5)
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)

    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图5-9_延迟并发曲线.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 图3: 国密算法性能对比柱状图
# ============================================================
def plot_gmsm_comparison():
    # 签名性能对比
    algorithms_sign = ['RSA-2048', 'RSA-3072', 'ECDSA-P256', 'SM2（国密）']
    sign_times = [8.92, 18.45, 0.89, 0.76]
    verify_times = [0.23, 0.31, 1.12, 0.94]

    # 哈希性能对比
    algorithms_hash = ['SHA-256', 'SM3（国密）']
    hash_times = [0.045, 0.038]

    fig, axes = plt.subplots(1, 2, figsize=(14, 6))

    # 子图1: 签名/验签性能
    x = np.arange(len(algorithms_sign))
    width = 0.35
    bars1 = axes[0].bar(x - width/2, sign_times, width, label='签名耗时(ms)',
                         color=COLORS['blue'], alpha=0.85, edgecolor='white')
    bars2 = axes[0].bar(x + width/2, verify_times, width, label='验签耗时(ms)',
                         color=COLORS['orange'], alpha=0.85, edgecolor='white')

    for bar in bars1:
        h = bar.get_height()
        axes[0].text(bar.get_x() + bar.get_width()/2., h + 0.3,
                     f'{h:.2f}', ha='center', va='bottom', fontsize=9)
    for bar in bars2:
        h = bar.get_height()
        axes[0].text(bar.get_x() + bar.get_width()/2., h + 0.3,
                     f'{h:.2f}', ha='center', va='bottom', fontsize=9)

    axes[0].set_xlabel('加密算法', fontsize=12)
    axes[0].set_ylabel('耗时（ms）', fontsize=12)
    axes[0].set_title('(a) 数字签名算法性能对比', fontsize=12, fontweight='bold')
    axes[0].set_xticks(x)
    axes[0].set_xticklabels(algorithms_sign, fontsize=10)
    axes[0].legend(fontsize=10)
    axes[0].grid(True, axis='y', linestyle='--', alpha=0.5)
    axes[0].spines['top'].set_visible(False)
    axes[0].spines['right'].set_visible(False)

    # 子图2: 哈希算法性能
    bar_colors = [COLORS['gray'], COLORS['green']]
    bars3 = axes[1].bar(algorithms_hash, hash_times, color=bar_colors, alpha=0.85,
                         width=0.4, edgecolor='white')
    for bar, val in zip(bars3, hash_times):
        axes[1].text(bar.get_x() + bar.get_width()/2., val + 0.001,
                     f'{val:.3f}ms', ha='center', va='bottom', fontsize=11, fontweight='bold')

    axes[1].set_xlabel('哈希算法', fontsize=12)
    axes[1].set_ylabel('单次哈希耗时（ms）', fontsize=12)
    axes[1].set_title('(b) 哈希算法性能对比', fontsize=12, fontweight='bold')
    axes[1].grid(True, axis='y', linestyle='--', alpha=0.5)
    axes[1].spines['top'].set_visible(False)
    axes[1].spines['right'].set_visible(False)

    # 添加性能提升标注
    axes[1].annotate('SM3比SHA-256\n快约15.6%',
                     xy=(1, hash_times[1]), xytext=(0.5, 0.04),
                     fontsize=10, color=COLORS['green'],
                     arrowprops=dict(arrowstyle='->', color=COLORS['green']),
                     ha='center')

    fig.suptitle('图5-10  国密算法与国际算法性能对比', fontsize=14, fontweight='bold', y=1.02)
    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图5-10_国密算法性能对比.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 图4: 系统资源占用折线图（CPU+内存）
# ============================================================
def plot_resource_usage():
    df = pd.read_csv(os.path.join(os.path.dirname(__file__), 'data/caliper_performance.csv'))
    # 取"溯源记录上链"场景
    sub = df[df['测试场景'] == '溯源记录上链'].sort_values('并发数')

    fig, ax1 = plt.subplots(figsize=(10, 6))
    ax2 = ax1.twinx()

    line1 = ax1.plot(sub['并发数'], sub['CPU使用率(%)'],
                     marker='o', color=COLORS['blue'], linewidth=2.5,
                     markersize=8, label='CPU使用率(%)')
    line2 = ax2.plot(sub['并发数'], sub['内存使用(MB)'],
                     marker='s', color=COLORS['orange'], linewidth=2.5,
                     markersize=8, linestyle='--', label='内存使用(MB)')

    ax1.set_xlabel('并发用户数（个）', fontsize=13)
    ax1.set_ylabel('CPU使用率（%）', fontsize=13, color=COLORS['blue'])
    ax2.set_ylabel('内存使用量（MB）', fontsize=13, color=COLORS['orange'])
    ax1.set_title('图5-11  系统资源占用随并发数变化曲线（溯源上链场景）',
                  fontsize=14, fontweight='bold', pad=15)
    ax1.set_xticks(sub['并发数'])
    ax1.tick_params(axis='y', labelcolor=COLORS['blue'])
    ax2.tick_params(axis='y', labelcolor=COLORS['orange'])
    ax1.set_ylim(0, 110)
    ax2.set_ylim(0, 1800)

    lines = line1 + line2
    labels = [l.get_label() for l in lines]
    ax1.legend(lines, labels, fontsize=11, loc='upper left')
    ax1.grid(True, linestyle='--', alpha=0.4)
    ax1.spines['top'].set_visible(False)

    # 添加警戒线
    ax1.axhline(y=80, color='red', linestyle=':', alpha=0.6, linewidth=1.5)
    ax1.text(200, 81, 'CPU 80% 警戒线', color='red', fontsize=9, ha='right')

    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图5-11_系统资源占用.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 图5: 区块链上链开销分析饼图
# ============================================================
def plot_blockchain_overhead():
    labels = ['Fabric共识耗时', '国密SM2签名', '国密SM3哈希', 'IPFS存储', 'MySQL写入', '网络传输']
    sizes  = [38.5, 22.3, 8.7, 15.6, 8.2, 6.7]
    colors_pie = [COLORS['blue'], COLORS['orange'], COLORS['green'],
                  COLORS['purple'], COLORS['yellow'], COLORS['gray']]
    explode = (0.05, 0.05, 0, 0, 0, 0)

    fig, ax = plt.subplots(figsize=(9, 7))
    wedges, texts, autotexts = ax.pie(
        sizes, explode=explode, labels=labels, colors=colors_pie,
        autopct='%1.1f%%', startangle=140,
        textprops={'fontsize': 11},
        wedgeprops={'edgecolor': 'white', 'linewidth': 2}
    )
    for autotext in autotexts:
        autotext.set_fontsize(10)
        autotext.set_fontweight('bold')

    ax.set_title('图5-12  溯源上链操作各环节耗时占比分析\n（平均总耗时: 320ms @ 50并发）',
                 fontsize=13, fontweight='bold', pad=20)
    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图5-12_上链耗时占比.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 图6: 综合性能雷达图
# ============================================================
def plot_radar_comparison():
    categories = ['吞吐量\n(TPS)', '响应速度\n(低延迟)', '安全性\n(国密)', '可靠性\n(成功率)', '可扩展性', '数据完整性']
    N = len(categories)

    # 本系统 vs 传统中心化系统
    values_ours = [7.8, 7.2, 9.5, 9.2, 8.0, 9.8]
    values_trad = [9.5, 9.2, 6.0, 7.5, 6.5, 6.8]

    values_ours += values_ours[:1]
    values_trad += values_trad[:1]

    angles = [n / float(N) * 2 * np.pi for n in range(N)]
    angles += angles[:1]

    fig, ax = plt.subplots(figsize=(8, 8), subplot_kw=dict(polar=True))

    ax.plot(angles, values_ours, 'o-', linewidth=2.5, color=COLORS['blue'], label='本系统（区块链+国密）')
    ax.fill(angles, values_ours, alpha=0.2, color=COLORS['blue'])
    ax.plot(angles, values_trad, 's--', linewidth=2.5, color=COLORS['orange'], label='传统中心化溯源系统')
    ax.fill(angles, values_trad, alpha=0.15, color=COLORS['orange'])

    ax.set_xticks(angles[:-1])
    ax.set_xticklabels(categories, fontsize=12)
    ax.set_ylim(0, 10)
    ax.set_yticks([2, 4, 6, 8, 10])
    ax.set_yticklabels(['2', '4', '6', '8', '10'], fontsize=9)
    ax.grid(True, linestyle='--', alpha=0.5)
    ax.set_title('图5-13  本系统与传统溯源系统综合性能雷达图',
                 fontsize=13, fontweight='bold', pad=30)
    ax.legend(loc='upper right', bbox_to_anchor=(1.35, 1.1), fontsize=11)

    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图5-13_综合性能雷达图.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 图7: 苹果产量趋势图（参考图2样式）
# ============================================================
def plot_apple_production():
    years = ['2014年','2015年','2016年','2017年','2018年','2019年','2020年','2021年','2022年','2023年','2024年']
    production = [3735.39, 3889.90, 4039.33, 4139.00, 3923.34, 4242.54, 4406.61, 4597.34, 4757.18, 4960.17, 5128.51]
    growth_rate = [None, 4.14, 3.84, 2.47, -5.21, 8.14, 3.87, 4.33, 3.48, 4.27, 3.39]

    fig, ax1 = plt.subplots(figsize=(13, 7))
    ax2 = ax1.twinx()

    bars = ax1.bar(years, production, color='#F9B8C8', edgecolor='#e8a0b0', linewidth=0.8, width=0.6, label='苹果产量（万吨）')
    for bar, val in zip(bars, production):
        ax1.text(bar.get_x() + bar.get_width()/2., bar.get_height() + 20,
                 f'{val}', ha='center', va='bottom', fontsize=9, fontweight='bold', color='#333')

    valid_x = [years[i] for i in range(1, len(years))]
    valid_y = [growth_rate[i] for i in range(1, len(growth_rate))]
    ax2.plot(valid_x, valid_y, 'o-', color='#3CB371', linewidth=2.5, markersize=8,
             label='产量增速（%）', zorder=5)
    for x, y in zip(valid_x, valid_y):
        offset = 0.15 if y >= 0 else -0.4
        ax2.text(x, y + offset, f'{y}%', ha='center', va='bottom', fontsize=9,
                 color='#2d8653', fontweight='bold')

    ax2.axhline(y=0, color='#999', linestyle='--', linewidth=1, alpha=0.6)

    ax1.set_xlabel('年份', fontsize=13)
    ax1.set_ylabel('苹果产量（万吨）', fontsize=13, color='#c0607a')
    ax2.set_ylabel('产量增速（%）', fontsize=13, color='#2d8653')
    ax1.set_title('图2-1  2014—2024年中国苹果产量及增速变化趋势', fontsize=14, fontweight='bold', pad=15)
    ax1.tick_params(axis='y', labelcolor='#c0607a')
    ax2.tick_params(axis='y', labelcolor='#2d8653')
    ax1.set_ylim(0, 6000)
    ax2.set_ylim(-8, 12)

    lines1, labels1 = ax1.get_legend_handles_labels()
    lines2, labels2 = ax2.get_legend_handles_labels()
    ax1.legend(lines1 + lines2, labels1 + labels2, loc='upper left', fontsize=11)
    ax1.grid(True, axis='y', linestyle='--', alpha=0.3)
    ax1.spines['top'].set_visible(False)

    plt.tight_layout()
    path = os.path.join(OUTPUT_DIR, '图2-1_苹果产量趋势.png')
    plt.savefig(path, dpi=150, bbox_inches='tight')
    plt.close()
    print(f'[OK] 已生成: {path}')

# ============================================================
# 主函数
# ============================================================
if __name__ == '__main__':
    print('=' * 60)
    print('  阿克苏苹果区块链溯源系统 - 性能图表生成脚本')
    print('=' * 60)
    plot_tps_vs_concurrency()
    plot_latency_vs_concurrency()
    plot_gmsm_comparison()
    plot_resource_usage()
    plot_blockchain_overhead()
    plot_radar_comparison()
    plot_apple_production()
    print('=' * 60)
    print(f'  全部图表已生成至: {OUTPUT_DIR}')
    print('=' * 60)
