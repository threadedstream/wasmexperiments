cmake_minimum_required(VERSION 3.5)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} --target=wasm32 --no-standard-libraries -Wl,--export-all -Wl,--no-entry")
set(CMAKE_C_COMPILER "clang")

project(casm)

add_executable(casm wasm.c)
