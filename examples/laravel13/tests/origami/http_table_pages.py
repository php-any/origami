import html as htmlmod
import json
import os
import re
import sys
import urllib.error
import urllib.request
from http.cookiejar import CookieJar

BASE = "http://127.0.0.1:18086"
DBG = os.path.join(os.path.dirname(__file__), "..", "..", "storage", "origami-debug")
os.makedirs(DBG, exist_ok=True)

cj = CookieJar()
opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))


def fetch(url, data=None, headers=None, timeout=40):
    h = {"User-Agent": "origami-e2e"}
    if headers:
        h.update(headers)
    req = urllib.request.Request(url, data=data, headers=h, method="POST" if data is not None else "GET")
    try:
        resp = opener.open(req, timeout=timeout)
        return resp.status, dict(resp.headers), resp.read()
    except urllib.error.HTTPError as e:
        return e.code, dict(e.headers), e.read()


def dump(name, body: bytes):
    path = os.path.join(DBG, name)
    with open(path, "wb") as f:
        f.write(body)
    return path


st, hdr, body = fetch(BASE + "/admin/login")
html = body.decode("utf-8", "replace")
snap_m = re.search(r'wire:snapshot="([^"]+)"', html)
csrf_m = re.search(r'data-csrf="([^"]*)"', html)
uri_m = re.search(r'data-update-uri="([^"]+)"', html)
snapshot = htmlmod.unescape(snap_m.group(1))
csrf = htmlmod.unescape(csrf_m.group(1) if csrf_m else "")
update_uri = htmlmod.unescape(uri_m.group(1))
if update_uri.startswith("/"):
    update_uri = BASE + update_uri

payload = json.dumps(
    {
        "components": [
            {
                "snapshot": snapshot,
                "updates": {
                    "data.email": "admin@example.com",
                    "data.password": "password",
                },
                "calls": [{"method": "authenticate", "params": []}],
            }
        ]
    }
).encode("utf-8")
st, hdr, body = fetch(
    update_uri,
    data=payload,
    headers={
        "Content-Type": "application/json",
        "Accept": "application/json",
        "X-Livewire": "true",
        "X-CSRF-TOKEN": csrf,
    },
)
print("login_post", st)

pages = [
    "/admin/roles",
    "/admin/permissions",
    "/admin/users",
    "/admin/admins",
    "/admin/products",
    "/admin/categories",
    "/admin/activities",
]
failed = False
for path in pages:
    st, hdr, body = fetch(BASE + path)
    text = body.decode("utf-8", "replace")
    dump("table-" + path.strip("/").replace("/", "_") + ".html", body)
    err = (
        "不支持调用函数" in text
        or "ViewException" in text
        or "Internal Server Error" in text
        or "当前值" in text
    )
    has_table = "fi-ta-table" in text or "fi-ta-row" in text
    title = ""
    tm = re.search(r"<title[^>]*>([^<]+)</title>", text, re.I)
    if tm:
        title = re.sub(r"\s+", " ", tm.group(1)).strip()
    print(f"{path} status={st} len={len(body)} err={err} table={has_table} title={title}")
    if st >= 500 or err:
        failed = True
        snippet = re.sub(r"\s+", " ", re.sub(r"<[^>]+>", " ", text))[:240]
        print("  snippet", snippet)

if failed:
    sys.exit(1)
print("OK table pages")
