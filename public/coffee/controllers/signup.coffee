angular = require 'angular'

angular.module('app').controller 'SignupCtrl', ($scope, User) ->
  $scope.signup = ->
    return if $scope.signupForm.$invalid

    User.post
      name: $scope.name
      email: $scope.email
      password: $scope.password
    .then ->
      console.log arguments
