/**
 * Created by Yinxiong on 2014-10-22.
 */
define([
'kernel',
'angular',
'ui/Toast',
'directive/svn',
'directive/client',
'directive/system',
'service/ClientService',
'service/SvnService',
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
        Connecting: 1,
        Alive: 2,
        Busy: 3
    })
    .controller('svnManagerCtrl', function ($scope, Status, Action, ClientService, SvnService, SocketInstance, Helper) {

        $scope.version = {
            Version: 0,
            Time: 0
        };

        $scope.deployEnable = false;
        $scope.upFileList = [];
        $scope.groupList = [];

        $scope.ACTION = Action;
        $scope.STATUS = Status;

        SocketInstance.setScope( $scope );

        SocketInstance.on('getClientList', function(data){
            Helper.data(data).then(function( data ){
                $scope.groupList = data.result;
                $scope.$apply();
            })
        });

        SocketInstance.on('heartbeat', function( data ){
            var client;
            for(var i in data.result) {
                if( client = $scope.findClient( i ) ){
                    client.Status = data.result[i]
                }
            }
            $scope.$$parse && $scope.$apply();
            core.delay(function(){
                SocketInstance.emit('heartbeat');
            }, 5000)
        });

        SocketInstance.on('deploy', function( data ){
            console.log(data)
        });

        var stillEmptyProcstat = false;
        SocketInstance.on('procstat', function( data ){
            if( $.isEmptyObject(data.result) ){
                if(!stillEmptyProcstat){
                    stillEmptyProcstat = true;
                    $scope.mapClients(function(client){
                        client.Proc = {
                            CPUPercent: 0,
                            MEMPercent: 0
                        }
                    })
                }
            } else {
                stillEmptyProcstat = false;
                for(var i in data.result) {
                    if( client = $scope.findClient( i ) ){
                        client.Proc = data.result[i]
                    }
                }
            }
            $scope.$$parse && $scope.$apply();
            core.delay(function(){
                SocketInstance.emit('procstat');
            }, 5000)
        });

        SocketInstance.on('lock', function(){
            Toast.makeText('Lock control!').show();
            $scope.setLockControl(true);
        });

        SocketInstance.on('unlock', function(){
            Toast.makeText('Unlock control!').show();
            $scope.setLockControl(false);
        });

        SocketInstance.on('svnup', function( data ){
            Helper.data(data).then(function( data ){
                var version = data.result;
                $scope.upgradeVersion({
                    Version: version.Version,
                    Time: version.Time
                });
            });
        });

        SocketInstance.emit('heartbeat');
        SocketInstance.emit('procstat');
        SocketInstance.emit('getClientList');

        $scope.upgradeVersion = function( version ){
            $scope.version = version;
        };

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
                $.each(org.Clients, function(index, client ){
                    if( client.Id == id ){
                        delete org.Clients[id];
                        target.Clients[id] = client;
                        $scope.$apply();
                        return false;
                    }
                });
            }
        };

        $scope.addClientToGroup = function( client, gid ){
            var group;
            if( group = $scope.findGroup(gid) ){
                group.Clients[client.Id] = client;
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
            $.each($scope.groupList, function(gid, group){
                var notFound = true;
                $.each(group.Clients, function(cid, host) {
                    if( cid == id ) {
                        result = host;
                        notFound = false;
                        return false
                    }
                });
                return notFound;
            });
            return result;
        };

        $scope.mapClients = function(fn){
            $.each($scope.groupList, function(index, group){
                $.each(group.Clients, function(index, host) {
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
                delete $scope.groupList[client.Group]["Clients"][client.Id];
            }
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

        SvnService.getLastVersion().then(function( data ){
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
    })


    ng.bootstrap(document, ['App']);
});