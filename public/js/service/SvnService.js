/**
 * Created by languid on 11/26/14.
 */

define(['kernel', 'angular', './module'],
function( core, ng, service){
    service.factory('SvnService', ['$http', 'Helper',function( $http, Helper ){
        return {
            svnup: function () {
                return Helper.result(
                    $http.post("/aj/svn/up", {
                        paths: ['_res', 'manage']
                    })
                )
            },
            getLastVersion: function(){
                return Helper.result(
                    $http.get('/aj/svn/lastVersion')
                );
            },
            deploy: function( files ){
                return Helper.result(
                    $http.post('/aj/deploy', {
                        files: files || []
                    })
                )
            },
            getUndeployFileList: function(){
                return Helper.result(
                    $http.get('/aj/svn/undeploy/files')
                )
            },
            getDeployLock: function(){
                return Helper.result(
                    $http.get('/aj/deploy/lock')
                )
            }
        }
    }])
});