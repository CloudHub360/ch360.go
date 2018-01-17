Describe "surf --version" {
  It "should output date of release" {
    surf --version | Should -Match "\d{4}\.\d{2}\.\d{2}"
  }

  It "should output short git hash" {
    surf --version | Should -Match ".*-([a-f0-9]{7,})"
  }

  It "should output build number" {
    surf --version | Should -Match ".*-(\d+)"
  }

  It "should output the supported format" {
    surf --version | Should -Match "\d{4}\.\d{2}\.\d{2}-([a-f0-9]{7,})-(\d+)"
  }
}
