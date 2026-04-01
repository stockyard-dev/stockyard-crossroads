package server

var dashboardHTML = []byte(`<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Stockyard Crossroads</title><style>:root{--bg:#1a1410;--surface:#241c15;--border:#3d2e1e;--rust:#c4622d;--cream:#f5e6c8;--muted:#7a6550;--text:#e8d5b0}*{box-sizing:border-box;margin:0;padding:0}body{background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace,sans-serif}header{background:var(--surface);border-bottom:1px solid var(--border);padding:1rem 2rem;display:flex;align-items:center;gap:1rem}.logo{color:var(--rust);font-size:1.25rem;font-weight:700}.badge{background:var(--rust);color:var(--cream);font-size:0.65rem;padding:0.2rem 0.5rem;border-radius:3px;font-weight:600;text-transform:uppercase}main{max-width:1100px;margin:0 auto;padding:2rem}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;margin-bottom:2rem}.stat{background:var(--surface);border:1px solid var(--border);border-radius:6px;padding:1.25rem;text-align:center}.stat-value{font-size:1.75rem;font-weight:700;color:var(--rust)}.stat-label{font-size:0.75rem;color:var(--muted);margin-top:0.25rem;text-transform:uppercase}.layout{display:grid;grid-template-columns:260px 1fr;gap:1rem}.card{background:var(--surface);border:1px solid var(--border);border-radius:6px;padding:1.5rem;margin-bottom:1rem}.card h2{font-size:0.85rem;color:var(--muted);text-transform:uppercase;letter-spacing:0.08em;margin-bottom:1rem}.form-row{display:flex;gap:0.5rem;margin-bottom:0.75rem;flex-wrap:wrap}select,input{background:var(--bg);border:1px solid var(--border);color:var(--text);padding:0.5rem 0.75rem;border-radius:4px;font-family:inherit;font-size:0.85rem;flex:1}.btn{background:var(--rust);color:var(--cream);border:none;padding:0.5rem 1rem;border-radius:4px;cursor:pointer;font-family:inherit;font-size:0.85rem;font-weight:600}.btn:hover{opacity:0.85}.btn-sm{padding:0.25rem 0.6rem;font-size:0.75rem}.btn-danger{background:#7a2020}.prof-item{padding:0.5rem 0.75rem;cursor:pointer;border-radius:4px;margin-bottom:0.25rem}.prof-item:hover,.prof-item.active{background:rgba(196,98,45,0.15)}.prof-item.active{border-left:3px solid var(--rust)}.link-row{display:flex;align-items:center;gap:0.5rem;padding:0.5rem 0;border-bottom:1px solid var(--border)}.link-icon{width:24px;text-align:center}.link-info{flex:1;min-width:0}.link-title{font-weight:600;color:var(--cream);font-size:0.85rem}.link-url{font-size:0.72rem;color:var(--muted);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.link-clicks{font-size:0.72rem;color:var(--muted)}.empty{color:var(--muted);font-size:0.85rem;padding:1rem 0;text-align:center}.dim{opacity:0.4}</style></head>
<body>
<header><span class="logo">&#x2B21; Stockyard</span><span style="color:var(--muted)">/</span><span style="color:var(--cream);font-weight:600">Crossroads</span><span class="badge">Link-in-Bio</span></header>
<main>
<div class="stats"><div class="stat"><div class="stat-value" id="s1">0</div><div class="stat-label">Profiles</div></div><div class="stat"><div class="stat-value" id="s2">0</div><div class="stat-label">Total Clicks</div></div><div class="stat"><div class="stat-value" id="s3">FREE</div><div class="stat-label">Tier</div></div></div>
<div class="layout">
<div>
<div class="card"><h2>New Profile</h2>
<div class="form-row"><input id="f-pslug" placeholder="Slug (e.g. jane)"><input id="f-pname" placeholder="Display name"></div>
<button class="btn btn-sm" onclick="addProfile()">Create</button>
<div id="prof-list" style="margin-top:1rem"><div class="empty">No profiles</div></div></div>
</div>
<div>
<div class="card" id="link-card">
<h2>Links: <span id="cur-prof" style="color:var(--cream)">—</span> <a id="pub-link" href="#" target="_blank" style="font-size:0.75rem;color:var(--rust);text-decoration:none;display:none">&#x1F517; Public page</a></h2>
<div id="link-form" style="display:none">
<div class="form-row"><input id="f-ltitle" placeholder="Title"><input id="f-lurl" placeholder="URL"><input id="f-licon" placeholder="Icon (emoji)" style="max-width:80px"></div>
<button class="btn btn-sm" onclick="addLink()">Add Link</button>
</div>
<div id="link-list"><div class="empty">Select a profile</div></div>
</div>
</div>
</div>
</main>
<script>
var curProf=null;
function load(){fetch('/api/stats').then(function(r){return r.json()}).then(function(d){document.getElementById('s1').textContent=d.profiles||0;document.getElementById('s2').textContent=d.total_clicks||0})}
function loadProfiles(){fetch('/api/profiles').then(function(r){return r.json()}).then(function(list){var el=document.getElementById('prof-list');el.innerHTML=list.length?list.map(function(p){return'<div class="prof-item'+(curProf===p.id?' active':'')+'" onclick="selectProfile('+p.id+',\''+p.slug+'\',\''+p.name+'\')"><div style="display:flex;justify-content:space-between">'+p.name+'<button class="btn btn-sm btn-danger" onclick="event.stopPropagation();delProfile('+p.id+')">x</button></div><div style="font-size:0.72rem;color:var(--muted)">@'+p.slug+'</div></div>'}).join(''):'<div class="empty">No profiles</div>'})}
function selectProfile(id,slug,name){curProf=id;document.getElementById('cur-prof').textContent=name;var pl=document.getElementById('pub-link');pl.href='/p/'+slug;pl.style.display='inline';document.getElementById('link-form').style.display='block';loadLinks(id);loadProfiles()}
function loadLinks(id){fetch('/api/profiles/'+id+'/links').then(function(r){return r.json()}).then(function(list){var el=document.getElementById('link-list');el.innerHTML=list.length?list.map(function(l){return'<div class="link-row'+(l.active?'':' dim')+'"><span class="link-icon">'+(l.icon||'&#x1F517;')+'</span><div class="link-info"><div class="link-title">'+l.title+'</div><div class="link-url">'+l.url+'</div></div><span class="link-clicks">&#x1F5B1; '+l.clicks+'</span><button class="btn btn-sm" onclick="toggleLink('+l.id+')">'+(l.active?'Hide':'Show')+'</button><button class="btn btn-sm btn-danger" onclick="delLink('+l.id+')">x</button></div>'}).join(''):'<div class="empty">No links yet</div>'})}
function addProfile(){var d={slug:document.getElementById('f-pslug').value.trim(),name:document.getElementById('f-pname').value.trim()};if(!d.slug||!d.name)return;fetch('/api/profiles',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(d)}).then(function(){loadProfiles();load()})}
function delProfile(id){fetch('/api/profiles/'+id,{method:'DELETE'}).then(function(){if(curProf===id)curProf=null;loadProfiles();load()})}
function addLink(){if(!curProf)return;var d={title:document.getElementById('f-ltitle').value.trim(),url:document.getElementById('f-lurl').value.trim(),icon:document.getElementById('f-licon').value.trim()};if(!d.title||!d.url)return;fetch('/api/profiles/'+curProf+'/links',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(d)}).then(function(){document.getElementById('f-ltitle').value='';document.getElementById('f-lurl').value='';loadLinks(curProf);load()})}
function toggleLink(id){fetch('/api/links/'+id+'/toggle',{method:'POST'}).then(function(){if(curProf)loadLinks(curProf)})}
function delLink(id){fetch('/api/links/'+id,{method:'DELETE'}).then(function(){if(curProf)loadLinks(curProf);load()})}
load();loadProfiles();
</script></body></html>`)
