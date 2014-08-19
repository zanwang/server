angular = require 'angular'
require 'angular-ui-router'
require 'angular-cookies'
require '../../bower_components/angular-resource/angular-resource'

angular.module 'app', ['ui.router', 'ngCookies', 'ngResource']

require './services'
require './directives'
require './controllers'
require './routes'
