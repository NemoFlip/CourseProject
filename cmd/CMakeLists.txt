
cmake_minimum_required (VERSION 3.8)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CURRENT_SOURCE_DIR D:/) 
set(CMAKE_TOOLCHAIN_FILE ${CMAKE_CURRENT_SOURCE_DIR}/vcpkg/scripts/buildsystems/vcpkg.cmake CACHE STRING "Vcpkg toolchain file")

if (POLICY CMP0141)
  cmake_policy(SET CMP0141 NEW)
  set(CMAKE_MSVC_DEBUG_INFORMATION_FORMAT "$<IF:$<AND:$<C_COMPILER_ID:MSVC>,$<CXX_COMPILER_ID:MSVC>>,$<$<CONFIG:Debug,RelWithDebInfo>:EditAndContinue>,$<$<CONFIG:Debug,RelWithDebInfo>:ProgramDatabase>>")
endif()

project(cmd)

find_package(Boost CONFIG REQUIRED)
find_package(Drogon CONFIG REQUIRED)


add_executable(${PROJECT_NAME} app/main.cpp
              src/LoginController.cpp
              src/ProfileController.cpp
              include/login/LoginController.hpp
              include/profile/ProfileController.hpp
)

target_link_libraries(${PROJECT_NAME} Drogon::Drogon ${BOOST_LIBS})

include_directories(${Boost_INCLUDE_DIRS})
