angular = require 'angular'

angular.module('app').directive 'equalTo', ->
  require: 'ngModel'
  scope:
    equalTo: '='
  link: (scope, elem, attrs, ctrl) ->
    scope.$watch ->
      scope.equalTo is ctrl.$modelValue
    , (value) ->
      ctrl.$setValidity 'equal', value