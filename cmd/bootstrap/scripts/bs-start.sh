#!/bin/bash

# Check if the number of arguments provided is not exactly one
if [ $# -ne 1 ]; then
  echo "USAGE: $0 <geoinfo-api-bin>"
  exit 1
fi

# Define an associative array mapping environment variable names to their corresponding command line arguments
declare -A arg_map=(
  ["PG_DB_USER"]="--db.user"
  ["PG_DB_PASSWORD"]="--db.password"
  ["PG_DB_DBNAME"]="--db.dbname"
  ["PG_DB_PORT"]="--db.port"
  ["PG_DB_HOST"]="--db.host"
  ["PG_DB_SSLCERT"]="--db.sslcert"
  ["PG_DB_SSLKEY"]="--db.sslkey"
  ["PG_DB_SSLROOTCERT"]="--db.sslrootcert"
  ["PG_DB_SSLMODE"]="--db.sslmode"
)

# Initialize an empty string to store the command line arguments
ARGS=""

# Loop through each environment variable in the array
for envar in "${!arg_map[@]}"; do 
  # Check if the environment variable is set to a non-empty value
  if [ -n "${!envar}" ]; then
    # Build the command line argument by concatenating the argument name and its value
    arg="${arg_map[$envar]}=${!envar}"
    # Append the argument to the ARGS string
    ARGS="$ARGS $arg"
  fi
done

if [ -n "$ARGS" ]; then
  # Call the binary with the generated command line arguments
  $1 $ARGS
else
  echo "Hint: You can either set each command line option manually, or populate the 'config.toml' and 'secret.toml' files and then pass the appropriate command line options. 

Example usage:
- Set all options manually:
  $ geoinfo-api --arg1=value1 --arg2=value2 --arg3=value3

- Populate 'config.toml' and 'secret.toml' files with your values and pass the appropriate command line options:
  $ geoinfo-api --config=config.toml --secret=secret.toml
"
  exit 1
fi
