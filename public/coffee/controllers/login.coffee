angular = require 'angular'

angular.module('app').controller 'LoginCtrl', ($scope, $http) ->
  $scope.login = ->
    return if $scope.loginForm.$invalid

    $http.post '/api/v1/tokens',
      login: $scope.name
      password: $scope.password
    .success (data) ->
      console.log data
