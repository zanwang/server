angular = require 'angular'
require 'angular-ui-router'

angular.module 'app', ['ui.router']

require './services'
require './directives'
require './controllers'
require './routes'
