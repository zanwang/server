angular = require 'angular'

angular.module('app').directive 'equalTo', ->
  require: 'ngModel'
  scope:
    equalTo: '='
  link: (scope, elem, attrs, ctrl) ->
    ctrl.$parsers.unshift (viewValue) ->
      if viewValue is scope.equalTo
        ctrl.$setValidity 'equal', true
        viewValue
      else
        ctrl.$setValidity 'equal', false
