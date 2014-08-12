angular = require 'angular'

angular.module('app').controller 'SignupCtrl', ($scope, $http) ->
  $scope.signup = ->
    return if $scope.signupForm.$invalid

    $http.post '/api/v1/users',
      name: $scope.name
      password: $scope.password
      email: $scope.email
    .success (data) ->
      console.log data
