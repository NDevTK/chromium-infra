create {
  source {
    script {
      name: "fetch.py"
      name: "prebuilt"
    }
    unpack_archive: true
    no_archive_prune: true
  }
  build {
    install: "install_official.sh"
  }
  package {
    disable_latest_ref: true
  }
}

create {
  platform_re: "linux-mips.*"
  source {
    script {
      name: "fetch.py"
      name: "source"
    }
  }
  build {
    install: "install_source.sh"
    tool: "tools/go116@1.16.7"  # depend on the prebuilt version in $PATH
  }
}

upload {
  pkg_prefix: "tools"
  pkg_name_override: "go"
}
