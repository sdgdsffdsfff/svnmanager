/**
 * Created by Yinxiong on 2014-10-22.
 */
define([
'kernel',
'angular',
'directive/svn',
'directive/group',
'service/ClientService',
'service/SvnService',
'service/SocketInstance',
'ngSanitize'
],
function (core, ng) {

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

        SocketInstance.setScope( $scope );
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

        SocketInstance.on('procstat', function( data ){
            for(var i in data.result) {
                if( client = $scope.findClient( i ) ){
                    client.Proc = data.result[i]
                }
            }
            $scope.$$parse && $scope.$apply();
            core.delay(function(){
                SocketInstance.emit('procstat');
            }, 5000)
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


        $scope.version = {
            Version: 0,
            Time: 0
        };
        $scope.deployEnable = false;
        $scope.upFileList = [];
        $scope.groupList = [];

        $scope.ACTION = Action;
        $scope.STATUS = Status;

        ($scope.refresh = function(){
            return ClientService.list().then(function (data) {
                $scope.groupList = data.result;
            });
        })();

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
                        org.Clients.splice(index, 1);
                        target.Clients.push(client);
                        $scope.$apply();
                        return false;
                    }
                });
            }
        };

        $scope.addClientToGroup = function( client, gid ){
            var group;
            if( group = $scope.findGroup(gid) ){
                group.Clients.push(client);
                return true;
            }
            return false;
        };

        $scope.findGroup = function( gid ){
            var result = null;
            if( !ng.isUndefined(gid) ){
                $.each($scope.groupList, function(index, group){
                    if( ng.isNumber(gid) ){
                        if( group.Id == gid ){
                            result = group;
                            return false;
                        }
                    }else if( ng.isFunction(gid) ){
                        if( gid(index, group) ){
                            result = group;
                            return false;
                        }
                    }
                });
            }
            return result;
        };

        $scope.findClientFromGroup = function( cid, gid ){
            var client = null,
                group = ng.isNumber(gid) ? $scope.findGroup(gid) :
                    ng.isObject(gid) ? gid : null;

            if( group ){
                $.each(group.Clients, function(index, host) {
                    if( host.Id == cid){
                        client = host;
                        return false;
                    }
                })
            }

            return client;
        };

        $scope.findClient = function( cid ){
            var client = null, host;
            $scope.findGroup(function(index, group){
                if( host = $scope.findClientFromGroup( cid, group ) ){
                    client = host;
                    return true;
                }
            });
            return client;
        };

        $scope.mapClients = function(fn){
            $.each($scope.groupList, function(index, group){
                $.each(group.Clients, function(index, host) {
                    fn(host)
                })
            })
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
        })
    });


    ng.bootstrap(document, ['App']);
});