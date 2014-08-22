angular = require 'angular'

angular.module('app').controller 'LoginCtrl', ($scope, Token, $state, Auth, Facebook) ->
  $scope.submitted = false
  $scope.submitting = false
  $scope.token = new Token()

  $scope.submit = ->
    return if $scope.submitting

    $scope.submitted = true
    return if $scope.loginForm.$invalid

    $scope.submitting = true

    $scope.token.$save().then (token) ->
      Auth.create token
      $state.go 'app.domains'
    , (err) ->
      $scope.errors = err.data.errors
      $scope.submitting = false

  $scope.facebookLogin = ->
    Facebook.then (FB) ->
      FB.login (res) ->
        if res.status is 'connected'
          $scope.token.user_id = res.authResponse.userID
          $scope.token.access_token = res.authResponse.accessToken

          $scope.token.$facebook().then (token) ->
            Auth.create token
            $state.go 'app.domains'
      , scope: 'public_profile,email'