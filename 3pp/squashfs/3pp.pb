create {
  platform_re: "linux-.*"
  source {
    url {
      download_url: "https://github.com/plougher/squashfs-tools/archive/4.4.tar.gz"
      version: "4.4"
    }
    unpack_archive: true
    cpe_base_address: "cpe:/a:phillip_lougher:squashfs:4.4"
    patch_dir: "patches"
    patch_version: "chromium.2"
  }
  build {
    dep: "static_libs/libzstd"
    dep: "static_libs/zlib"
    tool: "tools/sed"
  }
}

upload { pkg_prefix: "tools" }
