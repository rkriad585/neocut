cmake_minimum_required(VERSION 3.20)
project(neocut
    VERSION 1.0.2
    LANGUAGES C
    DESCRIPTION "Remove silence from MP3 audio files"
)

# Go build wrapper — CMake is not required for Go projects,
# this file exists for IDEs and toolchains that expect CMake.
#
# To build:  mkdir build && cd build && cmake .. && cmake --build .
# Or simply: go build ./cmd/neocut/

set(CMAKE_INSTALL_PREFIX ${CMAKE_BINARY_DIR}/install)

set(GO_FLAGS -ldflags "-X neocut/internal/config.Commit=${CMAKE_GIT_COMMIT}")

add_custom_target(neocut ALL
    COMMAND ${CMAKE_COMMAND} -E env
        GOFLAGS=${GO_FLAGS}
        go build -o ${CMAKE_BINARY_DIR}/neocut ./cmd/neocut/
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Building neocut..."
)

add_custom_target(test
    COMMAND go test -count=1 ./...
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Running tests..."
)

add_custom_target(clean
    COMMAND ${CMAKE_COMMAND} -E remove
        ${CMAKE_BINARY_DIR}/neocut
        ${CMAKE_BINARY_DIR}/neocut.exe
    COMMENT "Cleaning..."
)

add_custom_target(install
    COMMAND ${CMAKE_COMMAND} -E copy
        ${CMAKE_BINARY_DIR}/neocut
        ${CMAKE_INSTALL_PREFIX}/bin/neocut
    COMMENT "Installing..."
)

add_custom_target(uninstall
    COMMAND ${CMAKE_COMMAND} -E remove
        ${CMAKE_INSTALL_PREFIX}/bin/neocut
    COMMENT "Uninstalling..."
)
