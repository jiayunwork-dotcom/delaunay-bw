// delaunay-bw web console. Every number drawn here comes from the backend
// JSON API; the SVG mesh is rebuilt from /api/triangulate's triangle indices.
(() => {
  "use strict";

  const pointsInput = document.getElementById("points");
  const resultBox = document.getElementById("result");
  const errorBox = document.getElementById("error");
  const tolInput = document.getElementById("tol");
  const svg = document.getElementById("mesh");

  const NS = "http://www.w3.org/2000/svg";

  function showError(msg) {
    errorBox.textContent = msg;
    errorBox.hidden = false;
  }

  function clearError() {
    errorBox.hidden = true;
  }

  function parseTolerance() {
    const v = parseFloat(tolInput.value);
    return Number.isFinite(v) && v > 0 ? v : 1e-9;
  }

  function readPoints() {
    const raw = pointsInput.value.trim();
    if (!raw) throw new Error("点集为空");
    let data;
    try {
      data = JSON.parse(raw);
    } catch (e) {
      throw new Error("JSON 解析失败: " + e.message);
    }
    let pts = Array.isArray(data) ? data : data.points;
    if (!Array.isArray(pts) || pts.length === 0) throw new Error("缺少 points 数组");
    return pts;
  }

  async function post(path, body) {
    const resp = await fetch(path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    const text = await resp.text();
    let data;
    try {
      data = JSON.parse(text);
    } catch (e) {
      throw new Error("后端返回非 JSON: " + text.slice(0, 200));
    }
    if (!resp.ok) {
      // 失败路径：直接展示后端错误体
      throw new Error("HTTP " + resp.status + " 后端: " + (data.error || text));
    }
    return data;
  }

  function setPointsFromExample(data) {
    const pts = data.points || data;
    pointsInput.value = JSON.stringify(pts);
  }

  function drawMesh(pts, tris) {
    svg.innerHTML = "";
    const xs = pts.map((p) => p.x);
    const ys = pts.map((p) => p.y);
    const minX = Math.min(...xs), maxX = Math.max(...xs);
    const minY = Math.min(...ys), maxY = Math.max(...ys);
    const spanX = (maxX - minX) || 1, spanY = (maxY - minY) || 1;
    const pad = 24;
    const scale = Math.min(
      (500 - 2 * pad) / spanX,
      (380 - 2 * pad) / spanY,
    );
    const px = (p) => pad + (p.x - minX) * scale;
    const py = (p) => 380 - pad - (p.y - minY) * scale;

    const edges = new Set();
    for (const t of tris) {
      const a = pts[t[0]], b = pts[t[1]], c = pts[t[2]];
      const poly = document.createElementNS(NS, "polygon");
      poly.setAttribute("points",
        `${px(a)},${py(a)} ${px(b)},${py(b)} ${px(c)},${py(c)}`);
      poly.setAttribute("fill", "rgba(80,140,220,0.18)");
      poly.setAttribute("stroke", "#2c5aa0");
      poly.setAttribute("stroke-width", "1.2");
      svg.appendChild(poly);
      const es = [[t[0], t[1]], [t[1], t[2]], [t[2], t[0]]];
      for (const [u, v] of es) {
        edges.add(u < v ? `${u}-${v}` : `${v}-${u}`);
      }
    }
    for (let i = 0; i < pts.length; i++) {
      const dot = document.createElementNS(NS, "circle");
      dot.setAttribute("cx", px(pts[i]));
      dot.setAttribute("cy", py(pts[i]));
      dot.setAttribute("r", "2.6");
      dot.setAttribute("fill", "#b03030");
      svg.appendChild(dot);
    }
  }

  async function runTriangulate() {
    clearError();
    let pts;
    try {
      pts = readPoints();
    } catch (e) {
      showError(e.message);
      return;
    }
    try {
      const resp = await post("/api/triangulate", {
        points: pts,
        tolerance: parseTolerance(),
      });
      resultBox.textContent = JSON.stringify(resp, null, 2);
      drawMesh(pts, resp.triangles);
    } catch (e) {
      showError(e.message);
      resultBox.textContent = "计算失败";
    }
  }

  async function runVoronoi() {
    clearError();
    let pts;
    try {
      pts = readPoints();
    } catch (e) {
      showError(e.message);
      return;
    }
    try {
      const resp = await post("/api/voronoi", {
        points: pts,
        tolerance: parseTolerance(),
      });
      resultBox.textContent = JSON.stringify(resp, null, 2);
    } catch (e) {
      showError(e.message);
      resultBox.textContent = "计算失败";
    }
  }

  document.getElementById("load-example").addEventListener("click", async () => {
    clearError();
    try {
      const resp = await fetch("/example/grid-jitter.json");
      if (!resp.ok) throw new Error("加载示例失败: HTTP " + resp.status);
      const data = await resp.json();
      setPointsFromExample(data);
      resultBox.textContent = "示例已加载（" + data.points.length + " 点），点「计算三角剖分」。";
    } catch (e) {
      showError(e.message);
    }
  });

  document.getElementById("triangulate").addEventListener("click", runTriangulate);
  document.getElementById("voronoi").addEventListener("click", runVoronoi);
})();
