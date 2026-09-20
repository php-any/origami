# Dump roles/permissions HTML and extract mountAction / wire:init fragments
import html as htmlmod
import json
import os
import re
import sys
import urllib.error
import urllib.request
from http.cookiejar import CookieJar

BASE = "http://127.0.0.1:18086"
DBG = os.path.join(os.path.dirname(__file__), "..", "storage", "origami-debug")
os.makedirs(DBG, exist_ok=True)

cj = CookieJar()
opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))


def fetch(url, data=None, headers=None, timeout=45):
    h = {"User-Agent": "origami-e2e"}
    if headers:
        h.update(headers)
    req = urllib.request.Request(url, data=data, headers=h, method="POST" if data is not None else "GET")
    try:
        resp = opener.open(req, timeout=timeout)
        return resp.status, dict(resp.headers), resp.read()
    except urllib.error.HTTPError as e:
        return e.code, dict(e.headers), e.read()


st, hdr, body = fetch(BASE + "/admin/login")
html = body.decode("utf-8", "replace")
snap_m = re.search(r'wire:snapshot="([^"]*)"', html)
csrf_m = re.search(r'data-csrf="([^"]*)"', html)
uri_m = re.search(r'data-update-uri="([^"]+)"', html)
if not snap_m or not uri_m:
    print("FAIL login snapshot")
    sys.exit(1)
snapshot = htmlmod.unescape(snap_m.group(1))
csrf = htmlmod.unescape(csrf_m.group(1) if csrf_m else "")
update_uri = htmlmod.unescape(uri_m.group(1))
if update_uri.startswith("/"):
    update_uri = BASE + update_uri
payload = json.dumps({
    "components": [{
        "snapshot": snapshot,
        "updates": {"data.email": "admin@example.com", "data.password": "password"},
        "calls": [{"method": "authenticate", "params": []}],
    }]
}).encode("utf-8")
st, hdr, body = fetch(update_uri, data=payload, headers={
    "Content-Type": "application/json", "Accept": "application/json",
    "X-Livewire": "true", "X-CSRF-TOKEN": csrf,
})
print("login post", st)

for path in ["/admin/roles", "/admin/permissions", "/admin/users"]:
    st, hdr, body = fetch(BASE + path)
    name = "mount-" + path.strip("/").replace("/", "_") + ".html"
    out = os.path.join(DBG, name)
    with open(out, "wb") as f:
        f.write(body)
    text = body.decode("utf-8", "replace")
    print(f"== {path} status={st} len={len(body)}")
    # raw wire:init
    inits = re.findall(r'wire:init="([^"]*)"', text)
    print("  wire:init count", len(inits))
    for i, v in enumerate(inits[:8]):
        print("  init", i, v[:300])
    # unescaped mountAction occurrences
    clicks = re.findall(r'(?:wire:click|wire:init|x-on:click)="([^"]*mountAction[^"]*)"', text)
    print("  attr mountAction count", len(clicks))
    for i, v in enumerate(clicks[:12]):
        print("  attr", i, v[:400])
    # also look at truncated: mountAction( without closing in same attr
    broken = [v for v in clicks if v.count("(") > v.count(")")]
    print("  broken paren attrs", len(broken))
    for v in broken[:5]:
        print("  BROKEN", v[:400])
    # JSON.parse with raw quotes
    if "JSON.parse('" in text:
        idx = text.find("JSON.parse('")
        print("  JSON.parse ctx", text[max(0,idx-40):idx+120].replace("\n"," "))
