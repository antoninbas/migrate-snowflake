# migrate-snowflake

This is a thin wrapper around
[golang-migrate](https://github.com/golang-migrate/migrate), which only supports
Snowflake. This project consumes golang-migrate as a library.

The upstream golang-migrate CLI has the following issues;

 * no support for specifying the Snowflake virtual warehouse to use for queries
   (this is an issue specific to the provided driver for Snowflake)
 * credentials must be provided on the command-line as part of the Database URL,
   with no support for providing them as environment variables
