angular = require 'angular'

rDomainName = /^[a-zA-Z]+[a-zA-Z\d\-]*$/

angular.module('app').directive 'domainName', ->
  require: 'ngModel'
  link: (scope, elem, attrs, ctrl) ->
    scope.$watch ->
      ctrl.$modelValue
    , (value) ->
      ctrl.$setValidity 'domainName', if value then rDomainName.test value else true