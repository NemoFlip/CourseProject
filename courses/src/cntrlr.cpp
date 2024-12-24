#include "../include/cntrlr/cntrlr.hpp"
#include "../datastructures/Database.hpp"
#include "../datastructures/Logclass.hpp"

void Cntrlr::get_courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Json::Value response_body;

  Database::getInstance().select_query("courses", basic_values::condition, basic_values::sorting,
    [callback = std::move(callback), &response_body](const drogon::orm::Result& res) {
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
    [callback = std::move(callback), &response_body](const drogon::orm::DrogonDbException &exc) {
      response_body["Error"] = exc.base().what();
      auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
      callback(resp);
    });
}

void Cntrlr::join_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  auto body = req->getJsonObject();

  if (!body || !body->isMember("course_id") || !body->isMember("user_id")) {
    auto resp = drogon::HttpResponse::newHttpJsonResponse({ {"Error", "Invalid request"}});
    resp->setStatusCode(drogon::k400BadRequest);
    callback(resp);
  }
  else {
    const std::string course_id = (*body)["course_id"].asString();
    const std::string user_id = (*body)["user_id"].asString();

    Json::Value response_body;
    Database::getInstance().insert_query("courses_participation",
      "user_id, course_id",
      "?,?",
      [callback = std::move(callback)](const drogon::orm::Result& result) {
        auto resp = drogon::HttpResponse::newHttpResponse();
        callback(resp);
      },
      [callback = std::move(callback), &response_body](const drogon::orm::DrogonDbException& exc) {
        response_body["Error"] = exc.base().what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(response_body);
        callback(resp);
      },
      user_id,
      course_id);
  }
}

void Cntrlr::leave_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {

}

void Cntrlr::assess_course(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {

}