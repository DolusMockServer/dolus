class Dolus < Formula
  desc "A configurable mock server"
  url "https://github.com/DolusMockServer/dolus.git"
  license "MIT"
  head "https://github.com/DolusMockServer/dolus.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", cmd/dolus
  end

end
