#include "../include/cntrlr/cntrlr.hpp"
#include "../datastructures/Database.hpp"
#include "../datastructures/Logclass.hpp"

enum class CourseStates{
  JoinCourse,
  LeaveCourse
};

void sendJoinCourseQuery(Json::Value& response_body, const std::string& user_id, const std::string& course_id,
  std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Database::getInstance().insert_query("courses_participation",
    "user_id, course_id",
    "?,?",
    [callback = std::forward<decltype(callback)>(callback)](const drogon::orm::Result& result) {
      auto resp = drogon::HttpResponse::newHttpResponse();
      resp->setStatusCode(drogon::k204NoContent);
      callback(resp);
    },
    [callback = std::forward<decltype(callback)>(callback), &response_body](const drogon::orm::DrogonDbException& exc) {
      response_body["Error"] = exc.base().what();
      auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
      resp->setStatusCode(drogon::k500InternalServerError);
      callback(resp);
    },
    user_id,
    course_id);
}

void sendLeaveCourseQuery(Json::Value& response_body, const std::string& user_id, const std::string& course_id,
  std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Database::getInstance().delete_query("courses_participation",
    "user_id = ? and course_id = ?",
    [callback = std::forward<decltype(callback)>(callback)](const drogon::orm::Result& result) {
      auto resp = drogon::HttpResponse::newHttpResponse();
      resp->setStatusCode(drogon::k204NoContent);
      callback(resp);
    },
    [callback = std::forward<decltype(callback)>(callback), &response_body](const drogon::orm::DrogonDbException& exc) {
      response_body["Error"] = exc.base().what();
      auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
      resp->setStatusCode(drogon::k500InternalServerError);
      callback(resp);
    },
    user_id,
    course_id);
}

void courseParticipation(const CourseStates& State, const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  auto body = req->getJsonObject();

  if (body->empty() || !body->isMember("course_id") || !body->isMember("user_id") ||
    (*body)["course_id"].asString().empty() || (*body)["user_id"].asString().empty()) {
    auto resp = drogon::HttpResponse::newHttpJsonResponse({ {"Error", "Invalid request"} });
    resp->setStatusCode(drogon::k400BadRequest);
    callback(resp);
    return;
  }
  const std::string course_id = (*body)["course_id"].asString();
  const std::string user_id = (*body)["user_id"].asString();

  Json::Value response_body;
  switch (State)
  {
  case(CourseStates::JoinCourse):
    sendJoinCourseQuery(response_body, user_id, course_id, std::forward<decltype(callback)>(callback));
    break;
  case(CourseStates::LeaveCourse):
    sendLeaveCourseQuery(response_body, user_id, course_id, std::forward<decltype(callback)>(callback));
    break;
  default:
    break;
  }
}

void Cntrlr::getCourses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Json::Value response_body;

  Database::getInstance().select_query("courses", basic_values::condition, basic_values::sorting,
    [callback = std::forward<decltype(callback)>(callback), &response_body](const drogon::orm::Result& res) {
      if (res.empty()) {
        response_body["Error"] = "No available courses";
        return;
      }

      for (const auto& rows : res) {
        Json::Value tmp;
          
        for (const auto& column : rows) {
          if (!column.isNull()) {
            tmp[column.name()] = column.as<std::string>();
          }
          else {
            tmp[column.name()] = "Null";
          }
        }

        response_body.append(tmp);
      }
      auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
      callback(resp);
    },
    [callback = std::forward<decltype(callback)>(callback), &response_body](const drogon::orm::DrogonDbException &exc) {
      response_body["Error"] = exc.base().what();
      auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
      resp->setStatusCode(drogon::k500InternalServerError);
      callback(resp);
    });
}

void Cntrlr::joinCourse(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  courseParticipation(CourseStates::JoinCourse, req, std::forward<decltype(callback)>(callback));
}

void Cntrlr::leaveCourse(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  courseParticipation(CourseStates::LeaveCourse, req, std::forward<decltype(callback)>(callback));
}

void Cntrlr::rateCourse(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {

}

void Cntrlr::changeCourseRating(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {

}