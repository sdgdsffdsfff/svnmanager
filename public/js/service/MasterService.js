/**
 * Created by languid on 11/26/14.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){
    service.factory('MasterService', ['$http', 'Helper',function( $http, Helper ){
        return {
            svnup: function () {
                return Helper.result(
                    $http.post("/aj/update", {
                        paths: ['_res', 'manage']
                    })
                )
            },
            getLastVersion: function(){
                return Helper.result(
                    $http.get('/aj/lastVersion')
                );
            },
            deploy: function( filesId, clientsId, msg ){

                return Helper.result(
                    $http.post('/aj/deploy', {
                        filesId: filesId,
                        clientsId: clientsId,
                        message: msg
                    })
                )
            },
            compile: function(){
                return Helper.result(
                    $http.post('/aj/compile')
                )
            },
            getError: function(){
                return Helper.result(
                    $http.get('/aj/error')
                )
            },
            getUndeployFileList: function(){
                return Helper.result(
                    $http.get('/aj/undeploy/files')
                )
            }
        }
    }])
});