angular = require 'angular'

angular.module('app')
.constant('API_BASE_URL', '/api/v1/')
.constant('Config', require '../config/development.json')
.config ($stateProvider, $urlRouterProvider, $locationProvider, $httpProvider) ->
  $urlRouterProvider.otherwise '/app'
  $locationProvider.html5Mode true

  $httpProvider.interceptors.push 'TokenInterceptor'

  $stateProvider
    .state 'app',
      url: ''
      abstract: true
      data:
        protected: true
      template: require '../views/app'
    .state 'app.domains',
      url: '/app'
      data:
        title: 'Domains'
      template: require '../views/domains'
      controller: 'DomainListCtrl'
    .state 'app.settings',
      url: '/settings'
      data:
        title: 'Settings'
      template: require '../views/settings'
      controller: 'SettingsCtrl'
    .state 'login',
      url: '/login'
      template: require '../views/login'
      controller: 'LoginCtrl'
      data:
        skipIfAuthorized: true
        title: 'Log in'
    .state 'signup',
      url: '/signup'
      template: require '../views/signup'
      controller: 'SignupCtrl'
      data:
        skipIfAuthorized: true
        title: 'Sign up'
    .state 'password',
      url: '/forgot_password'
      template: require '../views/password'
      controller: 'PasswordCtrl'
      data:
        title: 'Forgot password'

.run ($rootScope, Auth, $state) ->
  $rootScope.$on '$stateChangeStart', (event, toState, toParams, fromState, fromParams) ->
    # Check authorization
    if Auth.token
      if toState.data and toState.data.skipIfAuthorized
        event.preventDefault()
        $state.go 'app.domains', {}, location: 'replace'
    else
      if toState.data and toState.data.protected
        event.preventDefault()
        $state.go 'login', {}, location: 'replace'

  $rootScope.$on '$stateChangeSuccess', (event, toState, toParams, fromState, fromParams) ->
    # Change page title
    $rootScope.$emit 'titleChanged', toState.data.title