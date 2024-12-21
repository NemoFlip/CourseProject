#include "../include/cntrlr/cntrlr.hpp"
#include "../datastructures/Database.hpp"
#include "../datastructures/Logclass.hpp"

void Cntrlr::get_courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  Json::Value courses_list;


  Database::getInstance().select_query("courses", basic_values::condition, basic_values::sorting,
    [callback = std::move(callback), &courses_list](const drogon::orm::Result& res) {
      if (res.empty()) {
        courses_list["Error"] = "No available courses";
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

          courses_list.append(tmp);
        }
      }
      auto req = drogon::HttpResponse::newHttpJsonResponse(courses_list);
      callback(req);
    },
    [callback = std::move(callback), &courses_list](const drogon::orm::DrogonDbException exc) {
      courses_list["Error"] = exc.base().what();
      auto req = drogon::HttpResponse::newHttpJsonResponse(courses_list);
      callback(req);
    });
}