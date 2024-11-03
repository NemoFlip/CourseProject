#include <drogon/drogon.h>
#include "../include/login/login.hpp"

int main() {

  /*drogon::app().registerHandler(
    "/login",
    [](const drogon::HttpRequestPtr& request,
      std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
        auto resp = drogon::HttpResponse::newHttpResponse();
        resp->setBody("Hello, World!");
        callback(resp);
    },
    { drogon::Get });
  ;*/
  drogon::app().setLogPath("./")
    .addListener("127.0.0.1", 8048)
    .setThreadNum(16)
    .run();
}