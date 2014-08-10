angular = require 'angular'
global._ = require 'lodash'
require 'angular-ui-router'
require 'restangular'

angular.module 'app', ['ui.router', 'restangular']

require './services'
require './directives'
require './controllers'
require './routes'
