/**
 * Created by Yinxiong on 2014-10-22.
 */
define([
'kernel',
'angular',
'ui/Toast',
'directive/master',
'directive/client',
'directive/system',
'service/ClientService',
'service/MasterService',
'service/SocketInstance',
'ngSanitize'
],
function (core, ng, Toast){

    var App = ng.module('App', ['App.services', 'App.directives', 'ngSanitize']);

    App
    .value('Action', {
        '1': 'Add',
        '2': 'Update',
        '3': 'Delete',
        '4': 'Wait'
    })
    .value('Status', {
        Die: 0,
        Free: 1,
        Alive: 2,
        Busy: 3
    })
    .value('DeployMessage', {
        Start: 0,
        Message: 1,
        Error: 2,
        Finish: 3
    })
    .controller('masterCtrl', function ($scope, Status, Action, DeployMessage, ClientService, MasterService, SocketInstance, Helper) {

        $scope.version = {
            Version: 0,
            Time: 0
        };

        $scope.master = {
            Error: false,
            Status: Status.Busy,
            Message: ''
        };

        $scope.deployEnable = false;
        $scope.upFileList = [];
        $scope.groupList = [];

        $scope.Action = Action;
        $scope.Status = Status;

        SocketInstance.setScope( $scope );

        SocketInstance.on('master', function( data ){
            Helper.data(data).then(function( data ){
                $scope.master = data.result;
            });

            !$scope.$$parse && $scope.$apply();
            core.delay(function(){
                SocketInstance.emit('master');
            }, 5000)
        });

        SocketInstance.on('heartbeat', function( data ){
            var client, status;
            for(var i in data.result) {
                if( client = $scope.findClient( i ) ){
                    status = data.result[i];
                    if( client.status != Status.Die && status == Status.Die ){
                        $scope.notify(client, "offline", true);
                    }
                    var d = data.result[i];
                    client.status = d.status;
                    client.message = d.message;
                    client.error = d.error;
                }
            }
            !$scope.$$parse && $scope.$apply();
            core.delay(function(){
                SocketInstance.emit('heartbeat');
            }, 5000)
        });

        SocketInstance.on('deploy', function( data ){
            var id = data.Id,
                message = data.Message;

            var client = $scope.setStatus(id, Status.Busy);

            switch( data.What ){
                case DeployMessage.Start:
                    $scope.notify(client, "deploying");
                    break;
                case DeployMessage.Message:
                    $scope.notify(client, message);
                    break;
                case DeployMessage.Error:
                    $scope.notify(client, message, true);
                    break;
                case DeployMessage.End:
                    $scope.notify(client, "deploy complete", true);
                    break;
            }

            $scope.$apply();
        });

        var stillEmptyProcstat = false;
        SocketInstance.on('procstat', function( data ){
            if( $.isEmptyObject(data.result) ){
                if(!stillEmptyProcstat){
                    stillEmptyProcstat = true;
                    $scope.mapClients(function(client){
                        client.proc = {
                            cpu_percent: 0,
                            mem_percent: 0
                        }
                    })
                }
            } else {
                stillEmptyProcstat = false;
                for(var i in data.result) {
                    if( client = $scope.findClient( i ) ){
                        client.proc = data.result[i]
                    }
                }
            }
            $scope.$$parse && $scope.$apply();
            core.delay(function(){
                SocketInstance.emit('procstat');
            }, 2000)
        });

        SocketInstance.on('lock', function(){
            Toast.makeText('Lock control!').show();
            $scope.setLockControl(true);
        });

        SocketInstance.on('unlock', function(){
            $scope.setLockControl(false);
        });

        SocketInstance.on('svnup', function( data ){
            Helper.data(data).then(function( data ){
                $scope.version = data.result;
            });
        });

        SocketInstance.on('syncDeployList', function( data ){

            if( data ){
                console.log( data )
            } else {
                $scope.mapClients(function( client ){
                    $scope.notify(client, "Synchronizing file list..");
                });
            }
        });

        SocketInstance.emit('heartbeat');
        SocketInstance.emit('procstat');
        SocketInstance.emit('master');

        $scope.fillList = function( data ){
            $scope.groupList = data;
            $scope.mapClients(function(client){
                client._lock = false;
                client._error = false;
                if( client.status == Status.Die ){
                    $scope.notify(client, "connecting..")
                }
            });
        };

        $scope.refresh = function(){
            return ClientService.refresh().then(function( data ){
                $scope.fillList(data.result);
            });
        };

        ($scope.getList = function(){
            return ClientService.list().then(function( data ){
                $scope.fillList(data.result);
            });
        })();

        /**
         * 移动主机到其他组并更新UI
         *
         * @param {number} id 客户机id
         * @param {number} ex 当前组
         * @param {number} to 目标组
         */
        $scope.moveClientToGroup = function (id, ex, to) {
            var target = $scope.findGroup(to),
                org = $scope.findGroup(ex);

            if( org && target ){
                $.each(org.clients, function(index, client ){
                    if( client.id == id ){
                        delete org.clients[id];
                        target.clients[id] = client;
                        $scope.$apply();
                        return false;
                    }
                });
            }
        };

        $scope.addClientToGroup = function( client, gid ){
            var group;
            if( group = $scope.findGroup(gid) ){
                group.clients[client.id] = client;
                return true;
            }
            return false;
        };

        $scope.findGroup = function( gid ){
            if( !ng.isUndefined(gid) ){
                if( ng.isNumber(gid) && gid in $scope.groupList ){
                    return $scope.groupList[gid]
                } else if( ng.isFunction(gid) ){
                    for( var i in $scope.groupList ){
                        if( gid(i, $scope.groupList[i] ) ){
                            return $scope.groupList[i]
                        }
                    }
                }
            }
            return null;
        };

        $scope.findClient = function( id ){
            var result = null;

            if( $.isPlainObject(id) ){
                result = id;
            } else {
                $.each($scope.groupList, function(gid, group){
                    var notFound = true;
                    $.each(group.clients, function(cid, host) {
                        if( cid == id ) {
                            result = host;
                            notFound = false;
                            return false
                        }
                    });
                    return notFound;
                });
            }
            return result;
        };

        $scope.mapClients = function(fn){
            $.each($scope.groupList, function(index, group){
                $.each(group.clients, function(index, host) {
                    fn(host)
                })
            })
        };

        $scope.mapGroup = function( fn ){
            $.each($scope.groupList, function(id, group){
                fn(id, group)
            });
        };

        $scope.delClient = function(id) {
            var client = $scope.findClient(id);
            if( client ){
                delete $scope.groupList[client.Group].clients[client.Id];
            }
        };

        $scope.lockClient = function( id, lock ){
            var client = $scope.findClient(id);
            if( client ){
                client._lock = $.isUndefined(lock) ? true : lock;
            }
        };

        $scope.setStatus = function(id, status) {
            var client = $scope.findClient(id);
            if( client ){
                client.status = status;
            }

            return client
        };

        $scope.onGroupChange = function () {
            var callbacks = [], lastValues = [];
            $scope.$watch('groupList', function (newValue) {
                var arr = [];
                $.each(newValue, function(id, g){
                    arr.push({
                        value: id,
                        text: g.Name
                    })
                });
                lastValues = arr;
                callbacks.map(function (fn) {
                    fn(arr);
                })
            });

            return function (fn) {
                callbacks.push(fn);
                fn(lastValues);
            }
        }();

        MasterService.getLastVersion().then(function( data ){
            $scope.version = data.result;
        });

        $scope.isEmptyObject = function( obj ){
            return $.isEmptyObject(obj)
        };

        $scope.lockControl = false;
        $scope.setLockControl = function( enable ){
            $scope.lockControl = enable;
            $scope.$digest();
        }
    });


    ng.bootstrap(document, ['App']);
});