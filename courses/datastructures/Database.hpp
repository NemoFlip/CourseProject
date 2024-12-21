#include <drogon/drogon.h>
#include <mutex>
#include <shared_mutex>
#include <vector>

struct basic_values {
  std::string sorting = "1";
  std::string condition = "1=1";
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

  template<typename FUNCTION1, typename FUNCTION2, typename... Args>
  void select_query(const std::string& from,
    std::string& conditinon,
    std::string& sorting,
    FUNCTION1 rCallBack,
    FUNCTION2 eCallback,
    Args&&... args) {
    std::shared_lock<std::shared_mutex> lock(shmutex_);
    instance->commit(std::format("select * from {} where {} order by {} limit ? offset ?", from, condition, sorting),
      std::forward<FUNCTION1>(rCallBack),
      std::forward<FUNCTION2>(eCallback),
      std::forward<Args...>(args));
  }
   
  template<typename FUNCTION1, typename FUNCTION2, typename... Args>
  void delete_query(const std::string& from,
    std::string& conditinon,
    FUNCTION1 rCallBack,
    FUNCTION2 eCallback,
    Args&&... args) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);
    instance->commit(std::format("delete from {} where {}", from, condition),
      std::forward<FUNCTION1>(rCallBack),
      std::forward<FUNCTION2>(eCallback),
      std::forward<Args...>(args));
  } 

  template<typename FUNCTION1, typename FUNCTION2, typename... Args>
  void insert_query(const std::string& into,
    std::string& parametrs,
    std::string& value,
    FUNCTION1 rCallBack,
    FUNCTION2 eCallback,
    Args&&... args) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);
    instance->commit(std::format("insert into {} ({}) values ({})", into, parametrs, value),
      std::forward<FUNCTION1>(rCallBack),
      std::forward<FUNCTION2>(eCallback),
      std::forward<Args...>(args));
  }

  template<typename FUNCTION1, typename FUNCTION2, typename... Args>
  void update_query(const std::string& to,
    std::string& upd_data,
    std::string& condirion,
    FUNCTION1 rCallBack,
    FUNCTION2 eCallback,
    Args&&... args) {
    std::unique_lock<std::shared_mutex> lock(shmutex_);
    instance->commit(std::format("update {} set {} where {}", to, upd_data, condition),
      std::forward<FUNCTION1>(rCallBack),
      std::forward<FUNCTION2>(eCallback),
      std::forward<Args...>(args));
  }


private:
  Database() {
    db_client = drogon::app().getDbClient("courses");
  };
  ~Database() = default;

  template<typename FUNCTION1, typename FUNCTION2, typename... Args>
  void commit(const std::string& query, FUNCTION1&& rCallBack, FUNCTION2&& eCallBack, Args&&... args) {
    instance->db_client->execSqlAsync(query, rCallBack, eCallBack, args);
  }

  static std::shared_ptr<Database> instance;
  static std::shared_ptr<drogon::orm::DbClient> db_client;
  static std::shared_mutex shmutex_;
};
