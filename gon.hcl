source = ["./qbee-cli"]
bundle_id = "io.qbee.qbee-cli"

apple_id {
  username = "mitchell@example.com"
  password = "@env:AC_PASSWORD"
  provider = "UL304B4VGY"
}

sign {
  application_identity = "Developer ID Application: Mitchell Hashimoto"
}

dmg {
  output_path = "qbee-cli.dmg"
  volume_name = "qbee-cli"
}

zip {
  output_path = "qbee-cli.zip"
}