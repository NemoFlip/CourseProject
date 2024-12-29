#include <drogon/drogon.h>
#include <shared_mutex>

struct basic_values {
  const static std::string sorting;
  const static std::string condition;
};
auto basic_values::condition = "1=1";
auto basic_values::sorting = "1";

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
    const std::string& conditinon,
    const std::string& sorting,
    FUNCTION1 rCallBack,
    FUNCTION2 eCallback,
    Args&&... args) {
    std::shared_lock<std::shared_mutex> lock(shmutex_);
    instance->commit(std::format("select * from {} where {} order by {}", from, condition, sorting),
      std::forward<FUNCTION1>(rCallBack),
      std::forward<FUNCTION2>(eCallback),
      std::forward<Args...>(args));
  }
   
  template<typename FUNCTION1, typename FUNCTION2, typename... Args>
  void delete_query(const std::string& from,
    const std::string& conditinon,
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
    const std::string& parametrs,
    const std::string& value,
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
    const std::string& upd_data,
    const std::string& condition,
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
