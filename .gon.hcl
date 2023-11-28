source = ["./qbee-cli"]
bundle_id = "io.qbee.qbee-cli"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "@env:AC_APPLICATION_IDENTITY"
}

dmg {
  output_path = "qbee-cli.dmg"
  volume_name = "qbee-cli"
}
