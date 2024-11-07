#include "../include/login/login.hpp"
#include "../include/servisesURIs.hpp"
#include <drogon/HttpClient.h>

void resendRequest(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback, 
  const std::string& path) 
{
  auto client = drogon::HttpClient::newHttpClient(servisesURI::auth_service);
  auto newReq = drogon::HttpRequest::newHttpRequest();

  newReq->setMethod(req->getMethod());
  newReq->setPath(path);
  newReq->setBody(std::string(req->getBody()));

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

void LoginController::login(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  resendRequest(req, std::move(callback), "/login");
}

void LoginController::logout(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  resendRequest(req, std::move(callback), "/logout");
}

void LoginController::registration(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  resendRequest(req, std::move(callback), "/registration");
}

void LoginController::passwordRecovery(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  resendRequest(req, std::move(callback), "/passwordRecovery");
}