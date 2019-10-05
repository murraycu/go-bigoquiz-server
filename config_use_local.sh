# Use the .local versions of the config files.
# (These cannot be added to the git repository.)

prefix="./config_oauth2/"
suffix="_credentials_secret.json"
end=".local"

function copy_file() {
  org=$1
  filename=${prefix}${org}${suffix}
  cp ${filename}${end} ${filename}

}

copy_file "google"
copy_file "github"
copy_file "facebook"
