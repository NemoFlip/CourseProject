
#include <drogon/drogon.h>
#include <drogon/HttpController.h>

class Cntrlr : public drogon::HttpController<Cntrlr> {
public:
  METHOD_LIST_BEGIN

  METHOD_ADD(Cntrlr::getCourses, "/courses/list", drogon::Get);
  METHOD_ADD(Cntrlr::joinCourse, "/courses/join_course", drogon::Post);
  METHOD_ADD(Cntrlr::leaveCourse, "/courses/leave_course", drogon::Delete);
  METHOD_ADD(Cntrlr::rateCourse, "/courses/rate_course", drogon::Post);
  METHOD_ADD(Cntrlr::changeCourseRating, "/courses/change_rating", drogon::Put);

  METHOD_LIST_END
protected:

  void getCourses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void joinCourse(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void leaveCourse(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void rateCourse(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);
  void changeCourseRating(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback);

};