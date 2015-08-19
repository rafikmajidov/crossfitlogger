    angular.module('project', ['ngRoute', 'ngResource'])

        .factory('Member', ['$resource',
            function($resource) {
                return $resource('/member/:entryId', {}, {
                    query: {method:'GET', params: { entryId : '' }, isArray: false},
                    authenticate: {method: 'POST', url: '/member/authenticate', isArray: false},
                    save: {method: 'POST', url: '/member/register'},
                    forgetPassword: {method: 'POST', url: '/member/forgetpassword'},
                    resetPassword: {method: 'POST', url: '/member/resetpassword'},
                    logout: {method: 'POST', url: '/member/logout'}
                });
            }]
        )

        .factory('Movement', 
            function($resource) {
                return $resource('/movement/:entryId', {}, {
                    query: {method:'GET', url: '/movement/all', isArray: true}
                    
                });
            }
        )

        // .factory('SessionMovement', 
        //     function($resource) {
        //         return $resource('/sessionmovement/:entryId', {}, {
        //             query: {method:'GET', url: '/sessionmovement/all', isArray: true},
        //             save: {method: 'POST', url: '/sessionmovement/save'},
                    
        //         });
        //     }
        // )

        .factory('SessionAllMovements', 
            function($resource) {
                return $resource('/sessionallmovements/:entryId', {}, {
                    query: {method:'GET', url: '/sessionallmovements/all', isArray: true},
                    save: {method: 'POST', url: '/sessionallmovements/save'},
                    getone: {method:'GET', url: '/sessionallmovements/get/:sessionId'},
                });
            }
        )


        

        // .factory('AllMovements', ['Movement', '$q',
        //     function(Movement, $q) {
        //         return function() {
        //             var delay = $q.defer();
        //             Movement.query(function(movements) {
        //                 delay.resolve(movements);
        //             }, function() {
        //                 delay.reject('Unable to fetch all movements');
        //             });


        //             return delay.promise;
        //         }
        //     }]
        // )

        .factory('Authentication', 
            function($window) {
                return {
                    getA: function() {
                        var a = '';
                        if ($window.localStorage['user']) {
                            var user = JSON.parse($window.localStorage["user"]);
                            if (user.a != undefined) {
                                a = user.a;
                            }
                        }
                        return a;
                    },
                    getB: function() {
                        var b = '';
                        if ($window.localStorage['user']) {
                            var user = JSON.parse($window.localStorage["user"]);
                            if (user.b != undefined) {
                                b = user.b;
                            }
                        }
                        return b;
                    }
                };
            })

        // register the interceptor as a service
        .factory('authInterceptor', function($q, $window, $location, Authentication) {
            return {
                // optional method
                'request': function(config) {
                    config.headers = config.headers || {};
                    if (Authentication.getA()
                        && Authentication.getB()) {
                        config.headers.Authorization = Authentication.getA() +
                            ' ' +
                            Authentication.getB() ; 
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
                .when('/member/login', {
                    controller:'LoginCtrl',
                    templateUrl:'./templates/login'
                })
                .when('/member/logout', {
                    controller:'LogoutCtrl',
                    templateUrl:'./templates/logout'
                })
                .when('/member/forget-password', {
                    controller:'ForgetPasswordCtrl',
                    templateUrl:'./templates/forgetpassword'
                })
                .when('/member/reset-password', {
                    controller:'ResetPasswordCtrl',
                    templateUrl:'./templates/resetpassword'
                })
                .when('/session/add', {
                    controller:'SessionAddEditCtrl',
                    resolve: {
                        Movements: function(Movement) {
                            return Movement.query();
                        },
                        Session: function() {}
                    },
                    templateUrl: './templates/sessionadd'                    
                })
                .when('/session/viewall', {
                    controller:'SessionViewallCtrl',
                    resolve: {
                        Sessions: function(SessionAllMovements) {
                            return SessionAllMovements.query();
                        }
                    },
                    templateUrl: './templates/sessionviewall'                    
                })
                .when('/session/:sessionId', {
                    controller:'SessionAddEditCtrl',
                    resolve: {
                        Movements: function(Movement) {
                            return Movement.query();
                        },
                        Session: function($q, $route, SessionAllMovements) {
                            var defer = $q.defer();
                            var data = SessionAllMovements.getone({sessionId:$route.current.params.sessionId})
                                .$promise.then(function(Session) {
                                    console.log('mi zdes')
                                    console.log(Session)
                                    return Session
                                });
                            defer.resolve(data);
                            return defer.promise
                            
                            // SessionAllMovements.getone({sessionId:$route.current.params.sessionId}, function(Session) {
                            //     console.log('team');
                            //     console.log(Session);
                            //     return Session;
                            // })

                        }
                    },
                    templateUrl: './templates/sessionadd'                    
                })
                .otherwise({
                    redirectTo:'/'
                });
        })

        .config(function ($httpProvider) {
            $httpProvider.interceptors.push('authInterceptor');
        })

        .controller('HomeCtrl', function($window, $location, $rootScope) {
            console.log('mi zdes')
            $rootScope.pageTitle = 'Home';
            if ($window.localStorage['member']) {
                $rootScope.member = $window.localStorage['member'];
                console.log('member present');
                console.log($rootScope.member);
            } else {
                console.log('we have to authenticate that member');
                console.log($rootScope)
                $location.path(routes.memberLogin);
            }
        })

        .controller('LoginCtrl', function($window, $location, $scope, $rootScope, User) {
            $rootScope.pageTitle = 'Login | Register';
            if ($window.localStorage['user']) {
                $location.path('/')
            } else {
                $scope.login = function(visitor) {
                    if (visitor != undefined && visitor.email != undefined && visitor.pass != undefined) {
                        console.log("login submit")
                        User.authenticate({email:visitor.email, pass: visitor.pass})
                            .$promise.then(function(user) {
                                if (user.status == 1 && user.content != undefined) {
                                    $window.localStorage["user"] = JSON.stringify(user.content)
                                    $location.path('/')
                                } else {
                                    $scope.errorMessage = user.content;
                                }
                            });
                    };
                }
                $scope.register = function(visitor) {
                    if (visitor != undefined && visitor.email != undefined && visitor.pass != undefined) {
                        console.log("register submit")
                        User.save({email:visitor.email, pass: visitor.pass})
                            .$promise.then(function(user) {
                                if (user.status == 1 && user.content != undefined) {
                                    $window.localStorage["user"] = JSON.stringify(user.content)
                                    $location.path('/')
                                } else {
                                    $scope.errorMessage = user.content;
                                }
                            });
                    };
                }
            }
        })
     
        .controller('LogoutCtrl', function($window, $location, $rootScope, User) {
            if ($window.localStorage['user']) {
                console.log('in logout')
                
                User.logout()
                    .$promise.then(function(user) {
                        if (user.status == 1 && user.content != undefined) {
                            console.log('success after logout')
                            $window.localStorage.removeItem("user");
                            $rootScope.user = null;
                            $location.path('/')
                        } else {
                            //$scope.errorMessage = user.content;
                            console.log('failure after logout')
                        }
                    });

            } else {
                console.log('user was logged out before');
                $location.path('/');
            }
        })

        .controller('ForgetPasswordCtrl', function($window, $location, $scope, $rootScope, User) {
            $rootScope.pageTitle = 'Forget password';
            $scope.forgetPassword = function(visitor) {
                    if (visitor != undefined && visitor.email != undefined ) {
                        console.log("forget password submit")
                        User.forgetPassword({visitor:visitor})
                            .$promise.then(function(user) {
                                if (user.status == 1 && user.content != undefined) {
                                    $location.path('/')
                                } else {
                                    $scope.errorMessage = user.content;
                                }
                            });
                    };
                }    
            
        })

        .controller('ResetPasswordCtrl', function($window, $location, $rootScope, $scope, User) {
            $rootScope.pageTitle = 'Reset password';
            var search = $location.search()

            if (search == undefined || search.token == undefined) {
                $location.path('/');
            } 

            $scope.resetPassword = function(visitor) {
                if (visitor != undefined 
                    && visitor.pass != undefined
                    && visitor.pass2 != undefined 
                    && search != undefined
                    && search.token != undefined) {
                    console.log("reset password submit")
                    User.resetPassword({visitor:visitor, token: search.token})
                        .$promise.then(function(user) {
                            if (user.status == 1 && user.content != undefined) {
                                $window.localStorage["user"] = JSON.stringify(user.content)
                                $location.path('/')
                            } else {
                                $scope.errorMessage = user.content;
                            }
                        });
                };
            }
            
        })

        .controller('SessionEditCtrl', function($window, $location, $rootScope, $scope, Movements, Session) {
            if ($window.localStorage['user']) {
                $rootScope.user = JSON.parse($window.localStorage["user"]);
                console.log(Movements)
                console.log(Session)
            } else {
                console.log('we have to authenticate that user');
                console.log($rootScope)
                $location.path('/user/login')
            }
        })
        
        .controller('SessionAddEditCtrl', function($window, $location, $rootScope, $scope, Movements, Session, SessionAllMovements) {
            if ($window.localStorage['user']) {
                $rootScope.user = JSON.parse($window.localStorage["user"]);
                console.log('movements')
                console.log(Movements);
                console.log('sessionallmovements')
                console.log(SessionAllMovements);
                console.log('session')
                console.log(Session)
                // update session
                if (Session != undefined) {
                    $rootScope.pageTitle = 'Update session';
                    $scope.sessionMovements = Session.sessionAllMovements;
                    $scope.sessionMovements[0].sessionId = Session.id;
                // add session
                } else {
                    $rootScope.pageTitle = 'Add session';
                    //$scope.sessionMovements = [{}];
                    var d = new Date();
                    $scope.sessionMovements[0].sessionDateTime = d.toString();
                    $scope.sessionMovements[0].sessionId = 0;
                    $scope.sessionAmraps = [];
                }
                $scope.movements = Movements;
                
                $scope.addSessionMovement = function() {
                    $scope.sessionMovements.push({'movement':'', 'weight': '', 'reps': '', done: false});
                }

                $scope.removeChecked = function() {
                    var oldList = $scope.sessionMovements;
                    var sessionDateTime = oldList[0].sessionDateTime;
                    var sessionId = oldList[0].sessionId;
                    $scope.sessionMovements = [];
                    angular.forEach(oldList, function(sessionMovement) {
                        if (!sessionMovement.done) {
                            $scope.sessionMovements.push(sessionMovement);
                        }
                    });
                    $scope.sessionMovements[0].sessionDateTime = sessionDateTime;
                    $scope.sessionMovements[0].sessionId = sessionId;
                }

                $scope.saveSession = function() {
                    SessionAllMovements.save({sessionAllMovements:$scope.sessionMovements})
                        .$promise.then(function(sessionAllMovements) {
                            if (sessionAllMovements.status == 1 && sessionAllMovements.content != undefined) {
                                console.log('success after save session movements')
                                
                                $location.path('/session/viewall')
                            } else {
                                //$scope.errorMessage = user.content;
                                console.log('failure after session save')
                            }
                        });
                    
                }

                $scope.addAmrap = function() {
                    //$scope.sessionAmrapMovements.push({'movement':'', 'weight': '', 'reps': '', done: false});
                    var pos = $scope.sessionAmraps.length + 1;
                    $scope.sessionAmraps.push({'name': 'Amrap ' + pos, 'pos': pos});
                }

            } else {
                console.log('we have to authenticate that user');
                console.log($rootScope)
                $location.path('/user/login')
            }
            
        })

        .controller('SessionViewallCtrl', function($window, $location, $rootScope, $scope, Sessions) {
            if ($window.localStorage['user']) {
                $rootScope.user = JSON.parse($window.localStorage["user"]);
                $rootScope.pageTitle = 'View all sessions';
                $scope.sessions = Sessions;
                
                
            } else {
                console.log('we have to authenticate that user');
                console.log($rootScope)
                $location.path('/user/login')
            }
            
        })
    }

