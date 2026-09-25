(() => {
  const API = location.origin;
  const WS_BASE = (location.protocol === 'https:' ? 'wss:' : 'ws:') + '//' + location.host;

  const state = {
    token: localStorage.getItem('pet_play_token') || '',
    user: null,
    pet: null,
    visitHostId: 0,
    ws: null,
  };

  const $ = (id) => document.getElementById(id);
  const show = (el, on) => el.classList.toggle('hidden', !on);

  async function api(method, path, body) {
    const headers = { 'Content-Type': 'application/json' };
    if (state.token) headers.Authorization = 'Bearer ' + state.token;
    const res = await fetch(API + path, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    });
    const json = await res.json();
    if (json.code !== 0) throw new Error(json.msg || '请求失败');
    return json.data;
  }

  function appendChat(box, msg) {
    const div = document.createElement('div');
    div.className = 'line';
    div.innerHTML = `<span class="name">${escapeHtml(msg.fromName || msg.fromId)}</span> ${escapeHtml(msg.content)}`;
    box.appendChild(div);
    box.scrollTop = box.scrollHeight;
  }

  function escapeHtml(s) {
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;');
  }

  function renderPet() {
    const pet = state.pet;
    if (!pet) return;
    $('petName').textContent = pet.name;
    $('petLevel').textContent = 'Lv.' + pet.level;
    $('txtHunger').textContent = pet.hunger;
    $('txtMood').textContent = pet.mood;
    $('barHunger').style.width = pet.hunger + '%';
    $('barMood').style.width = pet.mood + '%';
    let outfit = { hat: 'hat_basic', clothes: 'clothes_basic' };
    try { outfit = JSON.parse(pet.outfitJson || '{}'); } catch (_) {}
    const vis = $('petVisual');
    vis.dataset.hat = outfit.hat || 'hat_basic';
    vis.dataset.clothes = outfit.clothes || 'clothes_basic';
  }

  function renderUserChip() {
    const u = state.user;
    if (!u) return;
    $('userChip').textContent = `${u.nickname} · ID ${u.id}`;
    show($('userChip'), true);
  }

  function connectWs() {
    if (state.ws) {
      try { state.ws.close(); } catch (_) {}
    }
    const ws = new WebSocket(`${WS_BASE}/ws?token=${encodeURIComponent(state.token)}`);
    state.ws = ws;
    ws.onmessage = (ev) => {
      let msg;
      try { msg = JSON.parse(ev.data); } catch { return; }
      if (msg.type === 'chat.msg') {
        if (msg.channel === 'world') appendChat($('chatLog'), msg);
        if (msg.channel === 'friend') {
          const peer = Number($('dmPeerInput').value || 0);
          if (!peer || msg.fromId === peer || msg.toId === peer || msg.fromId === state.user.id) {
            appendChat($('dmLog'), msg);
          }
        }
      }
      if (msg.type === 'chat.rejected') {
        $('homeMsg').textContent = '消息被内容安全拦截: ' + (msg.reason || '');
      }
    };
    setInterval(() => {
      if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({ type: 'ping' }));
    }, 20000);
  }

  async function enterHome() {
    show($('viewLogin'), false);
    show($('viewHome'), true);
    renderUserChip();
    renderPet();
    connectWs();
    await refreshFriends();
    await loadWorldHistory();
  }

  async function ensurePet() {
    if (state.pet) return;
    const name = (state.user.nickname || '旅人') + '的宠物';
    state.pet = await api('POST', '/api/v1/pet', { name });
  }

  async function refreshFriends() {
    const list = await api('GET', '/api/v1/friends');
    const ul = $('friendList');
    ul.innerHTML = '';
    if (!list || !list.length) {
      ul.innerHTML = '<li>暂无好友，找另一个窗口互加吧</li>';
      return;
    }
    list.forEach((f) => {
      const li = document.createElement('li');
      li.textContent = `#${f.userId} ${f.nickname} [${f.status}]${f.online ? ' · 在线' : ''}`;
      li.style.cursor = 'pointer';
      li.onclick = () => {
        $('visitIdInput').value = f.userId;
        $('dmPeerInput').value = f.userId;
        $('friendIdInput').value = f.userId;
      };
      ul.appendChild(li);
    });
  }

  async function loadWorldHistory() {
    const list = await api('GET', '/api/v1/chat/history?channel=world&limit=40');
    $('chatLog').innerHTML = '';
    (list || []).forEach((m) => appendChat($('chatLog'), m));
  }

  $('btnLogin').onclick = async () => {
    $('loginErr').textContent = '';
    try {
      const code = $('codeInput').value.trim() || 'player';
      const nickname = $('nickInput').value.trim() || code;
      const data = await api('POST', '/api/v1/auth/login', {
        platform: 'dev',
        code,
        nickname,
      });
      state.token = data.token;
      state.user = data.user;
      state.pet = data.pet || null;
      localStorage.setItem('pet_play_token', state.token);
      await ensurePet();
      await enterHome();
    } catch (e) {
      $('loginErr').textContent = e.message;
    }
  };

  async function withHome(fn) {
    try {
      $('homeMsg').textContent = '';
      await fn();
      renderPet();
    } catch (e) {
      $('homeMsg').textContent = e.message;
    }
  }

  $('btnFeed').onclick = () => withHome(async () => {
    state.pet = await api('POST', '/api/v1/pet/feed', { foodCode: 'food_kibble' });
    $('homeMsg').textContent = '喂食成功';
  });
  $('btnTreat').onclick = () => withHome(async () => {
    state.pet = await api('POST', '/api/v1/pet/feed', { foodCode: 'food_treat' });
    $('homeMsg').textContent = '零食好吃！';
  });
  $('btnHat').onclick = () => withHome(async () => {
    state.pet = await api('POST', '/api/v1/pet/outfit', { hat: 'hat_star' });
    $('homeMsg').textContent = '换上星星帽';
  });
  $('btnCape').onclick = () => withHome(async () => {
    state.pet = await api('POST', '/api/v1/pet/outfit', { clothes: 'clothes_cape' });
    $('homeMsg').textContent = '披上小披风';
  });

  $('btnSendWorld').onclick = () => {
    const content = $('chatInput').value.trim();
    if (!content || !state.ws || state.ws.readyState !== WebSocket.OPEN) return;
    state.ws.send(JSON.stringify({ type: 'chat.send', channel: 'world', content }));
    $('chatInput').value = '';
  };

  $('btnReqFriend').onclick = async () => {
    try {
      await api('POST', '/api/v1/friends/request', { friendId: Number($('friendIdInput').value) });
      $('homeMsg').textContent = '已发送好友申请';
      await refreshFriends();
    } catch (e) { $('homeMsg').textContent = e.message; }
  };
  $('btnAcceptFriend').onclick = async () => {
    try {
      await api('POST', '/api/v1/friends/accept', { friendId: Number($('friendIdInput').value) });
      $('homeMsg').textContent = '已接受好友';
      await refreshFriends();
    } catch (e) { $('homeMsg').textContent = e.message; }
  };

  $('btnVisit').onclick = async () => {
    try {
      const id = Number($('visitIdInput').value);
      const home = await api('GET', `/api/v1/visit/${id}`);
      state.visitHostId = id;
      const pet = home.pet;
      $('visitInfo').innerHTML = `<b>${escapeHtml(home.host.nickname)}</b> 的家<br/>` +
        (pet ? `${escapeHtml(pet.name)} Lv.${pet.level} · 饱食 ${pet.hunger} · 心情 ${pet.mood}` : '对方还没宠物') +
        `<br/>${home.online ? '在线' : '离线'}`;
      $('btnVisitPet').disabled = !pet;
      $('btnVisitGift').disabled = !pet;
    } catch (e) {
      $('visitInfo').textContent = e.message;
    }
  };
  $('btnVisitPet').onclick = async () => {
    if (!state.visitHostId) return;
    const pet = await api('POST', `/api/v1/visit/${state.visitHostId}/interact`, { action: 'pet' });
    $('visitInfo').textContent = `抚摸成功，对方心情 ${pet.mood}`;
  };
  $('btnVisitGift').onclick = async () => {
    if (!state.visitHostId) return;
    const pet = await api('POST', `/api/v1/visit/${state.visitHostId}/interact`, { action: 'gift' });
    $('visitInfo').textContent = `送礼成功，饱食 ${pet.hunger} / 心情 ${pet.mood}`;
  };

  $('btnLoadDm').onclick = async () => {
    try {
      const peer = Number($('dmPeerInput').value);
      const list = await api('GET', `/api/v1/chat/history?channel=friend&peerId=${peer}&limit=40`);
      $('dmLog').innerHTML = '';
      (list || []).forEach((m) => appendChat($('dmLog'), m));
    } catch (e) { $('homeMsg').textContent = e.message; }
  };
  $('btnSendDm').onclick = () => {
    const content = $('dmInput').value.trim();
    const toId = Number($('dmPeerInput').value);
    if (!content || !toId || !state.ws || state.ws.readyState !== WebSocket.OPEN) return;
    state.ws.send(JSON.stringify({ type: 'chat.send', channel: 'friend', toId, content }));
    $('dmInput').value = '';
  };

  // Auto resume session
  (async () => {
    if (!state.token) return;
    try {
      const me = await api('GET', '/api/v1/me');
      state.user = me.user;
      state.pet = me.pet || null;
      await ensurePet();
      await enterHome();
    } catch (_) {
      localStorage.removeItem('pet_play_token');
      state.token = '';
    }
  })();
})();
