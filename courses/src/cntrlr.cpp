#include "../include/cntrlr/cntrlr.hpp"
#include "../datastructures/Database.hpp"
#include "../datastructures/Logclass.hpp"

void Cntrlr::get_courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Json::Value courses_list;

  Database::getInstance().do_query(
    "select * from courses_info",
    [callback = std::move(callback), &courses_list](const drogon::orm::Result& result) {

      if (result.empty()) {
        courses_list["error"] = "No courses available";
      }
      else {
        for (auto& rows : result) {
        
        }
      }

      auto resp = drogon::HttpResponse::newHttpJsonResponse(courses_list);
      callback(resp);
    },
    [callback = std::move(callback), &courses_list](const drogon::orm::DrogonDbException& e) {
      courses_list["error"] = e.base().what();

      auto resp = drogon::HttpResponse::newHttpJsonResponse(courses_list);
      callback(resp);
    }
  );
}