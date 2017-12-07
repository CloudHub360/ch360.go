Describe "ch360 --version" {
  It "should output date of release" {
    ch360 --version | Should -Match "\d{4}\.\d{2}\.\d{2}"
  }

  It "should output short git hash" {
    ch360 --version | Should -Match ".*-([a-f0-9]{7,})"
  }

  It "should output build number" {
    ch360 --version | Should -Match ".*-(\d+)"
  }

  It "should output the supported format" {
    ch360 --version | Should -Match "\d{4}\.\d{2}\.\d{2}-([a-f0-9]{7,})-(\d+)"
  }
}
