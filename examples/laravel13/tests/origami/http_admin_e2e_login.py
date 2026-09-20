# Real HTTP e2e against Origami serve: login then hit /admin pages.
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


def fetch(url, data=None, headers=None, timeout=40):
    h = {"User-Agent": "origami-e2e"}
    if headers:
        h.update(headers)
    req = urllib.request.Request(url, data=data, headers=h, method="POST" if data is not None else "GET")
    try:
        resp = opener.open(req, timeout=timeout)
        body = resp.read()
        return resp.status, dict(resp.headers), body
    except urllib.error.HTTPError as e:
        return e.code, dict(e.headers), e.read()


def dump(name, body: bytes):
    path = os.path.join(DBG, name)
    with open(path, "wb") as f:
        f.write(body)
    return path


def markers(html: str):
    return {
        "missing_root": "missing root tag" in html or "RootTagMissing" in html,
        "panic": "go作用域异常退出" in html or "nil pointer" in html,
        "引用参数": "引用参数只能传入变量" in html,
        "exception": "Internal Server Error" in html or "Server Error" in html,
        "login_page": "Sign in" in html or "登录" in html and "password" in html.lower(),
        "dashboard": "仪表盘" in html or "Origami Admin" in html,
        "wire": "wire:snapshot" in html,
    }


print("== GET /admin/login ==")
st, hdr, body = fetch(BASE + "/admin/login")
html = body.decode("utf-8", "replace")
dump("e2e-login.html", body)
print(f"status={st} len={len(body)}")
print("markers", markers(html))
if st != 200:
    print("FAIL login page")
    sys.exit(1)

snap_m = re.search(r'wire:snapshot="([^"]+)"', html)
csrf_m = re.search(r'data-csrf="([^"]*)"', html)
uri_m = re.search(r'data-update-uri="([^"]+)"', html)
if not snap_m or not uri_m:
    print("FAIL missing snapshot/update-uri")
    sys.exit(1)

snapshot = (
    snap_m.group(1)
    .replace("&quot;", '"')
    .replace("&amp;", "&")
    .replace("&#34;", '"')
)
# html.unescape
import html as htmlmod

snapshot = htmlmod.unescape(snap_m.group(1))
csrf = htmlmod.unescape(csrf_m.group(1) if csrf_m else "")
update_uri = htmlmod.unescape(uri_m.group(1))
if update_uri.startswith("/"):
    update_uri = BASE + update_uri
print("update_uri", update_uri)
print("csrf_len", len(csrf), "snap_len", len(snapshot))

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

print("== POST livewire authenticate ==")
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
dump("e2e-login-post.json", body)
text = body.decode("utf-8", "replace")
print(f"status={st} len={len(body)}")
print("head", text[:400].replace("\n", " "))
try:
    js = json.loads(text)
    redir = js["components"][0].get("effects", {}).get("redirect")
except Exception as e:
    redir = None
    print("json_err", e)
print("redirect", redir)
if st >= 500 or "missing root tag" in text or "go作用域异常退出" in text:
    print("FAIL login post")
    sys.exit(1)

pages = ["/admin", "/admin/users", "/admin/orders", "/admin/products", "/admin/settings"]
failed = False
for path in pages:
    print(f"== GET {path} ==")
    st, hdr, body = fetch(BASE + path)
    html = body.decode("utf-8", "replace")
    loc = hdr.get("Location") or hdr.get("location") or ""
    name = "e2e-" + path.strip("/").replace("/", "_") + ".html"
    dump(name, body)
    mk = markers(html)
    print(f"status={st} loc={loc} len={len(body)} markers={mk}")
    title = ""
    tm = re.search(r"<title[^>]*>([^<]+)</title>", html, re.I)
    if tm:
        title = tm.group(1).strip()
        print("title", title)
    bad = st >= 500 or mk["missing_root"] or mk["panic"] or mk["引用参数"]
    if st in (301, 302) and "login" in loc:
        print("FAIL bounced to login")
        bad = True
    if path == "/admin" and st == 200 and not mk["dashboard"] and not mk["wire"]:
        print("FAIL dashboard missing markers")
        bad = True
    if bad:
        failed = True
        print("snippet", re.sub(r"\s+", " ", re.sub(r"<[^>]+>", " ", html))[:280])

if failed:
    sys.exit(1)
print("OK e2e login + admin pages")
