
cmake_minimum_required (VERSION 3.8)

set(CMAKE_CXX_STANDARD 17)

if (POLICY CMP0141)
  cmake_policy(SET CMP0141 NEW)
  set(CMAKE_MSVC_DEBUG_INFORMATION_FORMAT "$<IF:$<AND:$<C_COMPILER_ID:MSVC>,$<CXX_COMPILER_ID:MSVC>>,$<$<CONFIG:Debug,RelWithDebInfo>:EditAndContinue>,$<$<CONFIG:Debug,RelWithDebInfo>:ProgramDatabase>>")
endif()

project(courses)

find_package(Drogon CONFIG REQUIRED)


add_executable(${PROJECT_NAME} app/main.cpp
               include/cntrlr/cntrlr.hpp
 "src/cntrlr.cpp")

target_link_libraries(${PROJECT_NAME} Drogon::Drogon ${BOOST_LIBS} spdlog)
