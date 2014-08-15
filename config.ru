# http://stackoverflow.com/questions/22261624/puts-output-not-displaying-in-heroku-logs-for-sinatra-app
$stdout.sync = true

# This file is used by Rack-based servers to start the application.

require ::File.expand_path('../config/environment',  __FILE__)
run Rails.application