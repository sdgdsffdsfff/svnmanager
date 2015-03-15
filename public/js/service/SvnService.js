/**
 * Created by languid on 11/26/14.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){
    service.factory('SvnService', ['$http', 'Helper',function( $http, Helper ){
        return {
            svnup: function () {
                return Helper.result(
                    $http.post("/aj/svn/update", {
                        paths: ['_res', 'manage']
                    })
                )
            },
            getLastVersion: function(){
                return Helper.result(
                    $http.get('/aj/svn/lastVersion')
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
            getUndeployFileList: function(){
                return Helper.result(
                    $http.get('/aj/svn/undeploy/files')
                )
            }
        }
    }])
});