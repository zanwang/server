angular = require 'angular'
require 'angular-ui-router'
require '../../bower_components/angular-resource/angular-resource'
require '../../bower_components/ngstorage/ngstorage'

angular.module 'app', ['ui.router', 'ngResource', 'ngStorage']

require './services'
require './directives'
require './controllers'
require './routes'
