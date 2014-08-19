angular = require 'angular'
Spinner = require 'spin.js'

angular.module('app').directive 'loading', ->
  restrict: 'A'
  scope:
    loading: '='
    loadingStyle: '='
  link: (scope, elem, attrs) ->
    options = angular.extend
      color: '#FFB064'
      width: 3
      lines: 10
    , scope.loadingStyle

    spinner = new Spinner options

    scope.$watch 'loading', (loading) ->
      if loading
        spinner.spin elem[0]
      else
        spinner.stop()