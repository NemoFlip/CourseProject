#include <drogon/drogon.h>


int main() {
  drogon::app().setLogPath("./")
    .addListener("127.0.0.1", 8048)
    .setThreadNum(16)
    .run();
}