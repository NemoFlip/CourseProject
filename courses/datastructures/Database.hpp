#include <drogon/drogon.h>
#include <mutex>
#include <shared_mutex>
#include <vector>

enum class QueryType{
  Select,
  Update,
  Delete,
  Insert
};

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

  template<typename OutputDataType, typename... Args>
  OutputDataType select_query(const std::string& from, Args... args) {
    std::shared_lock<std::shared_mutex> lock(shmutex_);
    try {
      instance->select(std::format("select * from {}"), from, args);
    }
    catch (std::format_error) {
      instance->select(std::format("select * from {} where {}"), from, args);
    }
  }
   
  template<typename OutputDataType, typename... Args>
  OutputDataType delete_query(const std::string& from, Args... args) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);

  }

  template<typename OutputDataType, typename... Args>
  OutputDataType insert_query(const std::string& to, Args... args) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);

  }

  template<typename OutputDataType, typename... Args>
  OutputDataType update_query(const std::string& to, Args... args) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);

  }

private:
  Database() {
    db_client = drogon::app().getDbClient("courses");
  };
  ~Database() = default;

  template<typename OutputDataType, typename... Args>
  void select(const std::string& query) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);

  }

  static std::shared_ptr<Database> instance;
  static std::shared_ptr<drogon::orm::DbClient> db_client;
  static std::shared_mutex shmutex_;
};
