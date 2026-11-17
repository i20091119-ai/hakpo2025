// 학급 대시보드 서비스 워커: 화면 파일만 저장해 두고, 인터넷이 끊겨도 마지막 화면을 띄운다.
// 날씨·급식·Firebase 같은 외부 요청은 건드리지 않는다.
const CACHE = 'class-dash-v1';
const SHELL = ['./dashboard.html', './manifest.json', './icon-192.png', './icon-512.png'];

self.addEventListener('install', e => {
  e.waitUntil(caches.open(CACHE).then(c => c.addAll(SHELL)).then(() => self.skipWaiting()));
});
self.addEventListener('activate', e => {
  e.waitUntil(caches.keys().then(ks => Promise.all(ks.filter(k => k !== CACHE).map(k => caches.delete(k)))).then(() => self.clients.claim()));
});
self.addEventListener('fetch', e => {
  const url = new URL(e.request.url);
  if (e.request.method !== 'GET' || url.origin !== location.origin) return;
  // 항상 최신 파일을 먼저 받고, 실패하면 저장본 사용
  e.respondWith(
    fetch(e.request).then(res => {
      if (res.ok) { const copy = res.clone(); caches.open(CACHE).then(c => c.put(e.request, copy)); }
      return res;
    }).catch(() => caches.match(e.request, {ignoreSearch: true}))
  );
});
