
#include <drogon/drogon.h>
#include <drogon/HttpController.h>

class Cntrlr : public drogon::HttpController<Cntrlr> {
public:
  METHOD_LIST_BEGIN

  METHOD_ADD(Cntrlr::get_courses, "/courses/list", drogon::Get);
  METHOD_ADD(Cntrlr::join_course, "/courses/join_course", drogon::Post);
  METHOD_ADD(Cntrlr::leave_course, "/courses/leave_course", drogon::Delete);
  METHOD_ADD(Cntrlr::assess_course, "/courses/assess_course", drogon::Put);

  METHOD_LIST_END
protected:

  void get_courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void join_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void leave_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void assess_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
};