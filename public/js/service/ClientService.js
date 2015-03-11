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
                update: function (cid, data) {
                    return Helper.result(
                        $http.post('/aj/client/'+cid+'/update', data)
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
                changeGroup: function (cid, gid) {
                    return Helper.result(
                        $http.post('/aj/client/' + cid + '/change/group/' + gid)
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