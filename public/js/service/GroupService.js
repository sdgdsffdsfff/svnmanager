/**
 * Created by languid on 12/10/14.
 */

define(['kernel', 'angular', './module'],
    function (core, ng, service) {
        service.factory('GroupService', function ($http, Helper) {
            return {
                add: function (name) {
                    return Helper.result(
                        $http.post('/aj/group/add', {
                            name: name
                        })
                    )
                },
                edit: function (id, name) {
                    return Helper.result(
                        $http.post('/aj/group/edit/' + id, {
                            name: name
                        })
                    )
                }
            }
        })
    });