/**
 * Created by languid on 12/10/14.
 */

define(['kernel', 'angular', './module'],
    function (core, ng, service) {
        service.factory('ClientService', ['$http', 'Helper', function ($http, Helper) {
            return {
                refresh: function () {
                    return Helper.result(
                        $http.get('/aj/client/refresh')
                    )
                },
                edit: function (cid, data) {
                    return Helper.result(
                        $http.post('/aj/client/'+cid+'/edit', data)
                    )
                },
                add: function( data ){
                    return Helper.result(
                        $http.post('/aj/client/add', data)
                    )
                },
                del: function( id ){
                    return Helper.result(
                        $http.post('/aj/client/'+id+'/del')
                    )
                },
                list: function () {
                    return Helper.result(
                        $http.get('/aj/client/list')
                    )
                },
                heartbeat: function(){
                    return Helper.result(
                        $http.get('/aj/client/heartbeat')
                    )
                },
                revert: function(){

                },
                changeGroup: function (cid, gid) {
                    return Helper.result(
                        $http.post('/aj/client/' + cid + '/change/group/' + gid)
                    )
                },
                update: function( id, fileIds){
                    return Helper.result(
                        $http.post('/aj/client/'+id+'/update', {
                            fileIds: fileIds
                        })
                    )
                },
                deploy: function(id) {
                    return Helper.result(
                        $http.post('/aj/client/'+id+'/deploy')
                    )
                },
                getBackupList: function(id){
                    return Helper.result(
                        $http.get('/aj/client/'+id+'/backuplist')
                    )
                },
                log: function( id ){
                    return Helper.result(
                        $http.get('/aj/client/'+id+'/log')
                    )
                },
                checkClientDeployable: function( ids ){
                    return Helper.result(
                        $http.post('/aj/client/check', {
                            clientsId: ids
                        })
                    )
                }
            }
        }])
    });