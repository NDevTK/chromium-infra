create {
  platform_re: ".*-386"
  unsupported: true
}

create {
  platform_re: "linux-.*|mac-.*"
  source {
    git {
      repo: "https://chromium.googlesource.com/external/github.com/Kitware/CMake"
      tag_pattern: "v%s"
    }
  }

  build {
    tool: "build_support/cmake_bootstrap"
    tool: "tools/ninja"
  }
}

upload { pkg_prefix: "tools" }
