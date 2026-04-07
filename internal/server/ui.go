package server

var dashboardHTML = []byte(`<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Stockyard Flume</title><style>:root{--bg:#1a1410;--surface:#241c15;--border:#3d2e1e;--rust:#c4622d;--cream:#f5e6c8;--muted:#7a6550;--text:#e8d5b0}*{box-sizing:border-box;margin:0;padding:0}body{background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace,sans-serif}header{background:var(--surface);border-bottom:1px solid var(--border);padding:1rem 2rem;display:flex;align-items:center;gap:1rem}.logo{color:var(--rust);font-size:1.25rem;font-weight:700}.badge{background:var(--rust);color:var(--cream);font-size:0.65rem;padding:0.2rem 0.5rem;border-radius:3px;font-weight:600;text-transform:uppercase}main{max-width:1200px;margin:0 auto;padding:2rem}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;margin-bottom:2rem}.stat{background:var(--surface);border:1px solid var(--border);border-radius:6px;padding:1.25rem;text-align:center}.stat-value{font-size:1.75rem;font-weight:700;color:var(--rust)}.stat-label{font-size:0.75rem;color:var(--muted);margin-top:0.25rem;text-transform:uppercase;letter-spacing:0.05em}.grid{display:grid;grid-template-columns:1fr 1fr;gap:1rem;margin-bottom:2rem}.card{background:var(--surface);border:1px solid var(--border);border-radius:6px;padding:1.5rem}.card h2{font-size:0.85rem;color:var(--muted);text-transform:uppercase;letter-spacing:0.08em;margin-bottom:1rem}.full{grid-column:1/-1}.form-row{display:flex;gap:0.5rem;margin-bottom:0.75rem;flex-wrap:wrap}select,input{background:var(--bg);border:1px solid var(--border);color:var(--text);padding:0.5rem 0.75rem;border-radius:4px;font-family:inherit;font-size:0.85rem;flex:1}.btn{background:var(--rust);color:var(--cream);border:none;padding:0.5rem 1rem;border-radius:4px;cursor:pointer;font-family:inherit;font-size:0.85rem;font-weight:600}.btn:hover{opacity:0.85}.btn-sm{padding:0.25rem 0.6rem;font-size:0.75rem}.btn-danger{background:#7a2020}.log-table{width:100%;border-collapse:collapse;font-size:0.78rem;font-family:'JetBrains Mono',monospace}.log-table th{text-align:left;color:var(--muted);padding:0.4rem;border-bottom:1px solid var(--border);font-size:0.7rem;text-transform:uppercase}.log-table td{padding:0.35rem 0.4rem;border-bottom:1px solid rgba(61,46,30,0.5);vertical-align:top}.log-debug{color:#7a6550}.log-info{color:#e8d5b0}.log-warn{color:#f0ad4e}.log-error{color:#d9534f}.log-fatal{color:#ff4444;font-weight:700}.ts{color:var(--muted);font-size:0.7rem;white-space:nowrap}.msg{word-break:break-all}.empty{color:var(--muted);font-size:0.85rem;padding:1rem 0;text-align:center}.token-box{background:var(--bg);border:1px solid var(--border);padding:0.4rem 0.75rem;border-radius:4px;font-size:0.75rem;color:var(--muted);word-break:break-all}</style></head>
<body>
<header><span class="logo">&#x2B21; Stockyard</span><span style="color:var(--muted)">/</span><span style="color:var(--cream);font-weight:600">Flume</span><span class="badge">Logs</span></header>
<main>
<div class="stats"><div class="stat"><div class="stat-value" id="s1">0</div><div class="stat-label">Total Logs</div></div><div class="stat"><div class="stat-value" id="s2">0</div><div class="stat-label">Streams</div></div><div class="stat"><div class="stat-value" id="s3">FREE</div><div class="stat-label">Tier</div></div></div>
<div class="grid">
<div class="card"><h2>New Stream</h2>
<div class="form-row"><input id="f-sname" placeholder="Stream name"><input id="f-ret" type="number" placeholder="Retention days" value="7" style="max-width:130px"></div>
<button class="btn btn-sm" onclick="addStream()">Create</button>
<div id="stream-list" style="margin-top:1rem"><div class="empty">No streams</div></div></div>
<div class="card"><h2>Search Logs</h2>
<div class="form-row"><select id="f-sid"><option value="">All Streams</option></select><select id="f-lvl"><option value="">All Levels</option><option>debug</option><option>info</option><option>warn</option><option>error</option><option>fatal</option></select></div>
<div class="form-row"><input id="f-search" placeholder="Search message..." oninput="loadLogs()"><input id="f-limit" type="number" placeholder="Limit" value="100" style="max-width:80px"></div>
<button class="btn btn-sm" onclick="loadLogs()">Refresh</button></div>
</div>
<div class="card full"><h2>Log Viewer <span style="color:var(--muted);font-size:0.75rem">(newest first)</span></h2><div id="log-view"><div class="empty">No logs yet. Ingest via POST /api/ingest/{stream_id}</div></div></div>
</main>
<script>
function load(){fetch('/api/stats').then(function(r){return r.json()}).then(function(d){document.getElementById('s1').textContent=d.total_logs||0})}
function loadStreams(){fetch('/api/streams').then(function(r){return r.json()}).then(function(list){document.getElementById('s2').textContent=list.length;var sel=document.getElementById('f-sid');sel.innerHTML='<option value="">All Streams</option>';list.forEach(function(s){sel.innerHTML+='<option value="'+s.id+'">'+s.name+'</option>'});var el=document.getElementById('stream-list');el.innerHTML=list.length?list.map(function(s){return'<div style="padding:0.5rem 0;border-bottom:1px solid var(--border)"><div style="display:flex;justify-content:space-between"><span style="color:var(--cream)">'+s.name+'</span><button class="btn btn-sm btn-danger" onclick="delStream('+s.id+')">x</button></div><div class="token-box">POST /api/ingest/'+s.id+' &nbsp;&bull;&nbsp; token: '+s.token+'</div></div>'}).join(''):'<div class="empty">No streams</div>'})}
function loadLogs(){var sid=document.getElementById('f-sid').value;var lvl=document.getElementById('f-lvl').value;var q=document.getElementById('f-search').value;var lim=document.getElementById('f-limit').value||100;var u='/api/logs?limit='+lim+(sid?'&stream_id='+sid:'')+(lvl?'&level='+lvl:'')+(q?'&q='+encodeURIComponent(q):'')
fetch(u).then(function(r){return r.json()}).then(function(list){var el=document.getElementById('log-view');el.innerHTML=list.length?'<table class="log-table"><thead><tr><th>Time</th><th>Stream</th><th>Level</th><th>Message</th></tr></thead><tbody>'+list.map(function(e){var cls='log-'+(e.level||'info');return'<tr><td class="ts">'+e.created_at+'</td><td style="color:var(--muted)">'+e.stream_name+'</td><td class="'+cls+'">'+e.level+'</td><td class="msg '+cls+'">'+e.message+(e.fields&&e.fields!=='{}'?'<span style="color:var(--muted)"> '+e.fields+'</span>':'')+'</td></tr>'}).join('')+"</tbody></table>":'<div class="empty">No logs match your filters</div>';load()})}
function addStream(){var n=document.getElementById('f-sname').value.trim();var r=parseInt(document.getElementById('f-ret').value)||7;if(!n)return;fetch('/api/streams',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:n,retention_days:r})}).then(function(){document.getElementById('f-sname').value='';loadStreams()})}
function delStream(id){fetch('/api/streams/'+id,{method:'DELETE'}).then(function(){loadStreams();loadLogs();load()})}
load();loadStreams();loadLogs();setInterval(function(){loadLogs();load()},10000);
fetch('/api/tier').then(r=>r.json()).then(j=>{if(j.tier==='free'){var b=document.getElementById('upgrade-banner');if(b)b.style.display='block'}}).catch(()=>{var b=document.getElementById('upgrade-banner');if(b)b.style.display='block'});
</script><script>
(function(){
  fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
    if(!cfg||typeof cfg!=='object')return;
    if(cfg.dashboard_title){
      document.title=cfg.dashboard_title;
      var h1=document.querySelector('h1');
      if(h1){
        var inner=h1.innerHTML;
        var firstSpan=inner.match(/<span[^>]*>[^<]*<\/span>/);
        if(firstSpan){h1.innerHTML=firstSpan[0]+' '+cfg.dashboard_title}
        else{h1.textContent=cfg.dashboard_title}
      }
    }
  }).catch(function(){});
})();
</script>
</body></html>`)
