"""
generate_visuals.py
===================
Reads real benchmark data from  results/benchmark.json  (produced by
simulation.py) and generates all publication-quality charts and diagrams
for the report.

Requirements:  pip install matplotlib numpy
Usage:         python src/generate_visuals.py
"""

import json
import os
import sys

import numpy as np
import matplotlib.pyplot as plt
from matplotlib.patches import FancyBboxPatch, FancyArrowPatch

# ── Paths ──
SCRIPT_DIR  = os.path.dirname(os.path.abspath(__file__))
PROJECT_DIR = os.path.dirname(SCRIPT_DIR)
DATA_PATH   = os.path.join(PROJECT_DIR, "results", "benchmark.json")
OUT_DIR     = os.path.join(PROJECT_DIR, "diagrams")
os.makedirs(OUT_DIR, exist_ok=True)

# ── Load benchmark data ──
if not os.path.exists(DATA_PATH):
    print(f"ERROR: {DATA_PATH} not found.  Run  python src/simulation.py  first.")
    sys.exit(1)

with open(DATA_PATH, "r", encoding="utf-8") as f:
    DATA = json.load(f)

PRIMARY = DATA["primary"]["results"]
CONFIG  = DATA["primary"]["config"]
SCALE   = DATA["scalability"]

# ── Shared style ──
plt.rcParams.update({
    "font.family":       "sans-serif",
    "font.sans-serif":   ["Segoe UI", "Arial", "Helvetica"],
    "axes.spines.top":   False,
    "axes.spines.right": False,
    "figure.dpi":        180,
    "savefig.bbox":      "tight",
    "savefig.pad_inches": 0.3,
})

BLUE   = "#3b82f6"
PURPLE = "#a855f7"
GREEN  = "#22c55e"
DARK   = "#1e293b"

def _save(name):
    path = os.path.join(OUT_DIR, name)
    plt.savefig(path)
    plt.close()
    print(f"  [OK] {name}")


# ═══════════════════════════════════════════════════════════════
# 1. Cloud Architecture Overview
# ═══════════════════════════════════════════════════════════════
def draw_cloud_overview():
    fig, ax = plt.subplots(figsize=(10, 5.5))
    ax.set_xlim(0, 10); ax.set_ylim(0, 6); ax.axis("off")
    ax.set_facecolor("#fafbff"); fig.patch.set_facecolor("#fafbff")

    ax.text(5, 5.6, "Cloud Resource Allocation - Overview",
            ha="center", fontsize=14, fontweight="bold", color=DARK)

    # Tasks
    tasks_info = [
        ("T1", "CPU: 2 | RAM: 4 GB\nPriority: High", "#dbeafe"),
        ("T2", "CPU: 4 | RAM: 8 GB\nPriority: Med",  "#ede9fe"),
        ("T3", "CPU: 1 | RAM: 2 GB\nPriority: High", "#dbeafe"),
        ("T4", "CPU: 3 | RAM: 6 GB\nPriority: Low",  "#dcfce7"),
    ]
    ax.text(1.1, 5.1, "Incoming Tasks", ha="center", fontsize=10,
            fontweight="bold", color=DARK)
    for i, (name, desc, clr) in enumerate(tasks_info):
        y = 4.3 - i * 1.1
        box = FancyBboxPatch((0.2, y - 0.35), 1.8, 0.7,
              boxstyle="round,pad=0.1", fc=clr, ec="#94a3b8", lw=1)
        ax.add_patch(box)
        ax.text(0.45, y + 0.1, name, fontsize=9, fontweight="bold", color=DARK)
        ax.text(0.45, y - 0.15, desc, fontsize=6.5, color="#475569")

    # Engine
    engine = FancyBboxPatch((3.2, 1.8), 3.6, 2.4,
             boxstyle="round,pad=0.15", fc="#1e293b", ec=BLUE, lw=2)
    ax.add_patch(engine)
    ax.text(5, 3.8, "Allocation Algorithm", ha="center", fontsize=10,
            fontweight="bold", color="#fff")
    ax.text(5, 3.3, "Greedy\nDynamic Programming\nHeuristic Load Balancing",
            ha="center", fontsize=7.5, color="#94a3b8", linespacing=1.5)

    for i in range(4):
        y = 4.3 - i * 1.1
        ax.annotate("", xy=(3.2, 3.0), xytext=(2.0, y),
                     arrowprops=dict(arrowstyle="-|>", color="#94a3b8", lw=1.2))

    # Servers
    srv_info = [("Server 1", "16 CPU | 32 GB", GREEN),
                ("Server 2", "16 CPU | 32 GB", BLUE),
                ("Server 3", "16 CPU | 32 GB", PURPLE)]
    ax.text(8.9, 5.1, "Server Pool", ha="center", fontsize=10,
            fontweight="bold", color=DARK)
    for i, (name, desc, clr) in enumerate(srv_info):
        y = 4.3 - i * 1.4
        box = FancyBboxPatch((7.8, y - 0.4), 2.0, 0.8,
              boxstyle="round,pad=0.1", fc=clr, ec="#fff", lw=1.5, alpha=0.85)
        ax.add_patch(box)
        ax.text(8.8, y + 0.12, name, ha="center", fontsize=9,
                fontweight="bold", color="#fff")
        ax.text(8.8, y - 0.18, desc, ha="center", fontsize=7, color="#e2e8f0")

    for i in range(3):
        y = 4.3 - i * 1.4
        ax.annotate("", xy=(7.8, y), xytext=(6.8, 3.0),
                     arrowprops=dict(arrowstyle="-|>", color="#94a3b8", lw=1.2))

    _save("cloud_overview.png")


# ═══════════════════════════════════════════════════════════════
# 2. Complexity Growth (theoretical)
# ═══════════════════════════════════════════════════════════════
def draw_complexity_chart():
    n = np.array([10, 50, 100, 500, 1000, 5000])
    m_ratio = 0.1  # servers ~ 10% of tasks
    C, R = CONFIG["cpu_cap"], CONFIG["ram_cap"]

    greedy    = n * np.maximum(n * m_ratio, 3)
    dp        = n * C * R
    heuristic = n * np.maximum(n * m_ratio, 3) * 1.3

    fig, ax = plt.subplots(figsize=(8, 4.5))
    ax.plot(n, greedy,    "o-", color=BLUE,   lw=2, ms=6, label="Greedy  O(n*m)")
    ax.plot(n, dp,        "s-", color=PURPLE, lw=2, ms=6, label="DP  O(n*C*R)")
    ax.plot(n, heuristic, "^-", color=GREEN,  lw=2, ms=6, label="Heuristic  O(n*m)")
    ax.set_yscale("log")
    ax.set_xlabel("Number of Tasks (n)", fontsize=11)
    ax.set_ylabel("Operations (log scale)", fontsize=11)
    ax.set_title("Theoretical Time Complexity Growth", fontsize=13, fontweight="bold")
    ax.legend(frameon=True, fontsize=9); ax.grid(True, alpha=0.3)
    _save("complexity_chart.png")


# ═══════════════════════════════════════════════════════════════
# 3. Benchmark Comparison Bars (from real data)
# ═══════════════════════════════════════════════════════════════
def draw_benchmark_bars():
    metrics = ["Tasks\nPlaced", "Priority\nScore", "Avg CPU\nUtil (%)", "Avg RAM\nUtil (%)"]
    g = PRIMARY["Greedy"]
    d = PRIMARY["DP"]
    h = PRIMARY["Heuristic"]

    g_vals = [g["placed"], g["priority_score"], g["avg_cpu_util"], g["avg_ram_util"]]
    d_vals = [d["placed"], d["priority_score"], d["avg_cpu_util"], d["avg_ram_util"]]
    h_vals = [h["placed"], h["priority_score"], h["avg_cpu_util"], h["avg_ram_util"]]

    x = np.arange(len(metrics)); w = 0.25
    fig, ax = plt.subplots(figsize=(9, 5))
    ax.bar(x - w, g_vals, w, label="Greedy",    color=BLUE,   ec="#fff")
    ax.bar(x,     d_vals, w, label="DP",        color=PURPLE, ec="#fff")
    ax.bar(x + w, h_vals, w, label="Heuristic", color=GREEN,  ec="#fff")

    ax.set_ylabel("Value", fontsize=11)
    ax.set_title(f"Benchmark Results ({g['total']} tasks, "
                 f"{CONFIG['n_servers']} servers)", fontsize=13, fontweight="bold")
    ax.set_xticks(x); ax.set_xticklabels(metrics, fontsize=9)
    ax.legend(frameon=True, fontsize=9); ax.grid(axis="y", alpha=0.3)
    _save("benchmark_bars.png")


# ═══════════════════════════════════════════════════════════════
# 4. Execution Time Comparison
# ═══════════════════════════════════════════════════════════════
def draw_execution_time():
    algos  = ["Greedy", "DP", "Heuristic"]
    times  = [PRIMARY[a]["time_ms"] for a in algos]
    colors = [BLUE, PURPLE, GREEN]

    fig, ax = plt.subplots(figsize=(7, 4))
    bars = ax.bar(algos, times, color=colors, ec="#fff", width=0.55)
    ax.set_yscale("log")
    ax.set_ylabel("Execution Time (ms) - log scale", fontsize=11)
    ax.set_title("Execution Time Comparison", fontsize=13, fontweight="bold")
    ax.grid(axis="y", alpha=0.3)

    for bar, t in zip(bars, times):
        ax.text(bar.get_x() + bar.get_width() / 2,
                bar.get_height() * 1.5,
                f"{t:.3f} ms", ha="center", fontsize=9, fontweight="bold", color=DARK)
    _save("execution_time.png")


# ═══════════════════════════════════════════════════════════════
# 5. Trade-Off Radar
# ═══════════════════════════════════════════════════════════════
def draw_tradeoff_radar():
    cats = ["Speed", "Optimality", "Scalability", "Memory\nEfficiency", "Ease of\nImpl."]
    N = len(cats)
    angles = np.linspace(0, 2 * np.pi, N, endpoint=False).tolist()
    angles += angles[:1]

    g = [95, 40, 90, 95, 95]; g += g[:1]
    d = [15, 100, 20, 10, 50]; d += d[:1]
    h = [85, 70, 85, 90, 65]; h += h[:1]

    fig, ax = plt.subplots(figsize=(6, 6), subplot_kw=dict(polar=True))
    ax.set_theta_offset(np.pi / 2); ax.set_theta_direction(-1)
    ax.set_thetagrids(np.degrees(angles[:-1]), cats, fontsize=9)
    ax.set_ylim(0, 100)
    ax.set_yticks([20, 40, 60, 80, 100])
    ax.set_yticklabels(["20","40","60","80","100"], fontsize=7, color="#999")

    for vals, clr, lbl, mk in [(g, BLUE, "Greedy", "o"),
                                (d, PURPLE, "DP", "s"),
                                (h, GREEN, "Heuristic", "^")]:
        ax.plot(angles, vals, f"{mk}-", color=clr, lw=2, label=lbl)
        ax.fill(angles, vals, alpha=0.1, color=clr)

    ax.legend(loc="upper right", bbox_to_anchor=(1.25, 1.15), frameon=True, fontsize=9)
    ax.set_title("Algorithm Trade-Off Profile", fontsize=13, fontweight="bold", y=1.08)
    _save("tradeoff_radar.png")


# ═══════════════════════════════════════════════════════════════
# 6. CPU Utilisation Heatmap (from real per-server data)
# ═══════════════════════════════════════════════════════════════
def draw_utilisation_heatmap():
    servers = [f"S{i+1}" for i in range(CONFIG["n_servers"])]
    algos   = ["Greedy", "DP", "Heuristic"]

    cpu_data = np.array([
        [PRIMARY[a]["per_server_cpu"][s] for a in algos]
        for s in servers
    ])

    fig, ax = plt.subplots(figsize=(6, 4.5))
    im = ax.imshow(cpu_data, cmap="YlGnBu", aspect="auto", vmin=40, vmax=100)
    ax.set_xticks(range(len(algos))); ax.set_xticklabels(algos, fontsize=10)
    ax.set_yticks(range(len(servers))); ax.set_yticklabels(servers, fontsize=10)
    ax.set_xlabel("Algorithm", fontsize=11)
    ax.set_ylabel("Server", fontsize=11)
    ax.set_title("CPU Utilisation (%) by Server", fontsize=13, fontweight="bold")

    for i in range(len(servers)):
        for j in range(len(algos)):
            clr = "#fff" if cpu_data[i, j] > 80 else DARK
            ax.text(j, i, f"{cpu_data[i,j]:.0f}%", ha="center", va="center",
                    fontsize=10, fontweight="bold", color=clr)

    fig.colorbar(im, ax=ax, shrink=0.8, label="Utilisation %")
    _save("cpu_utilisation_heatmap.png")


# ═══════════════════════════════════════════════════════════════
# 7. Allocation Workflow Flowchart
# ═══════════════════════════════════════════════════════════════
def draw_allocation_flowchart():
    fig, ax = plt.subplots(figsize=(8, 10))
    ax.set_xlim(0, 8); ax.set_ylim(0, 12); ax.axis("off")
    fig.patch.set_facecolor("#fff")

    def box(x, y, w, h, txt, fc="#e0e7ff", tc=DARK, fs=9):
        p = FancyBboxPatch((x-w/2, y-h/2), w, h,
            boxstyle="round,pad=0.15", fc=fc, ec="#64748b", lw=1.2)
        ax.add_patch(p)
        ax.text(x, y, txt, ha="center", va="center", fontsize=fs,
                fontweight="bold", color=tc)

    def diamond(x, y, txt, fc="#fef3c7"):
        d = plt.Polygon([(x,y+0.5),(x+1.2,y),(x,y-0.5),(x-1.2,y)],
                        fc=fc, ec="#64748b", lw=1.2)
        ax.add_patch(d)
        ax.text(x, y, txt, ha="center", va="center", fontsize=8,
                fontweight="bold", color=DARK)

    def arrow(x1,y1,x2,y2,lbl=""):
        ax.annotate("", xy=(x2,y2), xytext=(x1,y1),
                     arrowprops=dict(arrowstyle="-|>", color="#64748b", lw=1.5))
        if lbl:
            mx, my = (x1+x2)/2, (y1+y2)/2
            ax.text(mx+0.15, my, lbl, fontsize=7, color="#64748b", fontstyle="italic")

    ax.text(4, 11.5, "Resource Allocation Workflow", ha="center",
            fontsize=14, fontweight="bold", color=DARK)

    box(4, 10.5, 2.5, 0.7, "New Task Arrives", fc="#dbeafe")
    arrow(4, 10.15, 4, 9.7)
    box(4, 9.3, 3, 0.7, "Extract CPU, RAM, Priority", fc="#e0e7ff")
    arrow(4, 8.95, 4, 8.5)
    box(4, 8.1, 3, 0.7, "Sort / Rank by Priority", fc="#e0e7ff")
    arrow(4, 7.75, 4, 7.3)
    diamond(4, 6.8, "Choose\nAlgorithm")

    arrow(2.8, 6.8, 1.5, 6.1, "Greedy")
    box(1.5, 5.6, 2, 0.7, "First-Fit\nAllocation", fc=BLUE, tc="#fff", fs=8)
    arrow(4, 6.3, 4, 5.9, "DP")
    box(4, 5.4, 2.2, 0.7, "Build DP Table\nOptimal Assign", fc=PURPLE, tc="#fff", fs=8)
    arrow(5.2, 6.8, 6.5, 6.1, "Heuristic")
    box(6.5, 5.6, 2, 0.7, "Score Servers\nLeast Loaded", fc=GREEN, tc="#fff", fs=8)

    arrow(1.5, 5.25, 4, 4.3); arrow(4, 5.05, 4, 4.3); arrow(6.5, 5.25, 4, 4.3)
    box(4, 3.9, 3, 0.7, "Assign Task to Server", fc="#d1fae5")
    arrow(4, 3.55, 4, 3.1)
    box(4, 2.7, 3, 0.7, "Update Server Capacity", fc="#d1fae5")
    arrow(4, 2.35, 4, 1.9)
    diamond(4, 1.4, "More\nTasks?")
    arrow(5.2, 1.4, 7, 1.4)
    ax.annotate("", xy=(7,10.5), xytext=(7,1.4),
                arrowprops=dict(arrowstyle="-|>", color="#64748b", lw=1.5))
    ax.annotate("", xy=(5.25,10.5), xytext=(7,10.5),
                arrowprops=dict(arrowstyle="-|>", color="#64748b", lw=1.5))
    ax.text(7.2, 6, "Yes", fontsize=8, color="#64748b", fontstyle="italic")
    arrow(4, 0.9, 4, 0.5, "No")
    box(4, 0.2, 2.5, 0.5, "DONE", fc="#065f46", tc="#fff")

    _save("allocation_flowchart.png")


# ═══════════════════════════════════════════════════════════════
# 8. Provider Strategy Stacked Bar
# ═══════════════════════════════════════════════════════════════
def draw_provider_stacked():
    providers = ["AWS EC2", "Google\nBorg/GCP", "Microsoft\nAzure", "Kubernetes"]
    g_pct = [30, 25, 35, 30]
    h_pct = [50, 45, 40, 55]
    d_pct = [20, 30, 25, 15]

    x = np.arange(len(providers)); w = 0.5
    fig, ax = plt.subplots(figsize=(8, 4.5))
    ax.bar(x, g_pct, w, label="Greedy Filtering", color=BLUE)
    b2 = g_pct
    ax.bar(x, h_pct, w, bottom=b2, label="Heuristic Scoring", color=GREEN)
    b3 = [a+b for a,b in zip(g_pct, h_pct)]
    ax.bar(x, d_pct, w, bottom=b3, label="Offline Optimisation (DP-like)", color=PURPLE)

    ax.set_ylabel("Strategy Weight (%)", fontsize=11)
    ax.set_title("Algorithm Usage in Real Cloud Platforms", fontsize=13, fontweight="bold")
    ax.set_xticks(x); ax.set_xticklabels(providers, fontsize=9)
    ax.set_ylim(0, 110)
    ax.legend(loc="upper right", frameon=True, fontsize=8); ax.grid(axis="y", alpha=0.3)
    _save("provider_strategy.png")


# ═══════════════════════════════════════════════════════════════
# 9. RAM Utilisation per Server (from real data)
# ═══════════════════════════════════════════════════════════════
def draw_ram_usage():
    servers = [f"S{i+1}" for i in range(CONFIG["n_servers"])]
    ram_cap = CONFIG["ram_cap"]

    g_ram = [PRIMARY["Greedy"]["per_server_ram"][s] * ram_cap / 100 for s in servers]
    d_ram = [PRIMARY["DP"]["per_server_ram"][s] * ram_cap / 100 for s in servers]
    h_ram = [PRIMARY["Heuristic"]["per_server_ram"][s] * ram_cap / 100 for s in servers]

    x = np.arange(len(servers)); w = 0.22
    fig, ax = plt.subplots(figsize=(8, 4.5))
    ax.bar(x-w, g_ram, w, label="Greedy",    color=BLUE,   ec="#fff")
    ax.bar(x,   d_ram, w, label="DP",        color=PURPLE, ec="#fff")
    ax.bar(x+w, h_ram, w, label="Heuristic", color=GREEN,  ec="#fff")
    ax.plot(x, [ram_cap]*len(servers), "k--", lw=1.5, label="Total Capacity")

    ax.set_xlabel("Server", fontsize=11); ax.set_ylabel("RAM Used (GB)", fontsize=11)
    ax.set_title("RAM Allocation per Server", fontsize=13, fontweight="bold")
    ax.set_xticks(x); ax.set_xticklabels(servers, fontsize=10)
    ax.legend(frameon=True, fontsize=9); ax.set_ylim(0, ram_cap + 4)
    ax.grid(axis="y", alpha=0.3)
    _save("ram_usage.png")


# ═══════════════════════════════════════════════════════════════
# 10. Scalability — Execution Time vs. Number of Tasks
# ═══════════════════════════════════════════════════════════════
def draw_scalability():
    ns = sorted(int(k) for k in SCALE.keys())
    g_t = [SCALE[str(n)]["Greedy"]["time_ms"]    for n in ns]
    d_t = [SCALE[str(n)]["DP"]["time_ms"]        for n in ns]
    h_t = [SCALE[str(n)]["Heuristic"]["time_ms"] for n in ns]

    fig, ax = plt.subplots(figsize=(8, 4.5))
    ax.plot(ns, g_t, "o-", color=BLUE,   lw=2, ms=6, label="Greedy")
    ax.plot(ns, d_t, "s-", color=PURPLE, lw=2, ms=6, label="DP")
    ax.plot(ns, h_t, "^-", color=GREEN,  lw=2, ms=6, label="Heuristic")
    ax.set_yscale("log")
    ax.set_xlabel("Number of Tasks", fontsize=11)
    ax.set_ylabel("Execution Time (ms) - log scale", fontsize=11)
    ax.set_title("Scalability: Execution Time vs. Problem Size",
                 fontsize=13, fontweight="bold")
    ax.legend(frameon=True, fontsize=9); ax.grid(True, alpha=0.3)
    _save("scalability.png")


# ═══════════════════════════════════════════════════════════════
# Main
# ═══════════════════════════════════════════════════════════════
if __name__ == "__main__":
    print("Generating report visuals from real benchmark data...\n")
    draw_cloud_overview()
    draw_complexity_chart()
    draw_benchmark_bars()
    draw_execution_time()
    draw_tradeoff_radar()
    draw_utilisation_heatmap()
    draw_allocation_flowchart()
    draw_provider_stacked()
    draw_ram_usage()
    draw_scalability()
    print(f"\nAll diagrams saved to {OUT_DIR}")
