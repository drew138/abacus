htmx.config.responseHandling = [
  { code: "204", swap: false },
  { code: "[23]..", swap: true },
  { code: "[45]..", swap: true, error: true }
];

const byId = id => document.getElementById(id);
const closeModal = () => { byId("modal").innerHTML = ""; };
const closePop = () => { byId("pop").innerHTML = ""; };

document.addEventListener("click", e => {
  if (e.target.matches("[data-close]") || e.target.closest("[data-close-btn]")) {
    e.preventDefault();
    closeModal();
  }
  const pop = byId("pop");
  if (pop.innerHTML && !pop.contains(e.target) && !e.target.closest("[data-pop]")) closePop();
});

document.addEventListener("keydown", e => {
  if (e.key === "Escape") { closeModal(); closePop(); }
});

document.body.addEventListener("htmx:beforeRequest", e => {
  if (byId("pop").contains(e.detail.elt)) closePop();
});

document.body.addEventListener("htmx:afterSwap", e => {
  const target = e.detail.target, elt = e.detail.requestConfig && e.detail.requestConfig.elt;
  if (target.id === "pop" && elt) {
    const r = elt.getBoundingClientRect(), pop = byId("pop");
    pop.style.left = Math.max(8, Math.min(r.left, innerWidth - 238)) + "px";
    const h = Math.min(320, pop.scrollHeight), below = r.bottom + 6;
    pop.style.top = (below + h > innerHeight - 8 ? Math.max(8, r.top - h - 6) : below) + "px";
  }
  if (target.id === "exp-list" && elt && elt.hasAttribute("data-scroll")) {
    byId("exp").scrollIntoView({ behavior: matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth", block: "start" });
  }
  if (target.id === "main") { byId("main").scrollTop = 0; }
});

let toastTimer;
new MutationObserver(() => {
  clearTimeout(toastTimer);
  const area = byId("toast-area");
  if (!area.innerHTML.trim()) return;
  toastTimer = setTimeout(() => { area.innerHTML = ""; }, area.querySelector("[hx-post]") ? 6000 : 3500);
}).observe(byId("toast-area"), { childList: true, subtree: true });

function hiSlice(k) {
  const svg = document.querySelector(".donut");
  if (!svg) return;
  const slices = [...svg.querySelectorAll(".sl")];
  const sl = k ? slices.find(x => x.dataset.k === k) : null;
  slices.forEach(x => { x.classList.toggle("dim", !!sl && x !== sl); x.classList.toggle("hi", x === sl); });
  const [a, b, c] = ["dc1", "dc2", "dc3"].map(byId);
  const box = a.parentNode, key = sl ? k : null;
  if (box._k === key) return;
  box._k = key;
  clearTimeout(box._t);
  box.classList.add("fade");
  box._t = setTimeout(() => {
    box.classList.toggle("on", !!sl);
    if (sl) { a.textContent = sl.dataset.n; b.textContent = sl.dataset.v; c.textContent = sl.dataset.p; a.style.color = sl.dataset.c; }
    else { [a, b, c].forEach(x => { x.textContent = x.dataset.d; }); a.style.color = ""; }
    box.classList.remove("fade");
  }, 120);
}

function el(tag, cls, text) {
  const n = document.createElement(tag);
  if (cls) n.className = cls;
  if (text != null) n.textContent = text;
  return n;
}

function bindBars(root) {
  const tip = root.querySelector ? root.querySelector("#ctip") : null;
  if (!tip) return;
  const wrap = tip.parentNode, inner = tip.firstElementChild, segs = [...root.querySelectorAll(".bseg")];
  let hideTimer = null, current = null;
  const hide = () => { current = null; tip.classList.remove("on"); segs.forEach(x => x.classList.remove("dim")); };
  const fill = r => {
    const d = r.dataset, frag = document.createDocumentFragment();
    const h = el("div", "h"), dot = el("span", "dot");
    dot.style.background = d.c; h.append(dot, document.createTextNode(d.k));
    const s = el("div", "s");
    if (d.s) s.append(el("span", null, d.s));
    if (d.d) { const line = el("span"); line.append(el("span", d.dc, d.d), document.createTextNode(" " + d.vs)); s.append(line); }
    frag.append(h, el("div", "m", d.m), el("div", "v", d.v), s);
    inner.replaceChildren(frag);
  };
  const show = r => {
    segs.forEach(x => x.classList.toggle("dim", x !== r));
    const wasOn = tip.classList.contains("on");
    if (current !== r) {
      if (wasOn) { inner.classList.add("fade"); clearTimeout(inner._t); inner._t = setTimeout(() => { fill(r); inner.classList.remove("fade"); }, 110); }
      else fill(r);
      current = r;
    }
    const wb = wrap.getBoundingClientRect(), rb = r.getBoundingClientRect(), half = 96;
    const left = Math.max(half, Math.min(wb.width - half, rb.left - wb.left + rb.width / 2)), top = Math.max(84, rb.top - wb.top);
    if (!wasOn) tip.classList.add("nomove");
    tip.style.left = left + "px";
    tip.style.top = top + "px";
    if (!wasOn) { void tip.offsetWidth; tip.classList.remove("nomove"); }
    tip.classList.add("on");
  };
  segs.forEach(r => {
    r.addEventListener("pointerenter", () => { clearTimeout(hideTimer); show(r); });
    r.addEventListener("pointerleave", () => { hideTimer = setTimeout(hide, 90); });
  });
}

htmx.onLoad(root => {
  root.querySelectorAll(".sl").forEach(s => {
    s.addEventListener("mouseenter", () => hiSlice(s.dataset.k));
    s.addEventListener("focus", () => hiSlice(s.dataset.k));
    s.addEventListener("mouseleave", () => hiSlice(null));
    s.addEventListener("blur", () => hiSlice(null));
  });
  root.querySelectorAll("[data-hk]").forEach(r => {
    if (!r.dataset.hk) return;
    r.addEventListener("mouseenter", () => hiSlice(r.dataset.hk));
    r.addEventListener("mouseleave", () => hiSlice(null));
  });
  bindBars(root);
  if (root.querySelector) {
    const first = root.querySelector(".modal input:not([type=hidden])");
    if (first) first.focus();
  }
});
