#include <drogon/drogon.h>
#include <mutex>
class Database {

public:
  Database(const Database&) = delete;
  Database& operator=(const Database&) = delete;
  Database(const Database&&) = delete;
  Database& operator=(const Database&&) = delete;

  static Database& getInstance() {
    static std::once_flag flag;
    std::call_once(flag, []() {
      instance = std::make_shared<Database>();
    });
    return *instance;
  };

  template <typename FUNCTION1, typename FUNCTION2, typename... Arguments>
  void do_query(const std::string& query, FUNCTION1&& callback, FUNCTION2&& exceptCallback, Arguments &&...args) noexcept {
    this->db_ptr->execSqlAsync(query, callback, exceptCallback, args);
    
  };

private:
  Database() {
    db_client = drogon::app().getDbClient("courses");
  };
  ~Database() = default;

  static std::shared_ptr<Database> instance;
  static std::shared_ptr<drogon::orm::DbClient> db_client;
  static std::mutex mutex_;
};
