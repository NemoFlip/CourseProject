
#ifndef NAVIGATOR_HPP
#define NAVIGATOR_HPP

#include <drogon/drogon.h>
#include <drogon/HttpController.h>

class LoginController : public drogon::HttpController<LoginController> {
public:
  METHOD_LIST_BEGIN

  METHOD_ADD(LoginController::login, "/login", drogon::Post);
  METHOD_ADD(LoginController::logout, "/logout", drogon::Post);
  METHOD_ADD(LoginController::registration, "/registration", drogon::Post);
  METHOD_ADD(LoginController::passwordRecovery, "/passwordRecovery", drogon::Put);

  METHOD_LIST_END
protected:

  void login(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

  void logout(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

  void registration(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

  void passwordRecovery(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

};

#endif // !NAVIGATOR_HPP



