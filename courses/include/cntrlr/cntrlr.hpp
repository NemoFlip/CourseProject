#include <drogon/drogon.h>
#include <drogon/HttpController.h>



class Cntrlr : public drogon::HttpController<Cntrlr> {
public:
  METHOD_LIST_BEGIN

  METHOD_ADD(Cntrlr::courses, "/courses", drogon::Get);
  METHOD_ADD(Cntrlr::join_course, "/join", drogon::Post);
  METHOD_ADD(Cntrlr::leave_course, "/leave", drogon::Delete);

  METHOD_LIST_END

protected:

  void courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void join_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void leave_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);


};