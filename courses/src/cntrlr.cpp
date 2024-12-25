#include "../include/cntrlr/cntrlr.hpp"
#include "../datastructures/Database.hpp"
#include "../datastructures/Logclass.hpp"

enum class CourseStates{
  JoinCourse,
  LeaveCourse
};

void CourseParticipation(const CourseStates& State, const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  auto body = req->getJsonObject();

  if (body->empty() || !body->isMember("course_id") || !body->isMember("user_id") ||
    (*body)["course_id"].asString().empty() || (*body)["user_id"].asString().empty()) {
    auto resp = drogon::HttpResponse::newHttpJsonResponse({ {"Error", "Invalid request"} });
    resp->setStatusCode(drogon::k400BadRequest);
    callback(resp);
  }
  else {
    const std::string course_id = (*body)["course_id"].asString();
    const std::string user_id = (*body)["user_id"].asString();

    Json::Value response_body;
    if (State == CourseStates::JoinCourse) {
      Database::getInstance().insert_query("courses_participation",
        "user_id, course_id",
        "?,?",
        [callback = std::move(callback)](const drogon::orm::Result& result) {
          auto resp = drogon::HttpResponse::newHttpResponse();
          resp->setStatusCode(drogon::k204NoContent);
          callback(resp);
        },
        [callback = std::move(callback), &response_body](const drogon::orm::DrogonDbException& exc) {
          response_body["Error"] = exc.base().what();
          auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
          resp->setStatusCode(drogon::k500InternalServerError);
          callback(resp);
        },
        user_id,
        course_id);
    }
    else {
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
  }
}

void Cntrlr::get_courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Json::Value response_body;

  Database::getInstance().select_query("courses", basic_values::condition, basic_values::sorting,
    [callback = std::forward<decltype(callback)>(callback), &response_body](const drogon::orm::Result& res) {
      if (res.empty()) {
        response_body["Error"] = "No available courses";
      }
      else {
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

void Cntrlr::join_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  CourseParticipation(CourseStates::JoinCourse, req, std::forward<decltype(callback)>(callback));
}

void Cntrlr::leave_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  CourseParticipation(CourseStates::LeaveCourse, req, std::forward<decltype(callback)>(callback));
}

void Cntrlr::assess_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {

}