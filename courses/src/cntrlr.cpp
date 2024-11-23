#include "../include/cntrlr/cntrlr.hpp"


void Cntrlr::courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  auto courses_table = drogon::app().getDbClient("courses");
  
  courses_table->execSqlAsync("select name, description, assess from courses",
    [callback = std::move(callback)](const drogon::orm::Result& result) {
      Json::Value courses_list;
      auto resp = drogon::HttpResponse::newHttpJsonResponse(courses_list);

      callback(resp);
    },

    [callback = std::move(callback)](const drogon::orm::DrogonDbException& e) {
      Json::Value courses_list;
      auto resp = drogon::HttpResponse::newHttpJsonResponse(courses_list);

      callback(resp);
    }
  );
}