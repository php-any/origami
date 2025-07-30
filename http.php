namespace App;

server = new Net\Http\Server(port: 8080);

obj = {}
obj->number = 100

server->get("/", (request, response) => {
    response->write("Hello World");
    // 每次请求后，number会递增
    obj->number += 200;
})

spawn for (i = 0; i < 100; i++) {
    sleep(1);
    echo "异步中...", i, ":", obj->number, "\n";
}

spawn server->start();

$http = new Net\Http\Server(port: 8081);

$http->get("/test", (request, response) => {
    response->write("Hello http");

    obj->number += 200;
})

$http->start();