#include "../include/login/LoginController.hpp"
#include "../include/servisesURIs.hpp"
#include <drogon/HttpClient.h>
#include <chrono>

void resendRequest(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback, const std::string& path) 
{
  auto client = drogon::HttpClient::newHttpClient(servisesURI::auth_service);
  auto newReq = drogon::HttpRequest::newHttpRequest();

  newReq->setMethod(req->getMethod());
  newReq->setPath(path);
  newReq->setBody(std::string(req->getBody()));

  LOG->trace("Sending request to the {} adress to the {} endpoint.", servisesURI::auth_service, path);

  auto start = std::chrono::high_resolution_clock::now();
  client->sendRequest(
    newReq, [callback = std::move(callback), &start](drogon::ReqResult result,const drogon::HttpResponsePtr& resp) {

      auto end = std::chrono::high_resolution_clock::now();
      auto duration = std::chrono::duration<float>(end - start);
      
      if (result == drogon::ReqResult::Ok) {
        LOG->trace("Received a response with {} status codó.", resp->getStatusCode());
        LOG->info("Ñonnection succeeded.");
      }
      else {
        switch (result) {
        case drogon::ReqResult::BadResponse:
          LOG->error("Receives a badresponse from server.");
          break;
        case drogon::ReqResult::NetworkFailure:
          LOG->error("Trouble with the Network.");
          break;
        case drogon::ReqResult::BadServerAddress:
          LOG->error("Bad server address specified.");
          break;
        case drogon::ReqResult::Timeout:
          LOG->error("Request timed out.");
          break;
        case drogon::ReqResult::HandshakeError:
          LOG->error("TLS handshake error occurred.");
          break;
        case drogon::ReqResult::InvalidCertificate:
          LOG->error("The invalid certificate from server.");
          break;
        case drogon::ReqResult::EncryptionFailure:
          LOG->error("The encryption failed.");
          break;
        default:
          break;
        }
        LOG->info("Connection failure.");
      }
      LOG->info("Request time: {} ms", duration);
      callback(resp);
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