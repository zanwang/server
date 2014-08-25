angular = require 'angular'

angular.module('app').controller 'LoginCtrl', ($scope, Token, $state, Auth, Facebook, $window) ->
  $scope.submitted = false
  $scope.submitting = false
  $scope.socialLogging = false
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
    return if $scope.socialLogging

    $scope.socialLogging = true

    Facebook.then (FB) ->
      FB.login (res) ->
        if res.status is 'connected'
          $scope.token.user_id = res.authResponse.userID
          $scope.token.access_token = res.authResponse.accessToken

          $scope.token.$facebook().then (token) ->
            Auth.create token
            $state.go 'app.domains'
          , (err) ->
            console.log err
            $scope.socialLogging = false
        else
          console.log res
          $scope.socialLogging = false
      , scope: 'public_profile,email'

  $scope.twitterLogin = ->
    return if $scope.socialLogging

    $scope.socialLogging = true

    $window.open '/oauth/twitter/login', 'twitterLogin', 'width=500,height=400'