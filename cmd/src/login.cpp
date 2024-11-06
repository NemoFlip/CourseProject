#include "../include/login/login.hpp"
#include "../include/servisesURIs.hpp"
#include <drogon/HttpClient.h>

drogon::HttpRequestPtr createNewRequest(const drogon::HttpRequestPtr& req, const std::string& path) {
  auto newReq = drogon::HttpRequest::newHttpRequest();

  newReq->setMethod(req->getMethod());
  newReq->setPath(path);
  newReq->setBody(std::string(req->getBody()));
  
  return std::move(newReq);
}

void LoginController::login(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  auto client = drogon::HttpClient::newHttpClient(servisesURI::auth_service);
  auto newReq = createNewRequest(req, "/login");

  client->sendRequest(
    newReq, [callback = std::move(callback)](drogon::ReqResult result,const drogon::HttpResponsePtr& resp) {
      if (result == drogon::ReqResult::Ok) {
        //add some logs later
        callback(resp);
      }
      else {
        //add some logs
        std::cerr << result << '\n';
      }
    });
}

