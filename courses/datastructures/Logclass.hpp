#include <spdlog/sinks/basic_file_sink.h>
#include <mutex>

class Logger {
public:
  Logger(const Logger&) = delete;
  Logger(Logger&&) = delete;
  ~Logger() = default;
  Logger& operator=(const Logger&) = delete;
  Logger& operator=(Logger&&) = delete;

  static Logger& getInstance() {
    static std::once_flag flag;
    std::call_once(flag, []() {
      instance = std::make_shared<Logger>();
    });
    return *instance;
  };


private:
  Logger() {
    LOG=spdlog::basic_logger_mt("courses_logger", "log.txt");
  };

  static std::shared_ptr<Logger> instance;
  static std::shared_ptr<spdlog::logger> LOG;
};


#define LOG Logger::getInstance()
