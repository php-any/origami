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

ERR_MARKS = (
    "Internal Server Error",
    "ViewException",
    "不支持调用函数",
    "RootTagMissing",
    "missing root tag",
    "go作用域异常退出",
    "引用参数只能传入变量",
    "Vite manifest",
    "Fatal error",
    "ZY Fatal",
)


def fetch(url, data=None, headers=None, timeout=45):
    h = {"User-Agent": "origami-sweep"}
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


def inspect(path, st, body: bytes):
    text = body.decode("utf-8", "replace")
    hits = [m for m in ERR_MARKS if m in text]
    title = ""
    tm = re.search(r"<title[^>]*>([^<]+)</title>", text, re.I)
    if tm:
        title = re.sub(r"\s+", " ", tm.group(1)).strip()
    flags = {
        "wire": "wire:snapshot" in text,
        "table": "fi-ta-table" in text or "fi-ta-row" in text,
        "form": "fi-fo" in text or "fi-sc" in text,
        "widget": "fi-wi" in text or "fi-stats-overview" in text,
        "warning": "Undefined " in text or "Warning:" in text,
    }
    print(
        f"{path} status={st} len={len(body)} title={title} "
        f"err={hits or 'no'} flags={flags}"
    )
    if hits or st >= 500:
        snippet = re.sub(r"\s+", " ", re.sub(r"<[^>]+>", " ", text))[:280]
        print("  snippet", snippet)
    return text, hits, st >= 500 or bool(hits)


st, hdr, body = fetch(BASE + "/admin/login")
html = body.decode("utf-8", "replace")
snap_m = re.search(r'wire:snapshot="([^"]+)"', html)
csrf_m = re.search(r'data-csrf="([^"]*)"', html)
uri_m = re.search(r'data-update-uri="([^"]+)"', html)
if not snap_m or not uri_m:
    print("FAIL login page missing snapshot")
    sys.exit(1)
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
    "/admin",
    "/admin/settings",
    "/admin/users/create",
    "/admin/roles/create",
    "/admin/products/create",
    "/admin/admins/create",
    "/admin/permissions/create",
]
failed = False
htmls = {}
for path in pages:
    st, hdr, body = fetch(BASE + path)
    dump("sweep-" + path.strip("/").replace("/", "_") + ".html", body)
    text, hits, bad = inspect(path, st, body)
    htmls[path] = text
    if bad and st != 403:
        failed = True

# Follow first edit links from list pages already dumped if present; else fetch lists.
edit_re = re.compile(r'href="(http://127\.0\.0\.1:18086/admin/[^"]+/\d+/edit)"')
view_re = re.compile(r'href="(http://127\.0\.0\.1:18086/admin/[^"]+/\d+)"')
for list_path in ["/admin/users", "/admin/roles", "/admin/products", "/admin/admins", "/admin/activities"]:
    st, hdr, body = fetch(BASE + list_path)
    text = body.decode("utf-8", "replace")
    dump("sweep-" + list_path.strip("/").replace("/", "_") + ".html", body)
    _, hits, bad = inspect(list_path, st, body)
    if bad and st != 403:
        failed = True
    m = edit_re.search(text) or view_re.search(text)
    if m:
        url = m.group(1)
        path = url.replace(BASE, "")
        st, hdr, body = fetch(url)
        dump("sweep-" + path.strip("/").replace("/", "_") + ".html", body)
        _, hits, bad = inspect(path, st, body)
        if bad and st != 403:
            failed = True

edit_pages = [
    "/admin/users/1/edit",
    "/admin/roles/1/edit",
    "/admin/admins/1/edit",
    "/admin/permissions/1/edit",
    "/admin/products/1/edit",
]
for path in edit_pages:
    st, hdr, body = fetch(BASE + path)
    dump("sweep-" + path.strip("/").replace("/", "_") + ".html", body)
    _, hits, bad = inspect(path, st, body)
    if bad and st not in (403, 404):
        failed = True

if failed:
    sys.exit(1)
print("OK sweep")
