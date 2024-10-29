#include <drogon/drogon.h>
#include <drogon/HttpController.h>

class Navigator : public drogon::HttpController<Navigator> {
public:
  METHOD_LIST_BEGIN

    METHOD_ADD(Navigator::login, "/login", drogon::Get);
  METHOD_ADD(Navigator::profile, "/profile", drogon::Get);
  METHOD_ADD(Navigator::courses_list, "/courses", drogon::Get);

  METHOD_LIST_END
protected:

  void login(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

  void profile(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

  void courses_list(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

};
