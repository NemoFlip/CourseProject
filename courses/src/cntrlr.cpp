#include "../include/cntrlr/cntrlr.hpp"


void Cntrlr::courses(const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
  auto courses_table = drogon::app().getDbClient("courses");
  
  courses_table->execSqlAsync("select name, description, assess from courses",
    [callback = std::move(callback)](const drogon::orm::Result& result) {
      Json::Value courses_list;

      if (result.empty()) {
        courses_list["error"] = "No courses availaible";
      } else for (auto& row : result) {
        //adding row by row to the end
        Json::Value tmp;
        tmp["name"] = row["name"].as<std::string>();
        tmp["description"] = row["description"].as<std::string>();
        tmp["assess"] = row["assess"].as<int>();
        courses_list.append(tmp);
      }
      
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