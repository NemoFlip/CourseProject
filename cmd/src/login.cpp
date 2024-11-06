#include "../include/login/login.hpp"
#include "../include/servisesURIs.hpp"
#include <drogon/HttpClient.h>

void login(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {  
  std::string loginURL = servisesURI::auth_service;

  auto client = drogon::HttpClient::newHttpClient(loginURL);

  auto nextReq = drogon::HttpRequest::newHttpRequest();
  nextReq->setMethod(req->getMethod());
  nextReq->setPath("/login");
}