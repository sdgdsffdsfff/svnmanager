/**
 * Created by Yinxiong on 2014-11-05.
 */

define([
'kernel',
'angular',
'./module',
'service/Websocket'
],
function( core, ng, service){

    service.factory('SocketInstance', function( Websocket, Helper ) {
        var socket = Websocket();
        socket.scope = null;
        socket.setScope = function( s ){
            this.scope = s;
        };

        socket.on('broadcast', function( data ){
            console.log( data.data.Version )
        });
        return socket.listen();
    })
});