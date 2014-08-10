angular = require 'angular'

angular.module('app').config ($stateProvider, $urlRouterProvider, $locationProvider, RestangularProvider) ->
  $urlRouterProvider.otherwise '/app'
  $locationProvider.html5Mode true

  RestangularProvider.setBaseUrl '/api/v1'

  $stateProvider
    .state 'app',
      url: '/app'
      template: require '../views/app'
    .state 'login',
      url: '/login'
      template: require '../views/login'
      controller: 'LoginCtrl'
    .state 'signup',
      url: '/signup'
      template: require '../views/signup'
      controller: 'SignupCtrl'
