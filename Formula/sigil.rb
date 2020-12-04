# This file was generated by GoReleaser. DO NOT EDIT.
class Sigil < Formula
  desc ""
  homepage ""
  version "0.7.2"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/zendesk/sigil/releases/download/0.7.2/sigil_0.7.2_Darwin_x86_64.tar.gz"
    sha256 "d17b078105c247b4fd5e2fde981fd7f25635525bf866168c88f8a2ccb86b6fdf"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/zendesk/sigil/releases/download/0.7.2/sigil_0.7.2_Linux_x86_64.tar.gz"
    sha256 "678b1b4c9bf46c7c6deb8f0b7be6f6f670318b0916d246fc8277f66d0e929f6e"
  end

  def install
    bin.install "sigil"
  end
end
