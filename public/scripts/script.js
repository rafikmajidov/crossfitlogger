    var routes = {
        memberLogin: '/member/login',
        memberAuthenticate: '/member/authenticate',
        memberAuthenticate: '/member/authenticate',
        memberRegister: '/member/register',
        memberForgetPassword: '/member/forgetpassword',
        memberResetPassword: '/member/resetpassword',
        memberLogout: '/member/logout'
    };
    console.log(routes)

    angular.module('project', ['ngRoute', 'ngResource'])

        .factory('Member', ['$resource',
            function($resource) {
                return $resource('/member/:entryId', {}, {
                    query: {method:'GET', params: { entryId : '' }, isArray: false},
                    authenticate: {method: 'POST', url: routes.memberAuthenticate, isArray: false},
                    save: {method: 'POST', url: routes.memberRegister},
                    forgetPassword: {method: 'POST', url: routes.memberForgetPassword},
                    resetPassword: {method: 'POST', url: routes.memberResetPassword},
                    logout: {method: 'POST', url: routes.memberLogout}
                });
            }]
        )

        // register the interceptor as a service
        .factory('authInterceptor', function($q, $window, $location) {
            return {
                // optional method
                'request': function(config) {
                    config.headers = config.headers || {};
                    if ($window.localStorage['member']) {
                        config.headers.Authorization = $window.localStorage['member'];
                    }
                    return config;
                },
                // optional method
                'responseError': function(rejection) {
                    // do something on error
                    if (rejection.status === 401) {
                        // handle the case where the member is not authenticated
                    }
                    return $q.reject(rejection);
                }
            }
        })

        .config(function($routeProvider) {
            $routeProvider
                .when('/', {
                    controller:'HomeCtrl',
                    templateUrl:'./templates/home'
                })
                .when(routes.memberLogin, {
                    controller:'LoginCtrl',
                    templateUrl:'./templates/login'
                })
                .when(routes.memberLogout, {
                    controller:'LogoutCtrl',
                    templateUrl:'./templates/logout'
                })
                .when(routes.memberForgetPassword, {
                    controller:'ForgetPasswordCtrl',
                    templateUrl:'./templates/forgetpassword'
                })
                .when(routes.memberResetPassword, {
                    controller:'ResetPasswordCtrl',
                    templateUrl:'./templates/resetpassword'
                })
                .otherwise({
                    redirectTo:'/'
                });
        })

        .config(function ($httpProvider) {
            $httpProvider.interceptors.push('authInterceptor');
        })

        .controller('HomeCtrl', function($window, $location, $rootScope) {
            console.log('in HomeCtrl')
            $rootScope.pageTitle = 'Home';
            if ($window.localStorage['member']) {
                $rootScope.member = $window.localStorage['member'];
                console.log('member present');
                console.log($rootScope.member);
            } else {
                console.log('we have to authenticate that member');
                $location.path(routes.memberLogin);
            }
        })
    
        .controller('LoginCtrl', function($window, $location, $scope, $rootScope, Member) {
            $rootScope.pageTitle = 'Login | Register';
            if ($window.localStorage['member']) {
                $location.path('/')
            } else {
                $scope.login = function(visitor) {
                    if (visitor != undefined && visitor.email != undefined && visitor.pass != undefined) {
                        console.log("login submit")
                        Member.authenticate({email:visitor.email, pass: visitor.pass})
                            .$promise.then(function(member) {
                                if (member.S == 1 && member.A != undefined) {
                                    $window.localStorage['member'] = member.A
                                    $location.path('/')
                                } else {
                                    $scope.errorMessage = member.A;
                                }
                            });
                    };
                }
                $scope.register = function(visitor) {
                    if (visitor != undefined && visitor.email != undefined && visitor.pass != undefined) {
                        console.log("register submit")
                        Member.save({email:visitor.email, pass: visitor.pass})
                            .$promise.then(function(member) {
                                if (member.S == 1 && member.A != undefined) {
                                    $window.localStorage['member'] = member.A
                                    $location.path('/')
                                } else {
                                    $scope.errorMessage = member.A;
                                }
                            });
                    };
                }
            }
        })

        .controller('LogoutCtrl', function($window, $location, $rootScope, Member) {
            if ($window.localStorage['member']) {
                console.log('LogoutCtrl')
                
                Member.logout()
                    .$promise.then(function(member) {
                        if (member.S == 1 && member.A != undefined) {
                            console.log('success after logout')
                            $window.localStorage.removeItem('member');
                            $rootScope.member = null;
                            $location.path('/')
                        } else {
                            $scope.errorMessage = member.A;
                            console.log('failure after logout')
                        }
                    });

            } else {
                console.log('member was logged out before');
                $location.path('/');
            }
        })

        .controller('ForgetPasswordCtrl', function($window, $location, $scope, $rootScope, Member) {
            $rootScope.pageTitle = 'Forget password';
            $scope.forgetPassword = function(visitor) {
                    if (visitor != undefined && visitor.email != undefined ) {
                        console.log("forget password submit")
                        Member.forgetPassword({email:visitor.email})
                            .$promise.then(function(member) {
                                    $scope.errorMessage = member.A;
                            });
                    };
                }    
            
        })

        .controller('ResetPasswordCtrl', function($window, $location, $rootScope, $scope, Member) {
            $rootScope.pageTitle = 'Reset password';
            var search = $location.search()

            if (search == undefined || search.t == undefined) {
                $location.path('/');
            } 

            $scope.resetPassword = function(visitor) {
                if (visitor != undefined 
                    && visitor.pass != undefined
                    && visitor.pass2 != undefined 
                    && search != undefined
                    && search.t!= undefined) {
                    console.log("reset password submit")
                    if (visitor.pass != visitor.pass2) {
                        $scope.errorMessage = "Passwords mismatch";
                        return ;
                    }
                    Member.resetPassword({pass:visitor.pass, t: search.t})
                        .$promise.then(function(member) {
                                if (member.S == 1 && member.A != undefined) {
                                    $window.localStorage['member'] = member.A
                                    $location.path('/')
                                } else {
                                    $scope.errorMessage = member.A;
                                }
                            });
                };
            }
            
            
        });

